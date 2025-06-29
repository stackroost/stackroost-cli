package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

var purgeDomainLogsCmd = &cobra.Command{
	Use:   "purge-domain-logs",
	Short: "Delete access and error logs for a specific domain",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		if internal.IsNilOrEmpty(domain) {
			logger.Error("Please provide a domain using --domain")
			os.Exit(1)
		}

		server := internal.DetectServerType(domain)
		if server == "" {
			logger.Error("Could not detect server type for domain: " + domain)
			os.Exit(1)
		}

		var accessLog, errorLog string

		switch server {
		case "apache":
			accessLog = filepath.Join("/var/log/apache2", domain+"-access.log")
			errorLog = filepath.Join("/var/log/apache2", domain+"-error.log")
		case "nginx":
			accessLog = filepath.Join("/var/log/nginx", domain+"-access.log")
			errorLog = filepath.Join("/var/log/nginx", domain+"-error.log")
		case "caddy":
			logger.Warn("Caddy does not maintain traditional per-domain logs")
			return
		default:
			logger.Error("Unsupported server type: " + server)
			return
		}

		for _, logFile := range []string{accessLog, errorLog} {
			if _, err := os.Stat(logFile); err == nil {
				logger.Info("Deleting log: " + logFile)
				if err := internal.RunCommand("sudo", "rm", "-f", logFile); err != nil {
					logger.Error("Failed to delete log file: " + err.Error())
				} else {
					logger.Success("Deleted: " + logFile)
				}
			} else {
				logger.Warn("Log not found: " + logFile)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(purgeDomainLogsCmd)
	purgeDomainLogsCmd.Flags().String("domain", "", "Domain name to purge logs for")
	purgeDomainLogsCmd.MarkFlagRequired("domain")
}
