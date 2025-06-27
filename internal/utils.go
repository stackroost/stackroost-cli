package internal

import (
    "fmt"
    "os/exec"
    "os"
    "path/filepath"

)

func RunCommand(name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("command %s %v failed: %w", name, args, err)
    }
    return nil
}

func IsNilOrEmpty(s string) bool {
    return s == "" || s == "<nil>"
}

func DetectServerType(domain string) string {
	filename := domain + ".conf"

	paths := map[string]string{
		"apache": "/etc/apache2/sites-available",
		"nginx":  "/etc/nginx/sites-available",
		"caddy":  "/etc/caddy/sites-available",
	}

	for server, dir := range paths {
		if _, err := os.Stat(filepath.Join(dir, filename)); err == nil {
			return server
		}
	}
	return ""
}