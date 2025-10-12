/*
Copyright © 2025 Stackroost CLI

*/
package server

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"stackroost-cli/cmd/internal/utils"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage web servers (Apache, Nginx, Caddy)",
	Long:  `Commands for managing web server services including start, stop, reload, and status checks.`,
}

var servers = map[string]string{
	"apache": "httpd",
	"nginx":  "nginx",
	"caddy":  "caddy",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed web servers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Installed web servers:")
		for name, service := range servers {
			if isServiceInstalled(service) {
				fmt.Printf("- %s (%s)\n", name, service)
			}
		}
	},
}

var startCmd = &cobra.Command{
	Use:   "start [server]",
	Short: "Start a web server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		server := args[0]
		service, ok := servers[server]
		if !ok {
			fmt.Printf("Unknown server: %s\n", server)
			return
		}
		utils.RunCommand("sudo", "systemctl", "start", service)
		fmt.Printf("Started %s\n", server)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop [server]",
	Short: "Stop a web server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		server := args[0]
		service, ok := servers[server]
		if !ok {
			fmt.Printf("Unknown server: %s\n", server)
			return
		}
		utils.RunCommand("sudo", "systemctl", "stop", service)
		fmt.Printf("Stopped %s\n", server)
	},
}

var reloadCmd = &cobra.Command{
	Use:   "reload [server]",
	Short: "Reload a web server configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		server := args[0]
		service, ok := servers[server]
		if !ok {
			fmt.Printf("Unknown server: %s\n", server)
			return
		}
		utils.RunCommand("sudo", "systemctl", "reload", service)
		fmt.Printf("Reloaded %s\n", server)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of web servers",
	Run: func(cmd *cobra.Command, args []string) {
		for name, service := range servers {
			if isServiceInstalled(service) {
				fmt.Printf("%s:\n", name)
				cmd := exec.Command("systemctl", "status", service, "--no-pager", "-l")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
				fmt.Println()
			}
		}
	},
}

func AddServerCmd(root *cobra.Command) {
	root.AddCommand(serverCmd)

	serverCmd.AddCommand(listCmd)
	serverCmd.AddCommand(startCmd)
	serverCmd.AddCommand(stopCmd)
	serverCmd.AddCommand(reloadCmd)
	serverCmd.AddCommand(statusCmd)
}

func isServiceInstalled(service string) bool {
	cmd := exec.Command("systemctl", "is-active", service)
	err := cmd.Run()
	return err == nil
}