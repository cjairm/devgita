package commands

import "fmt"

type DebianCommand struct{}

func (d *DebianCommand) ExecCommand(cmd CommandParams) error {
	// cmd := exec.Command("apt-get", "install", "-y", packageName)
	// return cmd.Run()
	return nil
}

func (d *DebianCommand) PrintSomething() {
	fmt.Println("Executing command on Debian")
}
