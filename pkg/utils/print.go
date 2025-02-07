package utils

import (
	"fmt"
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
