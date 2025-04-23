package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

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
	// TODO: Check if `.zshrc` file is present or any other type of file for
	// configuration
	return files.AddLineToFile(line, filepath.Join(paths.HomeDir, ".zshrc"))
}

func (b *BaseCommand) MaybeSetup(line, toSearch string) error {
	// TODO: Check if `.zshrc` file is present or any other type of file for
	// configuration
	isAlreadySetup, err := files.ContentExistsInFile(
		filepath.Join(paths.HomeDir, ".zshrc"),
		toSearch,
	)
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
