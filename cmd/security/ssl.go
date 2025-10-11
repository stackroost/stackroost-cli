/*
Copyright © 2025 Stackroost CLI

*/
package security

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"stackroost-cli/cmd/internal/utils"
)

// sslCmd represents the ssl command
var sslCmd = &cobra.Command{
	Use:   "ssl",
	Short: "Manage SSL certificates",
	Long:  `Commands for issuing, renewing, revoking SSL certificates via Let's Encrypt and manual uploads.`,
}

var sslIssueCmd = &cobra.Command{
	Use:   "issue [domain]",
	Short: "Issue SSL certificate for domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		email, _ := cmd.Flags().GetString("email")
		root := viper.GetString("domains." + domain + ".root")
		server := viper.GetString("domains." + domain + ".server")
		if server == "apache" {
			utils.RunCommand("sudo", "certbot", "--apache", "-d", domain, "--email", email, "--agree-tos", "--non-interactive")
		} else {
			utils.RunCommand("sudo", "certbot", "certonly", "--webroot", "-w", root, "-d", domain, "--email", email, "--agree-tos", "--non-interactive")
		}
		// Update vhost to include SSL
		addSSLToVhost(domain, server)
		fmt.Printf("SSL issued for %s\n", domain)
	},
}

var sslRenewCmd = &cobra.Command{
	Use:   "renew [domain]",
	Short: "Renew SSL certificate",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		utils.RunCommand("sudo", "certbot", "renew", "--cert-name", domain)
		fmt.Printf("SSL renewed for %s\n", domain)
	},
}

var sslRevokeCmd = &cobra.Command{
	Use:   "revoke [domain]",
	Short: "Revoke SSL certificate",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		utils.RunCommand("sudo", "certbot", "revoke", "-d", domain)
		fmt.Printf("SSL revoked for %s\n", domain)
	},
}

var sslUploadCmd = &cobra.Command{
	Use:   "upload [domain]",
	Short: "Upload manual SSL certificate",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		cert, _ := cmd.Flags().GetString("cert")
		key, _ := cmd.Flags().GetString("key")
		dir := "/etc/letsencrypt/live/" + domain
		utils.RunCommand("sudo", "mkdir", "-p", dir)
		utils.RunCommand("sudo", "cp", cert, dir+"/fullchain.pem")
		utils.RunCommand("sudo", "cp", key, dir+"/privkey.pem")
		server := viper.GetString("domains." + domain + ".server")
		addSSLToVhost(domain, server)
		fmt.Printf("SSL uploaded for %s\n", domain)
	},
}

var sslListCmd = &cobra.Command{
	Use:   "list",
	Short: "List SSL certificates",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SSL certificates:")
		utils.RunCommand("sudo", "ls", "/etc/letsencrypt/live/")
	},
}

func AddSSLCmd(root *cobra.Command) {
	root.AddCommand(sslCmd)

	sslCmd.AddCommand(sslIssueCmd)
	sslCmd.AddCommand(sslRenewCmd)
	sslCmd.AddCommand(sslRevokeCmd)
	sslCmd.AddCommand(sslUploadCmd)
	sslCmd.AddCommand(sslListCmd)

	sslIssueCmd.Flags().String("email", "", "Email for Let's Encrypt")
	sslUploadCmd.Flags().String("cert", "", "Path to certificate file")
	sslUploadCmd.Flags().String("key", "", "Path to key file")
}

func addSSLToVhost(domain, server string) {
	if server == "apache" {
		// Certbot handles it
	} else if server == "nginx" {
		file := fmt.Sprintf("/etc/nginx/sites-available/%s", domain)
		content, _ := ioutil.ReadFile(file)
		// Add SSL server block
		sslBlock := fmt.Sprintf(`
server {
    listen 443 ssl;
    server_name %s;
    root %s;
    ssl_certificate /etc/letsencrypt/live/%s/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/%s/privkey.pem;
    location / {
        try_files $uri $uri/ =404;
    }
}`, domain, viper.GetString("domains."+domain+".root"), domain, domain)
		newContent := string(content) + sslBlock
		ioutil.WriteFile(file, []byte(newContent), 0644)
		utils.RunCommand("sudo", "systemctl", "reload", "nginx")
	}
}