package commands

import (
	"fmt"
)

type MacOSCommand struct{}

func (m *MacOSCommand) ExecCommand(cmd CommandParams) error {
	// cmd := exec.Command("mv", source, destination)
	// return cmd.Run()
	return nil
}

func (m *MacOSCommand) PrintSomething() {
	fmt.Println("Executing command on MacOS")
}
