# StackRoost CLI

StackRoost is a powerful command-line tool to manage Linux servers with ease. It supports domain configuration, user management, SSL setup, logs, monitoring, and more.

## Version

v1.0.0

## Features

### Domain Management
- `create-domain` – Create a new domain config with Apache/Nginx/Caddy.
- `backup-domain` – Backup public_html and MySQL DB.
- `clone-domain` – Clone full domain configuration and data.
- `list-domains` – List all active domains and statuses.
- `monitor` – Interactive TUI to monitor all domains.
- `remove-domain` – Remove domain, user, database, and config.
- `restore-domain` – Restore from domain backup archive.
- `status-domain` – Inspect domain config, SSL, and user.
- `toggle-site` – Enable or disable a site's config.
- `update-domain-port` – Update the domain port and reload the web server.

### Email
- `test-email` – Check if the server can send mail (mail/sendmail/msmtp).

### Firewall
- `enable-firewall` – Enable UFW and allow common/custom ports.
- `disable-firewall` – Safely disable UFW.

### Logs
- `analyze-log-traffic` – Analyze access log traffic (IP, URL, requests).
- `logs-domain` – View domain logs (access and error).
- `purge-domain-logs` – Delete domain logs safely.

### Security
- `run-security-check` – Run a server hardening security check.
- `secure-server` – Enable UFW, SSH restrictions, and config hardening.

### Server Tools
- `check-port` – Check if a domain's port is open.
- `server-health` – View CPU, RAM, disk usage, uptime, web server status.
- `inspect-config` – View web server config file.
- `restart-server` – Restart Apache/Nginx/Caddy.
- `schedule-restart` – Schedule a restart after delay.
- `sync-time` – Sync time using systemd-timesyncd.

### SSL Certificates
- `enable` – Enable Let's Encrypt SSL for a domain.
- `disable` – Remove SSL config and certs.
- `renew` – Renew SSL certificates.
- `expiry` – Check SSL expiry date.
- `test` – Test domain SSL cert status.

### User Management
- `list-users` – List shell users (UID ≥ 1000).
- `delete-user` – Delete a system user and their home directory.

## Install

Clone and build manually:
```bash
git clone https://github.com/stackroost/stackroost-cli.git
cd stackroost-cli
go build -o stackroost main.go
```

## Usage

```bash
./stackroost --help
./stackroost create-domain --name example.com --server nginx --shelluser --pass mypass --useridr --ssl
./stackroost monitor
```

## License

MIT License