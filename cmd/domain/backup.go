package domain

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

func GetBackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup-domain",
		Short: "Backup public_html and MySQL DB of a domain",
		Run: func(cmd *cobra.Command, args []string) {
			domain, _ := cmd.Flags().GetString("domain")
			backupType, _ := cmd.Flags().GetString("type")

			if internal.IsNilOrEmpty(domain) {
				logger.Error("Please provide a domain using --domain")
				os.Exit(1)
			}
			if backupType == "" {
				backupType = "tar.gz"
			}

			username := strings.Split(domain, ".")[0]
			timestamp := time.Now().Format("20060102_150405")
			backupDir := "/var/backups"
			os.MkdirAll(backupDir, 0755)

			publicPath := fmt.Sprintf("/home/%s/public_html", username)
			sqlDump := fmt.Sprintf("%s/%s-db.sql", backupDir, domain)

			logger.Info(fmt.Sprintf("Dumping MySQL database for %s", username))
			dumpCmd := []string{"mysqldump", "-u", username, fmt.Sprintf("-p%s", username), username}
			sqlFile, err := os.Create(sqlDump)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to create dump file: %v", err))
				os.Exit(1)
			}
			defer sqlFile.Close()

			cmdDump := exec.Command("sudo", dumpCmd...)
			cmdDump.Stdout = sqlFile
			cmdDump.Stderr = os.Stderr
			if err := cmdDump.Run(); err != nil {
				logger.Warn("Could not dump MySQL DB (likely bad credentials)")
			}

			baseName := fmt.Sprintf("%s/%s-%s", backupDir, domain, timestamp)

			switch backupType {
			case "tar.gz":
				output := baseName + ".tar.gz"
				logger.Info("Creating tar.gz archive")
				internal.RunCommand("sudo", "tar", "-czf", output, publicPath, sqlDump)
				logger.Success(fmt.Sprintf("Backup created: %s", output))
			case "tar":
				output := baseName + ".tar"
				logger.Info("Creating tar archive")
				internal.RunCommand("sudo", "tar", "-cf", output, publicPath, sqlDump)
				logger.Success(fmt.Sprintf("Backup created: %s", output))
			case "zip":
				output := baseName + ".zip"
				logger.Info("Creating zip archive")
				internal.RunCommand("sudo", "zip", "-r", output, publicPath, sqlDump)
				logger.Success(fmt.Sprintf("Backup created: %s", output))
			default:
				logger.Error("Unsupported backup type. Use: tar.gz, tar, zip")
			}

			os.Remove(sqlDump)
		},
	}

	cmd.Flags().String("domain", "", "Domain name to back up")
	cmd.Flags().String("type", "tar.gz", "Backup type: tar.gz, zip, tar")
	cmd.MarkFlagRequired("domain")

	return cmd
}
