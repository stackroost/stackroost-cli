/*
Copyright © 2025 Stackroost CLI

*/
package domain

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"stackroost-cli/cmd/internal/utils"
)

// domainCmd represents the domain command
var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Manage domains and virtual hosts",
	Long:  `Commands for adding, listing, removing, enabling, and disabling domains with virtual host configurations.`,
}

var domainAddCmd = &cobra.Command{
	Use:   "add [domain]",
	Short: "Add a new domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		server, _ := cmd.Flags().GetString("server")
		if server == "" {
			server = "apache" // default
		}
		createVhost(domain, server)
		fmt.Printf("Added domain %s for %s\n", domain, server)
	},
}

var domainListCmd = &cobra.Command{
	Use:   "list",
	Short: "List domains",
	Run: func(cmd *cobra.Command, args []string) {
		listDomains()
	},
}

var domainRemoveCmd = &cobra.Command{
	Use:   "remove [domain]",
	Short: "Remove a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		removeVhost(domain)
		fmt.Printf("Removed domain %s\n", domain)
	},
}

var domainEnableCmd = &cobra.Command{
	Use:   "enable [domain]",
	Short: "Enable a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		server := viper.GetString("domains." + domain + ".server")
		if server == "apache" {
			utils.RunCommand("sudo", "a2ensite", domain)
			utils.RunCommand("sudo", "systemctl", "reload", "apache2")
		} else if server == "nginx" {
			utils.RunCommand("sudo", "ln", "-sf", "/etc/nginx/sites-available/"+domain, "/etc/nginx/sites-enabled/"+domain)
			utils.RunCommand("sudo", "systemctl", "reload", "nginx")
		} else if server == "caddy" {
			utils.RunCommand("sudo", "systemctl", "reload", "caddy")
		}
		fmt.Printf("Enabled domain %s\n", domain)
	},
}

var domainDisableCmd = &cobra.Command{
	Use:   "disable [domain]",
	Short: "Disable a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		server := viper.GetString("domains." + domain + ".server")
		if server == "apache" {
			utils.RunCommand("sudo", "a2dissite", domain)
			utils.RunCommand("sudo", "systemctl", "reload", "apache2")
		} else if server == "nginx" {
			utils.RunCommand("sudo", "rm", "-f", "/etc/nginx/sites-enabled/"+domain)
			utils.RunCommand("sudo", "systemctl", "reload", "nginx")
		} else if server == "caddy" {
			utils.RunCommand("sudo", "systemctl", "reload", "caddy")
		}
		fmt.Printf("Disabled domain %s\n", domain)
	},
}

var domainSetRootCmd = &cobra.Command{
	Use:   "set-root [domain] [path]",
	Short: "Set document root for a domain",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		root := args[1]
		setDocumentRoot(domain, root)
		fmt.Printf("Set document root for %s to %s\n", domain, root)
	},
}

func AddDomainCommands(root *cobra.Command) {
	root.AddCommand(domainCmd)

	domainCmd.AddCommand(domainAddCmd)
	domainCmd.AddCommand(domainListCmd)
	domainCmd.AddCommand(domainRemoveCmd)
	domainCmd.AddCommand(domainEnableCmd)
	domainCmd.AddCommand(domainDisableCmd)
	domainCmd.AddCommand(domainSetRootCmd)

	domainAddCmd.Flags().String("server", "", "Web server (apache, nginx, caddy)")
}

func createVhost(domain, server string) {
	root := "/var/www/" + domain
	os.MkdirAll(root, 0755)

	viper.Set("domains."+domain+".server", server)
	viper.Set("domains."+domain+".root", root)
	viper.WriteConfig()

	var content string
	var file string
	if server == "apache" {
		content = fmt.Sprintf(`<VirtualHost *:80>
    ServerName %s
    DocumentRoot %s
    <Directory %s>
        AllowOverride All
        Require all granted
    </Directory>
</VirtualHost>`, domain, root, root)
		file = fmt.Sprintf("/etc/apache2/sites-available/%s.conf", domain)
	} else if server == "nginx" {
		content = fmt.Sprintf(`server {
    listen 80;
    server_name %s;
    root %s;
    index index.html index.htm;
    location / {
        try_files $uri $uri/ =404;
    }
}`, domain, root)
		file = fmt.Sprintf("/etc/nginx/sites-available/%s", domain)
	} else if server == "caddy" {
		content = fmt.Sprintf(`%s {
    root * %s
    file_server
}`, domain, root)
		file = fmt.Sprintf("/etc/caddy/sites/%s.caddyfile", domain)
		os.MkdirAll("/etc/caddy/sites", 0755)
	}
	ioutil.WriteFile(file, []byte(content), 0644)
}

func listDomains() {
	domains := viper.GetStringMap("domains")
	for domain := range domains {
		server := viper.GetString("domains." + domain + ".server")
		root := viper.GetString("domains." + domain + ".root")
		fmt.Printf("%s (%s) -> %s\n", domain, server, root)
	}
}

func removeVhost(domain string) {
	server := viper.GetString("domains." + domain + ".server")
	var file string
	if server == "apache" {
		file = fmt.Sprintf("/etc/apache2/sites-available/%s.conf", domain)
		utils.RunCommand("sudo", "a2dissite", domain)
	} else if server == "nginx" {
		file = fmt.Sprintf("/etc/nginx/sites-available/%s", domain)
		utils.RunCommand("sudo", "rm", "-f", "/etc/nginx/sites-enabled/"+domain)
	} else if server == "caddy" {
		file = fmt.Sprintf("/etc/caddy/sites/%s.caddyfile", domain)
	}
	os.Remove(file)
	// Remove from config
	viper.Set("domains."+domain, nil)
	viper.WriteConfig()
}

func setDocumentRoot(domain, root string) {
	server := viper.GetString("domains." + domain + ".server")
	var file string
	if server == "apache" {
		file = fmt.Sprintf("/etc/apache2/sites-available/%s.conf", domain)
		content, _ := ioutil.ReadFile(file)
		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			if strings.Contains(line, "DocumentRoot") {
				lines[i] = fmt.Sprintf("    DocumentRoot %s", root)
			}
		}
		ioutil.WriteFile(file, []byte(strings.Join(lines, "\n")), 0644)
	} else if server == "nginx" {
		file = fmt.Sprintf("/etc/nginx/sites-available/%s", domain)
		content, _ := ioutil.ReadFile(file)
		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			if strings.Contains(line, "root ") {
				lines[i] = fmt.Sprintf("    root %s;", root)
			}
		}
		ioutil.WriteFile(file, []byte(strings.Join(lines, "\n")), 0644)
	} else if server == "caddy" {
		file = fmt.Sprintf("/etc/caddy/sites/%s.caddyfile", domain)
		content, _ := ioutil.ReadFile(file)
		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			if strings.Contains(line, "root * ") {
				lines[i] = fmt.Sprintf("    root * %s", root)
			}
		}
		ioutil.WriteFile(file, []byte(strings.Join(lines, "\n")), 0644)
	}
	viper.Set("domains."+domain+".root", root)
	viper.WriteConfig()
}