package devgita

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/apps/git"
	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/paths"
)

func init() {
	logger.Init(false)
}

// testPaths holds the original path values for restoration after tests
type testPaths struct {
	appDir             string
	configDir          string
	templatesConfigDir string
	configDirPath      string
	globalConfigPath   string
	zshConfigPath      string
}

// setupTestPaths overrides package-level path variables with test paths
// Returns a cleanup function to restore original paths
func setupTestPaths(t *testing.T, tempDir string) func() {
	t.Helper()

	configDir := filepath.Join(tempDir, "config")
	devgitaConfigDir := filepath.Join(configDir, "devgita")

	// Save original values
	original := testPaths{
		appDir:             paths.Paths.App.Root,
		configDir:          paths.Paths.Config.Root,
		templatesConfigDir: paths.Paths.App.Configs.Templates,
		configDirPath:      configDirPath,
		globalConfigPath:   globalConfigPath,
		zshConfigPath:      zshConfigPath,
	}

	// Override with test values
	paths.Paths.App.Root = tempDir
	paths.Paths.Config.Root = configDir
	configDirPath = devgitaConfigDir
	globalConfigPath = filepath.Join(devgitaConfigDir, "global_config.yaml")
	zshConfigPath = filepath.Join(devgitaConfigDir, "devgita.zsh")

	// Return cleanup function
	return func() {
		paths.Paths.App.Root = original.appDir
		paths.Paths.Config.Root = original.configDir
		paths.Paths.App.Configs.Templates = original.templatesConfigDir
		configDirPath = original.configDirPath
		globalConfigPath = original.globalConfigPath
		zshConfigPath = original.zshConfigPath
	}
}

// setupConfigTestPaths is similar to setupTestPaths but includes templates directory
func setupConfigTestPaths(t *testing.T, tempDir, appDir, templatesDir string) func() {
	t.Helper()

	configDir := filepath.Join(tempDir, "config")
	devgitaConfigDir := filepath.Join(configDir, "devgita")

	// Save original values
	original := testPaths{
		appDir:             paths.Paths.App.Root,
		configDir:          paths.Paths.Config.Root,
		templatesConfigDir: paths.Paths.App.Configs.Templates,
		configDirPath:      configDirPath,
		globalConfigPath:   globalConfigPath,
		zshConfigPath:      zshConfigPath,
	}

	// Override with test values
	paths.Paths.App.Root = appDir
	paths.Paths.Config.Root = configDir
	paths.Paths.App.Configs.Templates = templatesDir
	configDirPath = devgitaConfigDir
	globalConfigPath = filepath.Join(devgitaConfigDir, "global_config.yaml")
	zshConfigPath = filepath.Join(devgitaConfigDir, "devgita.zsh")

	// Return cleanup function
	return func() {
		paths.Paths.App.Root = original.appDir
		paths.Paths.Config.Root = original.configDir
		paths.Paths.App.Configs.Templates = original.templatesConfigDir
		configDirPath = original.configDirPath
		globalConfigPath = original.globalConfigPath
		zshConfigPath = original.zshConfigPath
	}
}

// createTestDirs creates common test directory structure
func createTestDirs(
	t *testing.T,
	tempDir string,
) (configDir, appDir, templatesDir, devgitaConfigDir string) {
	t.Helper()

	configDir = filepath.Join(tempDir, "config")
	appDir = filepath.Join(tempDir, "app")
	templatesDir = filepath.Join(tempDir, "templates")
	devgitaConfigDir = filepath.Join(configDir, "devgita")

	if err := os.MkdirAll(devgitaConfigDir, 0o755); err != nil {
		t.Fatalf("Failed to create devgita config dir: %v", err)
	}
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		t.Fatalf("Failed to create app dir: %v", err)
	}
	if err := os.MkdirAll(templatesDir, 0o755); err != nil {
		t.Fatalf("Failed to create templates dir: %v", err)
	}

	return
}

