# syoi-access

syoi-access is a CLI tool to access internal services of syoi.org. It serves as
a proxy to internal services. For example, it can proxy SSH connections to a
specific instance of code.syoi.org.

## Installation

### Binary Distribution

The appliation is a single binary file. You may use those compiled in the
Release page if the matches your operating system and architecture.

### Go install

You may use the go install command to download, compile and install the
application.

```bash
go install github.com/syoi-org/syoi-access@latest
```

## Getting Started

### Setting up SSH access

Please get the hostname of your SSH hostname from your network administrator. It
should be something like `alice.ssh.syoi.org`. Append the following to your SSH
config `~/.ssh/config`.

```ssh
Host alice.ssh.syoi.org
  ProxyCommand syoi-access ssh --hostname %h
```
