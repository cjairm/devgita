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
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

// MoveContents moves files or folders from the source path to the target directory.
func MoveContents(srcPath, targetDir string) error {
	// Check if the target directory exists, if not, create it
	if !FileAlreadyExist(targetDir) {
		if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
			return fmt.Errorf("error creating target directory: %w", err)
		}
	}
	// Check if the source path is a directory
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("error stating source path: %w", err)
	}
	if srcInfo.IsDir() {
		// If it's a directory, copy its contents
		if err := CopyDirectory(srcPath, targetDir); err != nil {
			return fmt.Errorf("error copying directory: %w", err)
		}
	} else {
		// If it's a file, copy it directly to the target directory
		targetFilePath := filepath.Join(targetDir, srcInfo.Name())
		if err := CopyFile(srcPath, targetFilePath); err != nil {
			return fmt.Errorf("error copying file: %w", err)
		}
	}
	return nil
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
