package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

var scheduleRestartCmd = &cobra.Command{
	Use:   "schedule-restart",
	Short: "Schedule a server restart after a delay",
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		delay, _ := cmd.Flags().GetInt("delay")

		server = strings.ToLower(server)
		if server != "apache" && server != "nginx" && server != "caddy" {
			logger.Error("Invalid server type. Use --server apache|nginx|caddy")
			os.Exit(1)
		}

		if delay <= 0 {
			delay = 5
		}

		logger.Info(fmt.Sprintf("Server: %s", server))
		logger.Info(fmt.Sprintf("Restart scheduled in %d seconds...", delay))
		time.Sleep(time.Duration(delay) * time.Second)

		var service string
		switch server {
		case "apache":
			service = "apache2"
		case "nginx":
			service = "nginx"
		case "caddy":
			service = "caddy"
		}

		logger.Info("Restarting now...")
		if err := internal.RunCommand("sudo", "systemctl", "restart", service); err != nil {
			logger.Error(fmt.Sprintf("Failed to restart %s: %v", service, err))
			os.Exit(1)
		}

		logger.Success(fmt.Sprintf("%s restarted successfully", strings.Title(server)))
	},
}

func init() {
	rootCmd.AddCommand(scheduleRestartCmd)
	scheduleRestartCmd.Flags().String("server", "", "Server to restart (apache|nginx|caddy)")
	scheduleRestartCmd.Flags().Int("delay", 5, "Delay in seconds before restart")
	scheduleRestartCmd.MarkFlagRequired("server")
}
