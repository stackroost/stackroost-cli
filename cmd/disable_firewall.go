package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
	"strings"
)

var disableFirewallCmd = &cobra.Command{
	Use:   "disable-firewall",
	Short: "Disable the system firewall (UFW) safely",
	Run: func(cmd *cobra.Command, args []string) {
		flush, _ := cmd.Flags().GetBool("flush")

		logger.Info("Checking firewall status...")
		statusOutput := CMDRUN("sudo", "ufw", "status")

		if statusOutput == "" || statusOutput == "Status: inactive\n" {
			logger.Warn("Firewall is already inactive")
			return
		}

		logger.Info("Disabling firewall (UFW)...")
		if err := internal.RunCommand("sudo", "ufw", "disable"); err != nil {
			logger.Error(fmt.Sprintf("Failed to disable firewall: %v", err))
			os.Exit(1)
		}

		if flush {
			logger.Warn("Flushing all UFW rules...")
			if err := internal.RunCommand("sudo", "ufw", "reset"); err != nil {
				logger.Error(fmt.Sprintf("Failed to flush firewall rules: %v", err))
				os.Exit(1)
			}
			logger.Success("All firewall rules flushed.")
		}

		logger.Success("Firewall disabled successfully.")
	},
}

func init() {
	rootCmd.AddCommand(disableFirewallCmd)
	disableFirewallCmd.Flags().Bool("flush", false, "Flush all UFW rules after disabling")
}

func CMDRUN(name string, args ...string) string {
	var out strings.Builder
	cmd := exec.Command(name, args...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	_ = cmd.Run()
	return out.String()
}