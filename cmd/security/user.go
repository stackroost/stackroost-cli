/*
Copyright © 2025 Stackroost CLI

*/
package security

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"stackroost-cli/cmd/internal/utils"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage system users",
	Long:  `Commands for creating, listing, modifying, and removing system users with SSH access.`,
}

var userAddCmd = &cobra.Command{
	Use:   "add [username]",
	Short: "Add a new user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		password, _ := cmd.Flags().GetString("password")
		utils.RunCommand("sudo", "useradd", "-m", username)
		if password != "" {
			utils.RunCommand("sudo", "chpasswd", fmt.Sprintf("%s:%s", username, password))
		}
		fmt.Printf("Added user %s\n", username)
	},
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	Run: func(cmd *cobra.Command, args []string) {
		execCmd := exec.Command("cut", "-d:", "-f1", "/etc/passwd")
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		execCmd.Run()
	},
}

var userPasswdCmd = &cobra.Command{
	Use:   "passwd [username]",
	Short: "Change user password",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		utils.RunCommand("sudo", "passwd", username)
		fmt.Printf("Password changed for %s\n", username)
	},
}

var userRemoveCmd = &cobra.Command{
	Use:   "remove [username]",
	Short: "Remove a user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		utils.RunCommand("sudo", "userdel", "-r", username)
		fmt.Printf("Removed user %s\n", username)
	},
}

var userSSHEnableCmd = &cobra.Command{
	Use:   "ssh-enable [username]",
	Short: "Enable SSH for user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		utils.RunCommand("sudo", "usermod", "-aG", "ssh", username)
		fmt.Printf("SSH enabled for %s\n", username)
	},
}

var userSSHDisableCmd = &cobra.Command{
	Use:   "ssh-disable [username]",
	Short: "Disable SSH for user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		utils.RunCommand("sudo", "gpasswd", "-d", username, "ssh")
		fmt.Printf("SSH disabled for %s\n", username)
	},
}

func AddUserCmd(root *cobra.Command) {
	root.AddCommand(userCmd)

	userCmd.AddCommand(userAddCmd)
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userPasswdCmd)
	userCmd.AddCommand(userRemoveCmd)
	userCmd.AddCommand(userSSHEnableCmd)
	userCmd.AddCommand(userSSHDisableCmd)

	userAddCmd.Flags().String("password", "", "User password")
}