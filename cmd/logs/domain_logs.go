package logs

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

// GetDomainLogsCmd returns the CLI command for viewing domain logs
func GetDomainLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs-domain",
		Short: "View recent access and error logs for a domain",
		Run: func(cmd *cobra.Command, args []string) {
			domain, _ := cmd.Flags().GetString("domain")
			lines, _ := cmd.Flags().GetInt("lines")

			if internal.IsNilOrEmpty(domain) {
				logger.Error("Please provide a domain using --domain")
				os.Exit(1)
			}

			if lines <= 0 {
				lines = 50
			}

			server := internal.DetectServerType(domain)
			if server == "" {
				logger.Error("Could not detect server type for the domain")
				os.Exit(1)
			}

			var accessLogPath, errorLogPath string
			switch server {
			case "apache":
				accessLogPath = fmt.Sprintf("/var/log/apache2/%s-access.log", domain)
				errorLogPath = fmt.Sprintf("/var/log/apache2/%s-error.log", domain)
			case "nginx":
				accessLogPath = fmt.Sprintf("/var/log/nginx/%s-access.log", domain)
				errorLogPath = fmt.Sprintf("/var/log/nginx/%s-error.log", domain)
			case "caddy":
				accessLogPath = fmt.Sprintf("/var/log/caddy/%s-access.log", domain)
				errorLogPath = "" // Caddy doesn't have default error logs per domain
			default:
				logger.Error("Unsupported server type")
				os.Exit(1)
			}

			// Access log
			if stat, err := os.Stat(accessLogPath); err == nil {
				logger.Info(fmt.Sprintf("Access Log (%s):", accessLogPath))
				if stat.Size() == 0 {
					logger.Warn("Access log is empty")
				} else {
					internal.RunCommand("sudo", "tail", "-n", fmt.Sprintf("%d", lines), accessLogPath)
				}
			} else {
				logger.Warn("Access log not found")
			}

			// Error log
			if errorLogPath != "" {
				if stat, err := os.Stat(errorLogPath); err == nil {
					logger.Info(fmt.Sprintf("Error Log (%s):", errorLogPath))
					if stat.Size() == 0 {
						logger.Warn("Error log is empty")
					} else {
						internal.RunCommand("sudo", "tail", "-n", fmt.Sprintf("%d", lines), errorLogPath)
					}
				} else {
					logger.Warn("Error log not found")
				}
			} else if server == "caddy" {
				logger.Info("Caddy does not maintain a separate error log by default")
			}
		},
	}

	cmd.Flags().String("domain", "", "Domain name to view logs for")
	cmd.Flags().Int("lines", 50, "Number of lines to show from each log file")
	cmd.MarkFlagRequired("domain")

	return cmd
}
