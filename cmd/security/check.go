package security

import (
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"stackroost/internal/logger"
)

// GetCheckCmd returns the Cobra command for security check
func GetCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run-security-check",
		Short: "Run basic security checks on the server",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Running security checks...")

			// SSH service check
			out := capture("systemctl", "is-active", "ssh")
			if strings.Contains(out, "active") {
				logger.Success("SSH service: Active")
			} else {
				logger.Warn("SSH service: Inactive or not installed")
			}

			// Check root login via SSH
			sshdConfig := capture("sudo", "grep", "^PermitRootLogin", "/etc/ssh/sshd_config")
			if strings.Contains(sshdConfig, "no") {
				logger.Success("SSH root login: Disabled")
			} else {
				logger.Warn("SSH root login: Possibly enabled (check PermitRootLogin)")
			}

			// Password authentication
			passAuth := capture("sudo", "grep", "^PasswordAuthentication", "/etc/ssh/sshd_config")
			if strings.Contains(passAuth, "no") {
				logger.Success("Password authentication: Disabled (good)")
			} else {
				logger.Warn("Password authentication: Enabled (not recommended)")
			}

			// Firewall check (UFW)
			ufwStatus := capture("sudo", "ufw", "status")
			if strings.Contains(ufwStatus, "Status: active") {
				logger.Success("Firewall (UFW): Active")
			} else {
				logger.Warn("Firewall (UFW): Inactive or not installed")
			}

			// Fail2ban check
			fail2banStatus := capture("systemctl", "is-active", "fail2ban")
			if strings.Contains(fail2banStatus, "active") {
				logger.Success("Fail2ban: Running")
			} else {
				logger.Warn("Fail2ban: Not active or not installed")
			}
		},
	}
	return cmd
}

func capture(name string, args ...string) string {
	var out strings.Builder
	cmd := exec.Command(name, args...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	_ = cmd.Run()
	return out.String()
}
