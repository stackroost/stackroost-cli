package config

import (
    "fmt"
    "stackroost/config/apache"
)

// WebServerConfig defines the interface for generating web server configurations
type WebServerConfig interface {
    Generate(domain, port string) (string, error)
    GetFileExtension() string
}

// NewWebServerConfig creates a new configuration generator based on the server type
func NewWebServerConfig(serverType string) (WebServerConfig, error) {
    switch serverType {
    case "apache":
        return &apache.ApacheConfig{}, nil
    default:
        return nil, fmt.Errorf("unsupported web server type: %s", serverType)
    }
}