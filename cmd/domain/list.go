package domain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"stackroost/internal/logger"
)

type DomainInfo struct {
	Server string `json:"server"`
	Domain string `json:"domain"`
	Status string `json:"status"`
	User   string `json:"user"`
	Port   string `json:"port,omitempty"`
}

func GetListCmd() *cobra.Command {
	var (
		filterServer  string
		outputJSON    bool
		showEnabled   bool
		showDisabled  bool
	)

	cmd := &cobra.Command{
		Use:   "list-domains",
		Short: "List all configured domains and their status",
		Run: func(cmd *cobra.Command, args []string) {
			var results []DomainInfo
			servers := []string{"apache", "nginx", "caddy"}
			if filterServer != "" {
				servers = []string{filterServer}
			}

			for _, server := range servers {
				availableDir := getSitesAvailableDir(server)
				enabledDir := getSitesEnabledDir(server)

				files, err := os.ReadDir(availableDir)
				if err != nil {
					logger.Warn(fmt.Sprintf("Skipping %s: %v", server, err))
					continue
				}

				for _, file := range files {
					if file.IsDir() || !strings.HasSuffix(file.Name(), ".conf") {
						continue
					}

					domain := strings.TrimSuffix(file.Name(), ".conf")
					username := strings.Split(domain, ".")[0]

					linkPath := filepath.Join(enabledDir, file.Name())
					enabled := isSymlink(linkPath)

					status := "DISABLED"
					if enabled {
						status = "ENABLED"
					}

					if showEnabled && !enabled {
						continue
					}
					if showDisabled && enabled {
						continue
					}

					results = append(results, DomainInfo{
						Server: server,
						Domain: domain,
						Status: status,
						User:   username,
					})
				}
			}

			if outputJSON {
				jsonOutput, _ := json.MarshalIndent(results, "", "  ")
				fmt.Println(string(jsonOutput))
			} else {
				for _, d := range results {
					logger.Info(fmt.Sprintf("[%s] %-20s %-9s user: %-10s", strings.ToUpper(d.Server), d.Domain, d.Status, d.User))
				}
			}
		},
	}

	cmd.Flags().StringVar(&filterServer, "server", "", "Filter by server type (apache, nginx, caddy)")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&showEnabled, "enabled", false, "Only show enabled domains")
	cmd.Flags().BoolVar(&showDisabled, "disabled", false, "Only show disabled domains")

	return cmd
}

func getSitesAvailableDir(server string) string {
	switch server {
	case "apache":
		return "/etc/apache2/sites-available"
	case "nginx":
		return "/etc/nginx/sites-available"
	case "caddy":
		return "/etc/caddy/sites-available"
	default:
		return ""
	}
}

func getSitesEnabledDir(server string) string {
	switch server {
	case "apache":
		return "/etc/apache2/sites-enabled"
	case "nginx":
		return "/etc/nginx/sites-enabled"
	case "caddy":
		return "/etc/caddy/sites-enabled"
	default:
		return ""
	}
}

func isSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}
