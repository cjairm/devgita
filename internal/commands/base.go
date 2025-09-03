package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/logger"
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

func (b *BaseCommand) ExecCommand(cmd CommandParams) (string, error) {
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
		return "", err
	}

	stderrPipe, err := execCommand.StderrPipe()
	if err != nil {
		logger.L().Errorf("failed to get stderr pipe: %v", err)
		return "", err
	}

	var stdoutBuf, stderrBuf strings.Builder

	stdoutScanner := bufio.NewScanner(stdoutPipe)
	stderrScanner := bufio.NewScanner(stderrPipe)

	// Start command
	if err := execCommand.Start(); err != nil {
		logger.L().Errorf("failed to start command: %v", err)
		return "", err
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

	return strings.TrimSpace(stdoutBuf.String()), err
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
		logger.L().Errorw("Could not load global config", "error", err)
		return err
	} else {
		if b.isTrackedInGlobalConfig(globalConfig, pkgToInstall, itemType) {
			logger.L().Debugw("Item already tracked as installed by devgita", "item", pkgToInstall, "type", itemType)
			return nil
		}
		if b.isTrackedInAlreadyInstalledConfig(globalConfig, pkgToInstall, itemType) {
			logger.L().Debugw("Item already tracked as pre-existing", "item", pkgToInstall, "type", itemType)
			return nil
		}
		if b.isIgnoredInGlobalConfig(globalConfig, pkgToInstall, itemType) {
			logger.L().Debugw("Item is ignored, skipping", "item", pkgToInstall, "type", itemType)
			return nil
		}
	}

	isInstalled, err = checkInstalled(pkgToInstall)
	if err != nil {
		return err
	}

	if isInstalled {
		if globalConfig != nil {
			b.addToAlreadyInstalledConfig(globalConfig, pkgToInstall, itemType)
			globalConfig.Save()
		}
		return nil
	}

	var installErr error
	if installURLFunc != nil {
		installErr = installURLFunc(pkgToInstall)
	} else {
		installErr = installFunc(pkgToInstall)
	}

	if installErr == nil && globalConfig != nil {
		b.addToGlobalConfig(globalConfig, pkgToInstall, itemType)
		if err := globalConfig.Save(); err != nil {
			logger.L().Errorw("Failed to update global config after installation", "error", err)
		}
	}

	return installErr
}

