package logger

import (
	"fmt"
	"time"
)

const (
	Reset  = "\033[0m"
	Gray   = "\033[38;2;160;160;160m"
	Green  = "\033[38;2;0;220;100m"
	Yellow = "\033[38;2;255;204;0m"
	Red    = "\033[38;2;255;80;80m"
	Cyan   = "\033[38;2;135;206;235m"
	Blue   = "\033[38;2;100;180;255m"
	Bold   = "\033[1m"
)

func timeStamp() string {
	return fmt.Sprintf("%s[%s]%s", Gray, time.Now().Format("15:04:05"), Reset)
}

func log(label string, labelColor string, message string) {
	fmt.Printf("%s %s%-8s%s %s\n", timeStamp(), labelColor, "["+label+"]", Reset, message)
}

func Info(msg string) {
	log("INFO", Cyan, msg)
}

func Success(msg string) {
	log("SUCCESS", Green, msg)
}

func Warn(msg string) {
	log("WARN", Yellow, msg)
}

func Error(msg string) {
	log("ERROR", Red, msg)
}

func Debug(msg string) {
	log("DEBUG", Blue, msg)
}