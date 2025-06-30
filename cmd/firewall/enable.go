package firewall

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

var firewallPorts []int

// GetEnableCmd returns the Cobra command to enable the firewall
func GetEnableCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable-firewall",
		Short: "Enable UFW and allow common and custom ports",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Enabling UFW (Uncomplicated Firewall)")

			// Install ufw if not installed
			if err := internal.RunCommand("sudo", "apt-get", "install", "-y", "ufw"); err != nil {
				logger.Error(fmt.Sprintf("Failed to install UFW: %v", err))
				os.Exit(1)
			}

			// Allow essential ports
			defaultPorts := []int{22, 80, 443}
			for _, port := range defaultPorts {
				logger.Info(fmt.Sprintf("Allowing port: %d", port))
				internal.RunCommand("sudo", "ufw", "allow", fmt.Sprintf("%d", port))
			}

			// Allow custom ports
			for _, port := range firewallPorts {
				logger.Info(fmt.Sprintf("Allowing custom port: %d", port))
				internal.RunCommand("sudo", "ufw", "allow", fmt.Sprintf("%d", port))
			}

			// Enable ufw
			logger.Info("Enabling UFW")
			internal.RunCommand("sudo", "ufw", "--force", "enable")

			// Show status
			logger.Info("Firewall status:")
			internal.RunCommand("sudo", "ufw", "status", "verbose")

			logger.Success("Firewall configured and enabled successfully")
		},
	}

	cmd.Flags().IntSliceVarP(&firewallPorts, "port", "p", []int{}, "Additional custom ports to allow (comma separated)")
	return cmd
}
