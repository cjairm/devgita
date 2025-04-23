package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

var LookPathFn = exec.LookPath
var CommandFn = exec.Command

type BaseCommand struct {
	Platform CustomizablePlatform
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
