package domain

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func GetMonitorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "monitor",
		Short: "Interactive dashboard to monitor all domains",
		Run: func(cmd *cobra.Command, args []string) {
			p := tea.NewProgram(newModel())
			if err := p.Start(); err != nil {
				fmt.Printf("Error running TUI: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

type domainRow struct {
	Domain     string
	User       string
	Enabled    string
	DiskUsed   string
	LastLogin  string
	ServerType string
}

type model struct {
	table table.Model
}

var columns = []table.Column{
	{Title: "Domain", Width: 22},
	{Title: "User", Width: 12},
	{Title: "Enabled", Width: 8},
	{Title: "Disk Used", Width: 10},
	{Title: "Last Login", Width: 14},
	{Title: "Server", Width: 10},
}

func newModel() model {
	rows := fetchDomainRows()
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)
	t.SetStyles(defaultTableStyles())
	return model{table: t}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			// Refresh data
			m.table.SetRows(fetchDomainRows())
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render("StackRoost Monitor")
	return fmt.Sprintf("\n%s\n\n%s\n\n[↑↓ to scroll] [r]efresh [q]uit\n", title, m.table.View())
}

func fetchDomainRows() []table.Row {
	servers := map[string]string{
		"apache2": "/etc/apache2/sites-available",
		"nginx":   "/etc/nginx/sites-available",
		"caddy":   "/etc/caddy/sites-available",
	}
	var rows []table.Row

	for server, basePath := range servers {
		files, err := os.ReadDir(basePath)
		if err != nil {
			continue
		}

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".conf") {
				continue
			}

			// Exclude common system/default configs
			lower := strings.ToLower(file.Name())
			if strings.HasPrefix(lower, "000") ||
				strings.HasPrefix(lower, "default") ||
				strings.HasPrefix(lower, "template") {
				continue
			}

			domain := strings.TrimSuffix(file.Name(), ".conf")
			username := strings.Split(domain, ".")[0]
			row := domainRow{
				Domain:     domain,
				User:       username,
				Enabled:    checkEnabled(server, file.Name()),
				DiskUsed:   getDiskUsage(username),
				LastLogin:  getLastLogin(username),
				ServerType: server,
			}
			rows = append(rows, table.Row{
				row.Domain, row.User, row.Enabled, row.DiskUsed, row.LastLogin, row.ServerType,
			})
		}
	}

	return rows
}

func checkEnabled(server, filename string) string {
	var path string
	switch server {
	case "apache2":
		path = "/etc/apache2/sites-enabled/" + filename
	case "nginx":
		path = "/etc/nginx/sites-enabled/" + filename
	case "caddy":
		path = "/etc/caddy/sites-enabled/" + filename
	}
	if _, err := os.Lstat(path); err == nil {
		return "Yes"
	}
	return "No"
}

func getDiskUsage(user string) string {
	home := "/home/" + user
	cmd := exec.Command("du", "-sh", home)
	out, err := cmd.Output()
	if err != nil {
		return "-"
	}
	return strings.Fields(string(out))[0]
}

func getLastLogin(user string) string {
	cmd := exec.Command("lastlog", "-u", user)
	out, err := cmd.Output()
	if err != nil {
		return "-"
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) > 1 && !strings.Contains(lines[1], "**Never logged in**") {
		fields := strings.Fields(lines[1])
		if len(fields) >= 5 {
			return fmt.Sprintf("%s %s %s", fields[3], fields[4], fields[5])
		}
	}
	return "Never"
}

func defaultTableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.Bold(true).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Bold(true)
	return s
}
