package user

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"stackroost/internal/logger"
)

func GetListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-users",
		Short: "List all regular system users (UID ≥ 1000)",
		Run: func(cmd *cobra.Command, args []string) {
			file, err := os.Open("/etc/passwd")
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to open /etc/passwd: %v", err))
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			logger.Info("Listing non-system shell users (UID ≥ 1000)")
			for scanner.Scan() {
				line := scanner.Text()
				parts := strings.Split(line, ":")
				if len(parts) < 7 {
					continue
				}

				username := parts[0]
				uidStr := parts[2]
				shell := parts[6]

				uid, err := strconv.Atoi(uidStr)
				if err != nil || uid < 1000 {
					continue
				}

				if shell == "/bin/bash" || shell == "/bin/sh" {
					u, _ := user.Lookup(username)
					logger.Info(fmt.Sprintf("User: %-15s UID: %-5d Shell: %s Home: %s", username, uid, shell, u.HomeDir))
				}
			}
		},
	}
}
