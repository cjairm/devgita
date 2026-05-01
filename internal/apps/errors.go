package apps

import "errors"

var (
	ErrUninstallNotSupported = errors.New("uninstall not supported")
	ErrUpdateNotSupported    = errors.New("update not supported")
	ErrConfigureNotSupported = errors.New("configure not supported")
	ErrExecuteNotSupported   = errors.New("execute not supported")
)
