package email

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"stackroost/internal/logger"
)


var TestEmailCmd = &cobra.Command{
	Use:   "test-email",
	Short: "Check if email capability is available on this system (mail/sendmail/msmtp)",
	Run: func(cmd *cobra.Command, args []string) {
		mailers := []string{"mail", "sendmail", "msmtp"}
		found := false

		for _, bin := range mailers {
			_, err := exec.LookPath(bin)
			if err == nil {
				logger.Success(fmt.Sprintf("Mailer available: %s", bin))
				found = true
			}
		}

		if !found {
			logger.Warn("No mail sending utilities found (mail/sendmail/msmtp)")
			logger.Info("You can install one, e.g., `sudo apt install mailutils` or `sendmail`")
			os.Exit(1)
		}

		logger.Info("Email sending capability appears to be available")
	},
}

func init() {
	TestEmailCmd.Flags().String("to", "", "Recipient email address")
	TestEmailCmd.MarkFlagRequired("to")
}