func (b *BaseCommand) InstallFontFromURL(url, fontFileName string, runCache bool) error {
	tmpPath := fmt.Sprintf("/tmp/%s.ttf", fontFileName)

	// 1. Download font
	if _, err := b.ExecCommand(CommandParams{
		PreExecMsg: fmt.Sprintf("Downloading %s...", fontFileName),
		Command:    "curl",
		Args:       []string{"-o", tmpPath, url},
	}); err != nil {
		return err
	}

	// 2. Move font
	if _, err := b.ExecCommand(CommandParams{
		PreExecMsg: "Installing font...",
		Command:    "mv",
		Args:       []string{tmpPath, filepath.Join(paths.UserFontsDir, fontFileName+".ttf")},
	}); err != nil {
		return err
	}

	// 3. Update font cache if needed
	if runCache {
		if _, err := b.ExecCommand(CommandParams{
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
func (b *BaseCommand) isTrackedInGlobalConfig(
	config *config.GlobalConfig,
	itemName, itemType string,
) bool {
	switch itemType {
	case "font":
		return slices.Contains(config.Installed.Fonts, itemName)
	case "package":
		return slices.Contains(config.Installed.Packages, itemName)
	case "desktop_app":
		return slices.Contains(config.Installed.DesktopApps, itemName)
	case "terminal_tool":
		return slices.Contains(config.Installed.TerminalTools, itemName)
	case "theme":
		return slices.Contains(config.Installed.Themes, itemName)
	case "dev_language":
		return slices.Contains(config.Installed.DevLanguages, itemName)
	case "database":
		return slices.Contains(config.Installed.Databases, itemName)
	default:
		return false
	}
}

func (b *BaseCommand) isIgnoredInGlobalConfig(
	config *config.GlobalConfig,
	itemName, itemType string,
) bool {
	switch itemType {
	case "font":
		return slices.Contains(config.Ignored.Fonts, itemName)
	case "package":
		return slices.Contains(config.Ignored.Packages, itemName)
	case "desktop_app":
		return slices.Contains(config.Ignored.DesktopApps, itemName)
	case "terminal_tool":
		return slices.Contains(config.Ignored.TerminalTools, itemName)
	case "theme":
		return slices.Contains(config.Ignored.Themes, itemName)
	case "dev_language":
		return slices.Contains(config.Ignored.DevLanguages, itemName)
	case "database":
		return slices.Contains(config.Ignored.Databases, itemName)
	default:
		return false
	}
}

func (b *BaseCommand) addToGlobalConfig(config *config.GlobalConfig, itemName, itemType string) {
	switch itemType {
	case "font":
		if !slices.Contains(config.Installed.Fonts, itemName) {
			config.Installed.Fonts = append(config.Installed.Fonts, itemName)
		}
	case "package":
		if !slices.Contains(config.Installed.Packages, itemName) {
			config.Installed.Packages = append(config.Installed.Packages, itemName)
		}
	case "desktop_app":
		if !slices.Contains(config.Installed.DesktopApps, itemName) {
			config.Installed.DesktopApps = append(config.Installed.DesktopApps, itemName)
		}
	case "terminal_tool":
		if !slices.Contains(config.Installed.TerminalTools, itemName) {
			config.Installed.TerminalTools = append(config.Installed.TerminalTools, itemName)
		}
	case "theme":
		if !slices.Contains(config.Installed.Themes, itemName) {
			config.Installed.Themes = append(config.Installed.Themes, itemName)
		}
	case "dev_language":
		if !slices.Contains(config.Installed.DevLanguages, itemName) {
			config.Installed.DevLanguages = append(config.Installed.DevLanguages, itemName)
		}
	case "database":
		if !slices.Contains(config.Installed.Databases, itemName) {
			config.Installed.Databases = append(config.Installed.Databases, itemName)
		}
	}
}

// Helper function to check if item is tracked in already installed config
func (b *BaseCommand) isTrackedInAlreadyInstalledConfig(
	config *config.GlobalConfig,
	itemName, itemType string,
) bool {
	switch itemType {
	case "font":
		return slices.Contains(config.AlreadyInstalledConfig.Fonts, itemName)
	case "package":
		return slices.Contains(config.AlreadyInstalledConfig.Packages, itemName)
	case "desktop_app":
		return slices.Contains(config.AlreadyInstalledConfig.DesktopApps, itemName)
	case "terminal_tool":
		return slices.Contains(config.AlreadyInstalledConfig.TerminalTools, itemName)
	case "theme":
		return slices.Contains(config.AlreadyInstalledConfig.Themes, itemName)
	case "dev_language":
		return slices.Contains(config.AlreadyInstalledConfig.DevLanguages, itemName)
	case "database":
		return slices.Contains(config.AlreadyInstalledConfig.Databases, itemName)
	default:
		return false
	}
}

// Helper function to add item to already installed config
func (b *BaseCommand) addToAlreadyInstalledConfig(
	config *config.GlobalConfig,
	itemName, itemType string,
) {
	switch itemType {
	case "font":
		if !slices.Contains(config.AlreadyInstalledConfig.Fonts, itemName) {
			config.AlreadyInstalledConfig.Fonts = append(
				config.AlreadyInstalledConfig.Fonts,
				itemName,
			)
		}
	case "package":
		if !slices.Contains(config.AlreadyInstalledConfig.Packages, itemName) {
			config.AlreadyInstalledConfig.Packages = append(
				config.AlreadyInstalledConfig.Packages,
				itemName,
			)
		}
	case "desktop_app":
		if !slices.Contains(config.AlreadyInstalledConfig.DesktopApps, itemName) {
			config.AlreadyInstalledConfig.DesktopApps = append(
				config.AlreadyInstalledConfig.DesktopApps,
				itemName,
			)
		}
	case "terminal_tool":
		if !slices.Contains(config.AlreadyInstalledConfig.TerminalTools, itemName) {
			config.AlreadyInstalledConfig.TerminalTools = append(
				config.AlreadyInstalledConfig.TerminalTools,
				itemName,
			)
		}
	case "theme":
		if !slices.Contains(config.AlreadyInstalledConfig.Themes, itemName) {
			config.AlreadyInstalledConfig.Themes = append(
				config.AlreadyInstalledConfig.Themes,
				itemName,
			)
		}
	case "dev_language":
		if !slices.Contains(config.AlreadyInstalledConfig.DevLanguages, itemName) {
			config.AlreadyInstalledConfig.DevLanguages = append(
				config.AlreadyInstalledConfig.DevLanguages,
				itemName,
			)
		}
	case "database":
		if !slices.Contains(config.AlreadyInstalledConfig.Databases, itemName) {
			config.AlreadyInstalledConfig.Databases = append(
				config.AlreadyInstalledConfig.Databases,
				itemName,
			)
		}
	}
}
