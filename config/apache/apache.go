package apache

import "fmt"

// ApacheConfig implements the WebServerConfig interface for Apache
type ApacheConfig struct{}

// Generate creates an Apache virtual host configuration
func (a *ApacheConfig) Generate(domain, port string) (string, error) {
    vhostTemplate := `<VirtualHost *:%s>
    ServerName %s
    ServerAlias www.%s
    DocumentRoot /var/www/%s
    ErrorLog ${APACHE_LOG_DIR}/%s-error.log
    CustomLog ${APACHE_LOG_DIR}/%s-access.log combined
    <Directory /var/www/%s>
        Options Indexes FollowSymLinks
        AllowOverride All
        Require all granted
    </Directory>
</VirtualHost>`

    config := fmt.Sprintf(vhostTemplate, port, domain, domain, domain, domain, domain, domain)
    return config, nil
}

// GetFileExtension returns the file extension for Apache config files
func (a *ApacheConfig) GetFileExtension() string {
    return ".conf"
}