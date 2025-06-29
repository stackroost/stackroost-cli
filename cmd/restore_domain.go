package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

var restoreDomainCmd = &cobra.Command{
	Use:   "restore-domain",
	Short: "Restore domain files and MySQL DB from backup archive",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		backupFile, _ := cmd.Flags().GetString("file")

		if internal.IsNilOrEmpty(domain) || internal.IsNilOrEmpty(backupFile) {
			logger.Error("Please provide both --domain and --file")
			os.Exit(1)
		}

		username := strings.Split(domain, ".")[0]
		restoreDir := fmt.Sprintf("/tmp/restore-%s", username)
		os.MkdirAll(restoreDir, 0755)

		ext := filepath.Ext(backupFile)
		if ext == ".gz" || strings.HasSuffix(backupFile, ".tar.gz") {
			logger.Info("Extracting tar.gz archive")
			internal.RunCommand("sudo", "tar", "-xzf", backupFile, "-C", restoreDir)
		} else if ext == ".tar" {
			logger.Info("Extracting tar archive")
			internal.RunCommand("sudo", "tar", "-xf", backupFile, "-C", restoreDir)
		} else if ext == ".zip" {
			logger.Info("Extracting zip archive")
			internal.RunCommand("sudo", "unzip", "-o", backupFile, "-d", restoreDir)
		} else {
			logger.Error("Unsupported file type. Use: tar.gz, tar, or zip")
			os.Exit(1)
		}

		// Restore public_html
		publicPath := fmt.Sprintf("/home/%s/public_html", username)
		logger.Info(fmt.Sprintf("Restoring files to %s", publicPath))
		internal.RunCommand("sudo", "cp", "-r", filepath.Join(restoreDir, "home", username, "public_html"), filepath.Join("/home", username))
		internal.RunCommand("sudo", "chown", "-R", fmt.Sprintf("%s:%s", username, username), publicPath)

		// Restore MySQL if .sql found
		sqlPath := ""
		filepath.Walk(restoreDir, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".sql") {
				sqlPath = path
			}
			return nil
		})

		if sqlPath != "" {
			logger.Info(fmt.Sprintf("Restoring MySQL DB from %s", sqlPath))
			restoreCmd := exec.Command("sudo", "mysql", "-u", username, fmt.Sprintf("-p%s", username), username)
			sqlFile, _ := os.Open(sqlPath)
			defer sqlFile.Close()
			restoreCmd.Stdin = sqlFile
			restoreCmd.Stdout = os.Stdout
			restoreCmd.Stderr = os.Stderr
			if err := restoreCmd.Run(); err != nil {
				logger.Warn(fmt.Sprintf("MySQL restore failed: %v", err))
			} else {
				logger.Success("MySQL database restored")
			}
		} else {
			logger.Warn("No SQL file found in backup. Skipping DB restore.")
		}

		logger.Success(fmt.Sprintf("Domain '%s' restored successfully", domain))
	},
}

func init() {
	rootCmd.AddCommand(restoreDomainCmd)
	restoreDomainCmd.Flags().String("domain", "", "Domain name to restore")
	restoreDomainCmd.Flags().String("file", "", "Path to backup archive file (.tar.gz, .zip, .tar)")
	restoreDomainCmd.MarkFlagRequired("domain")
	restoreDomainCmd.MarkFlagRequired("file")
}
