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

	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
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

// Returns XDG_CONFIG_HOME or fallback to ~/.config
func (b *BaseCommand) ConfigDir(subDirs ...string) (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(append([]string{base}, subDirs...)...), nil
}

// Returns XDG_DATA_HOME or fallback to ~/.local/share
func (b *BaseCommand) DataDir(subDirs ...string) (string, error) {
	base := os.Getenv("XDG_DATA_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(append([]string{base}, subDirs...)...), nil
}

func (b *BaseCommand) AppDir(subDirs ...string) (string, error) {
	appDir, err := b.DataDir(constants.AppName)
	if err != nil {
		return "", err
	}
	return filepath.Join(append([]string{appDir}, subDirs...)...), nil
}

// Returns XDG_CACHE_HOME or fallback to ~/.cache
func (b *BaseCommand) CacheDir(subDirs ...string) (string, error) {
	base := os.Getenv("XDG_CACHE_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".cache")
	}
	return filepath.Join(append([]string{base}, subDirs...)...), nil
}

func (b *BaseCommand) CopyAppConfigDirToLocalConfigDir(fromDevgita, toLocal []string) error {
	appConfigDir, err := b.AppDir(append([]string{"configs"}, fromDevgita...)...)
	if err != nil {
		return err
	}
	configDir, err := b.ConfigDir(toLocal...)
	if err != nil {
		return err
	}
	if err = files.CopyDir(appConfigDir, configDir); err != nil {
		return err
	}
	return nil
}

func (b *BaseCommand) CopyAppConfigFileToHomeDir(fromDevgita ...string) error {
	appConfigDir, err := b.AppDir(append([]string{"configs"}, fromDevgita...)...)
	if err != nil {
		return err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	if err = files.CopyFile(appConfigDir, homeDir); err != nil {
		return err
	}
	return nil
}

func (b *BaseCommand) LoadGlobalConfig() (*GlobalConfig, error) {
	filename, err := b.AppDir("configs", "bash", "global_config.json")
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(filename)
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
	filename, err := b.AppDir("configs", "bash", "global_config.json")
	if err != nil {
		return err
	}
	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, file, 0644)
}

func (b *BaseCommand) ResetGlobalConfig() error {
	filename, err := b.AppDir("configs", "bash", "global_config.json")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, []byte("{}"), 0644)
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