// createTemplateFiles creates required template files for configuration tests
func createTemplateFiles(t *testing.T, templatesDir string) {
	t.Helper()

	// Create global_config.yaml template
	sourceConfigPath := filepath.Join(templatesDir, "global_config.yaml")
	sourceContent := "# Global config template\n"
	if err := os.WriteFile(sourceConfigPath, []byte(sourceContent), 0o644); err != nil {
		t.Fatalf("Failed to create source config: %v", err)
	}

	// Create devgita.zsh.tmpl template
	zshTemplatePath := filepath.Join(templatesDir, "devgita.zsh.tmpl")
	zshTemplateContent := "# Devgita Shell Configuration\n"
	if err := os.WriteFile(zshTemplatePath, []byte(zshTemplateContent), 0o644); err != nil {
		t.Fatalf("Failed to create zsh template: %v", err)
	}
}

// createMockDevgita creates a Devgita instance with mock Git and Base for testing
func createMockDevgita() (*Devgita, *commands.MockBaseCommand) {
	mockBase := commands.NewMockBaseCommand()
	mockCmd := commands.NewMockCommand()

	// Create Git with mock base command to prevent real command execution
	gitInstance := &git.Git{
		Cmd:  mockCmd,
		Base: mockBase,
	}

	// Create a real BaseCommand for the Devgita instance
	// The actual command execution will be mocked at a different layer
	baseCmd := commands.NewBaseCommand()

	return &Devgita{
		Git:  *gitInstance,
		Base: *baseCmd,
	}, mockBase
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

	oldAppDir := paths.Paths.App.Root
	paths.Paths.App.Root = appDir
	t.Cleanup(func() {
		paths.Paths.App.Root = oldAppDir
	})

	dg, mockBase := createMockDevgita()

	// Mock successful git clone
	mockBase.SetExecCommandResult("Cloning into...", "", nil)

	err := dg.SoftInstall()
	if err != nil {
		t.Fatalf("SoftInstall() failed: %v", err)
	}

	// Verify git clone was called
	if mockBase.GetExecCommandCallCount() != 1 {
		t.Fatalf("Expected 1 git command call, got %d", mockBase.GetExecCommandCallCount())
	}

	lastCall := mockBase.GetLastExecCommandCall()
	if lastCall == nil {
		t.Fatal("No git command was executed")
	}

	if lastCall.Command != "git" {
		t.Fatalf("Expected git command, got %q", lastCall.Command)
	}

	if len(lastCall.Args) < 1 || lastCall.Args[0] != "clone" {
		t.Fatalf("Expected git clone command, got args: %v", lastCall.Args)
	}
}

