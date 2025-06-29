package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

var deleteUserCmd = &cobra.Command{
	Use:   "delete-user",
	Short: "Delete a system user and optionally remove their home directory",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("user")
		removeHome, _ := cmd.Flags().GetBool("remove-home")

		if internal.IsNilOrEmpty(username) {
			logger.Error("Please provide a username using --user")
			os.Exit(1)
		}

		// Check if user exists
		_, err := user.Lookup(username)
		if err != nil {
			logger.Warn(fmt.Sprintf("User '%s' does not exist", username))
			return
		}

		logger.Info(fmt.Sprintf("Deleting user: %s", username))

		// Build command args
		cmdArgs := []string{"userdel"}
		if removeHome {
			cmdArgs = append(cmdArgs, "-r")
		}
		cmdArgs = append(cmdArgs, username)

		// Execute command
		if err := internal.RunCommand("sudo", cmdArgs...); err != nil {
			logger.Error(fmt.Sprintf("Failed to delete user %s: %v", username, err))
			os.Exit(1)
		}

		logger.Success(fmt.Sprintf("User '%s' deleted", username))
	},
}

func init() {
	rootCmd.AddCommand(deleteUserCmd)
	deleteUserCmd.Flags().String("user", "", "Username to delete")
	deleteUserCmd.Flags().Bool("remove-home", false, "Remove the user's home directory")
	deleteUserCmd.MarkFlagRequired("user")
}
