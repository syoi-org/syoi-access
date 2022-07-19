/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/syoi-org/syoi-access/internal/app"
)

var config app.Config

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Connect to SSH services",
	Long: `The ssh command proxies SSH connection to internal servers.

This command can act as the ProxyCommand in SSH config. You will be prompted to
login to cloudflare and the SSH connection will be established afterwards.`,
	Run: func(*cobra.Command, []string) {
		fmt.Println("ssh called")
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)

	sshCmd.Flags().StringVar(&config.Hostname, "hostname", "", "specifiy the hostname of the server to connect to")
	sshCmd.MarkFlagRequired("hostname")
	sshCmd.Flags().IntVar(&config.LocalBindPort, "local-bind-port", 0, "specify the port of cloudflared binds to")
}
