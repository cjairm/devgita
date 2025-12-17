package devgita

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/paths"
)

func init() {
	logger.Init(false)
}

func TestNew(t *testing.T) {
	dg := New()
	if dg == nil {
		t.Fatal("New() returned nil")
	}
}

func TestSoftInstall_DirectoryDoesNotExist(t *testing.T) {
	tempDir := t.TempDir()
	appDir := filepath.Join(tempDir, "devgita")

	oldAppDir := paths.AppDir
	paths.AppDir = appDir
	t.Cleanup(func() {
		paths.AppDir = oldAppDir
	})

	dg := New()
	err := dg.SoftInstall()
	if err == nil {
		// If somehow it succeeded, verify the directory was created
		if _, statErr := os.Stat(appDir); os.IsNotExist(statErr) {
			t.Fatal("Expected Install() to be called when directory doesn't exist")
		}
	}
}

func TestSoftInstall_DirectoryExistsWithFiles(t *testing.T) {
	tempDir := t.TempDir()

	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	oldAppDir := paths.AppDir
	paths.AppDir = tempDir
	t.Cleanup(func() {
		paths.AppDir = oldAppDir
	})

	dg := New()
	err := dg.SoftInstall()
	if err != nil {
		t.Fatalf("SoftInstall() failed: %v", err)
	}

	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("Expected existing file to be preserved")
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}
	if string(content) != "content" {
		t.Fatal("Expected file content to be unchanged")
	}
}

func TestUninstall_DirectoryEmpty(t *testing.T) {
	tempDir := t.TempDir()
	emptyDir := filepath.Join(tempDir, "empty")
	if err := os.Mkdir(emptyDir, 0o755); err != nil {
		t.Fatalf("Failed to create empty directory: %v", err)
	}

	oldAppDir := paths.AppDir
	paths.AppDir = emptyDir
	t.Cleanup(func() {
		paths.AppDir = oldAppDir
	})

	dg := New()

	err := dg.Uninstall()
	if err != nil {
		t.Fatalf("Uninstall() failed: %v", err)
	}

	if _, err := os.Stat(emptyDir); !os.IsNotExist(err) {
		t.Fatal("Expected empty directory to be removed")
	}
}

