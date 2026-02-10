package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/paths"
	"github.com/cjairm/devgita/pkg/utils"
)

var LookPathFn = exec.LookPath
var CommandFn = exec.Command

// BaseCommandExecutor defines the interface for executing commands and managing system state
// This interface allows for dependency injection and mocking in tests
type BaseCommandExecutor interface {
	// Command execution
	ExecCommand(cmd CommandParams) (string, string, error)

	// Shell configuration
	Setup(line string) error
	MaybeSetup(line, toSearch string) error

	// System checks
	IsDesktopAppPresent(dirPath, appName string) (bool, error)
	IsPackagePresent(cmd *exec.Cmd, packageName string) (bool, error)
	IsFontPresent(fontName string) (bool, error)

	// Installation helpers
	MaybeInstall(
		itemName string,
		alias []string,
		checkInstalled func(string) (bool, error),
		installFunc func(string) error,
		installURLFunc func(string) error,
		itemType string,
	) error
	InstallFontFromURL(url, fontFileName string, runCache bool) error
}

type BaseCommand struct {
	Platform CustomizablePlatform
}

type CommandParams struct {
	PreExecMsg  string
	PostExecMsg string
	IsSudo      bool
	Command     string
	Args        []string
}

func NewBaseCommand() *BaseCommand {
	return &BaseCommand{
		Platform: NewPlatform(),
	}
}

func NewBaseCommandCustom(p CustomizablePlatform) *BaseCommand {
	return &BaseCommand{
		Platform: p,
	}
}

func (b *BaseCommand) Setup(line string) error {
	return files.AddLineToFile(line, paths.Files.ShellConfig)
}

func (b *BaseCommand) MaybeSetup(line, toSearch string) error {
	isAlreadySetup, err := files.ContentExistsInFile(paths.Files.ShellConfig, toSearch)
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
	fontDirs := []string{paths.Paths.User.Fonts, paths.Paths.System.Fonts}
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

func (b *BaseCommand) ExecCommand(cmd CommandParams) (string, string, error) {
	if cmd.PreExecMsg != "" {
		utils.Print(cmd.PreExecMsg, "")
	}

	command := cmd.Command
	if cmd.IsSudo {
		command = "sudo " + command
	}

	logger.L().
		Debugw("Executing command", "command", strings.Join(append([]string{command}, cmd.Args...), " "))

	execCommand := exec.Command(command, cmd.Args...)
	execCommand.Stdin = os.Stdin

	stdoutPipe, err := execCommand.StdoutPipe()
	if err != nil {
		logger.L().Errorf("failed to get stdout pipe: %v", err)
		return "", "", err
	}

	stderrPipe, err := execCommand.StderrPipe()
	if err != nil {
		logger.L().Errorf("failed to get stderr pipe: %v", err)
		return "", "", err
	}

	var stdoutBuf, stderrBuf strings.Builder

	stdoutScanner := bufio.NewScanner(stdoutPipe)
	stderrScanner := bufio.NewScanner(stderrPipe)

	// Start command
	if err := execCommand.Start(); err != nil {
		logger.L().Errorf("failed to start command: %v", err)
		return "", "", err
	}

	// Read stdout
	go func() {
		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			stdoutBuf.WriteString(line + "\n")
			logger.L().Debugw("stdout", "line", line)
		}
	}()

	// Read stderr
	go func() {
		for stderrScanner.Scan() {
			line := stderrScanner.Text()
			stderrBuf.WriteString(line + "\n")
			logger.L().Debugw("stderr", "line", line)
		}
	}()

	// Wait for command to complete
	err = execCommand.Wait()
	if err != nil {
		logger.L().Errorw("command finished with error", "error", err, "stderr", stderrBuf.String())
	}

	if cmd.PostExecMsg != "" && err == nil {
		utils.Print(cmd.PostExecMsg, "")
	}

	return strings.TrimSpace(stdoutBuf.String()), strings.TrimSpace(stderrBuf.String()), err
}

func (b *BaseCommand) MaybeInstall(
	itemName string,
	alias []string,
	checkInstalled func(string) (bool, error),
	installFunc func(string) error,
	installURLFunc func(string) error, // Optional: Handle URL-based installation (e.g., for fonts)
	itemType string,
) error {
	var isInstalled bool
	var err error
	pkgToInstall := itemName
	if len(alias) > 0 {
		pkgToInstall = alias[0]
	}

	globalConfig := &config.GlobalConfig{}
	if err := globalConfig.Load(); err != nil {
		// HACK: If global config doesn't exist and we're trying to install git,
		// we can assume it's a fresh install and create the global config
		if pkgToInstall == constants.Git {
			globalConfig.Create()
		} else {
			logger.L().Errorw("Could not load global config", "error", err)
			return err
		}
	} else {
		if globalConfig.IsInstalledByDevgita(pkgToInstall, itemType) {
			logger.L().
				Debugw("Item already tracked as installed by devgita", "item", pkgToInstall, "type", itemType)
			return nil
		}
		if globalConfig.IsAlreadyInstalled(pkgToInstall, itemType) {
			logger.L().
				Debugw("Item is already installed, skipping", "item", pkgToInstall, "type", itemType)
			return nil
		}
	}

	isInstalled, err = checkInstalled(pkgToInstall)
	if err != nil {
		return err
	}

	if isInstalled {
		logger.L().
			Debugw("Item is already installed, marking as such in global config", "item", pkgToInstall, "type", itemType)
		globalConfig.AddToAlreadyInstalled(pkgToInstall, itemType)
		globalConfig.Save()
		return nil
	}

	var installErr error
	if installURLFunc != nil {
		installErr = installURLFunc(pkgToInstall)
	} else {
		installErr = installFunc(pkgToInstall)
	}

	if installErr == nil {
		globalConfig.AddToInstalled(pkgToInstall, itemType)
		if err := globalConfig.Save(); err != nil {
			logger.L().Errorw("Failed to update global config after installation", "error", err)
		}
	} else {
		logger.L().
			Warnw("Installation failed", "item", pkgToInstall, "type", itemType, "error", installErr)
	}

	return installErr
}

func (b *BaseCommand) InstallFontFromURL(url, fontFileName string, runCache bool) error {
	tmpPath := fmt.Sprintf("/tmp/%s.ttf", fontFileName)

	// 1. Download font
	if _, _, err := b.ExecCommand(CommandParams{
		PreExecMsg: fmt.Sprintf("Downloading %s...", fontFileName),
		Command:    "curl",
		Args:       []string{"-o", tmpPath, url},
	}); err != nil {
		return err
	}

	// 2. Move font
	if _, _, err := b.ExecCommand(CommandParams{
		PreExecMsg: "Installing font...",
		Command:    "mv",
		Args:       []string{tmpPath, filepath.Join(paths.Paths.User.Fonts, fontFileName+".ttf")},
	}); err != nil {
		return err
	}

	// 3. Update font cache if needed
	if runCache {
		if _, _, err := b.ExecCommand(CommandParams{
			PreExecMsg: "Refreshing font cache...",
			Command:    "fc-cache",
			Args:       []string{"-fv"},
			IsSudo:     true,
		}); err != nil {
			return fmt.Errorf("failed to refresh font cache: %w", err)
		}
	}
	return nil
}
