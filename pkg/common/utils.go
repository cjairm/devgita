package common

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
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
	if cmdInfo.PreExecutionMessage != "" {
		fmt.Print(cmdInfo.PreExecutionMessage + "\n")
	}
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

	if cmdInfo.PostExecutionMessage != "" {
		fmt.Print(cmdInfo.PostExecutionMessage + "\n\n")
	}

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

func DirAlreadyExist(folderPath string) bool {
	info, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false // Directory does not exist
		}
		return false // An error occurred (e.g., permission denied)
	}
	isDir := info.IsDir()
	if isDir {
		return true
	}
	return false
}

// MoveContents moves files or folders from the source path to the target directory.
func MoveContents(srcPath, targetDir string) error {
	// Check if the target directory exists, if not, create it
	if !DirAlreadyExist(targetDir) {
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

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Function to remove an item from a slice
func removeItem(slice []string, item string) []string {
	for i, s := range slice {
		if s == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func MultiSelect(label string, options []string) ([]string, error) {
	var selectedLangs []string
	availableOptions := options

	for {
		// Remove "All" and "None" if all programming languages are selected
		if len(selectedLangs) == len(options)-3 { // Exclude "All", "None", "Done"
			availableOptions = removeItem(availableOptions, "All")
			availableOptions = removeItem(availableOptions, "None")
		}

		// Create a prompt for selecting languages
		prompt := promptui.Select{
			Label: label,
			Items: availableOptions,
		}

		// Show the prompt to the user
		_, result, err := prompt.Run()
		if err != nil {
			return nil, fmt.Errorf("prompt failed: %w", err)
		}

		// Handle special cases
		switch result {
		case "All":
			// Add all languages to selectedLangs
			selectedLangs = append(
				selectedLangs,
				availableOptions[3:]...)
			return selectedLangs, nil
		case "None":
			return []string{}, nil
		case "Done":
			return selectedLangs, nil
		default:
			// Add the selected language to the list of selected languages
			if !contains(selectedLangs, result) {
				selectedLangs = append(selectedLangs, result)
			}
			// Remove the selected item from the available options
			availableOptions = removeItem(availableOptions, result)
		}
	}
}

func UpdateFile(filePath, contentToReplace string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	// Convert the data to a string for manipulation
	fileContent := string(data)
	// Replace <HOME-PATH> with the actual home directory path
	homeDir, err := os.UserHomeDir()
	updatedContent := strings.ReplaceAll(fileContent, contentToReplace, homeDir)
	// Write the updated content back to the configuration file
	err = os.WriteFile(filePath, []byte(updatedContent), 0644)
	if err != nil {
		return err
	}
	return nil
}

func ContentExistInFile(filePath, substringToFind string) (error, bool) {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Error opening .zshrc file:", err), false
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), strings.ToLower(substringToFind)) {
			found = true
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading .zshrc file:", err), false
	}
	return nil, found
}

func Reboot() {
	fmt.Print("Ready to reboot for all settings to take effect? (yes/no): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)
	if input == "yes" {
		cmd := exec.Command("sudo", "reboot")
		if err := cmd.Run(); err != nil {
			fmt.Println("Error executing reboot:", err)
		}
	} else {
		fmt.Println("Reboot canceled.")
	}
}
