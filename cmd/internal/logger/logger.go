package logger

import "github.com/fatih/color"

// Info message (blue)
func Info(message string) {
	color.New(color.FgBlue, color.Bold).Printf("[INFO] ")
	color.New(color.FgWhite).Println(message)
}

// Success message (green)
func Success(message string) {
	color.New(color.FgGreen, color.Bold).Printf("[SUCCESS] ")
	color.New(color.FgWhite).Println(message)
}

// Error message (red)
func Error(message string) {
	color.New(color.FgRed, color.Bold).Printf("[ERROR] ")
	color.New(color.FgWhite).Println(message)
}
