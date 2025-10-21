/*
Copyright © 2025 Stackroost CLI

*/
package server

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"stackroost-cli/cmd/internal/utils"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage web servers (Apache, Nginx, Caddy)",
	Long:  `Commands for managing web server services including start, stop, reload, and status checks.`,
}

var servers = map[string]map[string]string{
	"ubuntu": {
		"apache": "apache2",
		"nginx":  "nginx",
		"caddy":  "caddy",
	},
	"centos": {
		"apache": "httpd",
		"nginx":  "nginx",
		"caddy":  "caddy",
	},
	"fedora": {
		"apache": "httpd",
		"nginx":  "nginx",
		"caddy":  "caddy",
	},
	"debian": {
		"apache": "apache2",
		"nginx":  "nginx",
		"caddy":  "caddy",
	},
	"rhel": {
		"apache": "httpd",
		"nginx":  "nginx",
		"caddy":  "caddy",
	},
	"sles": {
		"apache": "apache2",
		"nginx":  "nginx",
		"caddy":  "caddy",
	},
	"opensuse": {
		"apache": "apache2",
		"nginx":  "nginx",
		"caddy":  "caddy",
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed web servers",
	Run: func(cmd *cobra.Command, args []string) {
		distro := detectDistro()
		fmt.Println("Installed web servers:")
		if distroServers, ok := servers[distro]; ok {
			for name, service := range distroServers {
				if isServiceInstalled(service) {
					fmt.Printf("- %s (%s)\n", name, service)
				}
			}
		} else {
			fmt.Printf("Unsupported distribution: %s\n", distro)
		}
	},
}

var startCmd = &cobra.Command{
	Use:   "start [server]",
	Short: "Start a web server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		distro := detectDistro()
		server := args[0]
		if distroServers, ok := servers[distro]; ok {
			if service, ok := distroServers[server]; ok {
				utils.RunCommand("sudo", "systemctl", "start", service)
				fmt.Printf("Started %s\n", server)
			} else {
				fmt.Printf("Unknown server: %s\n", server)
			}
		} else {
			fmt.Printf("Unsupported distribution: %s\n", distro)
		}
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop [server]",
	Short: "Stop a web server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		distro := detectDistro()
		server := args[0]
		if distroServers, ok := servers[distro]; ok {
			if service, ok := distroServers[server]; ok {
				utils.RunCommand("sudo", "systemctl", "stop", service)
				fmt.Printf("Stopped %s\n", server)
			} else {
				fmt.Printf("Unknown server: %s\n", server)
			}
		} else {
			fmt.Printf("Unsupported distribution: %s\n", distro)
		}
	},
}

var reloadCmd = &cobra.Command{
	Use:   "reload [server]",
	Short: "Reload a web server configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		distro := detectDistro()
		server := args[0]
		if distroServers, ok := servers[distro]; ok {
			if service, ok := distroServers[server]; ok {
				utils.RunCommand("sudo", "systemctl", "reload", service)
				fmt.Printf("Reloaded %s\n", server)
			} else {
				fmt.Printf("Unknown server: %s\n", server)
			}
		} else {
			fmt.Printf("Unsupported distribution: %s\n", distro)
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of web servers",
	Run: func(cmd *cobra.Command, args []string) {
		distro := detectDistro()
		if distroServers, ok := servers[distro]; ok {
			for name, service := range distroServers {
				if isServiceInstalled(service) {
					fmt.Printf("%s:\n", name)
					cmd := exec.Command("systemctl", "status", service, "--no-pager", "-l")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.Run()
					fmt.Println()
				}
			}
		} else {
			fmt.Printf("Unsupported distribution: %s\n", distro)
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

func detectDistro() string {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "unknown"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ID=") {
			id := strings.TrimPrefix(line, "ID=")
			id = strings.Trim(id, "\"")
			return id
		}
	}
	return "unknown"
}

func isServiceInstalled(service string) bool {
	cmd := exec.Command("systemctl", "is-active", service)
	err := cmd.Run()
	return err == nil
}