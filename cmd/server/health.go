package server

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/spf13/cobra"
	"stackroost/internal/logger"
)

func GetHealthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server-health",
		Short: "Check system resource usage and web server status",
		Run: func(cmd *cobra.Command, args []string) {
			printHostname()
			printUptime()
			printCPU()
			printMemory()
			printDisk()
			printWebServerStatus("apache2")
			printWebServerStatus("nginx")
			printWebServerStatus("caddy")
		},
	}
}

func printHostname() {
	info, err := host.Info()
	if err != nil {
		logger.Warn(fmt.Sprintf("Hostname fetch error: %v", err))
		return
	}
	logger.Info(fmt.Sprintf("Hostname: %s", info.Hostname))
}

func printUptime() {
	uptime, err := host.Uptime()
	if err != nil {
		logger.Warn(fmt.Sprintf("Uptime fetch error: %v", err))
		return
	}
	logger.Info(fmt.Sprintf("Uptime: %s", time.Duration(uptime)*time.Second))
}

func printCPU() {
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		logger.Warn(fmt.Sprintf("CPU usage error: %v", err))
		return
	}
	if len(percentages) > 0 {
		logger.Info(fmt.Sprintf("CPU Usage: %.2f%%", percentages[0]))
	}
}

func printMemory() {
	v, err := mem.VirtualMemory()
	if err != nil {
		logger.Warn(fmt.Sprintf("Memory usage error: %v", err))
		return
	}
	logger.Info(fmt.Sprintf("Memory Usage: %.2f%% of %.2f GB", v.UsedPercent, float64(v.Total)/1e9))
}

func printDisk() {
	usage, err := disk.Usage("/")
	if err != nil {
		logger.Warn(fmt.Sprintf("Disk usage error: %v", err))
		return
	}
	logger.Info(fmt.Sprintf("Disk Usage: %.2f%% of %.2f GB", usage.UsedPercent, float64(usage.Total)/1e9))
}

func printWebServerStatus(service string) {
	cmd := exec.Command("systemctl", "is-active", service)
	output, err := cmd.Output()
	status := strings.TrimSpace(string(output))

	if err == nil && status == "active" {
		logger.Success(fmt.Sprintf("%s: Running", strings.Title(service)))
	} else {
		logger.Warn(fmt.Sprintf("%s: Inactive or not installed", strings.Title(service)))
	}
}
