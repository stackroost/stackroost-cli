package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

var sslDomain string

var enableSSLCmd = &cobra.Command{
	Use:   "enable-ssl",
	Short: "Enable Let's Encrypt SSL for a specific domain",
	Run: func(cmd *cobra.Command, args []string) {
		if sslDomain == "" {
			logger.Error("Please provide a domain using --domain")
			os.Exit(1)
		}

		serverType := internal.DetectServerType(sslDomain)
		if serverType == "" {
			logger.Error(fmt.Sprintf("Could not detect server type for domain: %s", sslDomain))
			os.Exit(1)
		}

		if serverType == "caddy" {
			logger.Info("Caddy automatically handles SSL — no need to enable manually.")
			return
		}

		logger.Info(fmt.Sprintf("Detected %s configuration for %s", serverType, sslDomain))

		err := internal.EnableSSLCertbot(sslDomain, serverType)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to enable SSL for %s: %v", sslDomain, err))
			os.Exit(1)
		}

		logger.Success(fmt.Sprintf("SSL enabled successfully for %s", sslDomain))
	},
}

func init() {
	rootCmd.AddCommand(enableSSLCmd)
	enableSSLCmd.Flags().StringVar(&sslDomain, "domain", "", "Domain name to enable SSL for")
	enableSSLCmd.MarkFlagRequired("domain")
}
