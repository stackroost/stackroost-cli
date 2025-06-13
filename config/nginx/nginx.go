package nginx

import "fmt"

type NginxConfig struct{}

func (n *NginxConfig) Generate(domain, port, username string) (string, error) {
	config := `
server {
    listen %s;
    server_name %s www.%s;

    root /home/%s/public_html;
    index index.html index.htm index.php;

    access_log /var/log/nginx/%s-access.log;
    error_log /var/log/nginx/%s-error.log;

    location / {
        try_files $uri $uri/ =404;
    }

    location ~ \.php$ {
        include snippets/fastcgi-php.conf;
        fastcgi_pass unix:/run/php/php8.1-fpm.sock; # Update as per your PHP version
    }

    location ~ /\.ht {
        deny all;
    }
}
`
	return fmt.Sprintf(config, port, domain, domain, username, domain, domain), nil
}

func (n *NginxConfig) GetFileExtension() string {
	return ".conf"
}
