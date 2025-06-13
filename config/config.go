package config

import (
    "fmt"
    "stackroost/config/apache"
    "stackroost/config/nginx"
    "stackroost/config/caddy"
)

type WebServerConfig interface {
    Generate(domain, port, username string) (string, error)
    GetFileExtension() string
}

func NewWebServerConfig(serverType string) (WebServerConfig, error) {
    switch serverType {
    case "apache":
        return &apache.ApacheConfig{}, nil
    case "nginx":
		return &nginx.NginxConfig{}, nil
    case "caddy":
        return &caddy.CaddyConfig{}, nil
    default:
        return nil, fmt.Errorf("unsupported web server type: %s", serverType)
    }
}
