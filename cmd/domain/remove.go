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

func GetRemoveCmd() *cobra.Command {
	var domain, serverType string
	var keepUser bool

	cmd := &cobra.Command{
		Use:   "remove-domain",
		Short: "Remove a domain configuration, user, and database",
		Run: func(cmd *cobra.Command, args []string) {
			if internal.IsNilOrEmpty(domain) {
				logger.Error("Domain name is required")
				os.Exit(1)
			}

			username := strings.Split(domain, ".")[0]
			filename := domain + ".conf"

			logger.Info(fmt.Sprintf("Removing domain: %s", domain))

			// Step 1: Disable web server site
			switch serverType {
			case "apache":
				logger.Info("Disabling Apache site")
				internal.RunCommand("sudo", "a2dissite", filename)
				internal.RunCommand("sudo", "systemctl", "reload", "apache2")
			case "nginx":
				link := filepath.Join("/etc/nginx/sites-enabled", filename)
				internal.RunCommand("sudo", "rm", "-f", link)
				internal.RunCommand("sudo", "systemctl", "reload", "nginx")
			case "caddy":
				link := filepath.Join("/etc/caddy/sites-enabled", filename)
				internal.RunCommand("sudo", "rm", "-f", link)
				internal.RunCommand("sudo", "systemctl", "reload", "caddy")
			default:
				logger.Error(fmt.Sprintf("Unsupported server type: %s", serverType))
				os.Exit(1)
			}

			// Step 2: Remove config file
			configPath := getConfigPath(serverType, domain)
			if err := os.Remove(configPath); err != nil {
				logger.Warn(fmt.Sprintf("Could not delete config file: %v", err))
			} else {
				logger.Success(fmt.Sprintf("Removed config file: %s", configPath))
			}

			// Step 3: Remove MySQL database and user
			if err := internal.DropMySQLUserAndDatabase(username); err != nil {
				logger.Warn(fmt.Sprintf("MySQL cleanup failed: %v", err))
			} else {
				logger.Success("MySQL user and database removed")
			}

			// Step 4: Remove system user
			if !keepUser {
				logger.Info(fmt.Sprintf("Removing Linux user: %s", username))
				if currentUser, _ := user.Current(); currentUser.Username == username {
					logger.Error("Refusing to delete the current executing user")
					os.Exit(1)
				}
				internal.RunCommand("sudo", "userdel", "-r", username)
				logger.Success(fmt.Sprintf("User '%s' and home directory removed", username))
			} else {
				logger.Info("Keeping shell user and home directory (per flag)")
			}

			logger.Success(fmt.Sprintf("Domain '%s' removed successfully", domain))
		},
	}

	cmd.Flags().StringVarP(&domain, "name", "n", "", "Domain name to remove")
	cmd.Flags().StringVarP(&serverType, "server", "s", "apache", "Server type (apache, nginx, caddy)")
	cmd.Flags().BoolVar(&keepUser, "keep-user", false, "Keep the Linux user and home directory")
	cmd.MarkFlagRequired("name")

	return cmd
}

func getConfigPath(serverType, domain string) string {
	filename := domain + ".conf"
	switch serverType {
	case "apache":
		return filepath.Join("/etc/apache2/sites-available", filename)
	case "nginx":
		return filepath.Join("/etc/nginx/sites-available", filename)
	case "caddy":
		return filepath.Join("/etc/caddy/sites-available", filename)
	default:
		return ""
	}
}
