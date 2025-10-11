/*
Copyright © 2025 Stackroost CLI

*/
package logs

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View server logs",
	Long:  `View logs from web servers and system services.`,
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		logType, _ := cmd.Flags().GetString("type")
		if server == "" {
			server = "apache"
		}
		if logType == "" {
			logType = "access"
		}
		logFile := getLogFile(server, logType)
		execCmd := exec.Command("tail", "-f", logFile)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		execCmd.Run()
	},
}

func AddLogsCmd(root *cobra.Command) {
	root.AddCommand(logsCmd)

	logsCmd.Flags().String("server", "", "Web server (apache, nginx, caddy)")
	logsCmd.Flags().String("type", "", "Log type (access, error)")
}

func getLogFile(server, logType string) string {
	switch server {
	case "apache":
		if logType == "error" {
			return "/var/log/apache2/error.log"
		}
		return "/var/log/apache2/access.log"
	case "nginx":
		if logType == "error" {
			return "/var/log/nginx/error.log"
		}
		return "/var/log/nginx/access.log"
	case "caddy":
		return "/var/log/caddy.log" // Caddy logs to journal or file, assuming /var/log/caddy.log
	default:
		return "/var/log/apache2/access.log"
	}
}