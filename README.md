# KSSH

Kssh is a terminal-based SSH manager with a slick, keyboard-driven interface. It lets you quickly browse and connect to your SSH servers using a simple YAML config file — no more memorizing IPs or typing long commands.


## Features

- Simple YAML-based configuration
- Interactive TUI with keyboard navigation
- Vim-style navigation
- Support for custom SSH keys, ports, and usernames
- Automatic configuration generation
- Colorful and intuitive interface

## Installation

### Prerequisites

- Go 1.16 or higher
- SSH client installed on your system

### Installing from source

1. Clone the repository:
```bash
git clone https://github.com/tergel-sama/kssh.git 
cd kssh
```

2. Build the application:
```bash
go build -o kssh cmd/main.go
```

3. Move the binary to your PATH:
```bash
sudo mv kssh /usr/local/bin/
```

### Using Go Install

not tested 
```bash
go install github.com/tergel-sama/kssh@latest
```

## Usage

Simply run the application:

```bash
kssh
```

On first run, the application will create a default configuration file at `~/.ssh_hosts.yaml` if it doesn't exist.

### Navigation

- Use `↑/↓` arrow keys or `j/k` to select a host
- Press `Enter` to connect to the selected host
- Press `q` or `Ctrl+C` to quit

## Configuration

The configuration file is located at `~/.ssh_hosts.yaml` and uses YAML format.

### Example Configuration

```yaml
hosts:
  - name: production-server
    hostname: example.com
    user: admin
    port: 22
    key: ~/.ssh/id_rsa
    
  - name: staging-server
    hostname: staging.example.com
    user: developer
    port: 2222
    
  - name: raspberry-pi
    hostname: 192.168.1.100
    user: pi
    key: ~/.ssh/pi_key
```

### Configuration Fields

Each host entry supports the following fields:

| Field     | Description                           | Required | Default    |
|-----------|---------------------------------------|----------|------------|
| `name`    | Display name for the connection       | Yes      | N/A        |
| `hostname`| Server hostname or IP address         | Yes      | N/A        |
| `user`    | SSH username                          | Yes      | N/A        |
| `port`    | SSH port                              | No       | 22         |
| `key`     | Path to SSH private key               | No       | None       |

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