func TestUninstall_RemovesRepositoryAndConfig(t *testing.T) {
	tempDir := t.TempDir()
	appDir := filepath.Join(tempDir, "app")
	configDir := filepath.Join(tempDir, "config")
	devgitaConfigDir := filepath.Join(configDir, "devgita")

	if err := os.MkdirAll(appDir, 0o755); err != nil {
		t.Fatalf("Failed to create app directory: %v", err)
	}
	testFile := filepath.Join(appDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if err := os.MkdirAll(devgitaConfigDir, 0o755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}
	testGlobalConfigPath := filepath.Join(devgitaConfigDir, "global_config.yaml")
	if err := os.WriteFile(testGlobalConfigPath, []byte("app_path: /test"), 0o644); err != nil {
		t.Fatalf("Failed to create global config: %v", err)
	}

	// Override global paths and package-level variables
	oldAppDir := paths.AppDir
	oldConfigDir := paths.ConfigDir
	oldConfigDirPath := configDirPath
	oldGlobalConfigPath := globalConfigPath

	paths.AppDir = appDir
	paths.ConfigDir = configDir
	configDirPath = devgitaConfigDir
	globalConfigPath = testGlobalConfigPath

	t.Cleanup(func() {
		paths.AppDir = oldAppDir
		paths.ConfigDir = oldConfigDir
		configDirPath = oldConfigDirPath
		globalConfigPath = oldGlobalConfigPath
	})

	dg := New()

	err := dg.Uninstall()
	if err != nil {
		t.Fatalf("Uninstall() failed: %v", err)
	}

	if _, err := os.Stat(appDir); !os.IsNotExist(err) {
		t.Fatal("Expected app directory to be removed")
	}

	if _, err := os.Stat(testGlobalConfigPath); !os.IsNotExist(err) {
		t.Fatal("Expected global config file to be removed")
	}

	if _, err := os.Stat(devgitaConfigDir); !os.IsNotExist(err) {
		t.Fatal("Expected devgita config directory to be removed")
	}
}

func TestUninstall_NoRepositoryNoConfig(t *testing.T) {
	tempDir := t.TempDir()
	appDir := filepath.Join(tempDir, "app")
	configDir := filepath.Join(tempDir, "config")
	devgitaConfigDir := filepath.Join(configDir, "devgita")
	testGlobalConfigPath := filepath.Join(devgitaConfigDir, "global_config.yaml")

	// Override global paths and package-level variables
	oldAppDir := paths.AppDir
	oldConfigDir := paths.ConfigDir
	oldConfigDirPath := configDirPath
	oldGlobalConfigPath := globalConfigPath

	paths.AppDir = appDir
	paths.ConfigDir = configDir
	configDirPath = devgitaConfigDir
	globalConfigPath = testGlobalConfigPath

	t.Cleanup(func() {
		paths.AppDir = oldAppDir
		paths.ConfigDir = oldConfigDir
		configDirPath = oldConfigDirPath
		globalConfigPath = oldGlobalConfigPath
	})

	dg := New()

	err := dg.Uninstall()
	if err != nil {
		t.Fatalf("Uninstall() failed: %v", err)
	}
}

func TestForceInstall(t *testing.T) {
	tempDir := t.TempDir()

	existingFile := filepath.Join(tempDir, "old.txt")
	if err := os.WriteFile(existingFile, []byte("old"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	oldAppDir := paths.AppDir
	paths.AppDir = tempDir
	t.Cleanup(func() {
		paths.AppDir = oldAppDir
	})

	dg := New()
	_ = dg.ForceInstall()

	if _, err := os.Stat(existingFile); !os.IsNotExist(err) {
		t.Fatal("Expected existing file to be removed by Uninstall")
	}
}

func TestForceConfigure_CreatesConfig(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	appDir := filepath.Join(tempDir, "app")
	bashConfigDir := filepath.Join(tempDir, "bash")

	// Create all necessary directories
	devgitaConfigDir := filepath.Join(configDir, "devgita")
	if err := os.MkdirAll(devgitaConfigDir, 0o755); err != nil {
		t.Fatalf("Failed to create devgita config dir: %v", err)
	}
	if err := os.MkdirAll(bashConfigDir, 0o755); err != nil {
		t.Fatalf("Failed to create bash config dir: %v", err)
	}

	// Create source bash config file
	sourceConfigPath := filepath.Join(bashConfigDir, "global_config.yaml")
	sourceContent := "# Global config template\n"
	if err := os.WriteFile(sourceConfigPath, []byte(sourceContent), 0o644); err != nil {
		t.Fatalf("Failed to create source config: %v", err)
	}

	// Override global paths and package-level variables
	oldConfigDir := paths.ConfigDir
	oldAppDir := paths.AppDir
	oldBashConfigAppDir := paths.BashConfigAppDir
	oldConfigDirPath := configDirPath
	oldGlobalConfigPath := globalConfigPath

	paths.ConfigDir = configDir
	paths.AppDir = appDir
	paths.BashConfigAppDir = bashConfigDir
	configDirPath = devgitaConfigDir
	globalConfigPath = filepath.Join(devgitaConfigDir, "global_config.yaml")

	t.Cleanup(func() {
		paths.ConfigDir = oldConfigDir
		paths.AppDir = oldAppDir
		paths.BashConfigAppDir = oldBashConfigAppDir
		configDirPath = oldConfigDirPath
		globalConfigPath = oldGlobalConfigPath
	})

	dg := New()

	err := dg.ForceConfigure()
	if err != nil {
		t.Fatalf("ForceConfigure() failed: %v", err)
	}

	// Verify config file was created
	configPath := filepath.Join(configDir, "devgita", "global_config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Expected config file to be created")
	}

	// Verify config content includes AppPath and ConfigPath
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "app_path:") {
		t.Fatal("Expected config to contain app_path field")
	}
	if !strings.Contains(contentStr, "config_path:") {
		t.Fatal("Expected config to contain config_path field")
	}
	if !strings.Contains(contentStr, appDir) {
		t.Fatalf("Expected config to contain AppDir path %q", appDir)
	}
}

func TestForceConfigure_OverwritesExisting(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	appDir := filepath.Join(tempDir, "app")
	bashConfigDir := filepath.Join(tempDir, "bash")

	// Create directories
	devgitaConfigDir := filepath.Join(configDir, "devgita")
	if err := os.MkdirAll(devgitaConfigDir, 0o755); err != nil {
		t.Fatalf("Failed to create devgita config dir: %v", err)
	}
	if err := os.MkdirAll(bashConfigDir, 0o755); err != nil {
		t.Fatalf("Failed to create bash config dir: %v", err)
	}

	// Create source config
	sourceConfigPath := filepath.Join(bashConfigDir, "global_config.yaml")
	sourceContent := "# Global config template\n"
	if err := os.WriteFile(sourceConfigPath, []byte(sourceContent), 0o644); err != nil {
		t.Fatalf("Failed to create source config: %v", err)
	}

	// Create existing config with old values
	configPath := filepath.Join(devgitaConfigDir, "global_config.yaml")
	oldContent := "app_path: /old/path\nconfig_path: /old/config\n"
	if err := os.WriteFile(configPath, []byte(oldContent), 0o644); err != nil {
		t.Fatalf("Failed to create existing config: %v", err)
	}

	// Override global paths and package-level variables
	oldConfigDir := paths.ConfigDir
	oldAppDir := paths.AppDir
	oldBashConfigAppDir := paths.BashConfigAppDir
	oldConfigDirPath := configDirPath
	oldGlobalConfigPath := globalConfigPath

	paths.ConfigDir = configDir
	paths.AppDir = appDir
	paths.BashConfigAppDir = bashConfigDir
	configDirPath = devgitaConfigDir
	globalConfigPath = configPath

	t.Cleanup(func() {
		paths.ConfigDir = oldConfigDir
		paths.AppDir = oldAppDir
		paths.BashConfigAppDir = oldBashConfigAppDir
		configDirPath = oldConfigDirPath
		globalConfigPath = oldGlobalConfigPath
	})

	dg := New()

	err := dg.ForceConfigure()
	if err != nil {
		t.Fatalf("ForceConfigure() failed: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	contentStr := string(content)
	if strings.Contains(contentStr, "/old/path") {
		t.Fatal("Expected old app_path to be replaced")
	}
	if strings.Contains(contentStr, "/old/config") {
		t.Fatal("Expected old config_path to be replaced")
	}
	if !strings.Contains(contentStr, appDir) {
		t.Fatalf("Expected config to contain new AppDir path %q", appDir)
	}
}

func TestSoftConfigure_PreservesExistingConfig(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	appDir := filepath.Join(tempDir, "app")
	bashConfigDir := filepath.Join(tempDir, "bash")

	// Create directories
	devgitaConfigDir := filepath.Join(configDir, "devgita")
	if err := os.MkdirAll(devgitaConfigDir, 0o755); err != nil {
		t.Fatalf("Failed to create devgita config dir: %v", err)
	}
	if err := os.MkdirAll(bashConfigDir, 0o755); err != nil {
		t.Fatalf("Failed to create bash config dir: %v", err)
	}

	// Create source config
	sourceConfigPath := filepath.Join(bashConfigDir, "global_config.yaml")
	sourceContent := "# Global config template\n"
	if err := os.WriteFile(sourceConfigPath, []byte(sourceContent), 0o644); err != nil {
		t.Fatalf("Failed to create source config: %v", err)
	}

	// Create existing config with custom values
	configPath := filepath.Join(devgitaConfigDir, "global_config.yaml")
	customContent := "app_path: /custom/path\nconfig_path: /custom/config\ncustom_field: custom_value\n"
	if err := os.WriteFile(configPath, []byte(customContent), 0o644); err != nil {
		t.Fatalf("Failed to create existing config: %v", err)
	}

	// Override global paths and package-level variables
	oldConfigDir := paths.ConfigDir
	oldAppDir := paths.AppDir
	oldBashConfigAppDir := paths.BashConfigAppDir
	oldConfigDirPath := configDirPath
	oldGlobalConfigPath := globalConfigPath

	paths.ConfigDir = configDir
	paths.AppDir = appDir
	paths.BashConfigAppDir = bashConfigDir
	configDirPath = devgitaConfigDir
	globalConfigPath = configPath

	t.Cleanup(func() {
		paths.ConfigDir = oldConfigDir
		paths.AppDir = oldAppDir
		paths.BashConfigAppDir = oldBashConfigAppDir
		configDirPath = oldConfigDirPath
		globalConfigPath = oldGlobalConfigPath
	})

	dg := New()

	err := dg.SoftConfigure()
	if err != nil {
		t.Fatalf("SoftConfigure() failed: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	contentStr := string(content)
	if contentStr != customContent {
		t.Fatalf(
			"Expected config to be unchanged, got:\n%s\nExpected:\n%s",
			contentStr,
			customContent,
		)
	}
	if !strings.Contains(contentStr, "custom_field: custom_value") {
		t.Fatal("Expected custom field to be preserved")
	}
}
