package logs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/spf13/cobra"
	"stackroost/internal"
	"stackroost/internal/logger"
)

var AnalyzeTrafficCmd = &cobra.Command{
	Use:   "analyze-log-traffic",
	Short: "Analyze traffic from a domain's access log (IP, URLs, request count)",
	Run:   runAnalyzeTraffic,
}

func init() {
	AnalyzeTrafficCmd.Flags().String("domain", "", "Domain to analyze")
	AnalyzeTrafficCmd.Flags().Int("lines", 1000, "Number of lines to scan from the log file")
	AnalyzeTrafficCmd.MarkFlagRequired("domain")
}

// Main logic
func runAnalyzeTraffic(cmd *cobra.Command, args []string) {
	domain, _ := cmd.Flags().GetString("domain")
	lines, _ := cmd.Flags().GetInt("lines")

	if internal.IsNilOrEmpty(domain) {
		logger.Error("Please provide a domain using --domain")
		os.Exit(1)
	}

	server := internal.DetectServerType(domain)
	if internal.IsNilOrEmpty(server) {
		logger.Warn("Could not detect server type for domain")
		os.Exit(1)
	}

	var accessLog string
	switch server {
	case "apache":
		accessLog = fmt.Sprintf("/var/log/apache2/%s-access.log", domain)
	case "nginx":
		accessLog = fmt.Sprintf("/var/log/nginx/%s.access.log", domain)
	case "caddy":
		accessLog = fmt.Sprintf("/var/log/caddy/%s.access.log", domain)
	default:
		logger.Warn("Unsupported server type for log analysis")
		os.Exit(1)
	}

	file, err := os.Open(filepath.Clean(accessLog))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to open access log: %v", err))
		os.Exit(1)
	}
	defer file.Close()

	ipCount := make(map[string]int)
	urlCount := make(map[string]int)

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		if lines > 0 && lineNum > lines {
			break
		}

		re := regexp.MustCompile(`^(\S+) .*?"\S+ (\S+) .*?" (\d{3})`)
		match := re.FindStringSubmatch(line)
		if len(match) >= 4 {
			ip := match[1]
			url := match[2]
			ipCount[ip]++
			urlCount[url]++
		}
	}

	logger.Info("Top 5 IPs:")
	printTopN(ipCount, 5)

	logger.Info("Top 5 URLs:")
	printTopN(urlCount, 5)
}

func printTopN(data map[string]int, n int) {
	type kv struct {
		Key   string
		Value int
	}

	var sorted []kv
	for k, v := range data {
		sorted = append(sorted, kv{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	for i := 0; i < n && i < len(sorted); i++ {
		fmt.Printf("  %s → %d\n", sorted[i].Key, sorted[i].Value)
	}
}

// Exportable getter
func GetAnalyzeCmd() *cobra.Command {
	return AnalyzeTrafficCmd
}
