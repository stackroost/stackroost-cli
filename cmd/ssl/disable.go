package ssl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

func GetDisableCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable and remove SSL certificate for a specific domain",
		Run: func(cmd *cobra.Command, args []string) {
			domain, _ := cmd.Flags().GetString("domain")
			if internal.IsNilOrEmpty(domain) {
				logger.Error("Please provide a domain using --domain")
				os.Exit(1)
			}

			serverType := internal.DetectServerType(domain)
			if serverType == "" {
				logger.Error(fmt.Sprintf("Could not detect server type for domain: %s", domain))
				os.Exit(1)
			}

			if serverType == "caddy" {
				logger.Info("Caddy auto-manages SSL — no need to disable manually.")
				return
			}

			logger.Info(fmt.Sprintf("Detected %s configuration for %s", serverType, domain))
			logger.Info("Removing SSL certificate using Certbot...")

			cmdArgs := []string{
				"delete",
				"--cert-name", domain,
				"--non-interactive",
				"--quiet",
				"--agree-tos",
			}

			if err := internal.RunCommand("sudo", append([]string{"certbot"}, cmdArgs...)...); err != nil {
				logger.Warn(fmt.Sprintf("Certbot failed to delete certificate: %v", err))
				os.Exit(1)
			}

			logger.Success(fmt.Sprintf("SSL certificate removed for %s", domain))
		},
	}

	cmd.Flags().String("domain", "", "Domain name to disable SSL for")
	cmd.MarkFlagRequired("domain")

	return cmd
}
