package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type GlobalConfig struct {
	LocalConfigPath      string            `json:"localConfigPath"`
	RemoteConfigPath     string            `json:"remoteConfigPath"`
	SelectedTheme        string            `json:"selectedTheme"`
	Font                 string            `json:"font"`
	InstalledPackages    []string          `json:"installedPackages"`
	InstalledDesktopApps []string          `json:"installedDesktopApps"`
	Shortcuts            map[string]string `json:"shortcuts"`
}

type BaseCommand struct{}

func NewBaseCommand() *BaseCommand {
	return &BaseCommand{}
}

func (b *BaseCommand) IsMac() bool {
	return runtime.GOOS == "darwin"
}

func (b *BaseCommand) IsLinux() bool {
	return runtime.GOOS == "linux"
}

func (b *BaseCommand) LoadGlobalConfig() (*GlobalConfig, error) {
	file, err := os.ReadFile(filepath.Join(paths.BashConfigAppDir, "global_config.json"))
	if err != nil {
		return nil, err
	}
	var config GlobalConfig
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (b *BaseCommand) SetGlobalConfig(config *GlobalConfig) error {
	filePath := filepath.Join(paths.BashConfigAppDir, "global_config.json")
	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, file, 0644)
}

func (b *BaseCommand) ResetGlobalConfig() error {
	filePath := filepath.Join(paths.BashConfigAppDir, "global_config.json")
	return os.WriteFile(filePath, []byte("{}"), 0644)
}

func (b *BaseCommand) Setup(line string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	zshConfigFile := filepath.Join(homeDir, ".zshrc")
	return files.AddLineToFile(line, zshConfigFile)
}

func (b *BaseCommand) MaybeSetup(line, toSearch string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	zshConfigFile := filepath.Join(homeDir, ".zshrc")
	isAlreadySetup, err := files.ContentExistsInFile(zshConfigFile, toSearch)
	if err != nil {
		return err
	}
	if isAlreadySetup == true {
		return nil
	}
	return b.Setup(line)
}

func (b *BaseCommand) FindPackageInCommandOutput(cmd *exec.Cmd, packageName string) (bool, error) {
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("Failed running brew command: %v", err)
	}
	for _, line := range bytes.Split(out.Bytes(), []byte{'\n'}) {
		if b.IsMac() {
			if string(line) == packageName {
				return true, nil
			}
		} else if b.IsLinux() {
			// The output of `dpkg -l` has a specific format, we need to check the package name in the right column
			if len(line) > 0 {
				// The package name is typically the second column in the output
				fields := bytes.Fields(line)
				if len(fields) > 1 && string(fields[1]) == packageName {
					return true, nil
				}
			}

		}
	}
	return false, nil
}

func (b *BaseCommand) CheckFileExistsInDirectory(dirPath, name string) (bool, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return false, fmt.Errorf("Failed to read directory: %v", err)
	}
	for _, file := range files {
		lowerCaseName := strings.ToLower(file.Name())
		if strings.Contains(lowerCaseName, name) {
			if b.IsLinux() && strings.HasSuffix(lowerCaseName, ".desktop") {
				return true, nil
			}
			return true, nil
		}
	}
	return false, nil
}

//Example of how to use the config package
// configFile := "./configs/bash/devgita_config.json"
//
// // Load the configuration
// c, err := config.LoadConfig(configFile)
// if err != nil {
// 	fmt.Println("Error loading config:", err)
// 	return
// }
//
// // Print the loaded configuration
// fmt.Printf("Loaded Config: %+v\n", c)
//
// // Modify the configuration
// c.SelectedTheme = "light"
// c.InstalledPackages = append(c.InstalledPackages, "new-package")
//
// // Save the updated configuration
// err = config.SaveConfig(configFile, c)
// if err != nil {
// 	fmt.Println("Error saving config:", err)
// 	return
// }
//
// fmt.Println("Configuration saved successfully.")
