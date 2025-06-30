package ssl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

// GetRenewCmd returns the command to renew SSL certificates
func GetRenewCmd() *cobra.Command {
	var renewAll bool
	var domain string
	var force bool

	cmd := &cobra.Command{
		Use:   "renew",
		Short: "Renew SSL certificates for all or a specific domain",
		Run: func(cmd *cobra.Command, args []string) {
			if renewAll {
				logger.Info("Renewing SSL certificates for all domains")
				err := internal.RunCommand("sudo", "certbot", "renew")
				if err != nil {
					logger.Error(fmt.Sprintf("SSL renewal failed: %v", err))
					os.Exit(1)
				}
				logger.Success("All certificates renewed successfully")
				return
			}

			if internal.IsNilOrEmpty(domain) {
				logger.Error("Either --all or --domain <domain> must be provided")
				os.Exit(1)
			}

			serverType := internal.DetectServerType(domain)
			if serverType == "" {
				logger.Error(fmt.Sprintf("Could not detect server type for domain: %s", domain))
				os.Exit(1)
			}

			logger.Info(fmt.Sprintf("Renewing certificate for %s (%s)", domain, serverType))

			cmdArgs := []string{
				fmt.Sprintf("--%s", serverType),
				"-d", domain,
				"-d", "www." + domain,
				"--non-interactive",
				"--agree-tos",
				"--register-unsafely-without-email",
			}

			if force {
				cmdArgs = append(cmdArgs, "--force-renewal")
			}

			if err := internal.RunCommand("sudo", append([]string{"certbot"}, cmdArgs...)...); err != nil {
				logger.Error(fmt.Sprintf("SSL renewal failed for %s: %v", domain, err))
				os.Exit(1)
			}

			logger.Success(fmt.Sprintf("Certificate renewed for %s", domain))
		},
	}

	cmd.Flags().BoolVar(&renewAll, "all", false, "Renew all certificates")
	cmd.Flags().StringVar(&domain, "domain", "", "Domain to renew certificate for")
	cmd.Flags().BoolVar(&force, "force", false, "Force renew the certificate")
	return cmd
}
