package ssl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

// GetEnableCmd returns the command to enable SSL for a domain
func GetEnableCmd() *cobra.Command {
	var domain string

	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable Let's Encrypt SSL for a specific domain",
		Run: func(cmd *cobra.Command, args []string) {
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
				logger.Info("Caddy automatically handles SSL — no need to enable manually.")
				return
			}

			logger.Info(fmt.Sprintf("Detected %s configuration for %s", serverType, domain))

			err := internal.EnableSSLCertbot(domain, serverType)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to enable SSL for %s: %v", domain, err))
				os.Exit(1)
			}

			logger.Success(fmt.Sprintf("SSL enabled successfully for %s", domain))
		},
	}

	cmd.Flags().StringVar(&domain, "domain", "", "Domain name to enable SSL for")
	cmd.MarkFlagRequired("domain")
	return cmd
}
