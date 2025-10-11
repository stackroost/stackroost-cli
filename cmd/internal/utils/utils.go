package utils

import (
	"fmt"
	"os/exec"
)

func RunCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil // or os.Stdout if want to print
	cmd.Stderr = nil
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	}
}