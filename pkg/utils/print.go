package utils

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func PrintError(errMsg string) {
	Print(errMsg, Red)
}

func PrintSuccess(errMsg string) {
	Print(errMsg, Green)
}

func PrintSecondary(msg string) {
	Print(msg, Gray)
}

func PrintInfo(msg string) {
	Print(msg, Blue)
}

func PrintWarning(msg string) {
	Print(msg, Yellow)
}

func PrintBold(msg string) {
	Print(msg, Bold)
}

func Print(msg, custom string) {
	if msg != "" {
		if custom == "" {
			fmt.Println(msg)
		} else {
			fmt.Printf("%s%s%s\n", custom, msg, Reset)
		}
	}
}

func PrompCustomHelp(cmd *cobra.Command, args []string) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Usage:\n  %s\n\n", cmd.Use))
	sb.WriteString(cmd.Long + "\n\n")
	PrintBold(sb.String())
}
