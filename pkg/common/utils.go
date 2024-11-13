package common

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
)

var Devgita = fmt.Sprintf(`
%s
    .___                .__  __          
  __| _/_______  ______ |__|/  |______   
 / __ |/ __ \  \/ / ___\|  \   __\__  \  
/ /_/ \  ___/\   / /_/  >  ||  |  / __ \_
\____ |\___  >\_/\___  /|__||__| (____  /
     \/    \/   /_____/               \/ 
@cjairm
%s`, "\033[1m", "\033[0m")

func ExecCommand(startMessage string, endMessage string, name string, args ...string) error {
	fmt.Printf(startMessage + "\n")
	cmd := exec.Command(name, args...)

	// Create pipes to capture standard output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// Create pipes to capture standard error
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Color("bold", "fgMagenta")
	s.Start()

	// Start the command
	if err := cmd.Start(); err != nil {
		return err
	}

	// Create a scanner for stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			s.Suffix = "\n" + scanner.Text()
		}
	}()

	// Create a scanner for stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			s.Suffix = "\n" + scanner.Text()
		}
	}()

	err = cmd.Wait()

	s.Stop()

	fmt.Printf(endMessage + "\n\n")

	return err
}

// copyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, os.ModePerm)
}

func MkdirOrCopyFile(
	filePath string,
	dirPath string,
	devgitaPathFile string,
	appName string,
) error {
	// Check if file already exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
		return CopyFile(devgitaPathFile, filePath)
	}
	fmt.Printf("File for %s already exists.\n\n", appName)
	return nil
}

func IsCommandInstalled(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
