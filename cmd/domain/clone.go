package domain

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

func GetCloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone-domain",
		Short: "Clone configuration, files, and database from one domain to another",
		Run: func(cmd *cobra.Command, args []string) {
			src, _ := cmd.Flags().GetString("source")
			dest, _ := cmd.Flags().GetString("target")
			cloneDB, _ := cmd.Flags().GetBool("clone-db")
			cloneUser, _ := cmd.Flags().GetBool("clone-user")

			if internal.IsNilOrEmpty(src) || internal.IsNilOrEmpty(dest) {
				logger.Error("Please provide both --source and --target domains")
				os.Exit(1)
			}

			logger.Info(fmt.Sprintf("Cloning domain from %s to %s", src, dest))

			serverType := internal.DetectServerType(src)
			if serverType == "" {
				logger.Error("Could not detect server type for source domain")
				os.Exit(1)
			}

			// Copy config
			srcConf := filepath.Join("/etc", serverType, "sites-available", src+".conf")
			destConf := filepath.Join("/etc", serverType, "sites-available", dest+".conf")
			logger.Info("Copying config...")
			internal.RunCommand("sudo", "cp", srcConf, destConf)
			internal.RunCommand("sudo", "sed", "-i", fmt.Sprintf("s/%s/%s/g", src, dest), destConf)

			// Copy public_html
			srcUser := strings.Split(src, ".")[0]
			destUser := strings.Split(dest, ".")[0]
			srcPath := fmt.Sprintf("/home/%s/public_html", srcUser)
			destPath := fmt.Sprintf("/home/%s/public_html", destUser)

			logger.Info("Copying website files...")
			internal.RunCommand("sudo", "mkdir", "-p", destPath)
			internal.RunCommand("sudo", "cp", "-r", srcPath+"/.", destPath)

			if cloneUser {
				logger.Info("Cloning user...")
				internal.RunCommand("sudo", "useradd", "-m", "-s", "/bin/bash", destUser)
				internal.RunCommand("sudo", "chown", "-R", fmt.Sprintf("%s:%s", destUser, destUser), "/home/"+destUser)
			}

			if cloneDB {
				logger.Info("Cloning MySQL database...")
				dumpFile := fmt.Sprintf("/tmp/%s.sql", srcUser)
				internal.RunCommand("sudo", "mysqldump", "-u", "root", srcUser, "-r", dumpFile)
				internal.CreateMySQLUserAndDatabase(destUser, "changeme123")
				internal.RunCommand("sudo", "mysql", "-u", "root", destUser, "-e", fmt.Sprintf("source %s", dumpFile))
				internal.RunCommand("sudo", "rm", "-f", dumpFile)
			}

			// Enable site
			logger.Info("Enabling cloned site...")
			switch serverType {
			case "apache":
				internal.RunCommand("sudo", "a2ensite", dest+".conf")
				internal.RunCommand("sudo", "systemctl", "reload", "apache2")
			case "nginx":
				link := filepath.Join("/etc/nginx/sites-enabled", dest+".conf")
				internal.RunCommand("sudo", "ln", "-s", destConf, link)
				internal.RunCommand("sudo", "systemctl", "reload", "nginx")
			case "caddy":
				link := filepath.Join("/etc/caddy/sites-enabled", dest+".conf")
				internal.RunCommand("sudo", "ln", "-s", destConf, link)
				internal.RunCommand("sudo", "systemctl", "reload", "caddy")
			}

			logger.Success(fmt.Sprintf("Domain %s cloned to %s", src, dest))
		},
	}

	cmd.Flags().String("source", "", "Source domain to clone from")
	cmd.Flags().String("target", "", "New domain name")
	cmd.Flags().Bool("clone-db", false, "Clone MySQL database")
	cmd.Flags().Bool("clone-user", false, "Clone shell user")
	cmd.MarkFlagRequired("source")
	cmd.MarkFlagRequired("target")

	return cmd
}
