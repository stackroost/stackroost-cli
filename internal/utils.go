package internal

import (
    "fmt"
    "os/exec"
)

// RunCommand runs a shell command with args and returns error if any
func RunCommand(name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("command %s %v failed: %w", name, args, err)
    }
    return nil
}

// IsNilOrEmpty returns true if the string is empty or "<nil>"
func IsNilOrEmpty(s string) bool {
    return s == "" || s == "<nil>"
}
