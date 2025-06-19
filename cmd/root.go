package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"stackroost/config"
	"stackroost/internal"
	"stackroost/internal/logger"
	"strings"
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
		logger.Info("Starting create-domain command execution")

		domain, _ := cmd.Flags().GetString("name")
		port, _ := cmd.Flags().GetString("port")
		serverType, _ := cmd.Flags().GetString("server")
		shellUser, _ := cmd.Flags().GetBool("shelluser")
		password, _ := cmd.Flags().GetString("pass")
		createDir, _ := cmd.Flags().GetBool("useridr")

		if internal.IsNilOrEmpty(domain) {
			logger.Error("Domain name is required and cannot be empty")
			os.Exit(1)
		}

		if internal.IsNilOrEmpty(port) {
			logger.Info("No port specified, defaulting to 80")
			port = "80"
		}

		username := strings.Split(domain, ".")[0]
		logger.Debug(fmt.Sprintf("Extracted username: %s from domain: %s", username, domain))

		ext := ".conf"
		var configPath string

		switch serverType {
		case "nginx":
			configPath = filepath.Join("/etc/nginx/sites-available", domain+ext)
		case "caddy":
			configPath = filepath.Join("/etc/caddy/sites-available", domain+ext)
		case "apache":
			configPath = filepath.Join("/etc/apache2/sites-available", domain+ext)
		default:
			logger.Error(fmt.Sprintf("Unsupported server type: %s. Supported types are: apache, nginx, caddy", serverType))
			os.Exit(1)
		}

		if _, err := os.Stat(configPath); err == nil {
			logger.Error(fmt.Sprintf("Configuration for '%s' already exists at %s", domain, configPath))
			os.Exit(1)
		}

		logger.Info(fmt.Sprintf("Initiating setup for domain: %s with server type: %s", domain, serverType))

		if shellUser {
			if internal.IsNilOrEmpty(password) {
				logger.Error("Password is required when shelluser is enabled")
				os.Exit(1)
			}

			logger.Info(fmt.Sprintf("Creating system user: %s", username))

			userAddCmd := fmt.Sprintf("id -u %s || useradd -m -s /bin/bash %s", username, username)
			setPassCmd := fmt.Sprintf("echo '%s:%s' | chpasswd", username, password)

			if err := internal.RunCommand("sudo", "bash", "-c", userAddCmd); err != nil {
				logger.Error(fmt.Sprintf("Failed to create user %s: %v", username, err))
				os.Exit(1)
			}

			if err := internal.RunCommand("sudo", "bash", "-c", setPassCmd); err != nil {
				logger.Error(fmt.Sprintf("Failed to set password for user %s: %v", username, err))
				os.Exit(1)
			}

			logger.Success(fmt.Sprintf("User '%s' created with shell access", username))
		}

		if createDir {
			logger.Info("Creating public_html directory for user")

			publicHtmlPath := fmt.Sprintf("/home/%s/public_html", username)
			if err := os.MkdirAll(publicHtmlPath, 0755); err != nil {
				logger.Error(fmt.Sprintf("Failed to create directory %s: %v", publicHtmlPath, err))
				os.Exit(1)
			}

			if err := internal.RunCommand("sudo", "chown", "-R", fmt.Sprintf("%s:%s", username, username), fmt.Sprintf("/home/%s", username)); err != nil {
				logger.Error(fmt.Sprintf("Failed to assign ownership for %s: %v", username, err))
				os.Exit(1)
			}

			logger.Success(fmt.Sprintf("Directory '%s' created and owned by '%s'", publicHtmlPath, username))
		}

		logger.Info(fmt.Sprintf("Generating %s configuration", serverType))

		configGen, err := config.NewWebServerConfig(serverType)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to create config generator: %v", err))
			os.Exit(1)
		}

		configContent, err := configGen.Generate(domain, port, username)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to generate config for %s: %v", domain, err))
			os.Exit(1)
		}

		if err := writeConfigFile(domain, configContent, configGen.GetFileExtension()); err != nil {
			logger.Error(fmt.Sprintf("Failed to write config file: %v", err))
			os.Exit(1)
		}

		if err := internal.CreateMySQLUserAndDatabase(username, password); err != nil {
			logger.Error(fmt.Sprintf("Database setup failed for %s: %v", username, err))
			os.Exit(1)
		}

		logger.Success("Configuration file created")

		filename := fmt.Sprintf("%s%s", domain, configGen.GetFileExtension())

		switch serverType {
		case "apache":
			logger.Info("Enabling site with a2ensite")
			if err := internal.RunCommand("sudo", "a2ensite", filename); err != nil {
				logger.Error(fmt.Sprintf("Failed to enable site %s: %v", filename, err))
				os.Exit(1)
			}

			logger.Info("Reloading Apache server")
			if err := internal.RunCommand("sudo", "systemctl", "reload", "apache2"); err != nil {
				logger.Error(fmt.Sprintf("Failed to reload Apache: %v", err))
				os.Exit(1)
			}

		case "nginx":
			sitePath := filepath.Join("/etc/nginx/sites-available", filename)
			linkPath := filepath.Join("/etc/nginx/sites-enabled", filename)
			logger.Info("Enabling Nginx site")
			if _, err := os.Stat(linkPath); os.IsNotExist(err) {
				if err := internal.RunCommand("sudo", "ln", "-s", sitePath, linkPath); err != nil {
					logger.Error(fmt.Sprintf("Failed to enable Nginx site %s: %v", filename, err))
					os.Exit(1)
				}
			}

			logger.Info("Reloading Nginx server")
			if err := internal.RunCommand("sudo", "systemctl", "reload", "nginx"); err != nil {
				logger.Error(fmt.Sprintf("Failed to reload Nginx: %v", err))
				os.Exit(1)
			}

		case "caddy":
			sitePath := filepath.Join("/etc/caddy/sites-available", filename)
			linkPath := filepath.Join("/etc/caddy/sites-enabled", filename)
			logger.Info("Enabling Caddy site")
			if _, err := os.Stat(linkPath); os.IsNotExist(err) {
				if err := internal.RunCommand("sudo", "ln", "-s", sitePath, linkPath); err != nil {
					logger.Error(fmt.Sprintf("Failed to enable Caddy site %s: %v", filename, err))
					os.Exit(1)
				}
			}

			logger.Info("Reloading Caddy server")
			if err := internal.RunCommand("sudo", "systemctl", "reload", "caddy"); err != nil {
				logger.Error(fmt.Sprintf("Failed to reload Caddy: %v", err))
				os.Exit(1)
			}
		}

		enableSSL, _ := cmd.Flags().GetBool("ssl")
		if enableSSL && (serverType == "apache" || serverType == "nginx") {
			err := internal.EnableSSLCertbot(domain, serverType)
			if err != nil {
				logger.Error(fmt.Sprintf("SSL setup failed: %v", err))
				os.Exit(1)
			}
		}

		logger.Success(fmt.Sprintf("%s configuration created and enabled for %s on port %s", serverType, domain, port))
	},
}

