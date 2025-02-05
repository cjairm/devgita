package common

import (
	"fmt"
	"runtime"
)

func GitCommand(gitArgs ...string) error {
	if len(gitArgs) == 0 {
		return fmt.Errorf("No command provided")
	}
	switch runtime.GOOS {
	case "darwin":
		cmd := CommandInfo{
			PreExecutionMessage:  "",
			PostExecutionMessage: "",
			IsSudo:               false,
			Command:              "git",
			Args:                 gitArgs,
		}
		return ExecCommand(cmd)
	case "linux":
		fmt.Println("Linux")
	default:
		fmt.Printf("Unsupported operating system: %s\n", runtime.GOOS)
	}
	return nil
}
