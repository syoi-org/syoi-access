package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

type Config struct {
	Hostname      string
	LocalBindPort int
}

func getListenAddr() (net.Addr, error) {
	// create listener on random available port on localhost
	listener, err := net.Listen("tcp", "[::1]:0")
	if err != nil {
		return nil, err
	}
	defer listener.Close()
	return listener.Addr(), nil
}

func dialWithRetry(ctx context.Context, network, addr string) (net.Conn, error) {
	var dialer net.Dialer
	// retry until context is cancelled
	for {
		conn, err := dialer.DialContext(ctx, network, addr)
		if err == nil {
			return conn, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(1 * time.Second): // retry after 1 second
		}
	}
}

func Ssh(config *Config) error {
	var addr net.Addr
	var err error

	if config.LocalBindPort == 0 {
		// Get the listen address which is a random available port on localhost
		addr, err = getListenAddr()
		if err != nil {
			return fmt.Errorf("failed to get listen address: %w", err)
		}
	} else {
		// Listen on the specified port
		addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("[::1]:%d", config.LocalBindPort))
		if err != nil {
			return fmt.Errorf("failed to resolve listen address: %w", err)
		}
	}

	// Run cloudflare tunnel in background
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Create the Cloudflare tunnel execute command
	cmd := exec.CommandContext(ctx, "cloudflared", "access", "ssh", "--hostname", "ssh.syoi.org", "--url", fmt.Sprintf("%s://%s", addr.Network(), addr.String()))
	log.Printf("Running cloudflared command: %s", strings.Join(cmd.Args, " "))

	// Setup pipes to capture stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	defer stdout.Close()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	defer stderr.Close()

	// Start the Cloudflare tunnel
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start cloudflared: %w", err)
	}

	// Copy stdout and stderr to the console
	go func() {
		io.Copy(os.Stderr, stdout)
	}()
	go func() {
		io.Copy(os.Stderr, stderr)
	}()

	// try to acquire tcp connection to listening port of Cloudflare tunnel
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	conn, err := dialWithRetry(dialCtx, addr.Network(), addr.String())
	if err != nil {
		return fmt.Errorf("failed to dial to listening port: %w", err)
	}

	// establish TLS connection
	tlsConn := tls.Client(conn, &tls.Config{ServerName: config.Hostname})
	defer tlsConn.Close()
	if err := tlsConn.HandshakeContext(ctx); err != nil {
		return fmt.Errorf("failed to establish TLS connection: %w", err)
	}

	// pipe stdin and stdout to the TLS connection
	go func() {
		io.Copy(os.Stdout, tlsConn)
	}()
	go func() {
		io.Copy(tlsConn, os.Stdin)
		stop() // stop after stdin is closed
	}()

	<-ctx.Done()

	return nil
}
