package cmd

import (
	"fmt"
	"os"
	"strings"
	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
	"os/exec"
)

var toggleSiteCmd = &cobra.Command{
	Use:   "toggle-site",
	Short: "Enable or disable a site's configuration",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		if internal.IsNilOrEmpty(domain) {
			logger.Error("Please provide a domain using --domain")
			os.Exit(1)
		}

		serverType := internal.DetectServerType(domain)
		if serverType == "" {
			logger.Warn("Could not detect server type")
			os.Exit(1)
		}
		logger.Info(fmt.Sprintf("Detected server: %s", serverType))

		filename := domain + ".conf"
		enabled := false

		switch serverType {
		case "apache":
			output := CaptureCommand("a2query", "-s", domain)
			enabled = strings.Contains(output, "is enabled")
			if enabled {
				logger.Info(fmt.Sprintf("Disabling Apache site: %s", domain))
				internal.RunCommand("sudo", "a2dissite", filename)
			} else {
				logger.Info(fmt.Sprintf("Enabling Apache site: %s", domain))
				internal.RunCommand("sudo", "a2ensite", filename)
			}
			internal.RunCommand("sudo", "systemctl", "reload", "apache2")

		case "nginx":
			sitesAvailable := "/etc/nginx/sites-available/" + filename
			sitesEnabled := "/etc/nginx/sites-enabled/" + filename
			if _, err := os.Stat(sitesEnabled); err == nil {
				logger.Info(fmt.Sprintf("Disabling Nginx site: %s", domain))
				internal.RunCommand("sudo", "rm", "-f", sitesEnabled)
			} else {
				logger.Info(fmt.Sprintf("Enabling Nginx site: %s", domain))
				internal.RunCommand("sudo", "ln", "-s", sitesAvailable, sitesEnabled)
			}
			internal.RunCommand("sudo", "systemctl", "reload", "nginx")

		case "caddy":
			sitesAvailable := "/etc/caddy/sites-available/" + filename
			sitesEnabled := "/etc/caddy/sites-enabled/" + filename
			if _, err := os.Stat(sitesEnabled); err == nil {
				logger.Info(fmt.Sprintf("Disabling Caddy site: %s", domain))
				internal.RunCommand("sudo", "rm", "-f", sitesEnabled)
			} else {
				logger.Info(fmt.Sprintf("Enabling Caddy site: %s", domain))
				internal.RunCommand("sudo", "ln", "-s", sitesAvailable, sitesEnabled)
			}
			internal.RunCommand("sudo", "systemctl", "reload", "caddy")
		}

		logger.Success(fmt.Sprintf("Site %s toggled successfully", domain))
	},
}

func init() {
	rootCmd.AddCommand(toggleSiteCmd)
	toggleSiteCmd.Flags().String("domain", "", "Domain name to toggle")
	toggleSiteCmd.MarkFlagRequired("domain")
}

func CaptureCommand(name string, args ...string) string {
	var out strings.Builder
	cmd := exec.Command(name, args...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	_ = cmd.Run()
	return out.String()
}