package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

var updatePortCmd = &cobra.Command{
	Use:   "update-domain-port",
	Short: "Update the port for a domain and reload the web server",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		newPort, _ := cmd.Flags().GetString("port")

		if internal.IsNilOrEmpty(domain) || internal.IsNilOrEmpty(newPort) {
			logger.Error("Both --domain and --port are required")
			os.Exit(1)
		}

		server := internal.DetectServerType(domain)
		if server == "" {
			logger.Error("Could not detect server type (no config found)")
			os.Exit(1)
		}

		var configPath string
		switch server {
		case "apache":
			configPath = filepath.Join("/etc/apache2/sites-available", domain+".conf")
		case "nginx":
			configPath = filepath.Join("/etc/nginx/sites-available", domain+".conf")
		case "caddy":
			configPath = filepath.Join("/etc/caddy/sites-available", domain+".conf")
		}

		logger.Info(fmt.Sprintf("Updating port in %s configuration", server))

		content, err := os.ReadFile(configPath)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to read config file: %v", err))
			os.Exit(1)
		}

		// Backup
		backupPath := configPath + ".bak"
		_ = os.WriteFile(backupPath, content, 0644)
		logger.Info(fmt.Sprintf("Backup created: %s", backupPath))

		updated := strings.ReplaceAll(string(content), ":80", ":"+newPort)
		if err := os.WriteFile(configPath, []byte(updated), 0644); err != nil {
			logger.Error(fmt.Sprintf("Failed to update config file: %v", err))
			os.Exit(1)
		}

		logger.Success(fmt.Sprintf("Port updated to %s for domain %s", newPort, domain))

		// Reload
		switch server {
		case "apache":
			internal.RunCommand("sudo", "systemctl", "reload", "apache2")
		case "nginx":
			internal.RunCommand("sudo", "systemctl", "reload", "nginx")
		case "caddy":
			internal.RunCommand("sudo", "systemctl", "reload", "caddy")
		}

		logger.Success(fmt.Sprintf("%s server reloaded successfully", strings.ToUpper(server)))
	},
}

func init() {
	rootCmd.AddCommand(updatePortCmd)
	updatePortCmd.Flags().String("domain", "", "Domain name to update")
	updatePortCmd.Flags().String("port", "", "New port number")
	updatePortCmd.MarkFlagRequired("domain")
	updatePortCmd.MarkFlagRequired("port")
}
