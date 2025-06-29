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

var syncTimeCmd = &cobra.Command{
	Use:   "sync-time",
	Short: "Sync server time with NTP using systemd-timesyncd",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Enabling NTP time sync...")
		if err := internal.RunCommand("sudo", "timedatectl", "set-ntp", "true"); err != nil {
			logger.Error(fmt.Sprintf("Failed to enable NTP sync: %v", err))
			os.Exit(1)
		}

		logger.Info("Restarting time sync service...")
		if err := internal.RunCommand("sudo", "systemctl", "restart", "systemd-timesyncd"); err != nil {
			logger.Error(fmt.Sprintf("Failed to restart timesync service: %v", err))
			os.Exit(1)
		}

		logger.Success("Time synchronization triggered successfully.")

		logger.Info("Current time sync status:")
		status := TimeCaptureCommand("timedatectl", "status")
		fmt.Println(status)
	},
}

func init() {
	rootCmd.AddCommand(syncTimeCmd)
}

func TimeCaptureCommand(name string, args ...string) string {
	var out strings.Builder
	cmd := exec.Command(name, args...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	_ = cmd.Run()
	return out.String()
}