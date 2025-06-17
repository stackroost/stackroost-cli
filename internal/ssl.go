package internal

import (
	"fmt"
	"stackroost/internal/logger"
)

func EnableSSLCertbot(domain string, serverType string) error {
	logger.Info(fmt.Sprintf("Requesting SSL certificate for %s using Certbot", domain))

	cmd := []string{
		fmt.Sprintf("--%s", serverType),
		"-d", domain,
		"-d", "www." + domain,
		"--non-interactive",
		"--agree-tos",
		"--register-unsafely-without-email", 
	}

	err := RunCommand("sudo", append([]string{"certbot"}, cmd...)...)
	if err != nil {
		return fmt.Errorf("certbot SSL generation failed: %v", err)
	}

	logger.Success(fmt.Sprintf("SSL certificate installed for %s", domain))
	return nil
}