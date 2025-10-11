/*
Copyright © 2025 Stackroost CLI

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"stackroost-cli/cmd/domain"
	"stackroost-cli/cmd/internal/logger"
	"stackroost-cli/cmd/logs"
	"stackroost-cli/cmd/remote"
	"stackroost-cli/cmd/security"
	"stackroost-cli/cmd/server"
)



var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stackroost",
	Short: "Cross-server CLI manager for web servers, domains, SSL, and user management",
	Long: `Stackroost is a powerful CLI tool for managing web servers (Apache, Nginx, Caddy),
domains, SSL certificates, user access, and multi-server deployments across local and remote systems.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	PrintBanner()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.stackroost.yaml)")

	domain.AddDomainCommands(rootCmd)
	server.AddServerCmd(rootCmd)
	security.AddSSLCmd(rootCmd)
	security.AddUserCmd(rootCmd)
	remote.AddRemoteCmd(rootCmd)
	logs.AddLogsCmd(rootCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".stackroost")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		logger.Info(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
	}
}


