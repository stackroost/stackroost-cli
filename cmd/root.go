package cmd

import (
    "fmt"
    "os"
    "path/filepath"
    "stackroost/config"
    "stackroost/internal"
    "github.com/spf13/cobra"
)

// rootCmd is the base command
var rootCmd = &cobra.Command{
    Use:   "stackroost",
    Short: "StackRoost CLI - manage your Linux servers with ease",
    Run: func(cmd *cobra.Command, args []string) {
        printWelcome()
    },
}

// createDomainCmd is the command to create a web server configuration
var createDomainCmd = &cobra.Command{
    Use:   "create-domain",
    Short: "Create a web server configuration for a domain",
    Run: func(cmd *cobra.Command, args []string) {
        domain, _ := cmd.Flags().GetString("name")
        port, _ := cmd.Flags().GetString("port")
        serverType, _ := cmd.Flags().GetString("server")

        if internal.IsNilOrEmpty(domain) {
            fmt.Println("Error: --name flag is required and cannot be empty")
            os.Exit(1)
        }
        if internal.IsNilOrEmpty(port) {
            port = "80" // Default port
        }

        // Create web server configuration generator
        configGen, err := config.NewWebServerConfig(serverType)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            os.Exit(1)
        }

        // Generate configuration
        configContent, err := configGen.Generate(domain, port)
        if err != nil {
            fmt.Printf("Error generating config: %v\n", err)
            os.Exit(1)
        }

        // Write configuration to file
        if err := writeConfigFile(domain, configContent, configGen.GetFileExtension()); err != nil {
            fmt.Printf("Error writing config file: %v\n", err)
            os.Exit(1)
        }

        filename := fmt.Sprintf("%s%s", domain, configGen.GetFileExtension())

        // Enable site using a2ensite
        if err := internal.RunCommand("sudo", "a2ensite", filename); err != nil {
            fmt.Printf("Failed to enable site: %v\n", err)
            os.Exit(1)
        }

        // Reload apache to apply changes
        if err := internal.RunCommand("sudo", "systemctl", "reload", "apache2"); err != nil {
            fmt.Printf("Failed to reload apache: %v\n", err)
            os.Exit(1)
        }

        fmt.Printf("%s configuration created and enabled for %s on port %s\n", serverType, domain, port)
    },
}

func init() {
    rootCmd.AddCommand(createDomainCmd)
    createDomainCmd.Flags().StringP("name", "n", "", "Domain name for the configuration (e.g., mahesh.spark.dev)")
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

// writeConfigFile writes the configuration to a file
func writeConfigFile(domain, content, extension string) error {
    outputDir := "/etc/apache2/sites-available"
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return fmt.Errorf("failed to create output directory: %v", err)
    }

    filename := fmt.Sprintf("%s%s", domain, extension)
    outputPath := filepath.Join(outputDir, filename)

    if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
        return fmt.Errorf("failed to write config file: %v", err)
    }

    return nil
}
