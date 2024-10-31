package debian

import "fmt"

func PreInstall() {
	fmt.Println("Make sure to install: ")
	fmt.Println("sudo apt-get update >/dev/null")
	fmt.Println("sudo apt-get install -y git >/dev/null")
}
