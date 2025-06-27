package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

var disableSSLCmd = &cobra.Command{
	Use:   "disable-ssl",
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
		logger.Info("Removing SSL certificate using Certbot")

		cmdArgs := []string{
			"delete",
			"--cert-name", domain,
			"--non-interactive",
			"--quiet",
			"--agree-tos",
		}

		err := internal.RunCommand("sudo", append([]string{"certbot"}, cmdArgs...)...)
		if err != nil {
			logger.Warn(fmt.Sprintf("Certbot failed to delete certificate: %v", err))
			os.Exit(1)
		}

		logger.Success(fmt.Sprintf("SSL certificate removed for %s", domain))
	},
}

func init() {
	rootCmd.AddCommand(disableSSLCmd)
	disableSSLCmd.Flags().String("domain", "", "Domain name to disable SSL for")
	disableSSLCmd.MarkFlagRequired("domain")
}
