package commands

import "runtime"

// cmdImpl := commands.NewCommand()
// cmdImpl.PrintSomething()
func NewCommand() Command {
	switch runtime.GOOS {
	case "darwin":
		return &MacOSCommand{}
	// TODO: Is it possible to detect the distribution of Linux?
	case "linux":
		return &DebianCommand{}
	default:
		panic("unsupported operating system")
	}
}
