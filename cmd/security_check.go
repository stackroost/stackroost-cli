package cmd

import (
	"os/exec"
	"strings"
	
	"github.com/spf13/cobra"
	"stackroost/internal/logger"
)

var securityCheckCmd = &cobra.Command{
	Use:   "run-security-check",
	Short: "Run basic security checks on the server",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Running security checks...")

		// SSH service check
		out := securityCaptureCommand("systemctl", "is-active", "ssh")
		if strings.Contains(out, "active") {
			logger.Success("SSH service: Active")
		} else {
			logger.Warn("SSH service: Inactive or not installed")
		}

		// Check root login via SSH
		sshdConfig := securityCaptureCommand("sudo", "grep", "^PermitRootLogin", "/etc/ssh/sshd_config")
		if strings.Contains(sshdConfig, "no") {
			logger.Success("SSH root login: Disabled")
		} else {
			logger.Warn("SSH root login: Possibly enabled (check PermitRootLogin)")
		}

		// Password authentication
		passAuth := securityCaptureCommand("sudo", "grep", "^PasswordAuthentication", "/etc/ssh/sshd_config")
		if strings.Contains(passAuth, "no") {
			logger.Success("Password authentication: Disabled (good)")
		} else {
			logger.Warn("Password authentication: Enabled (not recommended)")
		}

		// Firewall check (UFW)
		ufwStatus := securityCaptureCommand("sudo", "ufw", "status")
		if strings.Contains(ufwStatus, "Status: active") {
			logger.Success("Firewall (UFW): Active")
		} else {
			logger.Warn("Firewall (UFW): Inactive or not installed")
		}

		// fail2ban check
		fail2banStatus := securityCaptureCommand("systemctl", "is-active", "fail2ban")
		if strings.Contains(fail2banStatus, "active") {
			logger.Success("Fail2ban: Running")
		} else {
			logger.Warn("Fail2ban: Not active or not installed")
		}
	},
}

func init() {
	rootCmd.AddCommand(securityCheckCmd)
}

func securityCaptureCommand(name string, args ...string) string {
	var out strings.Builder
	cmd := exec.Command(name, args...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	_ = cmd.Run()
	return out.String()
}