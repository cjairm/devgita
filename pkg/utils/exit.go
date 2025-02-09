package utils

import (
	"os"
)

func MaybeExitWithError(err error) {
	if err == nil {
		return
	}
	PrintError(err.Error())
	os.Exit(1)
}
