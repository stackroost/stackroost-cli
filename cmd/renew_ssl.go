package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

var (
	renewAll   bool
	domainName string
	forceFlag  bool
)

var renewSSLCmd = &cobra.Command{
	Use:   "renew-ssl",
	Short: "Renew SSL certificates for all or specific domains",
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

		if domainName == "" {
			logger.Error("Either --all or --domain <domain> must be provided")
			os.Exit(1)
		}

		serverType := internal.DetectServerType(domainName)
		if serverType == "" {
			logger.Error(fmt.Sprintf("Could not detect server type for domain: %s", domainName))
			os.Exit(1)
		}

		logger.Info(fmt.Sprintf("Renewing certificate for %s (%s)", domainName, serverType))

		cmdArgs := []string{
			fmt.Sprintf("--%s", serverType),
			"-d", domainName,
			"-d", "www." + domainName,
			"--non-interactive",
			"--agree-tos",
			"--register-unsafely-without-email",
		}

		if forceFlag {
			cmdArgs = append(cmdArgs, "--force-renewal")
		}

		if err := internal.RunCommand("sudo", append([]string{"certbot"}, cmdArgs...)...); err != nil {
			logger.Error(fmt.Sprintf("SSL renewal failed for %s: %v", domainName, err))
			os.Exit(1)
		}

		logger.Success(fmt.Sprintf("Certificate renewed for %s", domainName))
	},
}

func init() {
	rootCmd.AddCommand(renewSSLCmd)
	renewSSLCmd.Flags().BoolVar(&renewAll, "all", false, "Renew all certificates")
	renewSSLCmd.Flags().StringVar(&domainName, "domain", "", "Domain to renew certificate for")
	renewSSLCmd.Flags().BoolVar(&forceFlag, "force", false, "Force renew the certificate")
}