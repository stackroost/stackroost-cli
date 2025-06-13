package internal

import (
    "fmt"
    "os/exec"
    "os"
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