func init() {
	rootCmd.AddCommand(createDomainCmd)
	createDomainCmd.Flags().StringP("name", "n", "", "Domain name for the configuration (e.g., example.com)")
	createDomainCmd.Flags().Bool("shelluser", false, "Create a shell user for the domain")
	createDomainCmd.Flags().String("pass", "", "Password for the shell user")
	createDomainCmd.Flags().Bool("useridr", false, "Create user directory /home/<user>/public_html")
	createDomainCmd.Flags().StringP("port", "p", "80", "Port for the configuration (default: 80)")
	createDomainCmd.Flags().StringP("server", "s", "apache", "Web server type (e.g., apache, nginx, caddy)")
	createDomainCmd.Flags().Bool("ssl", false, "Enable Let's Encrypt SSL (Apache/Nginx only)")
	createDomainCmd.MarkFlagRequired("name")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error(fmt.Sprintf("Command execution failed: %v", err))
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func printWelcome() {
	reset := "\033[0m"
	bold := "\033[1m"
	gray := "\033[38;2;180;180;180m"

	title := "WELCOME TO STACKROOST CLI"
	subTitle := "Smart Linux Server Manager"

	startR, startG, startB := 135, 206, 235
	endR, endG, endB := 255, 0, 0

	length := len(title)

	fmt.Println()
	fmt.Print(bold)

	for i, ch := range title {
		t := float64(i) / float64(length-1)
		r := int(float64(startR)*(1-t) + float64(endR)*t)
		g := int(float64(startG)*(1-t) + float64(endG)*t)
		b := int(float64(startB)*(1-t) + float64(endB)*t)
		fmt.Printf("\033[38;2;%d;%d;%dm%c", r, g, b, ch)
	}

	fmt.Println(reset)
	fmt.Println()

	fmt.Printf("%s\033[38;2;135;206;235m%s%s\n\n", bold, subTitle, reset)

	fmt.Println(gray + "Welcome to StackRoost — your powerful CLI for managing Linux domains," + reset)
	fmt.Println(gray + "creating shell users, and configuring Apache · Nginx · Caddy effortlessly." + reset)
	fmt.Println()

	fmt.Printf("%s\033[38;2;200;200;200mtry: %sstackroost --help%s\n\n", bold, reset, reset)
}

func writeConfigFile(domain, content, extension string) error {
	logger.Info(fmt.Sprintf("Writing configuration file for %s", domain))

	var configPath string
	if extension == ".conf" {
		if strings.HasPrefix(content, "server") {
			configPath = "/etc/nginx/sites-available"
			logger.Debug("Detected nginx configuration")
		} else if strings.Contains(content, "php_fastcgi") {
			configPath = "/etc/caddy/sites-available"
			logger.Debug("Detected caddy configuration")
		} else {
			configPath = "/etc/apache2/sites-available"
			logger.Debug("Detected apache configuration")
		}
	}

	if err := os.MkdirAll(configPath, 0755); err != nil {
		logger.Error(fmt.Sprintf("Failed to create output directory %s: %v", configPath, err))
		return err
	}

	filename := fmt.Sprintf("%s%s", domain, extension)
	outputPath := filepath.Join(configPath, filename)

	if _, err := os.Stat(outputPath); err == nil {
		logger.Error(fmt.Sprintf("Configuration already exists at %s for domain %s", outputPath, domain))
		return fmt.Errorf("configuration exists")
	}

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		logger.Error(fmt.Sprintf("Failed to write config file %s: %v", outputPath, err))
		return err
	}

	logger.Success(fmt.Sprintf("Configuration file written to %s", outputPath))
	return nil
}
