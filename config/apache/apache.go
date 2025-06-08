package apache

import "fmt"

// ApacheConfig implements the WebServerConfig interface for Apache
type ApacheConfig struct{}

// Generate creates an Apache virtual host configuration
func (a *ApacheConfig) Generate(domain, port, username string) (string, error) {
    vhostTemplate := `<VirtualHost *:%s>
    ServerName %s
    ServerAlias www.%s
    DocumentRoot /home/%s/public_html
    ErrorLog ${APACHE_LOG_DIR}/%s-error.log
    CustomLog ${APACHE_LOG_DIR}/%s-access.log combined
    <Directory /home/%s/public_html>
        Options Indexes FollowSymLinks
        AllowOverride All
        Require all granted
    </Directory>
</VirtualHost>`

    return fmt.Sprintf(vhostTemplate, port, domain, domain, username, domain, domain, username), nil
}


// GetFileExtension returns the file extension for Apache config files
func (a *ApacheConfig) GetFileExtension() string {
    return ".conf"
}