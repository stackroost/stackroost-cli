package security

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

var (
	ports            string
	disableRootLogin bool
	enforceKeyAuth   bool
)

var secureServerCmd = &cobra.Command{
	Use:   "secure-server",
	Short: "Secure the server by enabling firewall, restricting SSH, and hardening config",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Starting server hardening...")

		// Enable UFW and set defaults
		logger.Info("Enabling UFW...")
		internal.RunCommand("sudo", "ufw", "--force", "enable")
		internal.RunCommand("sudo", "ufw", "default", "deny", "incoming")
		internal.RunCommand("sudo", "ufw", "default", "allow", "outgoing")

		for _, port := range strings.Split(ports, ",") {
			port = strings.TrimSpace(port)
			if port != "" {
				internal.RunCommand("sudo", "ufw", "allow", port)
				logger.Success(fmt.Sprintf("Allowed port: %s", port))
			}
		}

		// SSH Hardening
		if disableRootLogin || enforceKeyAuth {
			sshConf := "/etc/ssh/sshd_config"

			if disableRootLogin {
				internal.RunCommand("sudo", "sed", "-i", "s/^#*PermitRootLogin.*/PermitRootLogin no/", sshConf)
				logger.Success("Root login disabled in SSH config")
			}

			if enforceKeyAuth {
				internal.RunCommand("sudo", "sed", "-i", "s/^#*PasswordAuthentication.*/PasswordAuthentication no/", sshConf)
				logger.Success("PasswordAuthentication disabled — key-based auth enforced")
			}

			internal.RunCommand("sudo", "systemctl", "restart", "ssh")
		}

		logger.Success("Server hardening completed.")
	},
}

func init() {
	secureServerCmd.Flags().StringVar(&ports, "allow-ports", "22,80,443", "Comma-separated ports to allow through UFW")
	secureServerCmd.Flags().BoolVar(&disableRootLogin, "disable-root-login", false, "Disable SSH root login")
	secureServerCmd.Flags().BoolVar(&enforceKeyAuth, "enforce-ssh-key-only", false, "Disable SSH password login")
}

// Exported for root registration
func GetSecureCmd() *cobra.Command {
	return secureServerCmd
}
