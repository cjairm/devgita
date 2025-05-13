package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
	"github.com/cjairm/devgita/pkg/utils"
)

var LookPathFn = exec.LookPath
var CommandFn = exec.Command

type BaseCommand struct {
	Platform CustomizablePlatform
}

type CommandParams struct {
	PreExecMsg  string
	PostExecMsg string
	IsSudo      bool
	Verbose     bool
	Command     string
	Args        []string
}

func NewBaseCommand() *BaseCommand {
	return &BaseCommand{
		Platform: *NewPlatform(),
	}
}

func NewBaseCommandCustom(p CustomizablePlatform) *BaseCommand {
	return &BaseCommand{
		Platform: p,
	}
}

func (b *BaseCommand) Setup(line string) error {
	return files.AddLineToFile(line, paths.ShellConfigFile)
}

func (b *BaseCommand) MaybeSetup(line, toSearch string) error {
	isAlreadySetup, err := files.ContentExistsInFile(paths.ShellConfigFile, toSearch)
	if err != nil {
		return err
	}
	if isAlreadySetup == true {
		return nil
	}
	return b.Setup(line)
}

func (b *BaseCommand) IsDesktopAppPresent(dirPath, appName string) (bool, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return false, fmt.Errorf("failed to read directory: %v", err)
	}
	for _, file := range files {
		filename := strings.ToLower(file.Name())
		if strings.Contains(filename, strings.ToLower(appName)) {
			if b.Platform.IsLinux() && strings.HasSuffix(filename, ".desktop") {
				return true, nil
			}
			if b.Platform.IsMac() {
				return true, nil
			}
		}
	}
	return false, nil
}

func (b *BaseCommand) IsPackagePresent(cmd *exec.Cmd, packageName string) (bool, error) {
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("failed running command: %w", err)
	}
	lines := bytes.Split(out.Bytes(), []byte{'\n'})
	if b.Platform.IsMac() {
		return findPackageInBrewOutput(lines, packageName), nil
	} else if b.Platform.IsLinux() {
		return findPackageInDpkgOutput(lines, packageName), nil
	}
	return false, nil
}

func findPackageInBrewOutput(lines [][]byte, packageName string) bool {
	for _, line := range lines {
		if string(line) == packageName {
			return true
		}
	}
	return false
}

func findPackageInDpkgOutput(lines [][]byte, packageName string) bool {
	for _, line := range lines {
		// The package name is typically the second column in the output
		fields := bytes.Fields(line)
		if len(fields) > 1 && string(fields[1]) == packageName {
			return true
		}
	}
	return false
}

func (b *BaseCommand) IsFontPresent(fontName string) (bool, error) {
	if _, err := LookPathFn("fc-list"); err == nil {
		cmd := CommandFn("fc-list", ":", "family")
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err == nil {
			lines := bytes.Split(out.Bytes(), []byte{'\n'})
			fontNameLower := strings.ToLower(fontName)
			for _, line := range lines {
				if strings.Contains(strings.ToLower(string(line)), fontNameLower) {
					return true, nil
				}
			}
			return false, nil
		}
	}
	// Fallback: scan known font directories
	fontDirs := []string{paths.UserFontsDir, paths.SystemFontsDir}
	for _, dir := range fontDirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue // ignore unreadable dirs
		}
		for _, file := range files {
			if strings.Contains(strings.ToLower(file.Name()), strings.ToLower(fontName)) &&
				hasFontExtension(strings.ToLower(file.Name())) {
				return true, nil
			}
		}
	}
	return false, nil
}

func hasFontExtension(filename string) bool {
	return strings.HasSuffix(filename, ".ttf") ||
		strings.HasSuffix(filename, ".otf") ||
		strings.HasSuffix(filename, ".woff") ||
		strings.HasSuffix(filename, ".woff2")
}

func (b *BaseCommand) ExecCommand(cmd CommandParams) error {
	if cmd.PreExecMsg != "" {
		utils.Print(cmd.PreExecMsg, "")
	}
	command := cmd.Command
	if cmd.IsSudo {
		command = "sudo " + command
	}
	execCommand := exec.Command(command, cmd.Args...)

	// Create pipes to capture standard output
	stdout, err := execCommand.StdoutPipe()
	if err != nil {
		return err
	}

	// Create pipes to capture standard error
	stderr, err := execCommand.StderrPipe()
	if err != nil {
		return err
	}

	// To allow user to interact with the command (password)
	execCommand.Stdin = os.Stdin

	// Start the command
	if err := execCommand.Start(); err != nil {
		return err
	}

	if cmd.Verbose == true {
		// Create a scanner for stdout
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				utils.PrintSecondary(scanner.Text())
			}
		}()

		// Create a scanner for stderr
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				utils.PrintError(scanner.Text())
			}
		}()
	}

	err = execCommand.Wait()
	if cmd.PostExecMsg != "" && err == nil {
		utils.Print(cmd.PostExecMsg, "")
	}

	return err
}

func (b *BaseCommand) MaybeInstall(
	itemName string,
	alias []string,
	checkInstalled func(string) (bool, error),
	installFunc func(string) error,
	installURLFunc func(string) error, // Optional: Handle URL-based installation (e.g., for fonts)
) error {
	var isInstalled bool
	var err error
	pkgToInstall := itemName
	if len(alias) > 0 {
		pkgToInstall = alias[0]
	}

	// Check if the item is installed
	isInstalled, err = checkInstalled(pkgToInstall)
	if err != nil {
		return err
	}

	// If installed, do nothing
	if isInstalled {
		return nil
	}

	if installURLFunc != nil {
		return installURLFunc(pkgToInstall)
	}
	return installFunc(pkgToInstall)
}

func (b *BaseCommand) InstallFontFromURL(url, fontFileName string, runCache bool) error {
	tmpPath := fmt.Sprintf("/tmp/%s.ttf", fontFileName)

	// 1. Download font
	if err := b.ExecCommand(CommandParams{
		PreExecMsg: fmt.Sprintf("Downloading %s...", fontFileName),
		Command:    "curl",
		Args:       []string{"-o", tmpPath, url},
	}); err != nil {
		return err
	}

	// 2. Move font
	if err := b.ExecCommand(CommandParams{
		PreExecMsg: "Installing font...",
		Command:    "mv",
		Args:       []string{tmpPath, filepath.Join(paths.UserFontsDir, fontFileName+".ttf")},
	}); err != nil {
		return err
	}

	// 3. Update font cache if needed
	if runCache {
		return b.ExecCommand(CommandParams{
			PreExecMsg: "Refreshing font cache...",
			Command:    "fc-cache",
			Args:       []string{"-fv"},
			IsSudo:     true,
		})
	}
	return nil
}