func TestSoftInstall_DirectoryExistsWithFiles(t *testing.T) {
	tempDir := t.TempDir()

	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	oldAppDir := paths.Paths.App.Root
	paths.Paths.App.Root = tempDir
	t.Cleanup(func() {
		paths.Paths.App.Root = oldAppDir
	})

	dg, mockBase := createMockDevgita()

	err := dg.SoftInstall()
	if err != nil {
		t.Fatalf("SoftInstall() failed: %v", err)
	}

	// Verify git clone was NOT called (directory already exists)
	if mockBase.GetExecCommandCallCount() != 0 {
		t.Fatalf(
			"Expected 0 git command calls (directory exists), got %d",
			mockBase.GetExecCommandCallCount(),
		)
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

	cleanup := setupTestPaths(t, emptyDir)
	t.Cleanup(cleanup)

	dg, _ := createMockDevgita()

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

	// Save original values
	oldAppDir := paths.Paths.App.Root
	oldConfigDir := paths.Paths.Config.Root
	oldConfigDirPath := configDirPath
	oldGlobalConfigPath := globalConfigPath
	oldZshConfigPath := zshConfigPath

	// Override with test values
	paths.Paths.App.Root = appDir
	paths.Paths.Config.Root = configDir
	configDirPath = devgitaConfigDir
	globalConfigPath = testGlobalConfigPath
	zshConfigPath = filepath.Join(devgitaConfigDir, "devgita.zsh")

	t.Cleanup(func() {
		paths.Paths.App.Root = oldAppDir
		paths.Paths.Config.Root = oldConfigDir
		configDirPath = oldConfigDirPath
		globalConfigPath = oldGlobalConfigPath
		zshConfigPath = oldZshConfigPath
	})

	dg, _ := createMockDevgita()

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
}

func TestUninstall_NoRepositoryNoConfig(t *testing.T) {
	tempDir := t.TempDir()
	appDir := filepath.Join(tempDir, "app")

	cleanup := setupTestPaths(t, appDir)
	t.Cleanup(cleanup)

	dg, _ := createMockDevgita()

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

	cleanup := setupTestPaths(t, tempDir)
	t.Cleanup(cleanup)

	dg, mockBase := createMockDevgita()

	// Mock successful git clone
	mockBase.SetExecCommandResult("Cloning into...", "", nil)

	err := dg.ForceInstall()
	if err != nil {
		t.Fatalf("ForceInstall() failed: %v", err)
	}

	// Verify old file was removed by Uninstall
	if _, err := os.Stat(existingFile); !os.IsNotExist(err) {
		t.Fatal("Expected existing file to be removed by Uninstall")
	}

	// Verify git clone was called
	if mockBase.GetExecCommandCallCount() != 1 {
		t.Fatalf("Expected 1 git clone call, got %d", mockBase.GetExecCommandCallCount())
	}
}

func TestForceConfigure_CreatesConfig(t *testing.T) {
	tempDir := t.TempDir()
	_, appDir, templatesDir, devgitaConfigDir := createTestDirs(t, tempDir)

	createTemplateFiles(t, templatesDir)

	cleanup := setupConfigTestPaths(t, tempDir, appDir, templatesDir)
	t.Cleanup(cleanup)

	dg, _ := createMockDevgita()

	err := dg.ForceConfigure()
	if err != nil {
		t.Fatalf("ForceConfigure() failed: %v", err)
	}

	// Verify config file was created
	configPath := filepath.Join(devgitaConfigDir, "global_config.yaml")
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

	// Verify shell config file (devgita.zsh) was created by RegenerateShellConfig
	// RegenerateShellConfig creates the file at paths.Paths.App.Root/devgita.zsh
	zshConfigPath := filepath.Join(appDir, "devgita.zsh")
	if _, err := os.Stat(zshConfigPath); os.IsNotExist(err) {
		t.Fatalf("Expected shell config file (devgita.zsh) to be created at %s", zshConfigPath)
	}

	// Verify zsh config content exists
	zshContent, err := os.ReadFile(zshConfigPath)
	if err != nil {
		t.Fatalf("Failed to read zsh config file: %v", err)
	}
	if len(zshContent) == 0 {
		t.Fatal("Expected zsh config file to have content")
	}
}

func TestForceConfigure_OverwritesExisting(t *testing.T) {
	tempDir := t.TempDir()
	_, appDir, templatesDir, devgitaConfigDir := createTestDirs(t, tempDir)

	createTemplateFiles(t, templatesDir)

	// Create existing config with old values
	configPath := filepath.Join(devgitaConfigDir, "global_config.yaml")
	oldContent := "app_path: /old/path\nconfig_path: /old/config\n"
	if err := os.WriteFile(configPath, []byte(oldContent), 0o644); err != nil {
		t.Fatalf("Failed to create existing config: %v", err)
	}

	cleanup := setupConfigTestPaths(t, tempDir, appDir, templatesDir)
	t.Cleanup(cleanup)

	dg, _ := createMockDevgita()

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
	_, appDir, templatesDir, devgitaConfigDir := createTestDirs(t, tempDir)

	createTemplateFiles(t, templatesDir)

	// Create existing config with custom values
	configPath := filepath.Join(devgitaConfigDir, "global_config.yaml")
	customContent := "app_path: /custom/path\nconfig_path: /custom/config\ncustom_field: custom_value\n"
	if err := os.WriteFile(configPath, []byte(customContent), 0o644); err != nil {
		t.Fatalf("Failed to create existing config: %v", err)
	}

	// Create zsh config file so SoftConfigure checks pass
	zshPath := filepath.Join(devgitaConfigDir, "devgita.zsh")
	if err := os.WriteFile(zshPath, []byte("# zsh config"), 0o644); err != nil {
		t.Fatalf("Failed to create zsh config: %v", err)
	}

	cleanup := setupConfigTestPaths(t, tempDir, appDir, templatesDir)
	t.Cleanup(cleanup)

	dg, _ := createMockDevgita()

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
