# Stackroost CLI

A powerful command-line interface for managing web servers, domains, SSL certificates, user access, and multi-server deployments across local and remote systems.

## Features

- **Domain Management**: Add, list, remove, enable, and disable domains with virtual host configurations for Apache, Nginx, and Caddy
- **Web Server Management**: Control Apache, Nginx, and Caddy services (start, stop, reload, status)
- **SSL Certificate Management**: Issue, renew, revoke SSL certificates via Let's Encrypt or upload manual certificates
- **User Management**: Create, modify, and manage system users with SSH access control
- **Remote Server Management**: Add remote servers and execute commands via SSH
- **Log Monitoring**: View real-time server logs (access and error logs)

## Installation

### Prerequisites

- Go 1.25.1 or later
- sudo access on target systems
- Web servers (Apache/Nginx/Caddy) installed on target systems
- Certbot for SSL certificate management

### Build from Source

```bash
git clone https://github.com/stackroost/stackroost-cli.git
cd stackroost-cli
go build -o stackroost .
sudo mv stackroost /usr/local/bin/
```

### Download Binary

Download the latest release from the [releases page](https://github.com/stackroost/stackroost-cli/releases) and make it executable:

```bash
chmod +x stackroost
sudo mv stackroost /usr/local/bin/
```

## Configuration

Stackroost uses a YAML configuration file stored at `~/.stackroost.yaml`. The configuration is automatically created and updated as you use the CLI.

You can also specify a custom config file:

```bash
stackroost --config /path/to/config.yaml
```

## Usage

### Domain Management

#### Add a new domain
```bash
stackroost domain add example.com --server apache
```

#### List domains
```bash
stackroost domain list
```

#### Remove a domain
```bash
stackroost domain remove example.com
```

#### Enable/Disable domains
```bash
stackroost domain enable example.com
stackroost domain disable example.com
```

#### Set document root
```bash
stackroost domain set-root example.com /var/www/example.com/public
```

### Web Server Management

#### List installed servers
```bash
stackroost server list
```

#### Control servers
```bash
stackroost server start apache
stackroost server stop nginx
stackroost server reload caddy
```

#### Check server status
```bash
stackroost server status
```

### SSL Certificate Management

#### Issue SSL certificate (Let's Encrypt)
```bash
stackroost ssl issue example.com --email admin@example.com
```

#### Renew SSL certificate
```bash
stackroost ssl renew example.com
```

#### Revoke SSL certificate
```bash
stackroost ssl revoke example.com
```

#### Upload manual SSL certificate
```bash
stackroost ssl upload example.com --cert /path/to/cert.pem --key /path/to/key.pem
```

#### List SSL certificates
```bash
stackroost ssl list
```

### User Management

#### Add a new user
```bash
stackroost user add john --password mypassword
```

#### List users
```bash
stackroost user list
```

#### Change user password
```bash
stackroost user passwd john
```

#### Remove a user
```bash
stackroost user remove john
```

#### Enable/Disable SSH access
```bash
stackroost user ssh-enable john
stackroost user ssh-disable john
```

### Remote Server Management

#### Add a remote server
```bash
stackroost remote add myserver user@192.168.1.100 --key ~/.ssh/id_rsa
```

#### List remote servers
```bash
stackroost remote list
```

#### Execute command on remote server
```bash
stackroost remote exec myserver "sudo apt update && sudo apt upgrade"
```

### Log Monitoring

#### View server logs
```bash
stackroost logs --server apache --type access
stackroost logs --server nginx --type error
```

## Supported Web Servers

- **Apache**: Full virtual host management, SSL integration
- **Nginx**: Virtual host configuration, SSL support
- **Caddy**: Basic virtual host setup, automatic SSL (limited CLI control)

## Security Notes

- SSL certificates are managed via Certbot (Let's Encrypt)
- SSH connections use key-based authentication
- All server management commands require sudo privileges
- Configuration files are stored securely in the user's home directory

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This tool requires administrative privileges and can modify system configurations. Use with caution and ensure you have backups of your configurations before making changes.