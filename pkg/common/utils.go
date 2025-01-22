package common

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	Reset = "\033[0m"
	Gray  = "\033[90m"
	Red   = "\033[31m"
)

type CommandInfo struct {
	PreExecutionMessage  string
	PostExecutionMessage string
	IsSudo               bool
	Command              string
	Args                 []string
}

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

func ExecCommand(cmdInfo CommandInfo) error {
	fmt.Print(cmdInfo.PreExecutionMessage + "\n")
	command := cmdInfo.Command
	if cmdInfo.IsSudo {
		command = "sudo " + command
	}
	cmd := exec.Command(command, cmdInfo.Args...)

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

	// To allow user to interact with the command (password)
	cmd.Stdin = os.Stdin

	// Start the command
	if err := cmd.Start(); err != nil {
		return err
	}

	// Create a scanner for stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Printf(Gray+"%s"+Reset+"\n", scanner.Text())
		}
	}()

	// Create a scanner for stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Printf(Red+"%s"+Reset+"\n", scanner.Text())
		}
	}()

	err = cmd.Wait()

	fmt.Print(cmdInfo.PostExecutionMessage + "\n\n")

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

// CopyDirectory copies a directory from src to dst
func CopyDirectory(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
				return err
			}
			if err := CopyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func MkdirOrCopyFile(
	filePath string,
	dirPath string,
	devgitaPathFile string,
	appName string,
) error {
	// Check if file already exists
	if !FileAlreadyExist(filePath) {
		// Create the directory if it doesn't exist
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
		return CopyFile(devgitaPathFile, filePath)
	}
	fmt.Printf("File for %s already exists.\n\n", appName)
	return nil
}

func FileAlreadyExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}

func IsCommandInstalled(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func DownloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
