package domain

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

func GetStatusCmd() *cobra.Command {
	var domain string

	cmd := &cobra.Command{
		Use:   "status-domain",
		Short: "Inspect the configuration, user, and SSL status of a domain",
		Run: func(cmd *cobra.Command, args []string) {
			if internal.IsNilOrEmpty(domain) {
				logger.Error("Please provide a domain using --domain")
				os.Exit(1)
			}

			logger.Info(fmt.Sprintf("Inspecting domain: %s", domain))

			serverType := internal.DetectServerType(domain)
			if serverType == "" {
				logger.Warn("Could not detect server type (no config found)")
			} else {
				logger.Info(fmt.Sprintf("Server: %s", serverType))
			}

			// Check if site is enabled
			enabled := false
			switch serverType {
			case "apache":
				out := internal.CaptureCommand("a2query", "-s", domain)
				enabled = strings.Contains(out, "is enabled")
			case "nginx":
				linkPath := filepath.Join("/etc/nginx/sites-enabled", domain+".conf")
				_, err := os.Stat(linkPath)
				enabled = err == nil
			case "caddy":
				linkPath := filepath.Join("/etc/caddy/sites-enabled", domain+".conf")
				_, err := os.Stat(linkPath)
				enabled = err == nil
			}

			if enabled {
				logger.Info("Status: ENABLED")
			} else {
				logger.Info("Status: DISABLED")
			}

			// Shell user check
			username := strings.Split(domain, ".")[0]
			if _, err := user.Lookup(username); err == nil {
				logger.Info(fmt.Sprintf("Shell User: %s ✔", username))
			} else {
				logger.Warn(fmt.Sprintf("Shell User: %s", username))
			}

			// public_html check
			htmlPath := fmt.Sprintf("/home/%s/public_html", username)
			if _, err := os.Stat(htmlPath); err == nil {
				logger.Info(fmt.Sprintf("Public HTML: %s ✔", htmlPath))
			} else {
				logger.Warn(fmt.Sprintf("Public HTML: %s", htmlPath))
			}

			// SSL cert status
			if serverType == "caddy" {
				logger.Info("SSL Certificate: Handled automatically by Caddy")
			} else {
				out := internal.CaptureCommand("sudo", "certbot", "certificates", "--cert-name", domain)
				if strings.Contains(out, domain) {
					logger.Info("SSL Certificate: Installed via Certbot")
				} else {
					logger.Warn("SSL Certificate: Not found")
				}
			}
		},
	}

	cmd.Flags().StringVar(&domain, "domain", "", "Domain name to inspect")
	cmd.MarkFlagRequired("domain")
	return cmd
}
