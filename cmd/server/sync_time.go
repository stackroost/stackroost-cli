package server

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

func GetSyncTimeCmd() *cobra.Command {
	cmd := &cobra.Command{
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
			status := internal.CaptureCommand("timedatectl", "status")
			fmt.Println(status)
		},
	}
	return cmd
}

