package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

var inspectConfigCmd = &cobra.Command{
	Use:   "inspect-config",
	Short: "View the web server configuration file of a domain",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		if internal.IsNilOrEmpty(domain) {
			logger.Error("Please provide a domain using --domain")
			os.Exit(1)
		}

		serverType := internal.DetectServerType(domain)
		if serverType == "" {
			logger.Error("Could not detect server type. No config file found.")
			os.Exit(1)
		}

		var configPath string
		switch serverType {
		case "apache":
			configPath = filepath.Join("/etc/apache2/sites-available", domain+".conf")
		case "nginx":
			configPath = filepath.Join("/etc/nginx/sites-available", domain+".conf")
		case "caddy":
			configPath = filepath.Join("/etc/caddy/sites-available", domain+".conf")
		default:
			logger.Error("Unsupported server type")
			os.Exit(1)
		}

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			logger.Error(fmt.Sprintf("Config file not found at: %s", configPath))
			os.Exit(1)
		}

		logger.Info(fmt.Sprintf("Showing config: %s", configPath))
		internal.RunCommand("sudo", "cat", configPath)
	},
}

func init() {
	rootCmd.AddCommand(inspectConfigCmd)
	inspectConfigCmd.Flags().String("domain", "", "Domain name to inspect config for")
	inspectConfigCmd.MarkFlagRequired("domain")
}
