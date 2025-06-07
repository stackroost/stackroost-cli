package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// rootCmd is the base command
var rootCmd = &cobra.Command{
	Use:   "stackrooy",
	Short: "StackRooy CLI - manage your Linux servers with ease",
	Run: func(cmd *cobra.Command, args []string) {
		printWelcome()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}

func printWelcome() {
	fmt.Println("Welcome to StackRoot CLI!")
	fmt.Println("Your terminal assistant for managing Linux servers.")
}
