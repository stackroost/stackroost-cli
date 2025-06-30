package server

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

func GetRestartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart-server",
		Short: "Restart a specific web server (apache, nginx, or caddy)",
		Run: func(cmd *cobra.Command, args []string) {
			server, _ := cmd.Flags().GetString("server")

			if internal.IsNilOrEmpty(server) {
				logger.Error("Please provide a server type using --server (apache, nginx, or caddy)")
				os.Exit(1)
			}

			server = strings.ToLower(server)
			validServers := map[string]string{
				"apache": "apache2",
				"nginx":  "nginx",
				"caddy":  "caddy",
			}

			systemName, ok := validServers[server]
			if !ok {
				logger.Error("Unsupported server type. Use one of: apache, nginx, or caddy")
				os.Exit(1)
			}

			logger.Info(fmt.Sprintf("Restarting %s server...", server))
			err := internal.RunCommand("sudo", "systemctl", "restart", systemName)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to restart %s: %v", server, err))
				os.Exit(1)
			}

			logger.Success(fmt.Sprintf("%s restarted successfully", strings.Title(server)))
		},
	}

	cmd.Flags().String("server", "", "Web server to restart (apache, nginx, or caddy)")
	cmd.MarkFlagRequired("server")
	return cmd
}
