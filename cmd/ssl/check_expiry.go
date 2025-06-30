package ssl

import (
	"fmt"
	"os"
	"crypto/tls"
	"time"
	"github.com/spf13/cobra"
	"stackroost/internal/logger"
)

var CheckSSLExpiryCmd = &cobra.Command{
	Use:   "check-ssl-expiry",
	Short: "Check the SSL certificate expiry for a domain",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		if domain == "" {
			logger.Error("Please provide a domain using --domain")
			os.Exit(1)
		}
		checkSSLExpiry(domain)
	},
}

func init() {
	CheckSSLExpiryCmd.Flags().String("domain", "", "Domain to check SSL expiry for")
	CheckSSLExpiryCmd.MarkFlagRequired("domain")
}

func checkSSLExpiry(domain string) {
	conn, err := tls.Dial("tcp", domain+":443", nil)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect: %v", err))
		os.Exit(1)
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		logger.Error("No SSL certificates found")
		os.Exit(1)
	}
	expiry := certs[0].NotAfter
	daysLeft := int(time.Until(expiry).Hours() / 24)

	logger.Info(fmt.Sprintf("SSL for %s expires on: %s (%d days left)", domain, expiry.Format(time.RFC1123), daysLeft))

	if daysLeft < 15 {
		logger.Warn("SSL certificate is expiring soon!")
	}
}
