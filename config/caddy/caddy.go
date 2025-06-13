package caddy

import "fmt"

type CaddyConfig struct{}

func (c *CaddyConfig) Generate(domain, port, username string) (string, error) {
	config := `
%s:%s {
    root * /home/%s/public_html
    file_server
    encode gzip

    php_fastcgi unix//run/php/php8.1-fpm.sock

    log {
        output file /var/log/caddy/%s-access.log
        format single_field common_log
    }

    handle_errors {
        respond "Something went wrong" 500
    }
}
`
	return fmt.Sprintf(config, domain, port, username, domain), nil
}

func (c *CaddyConfig) GetFileExtension() string {
	return ".conf"
}
