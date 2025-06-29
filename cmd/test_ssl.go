package cmd

import (
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

var testSSLCmd = &cobra.Command{
	Use:   "test-ssl",
	Short: "Check SSL certificate status for a domain",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		port, _ := cmd.Flags().GetString("port")

		if internal.IsNilOrEmpty(domain) {
			logger.Error("Please provide a domain using --domain")
			os.Exit(1)
		}

		address := fmt.Sprintf("%s:%s", domain, port)
		logger.Info(fmt.Sprintf("Testing SSL certificate for %s...", domain))

		conn, err := tls.Dial("tcp", address, nil)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to connect to %s: %v", address, err))
			os.Exit(1)
		}
		defer conn.Close()

		certs := conn.ConnectionState().PeerCertificates
		if len(certs) == 0 {
			logger.Error("No certificate found")
			os.Exit(1)
		}

		cert := certs[0]
		now := time.Now()
		if now.Before(cert.NotBefore) || now.After(cert.NotAfter) {
			logger.Error("SSL certificate is invalid or expired ")
		} else {
			logger.Success("SSL is valid ")
			logger.Info(fmt.Sprintf("Issuer: %s", cert.Issuer.CommonName))
			logger.Info(fmt.Sprintf("Expires: %s (in %d days)", cert.NotAfter.Format(time.RFC1123), int(cert.NotAfter.Sub(now).Hours()/24)))
		}
	},
}

func init() {
	rootCmd.AddCommand(testSSLCmd)
	testSSLCmd.Flags().String("domain", "", "Domain to test (required)")
	testSSLCmd.Flags().String("port", "443", "Port to test SSL (default: 443)")
	testSSLCmd.MarkFlagRequired("domain")
}
