package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"stackroost/config"
	"stackroost/internal"
)

var rootCmd = &cobra.Command{
	Use:   "stackroost",
	Short: "StackRoost CLI - manage your Linux servers with ease",
	Run: func(cmd *cobra.Command, args []string) {
		printWelcome()
	},
}

var createDomainCmd = &cobra.Command{
	Use:   "create-domain",
	Short: "Create a web server configuration for a domain",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("name")
		port, _ := cmd.Flags().GetString("port")
		serverType, _ := cmd.Flags().GetString("server")
		shellUser, _ := cmd.Flags().GetBool("shelluser")
		password, _ := cmd.Flags().GetString("pass")
		createDir, _ := cmd.Flags().GetBool("useridr")

		if internal.IsNilOrEmpty(domain) {
			fmt.Println("Error: --name flag is required and cannot be empty")
			os.Exit(1)
		}
		if internal.IsNilOrEmpty(port) {
			port = "80"
		}

		// Extract username from domain
		username := strings.Split(domain, ".")[0]

		// Check config existence first
		ext := ".conf"
		var configPath string

		switch serverType {
		case "nginx":
			configPath = filepath.Join("/etc/nginx/sites-available", domain+ext)
		default: // assume apache
			configPath = filepath.Join("/etc/apache2/sites-available", domain+ext)
		}

		if _, err := os.Stat(configPath); err == nil {
			fmt.Printf("  Configuration for '%s' already exists at %s\n", domain, configPath)
			fmt.Println(" Aborting to prevent overwriting existing configuration.")
			os.Exit(1)
		}

		fmt.Println(" Starting setup for domain:", domain)

		// Shell user creation
		if shellUser {
			if internal.IsNilOrEmpty(password) {
				fmt.Println(" Error: --pass is required when --shelluser is true")
				os.Exit(1)
			}

			fmt.Println("🔧 Creating system user:", username)

			userAddCmd := fmt.Sprintf("id -u %s || useradd -m -s /bin/bash %s", username, username)
			setPassCmd := fmt.Sprintf("echo '%s:%s' | chpasswd", username, password)

			if err := internal.RunCommand("sudo", "bash", "-c", userAddCmd); err != nil {
				fmt.Printf(" Failed to create user: %v\n", err)
				os.Exit(1)
			}

			if err := internal.RunCommand("sudo", "bash", "-c", setPassCmd); err != nil {
				fmt.Printf(" Failed to set password: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf(" User '%s' created with shell access\n", username)
		}

		// Create user directory
		if createDir {
			fmt.Println(" Creating public_html directory for user...")

			publicHtmlPath := fmt.Sprintf("/home/%s/public_html", username)
			if err := os.MkdirAll(publicHtmlPath, 0755); err != nil {
				fmt.Printf(" Failed to create directory: %v\n", err)
				os.Exit(1)
			}

			if err := internal.RunCommand("sudo", "chown", "-R", fmt.Sprintf("%s:%s", username, username), fmt.Sprintf("/home/%s", username)); err != nil {
				fmt.Printf(" Failed to assign ownership: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf(" Directory '%s' created and owned by '%s'\n", publicHtmlPath, username)
		}

		fmt.Println(" Generating Apache configuration...")

		configGen, err := config.NewWebServerConfig(serverType)
		if err != nil {
			fmt.Printf(" Error: %v\n", err)
			os.Exit(1)
		}

		// Use extracted username in DocumentRoot
		configContent, err := configGen.Generate(domain, port, username)
		if err != nil {
			fmt.Printf(" Error generating config: %v\n", err)
			os.Exit(1)
		}

		if err := writeConfigFile(domain, configContent, configGen.GetFileExtension()); err != nil {
			fmt.Printf(" Error writing config file: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(" Configuration file created.")

		filename := fmt.Sprintf("%s%s", domain, configGen.GetFileExtension())

		switch serverType {
		case "apache":
			fmt.Println(" Enabling site with a2ensite...")
			if err := internal.RunCommand("sudo", "a2ensite", filename); err != nil {
				fmt.Printf(" Failed to enable site: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Reloading Apache server...")
			if err := internal.RunCommand("sudo", "systemctl", "reload", "apache2"); err != nil {
				fmt.Printf(" Failed to reload apache: %v\n", err)
				os.Exit(1)
			}

		case "nginx":
			sitePath := filepath.Join("/etc/nginx/sites-available", filename)
			linkPath := filepath.Join("/etc/nginx/sites-enabled", filename)
			fmt.Println(" Enabling Nginx site...")
			if _, err := os.Stat(linkPath); os.IsNotExist(err) {
				if err := internal.RunCommand("sudo", "ln", "-s", sitePath, linkPath); err != nil {
					fmt.Printf(" Failed to enable nginx site: %v\n", err)
					os.Exit(1)
				}
			}

			fmt.Println("Reloading Nginx server...")
			if err := internal.RunCommand("sudo", "systemctl", "reload", "nginx"); err != nil {
				fmt.Printf(" Failed to reload nginx: %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Printf("🎉 %s configuration created and enabled for %s on port %s\n", serverType, domain, port)
	},
}

func init() {
	rootCmd.AddCommand(createDomainCmd)
	createDomainCmd.Flags().StringP("name", "n", "", "Domain name for the configuration (e.g., mahesh.spark.dev)")
	createDomainCmd.Flags().Bool("shelluser", false, "Create a shell user for the domain")
	createDomainCmd.Flags().String("pass", "", "Password for the shell user")
	createDomainCmd.Flags().Bool("useridr", false, "Create user directory /home/<user>/public_html")
	createDomainCmd.Flags().StringP("port", "p", "80", "Port for the configuration (default: 80)")
	createDomainCmd.Flags().StringP("server", "s", "apache", "Web server type (e.g., apache, nginx, caddy)")
	createDomainCmd.MarkFlagRequired("name")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func printWelcome() {
	fmt.Println("Welcome to StackRoost CLI!")
	fmt.Println("Your terminal assistant for managing Linux servers.")
}

func writeConfigFile(domain, content, extension string) error {
	var outputDir string
	if extension == ".conf" {
		if strings.HasPrefix(content, "server") {
			outputDir = "/etc/nginx/sites-available"
		} else {
			outputDir = "/etc/apache2/sites-available"
		}
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	filename := fmt.Sprintf("%s%s", domain, extension)
	outputPath := filepath.Join(outputDir, filename)

	if _, err := os.Stat(outputPath); err == nil {
		return fmt.Errorf("configuration for '%s' already exists at %s", domain, outputPath)
	}

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}
