package cmd

import (
	"fmt"
	"github.com/fatih/color"
)

func PrintBanner() {
	title := color.New(color.FgCyan, color.Bold).SprintFunc()
	sub   := color.New(color.FgWhite).SprintFunc()

	fmt.Println(title("\n*[{   Stackroost   }]*"))
	fmt.Println(sub("   Cross-Server CLI • Apache • Nginx • Caddy • SSL • Domains\n"))
}
