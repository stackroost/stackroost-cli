package cmd

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

var checkPortCmd = &cobra.Command{
	Use:   "check-port",
	Short: "Check if a specific port is open for a domain",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		port, _ := cmd.Flags().GetString("port")
		timeoutSec, _ := cmd.Flags().GetInt("timeout")

		if internal.IsNilOrEmpty(domain) || internal.IsNilOrEmpty(port) {
			logger.Error("Please provide both --domain and --port")
			os.Exit(1)
		}

		address := fmt.Sprintf("%s:%s", domain, port)
		timeout := time.Duration(timeoutSec) * time.Second

		logger.Info(fmt.Sprintf("Checking port %s on domain %s...", port, domain))

		conn, err := net.DialTimeout("tcp", address, timeout)
		if err != nil {
			logger.Error(fmt.Sprintf("Port %s is not reachable on %s (%v)", port, domain, err))
			os.Exit(1)
		}
		conn.Close()

		logger.Success(fmt.Sprintf("Port %s is open and reachable on %s", port, domain))
	},
}

func init() {
	rootCmd.AddCommand(checkPortCmd)
	checkPortCmd.Flags().String("domain", "", "Domain to check")
	checkPortCmd.Flags().String("port", "", "Port to check (e.g., 80, 443, 3000)")
	checkPortCmd.Flags().Int("timeout", 3, "Timeout in seconds (default: 3s)")
	checkPortCmd.MarkFlagRequired("domain")
	checkPortCmd.MarkFlagRequired("port")
}
