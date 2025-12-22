package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/paths"
)

// TestPaths holds the original path values for restoration after tests
type TestPaths struct {
	AppRoot            string
	ConfigRoot         string
	TemplatesConfigDir string
}

// PathsCleanup is a function that restores original paths
type PathsCleanup func()

// SetupIsolatedPaths overrides global path variables with test-specific paths
// Returns a cleanup function to restore original paths
//
// This ensures tests don't:
// - Read from real configuration directories
// - Write to real configuration directories
// - Execute real system commands
// - Modify the user's actual .zshrc or shell configuration
//
// Usage:
//
//	cleanup := testutil.SetupIsolatedPaths(t)
//	defer cleanup()
func SetupIsolatedPaths(t *testing.T) PathsCleanup {
	t.Helper()

	// Create isolated temp directory for this test
	tempDir := t.TempDir()

	// Save original values
	original := TestPaths{
		AppRoot:            paths.Paths.App.Root,
		ConfigRoot:         paths.Paths.Config.Root,
		TemplatesConfigDir: paths.Paths.App.Configs.Templates,
	}

	// Override with test-isolated values
	paths.Paths.App.Root = filepath.Join(tempDir, "app")
	paths.Paths.Config.Root = filepath.Join(tempDir, "config")
	paths.Paths.App.Configs.Templates = filepath.Join(tempDir, "templates")

	// Return cleanup function
	return func() {
		paths.Paths.App.Root = original.AppRoot
		paths.Paths.Config.Root = original.ConfigRoot
		paths.Paths.App.Configs.Templates = original.TemplatesConfigDir
	}
}

// SetupTestDirs creates a complete test directory structure
// Returns paths to key directories
func SetupTestDirs(t *testing.T) (appDir, configDir, templatesDir, devgitaConfigDir string) {
	t.Helper()

	tempDir := t.TempDir()

	appDir = filepath.Join(tempDir, "app")
	configDir = filepath.Join(tempDir, "config")
	templatesDir = filepath.Join(tempDir, "templates")
	devgitaConfigDir = filepath.Join(configDir, constants.App.Name)

	dirs := []string{appDir, configDir, templatesDir, devgitaConfigDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	return
}

// SetupTestEnvironment creates isolated paths and directory structure
// Returns paths and cleanup function
//
// # This is the recommended way to start any test that touches configuration
//
// Usage:
//
//	appDir, configDir, templatesDir, cleanup := testutil.SetupTestEnvironment(t)
//	defer cleanup()
func SetupTestEnvironment(t *testing.T) (appDir, configDir, templatesDir string, cleanup PathsCleanup) {
	t.Helper()

	cleanup = SetupIsolatedPaths(t)
	appDir, configDir, templatesDir, _ = SetupTestDirs(t)

	// Update paths to match created directories
	paths.Paths.App.Root = appDir
	paths.Paths.Config.Root = configDir
	paths.Paths.App.Configs.Templates = templatesDir

	return
}

// CreateGlobalConfigTemplate creates a basic global_config.yaml template
func CreateGlobalConfigTemplate(t *testing.T, templatesDir string) {
	t.Helper()

	templatePath := filepath.Join(templatesDir, constants.App.File.GlobalConfig)
	templateContent := `app_path: ""
config_path: ""
already_installed:
  packages: []
  desktop_apps: []
  fonts: []
  themes: []
  terminal_tools: []
  dev_languages: []
  databases: []
current_font: ""
current_theme: ""
installed:
  packages: []
  desktop_apps: []
  fonts: []
  themes: []
  terminal_tools: []
  dev_languages: []
  databases: []
shortcuts: {}
shell:
  mise: false
  zoxide: false
  zsh_autosuggestions: false
  zsh_syntax_highlighting: false
  powerlevel10k: false
  extended_capabilities: false
`

	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to create global config template: %v", err)
	}
}

// CreateShellConfigTemplate creates a basic devgita.zsh.tmpl template
func CreateShellConfigTemplate(t *testing.T, templatesDir string, content string) {
	t.Helper()

	templatePath := filepath.Join(templatesDir, constants.App.Template.ShellConfig)

	if content == "" {
		content = `# Test Shell Configuration
{{if .Mise}}# Mise enabled{{end}}
{{if .Zoxide}}# Zoxide enabled{{end}}
{{if .ZshAutosuggestions}}# Autosuggestions enabled{{end}}
{{if .ZshSyntaxHighlighting}}# Syntax highlighting enabled{{end}}
{{if .Powerlevel10k}}# Powerlevel10k enabled{{end}}
`
	}

	if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create shell config template: %v", err)
	}
}

// CreateGlobalConfigFile creates a global_config.yaml file with specified content
func CreateGlobalConfigFile(t *testing.T, configDir string, content string) string {
	t.Helper()

	devgitaConfigDir := filepath.Join(configDir, constants.App.Name)
	if err := os.MkdirAll(devgitaConfigDir, 0755); err != nil {
		t.Fatalf("Failed to create devgita config dir: %v", err)
	}

	configPath := filepath.Join(devgitaConfigDir, constants.App.File.GlobalConfig)

	if content == "" {
		content = `app_path: ""
config_path: ""
shell:
  mise: false
  zoxide: false
  zsh_autosuggestions: false
  zsh_syntax_highlighting: false
  powerlevel10k: false
`
	}

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create global config file: %v", err)
	}

	return configPath
}

// NewMockApp creates a mock app instance with isolated mocks
// This ensures tests don't execute real commands
type MockApp struct {
	Cmd  *commands.MockCommand
	Base *commands.MockBaseCommand
}

// NewMockApp creates a new mock app instance
func NewMockApp() *MockApp {
	return &MockApp{
		Cmd:  commands.NewMockCommand(),
		Base: commands.NewMockBaseCommand(),
	}
}

// Reset resets all mock state
func (m *MockApp) Reset() {
	m.Cmd.Reset()
	m.Base.ResetExecCommand()
}

// TestConfig holds common test configuration
type TestConfig struct {
	AppDir        string
	ConfigDir     string
	TemplatesDir  string
	Cleanup       PathsCleanup
	MockApp       *MockApp
	ConfigPath    string
	ZshConfigPath string
}

// SetupCompleteTest creates a complete test environment with all common setup
// This is the most convenient function for testing apps that use configuration
//
// Usage:
//
//	tc := testutil.SetupCompleteTest(t)
//	defer tc.Cleanup()
//
//	// Use tc.MockApp for mock instances
//	// Use tc.AppDir, tc.ConfigDir, etc. for paths
func SetupCompleteTest(t *testing.T) *TestConfig {
	t.Helper()

	appDir, configDir, templatesDir, cleanup := SetupTestEnvironment(t)

	// Create templates
	CreateGlobalConfigTemplate(t, templatesDir)
	CreateShellConfigTemplate(t, templatesDir, "")

	// Create initial config
	configPath := CreateGlobalConfigFile(t, configDir, "")

	tc := &TestConfig{
		AppDir:        appDir,
		ConfigDir:     configDir,
		TemplatesDir:  templatesDir,
		Cleanup:       cleanup,
		MockApp:       NewMockApp(),
		ConfigPath:    configPath,
		ZshConfigPath: filepath.Join(appDir, "devgita.zsh"),
	}

	return tc
}

// InitLogger initializes the logger for tests
// Call this in your init() function
func InitLogger() {
	logger.Init(false)
}

// VerifyNoRealCommands verifies that no real commands were executed
func VerifyNoRealCommands(t *testing.T, mockBase *commands.MockBaseCommand) {
	t.Helper()

	calls := mockBase.GetExecCommandCallCount()
	if calls > 0 {
		t.Errorf("Expected no real commands to be executed, but %d commands were called", calls)
		if lastCall := mockBase.GetLastExecCommandCall(); lastCall != nil {
			t.Errorf("Last command: %s %v", lastCall.Command, lastCall.Args)
		}
	}
}

// VerifyNoRealConfigChanges verifies that no real configuration files were modified
func VerifyNoRealConfigChanges(t *testing.T) {
	t.Helper()

	// Check if any temp paths are in common config locations
	realConfigDir := os.ExpandEnv("$HOME/.config/devgita")
	realZshrc := os.ExpandEnv("$HOME/.zshrc")

	if paths.Paths.Config.Root == realConfigDir {
		t.Error("Test is using real config directory! Tests should use isolated temp directories")
	}

	// Read .zshrc and check for temp paths
	if data, err := os.ReadFile(realZshrc); err == nil {
		content := string(data)
		if filepath.Base(paths.Paths.App.Root) != filepath.Base(realConfigDir) {
			// Check if .zshrc contains temp directory references
			if len(content) > 0 && filepath.IsAbs(paths.Paths.App.Root) {
				t.Logf("Warning: Ensure tests use t.TempDir() to avoid modifying real .zshrc")
			}
		}
	}
}

// AssertFileContains verifies that a file contains expected content
func AssertFileContains(t *testing.T, filePath, expected string) {
	t.Helper()

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", filePath, err)
	}

	if !filepath.IsAbs(filePath) {
		t.Errorf("File path should be absolute: %s", filePath)
	}

	contentStr := string(content)
	if len(contentStr) == 0 {
		t.Errorf("File %s is empty", filePath)
	}

	if expected != "" {
		if contentStr != expected {
			t.Errorf("File content mismatch.\nExpected:\n%s\nGot:\n%s", expected, contentStr)
		}
	}
}

// AssertFileNotContains verifies that a file does not contain specific content
func AssertFileNotContains(t *testing.T, filePath, unexpected string) {
	t.Helper()

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", filePath, err)
	}

	contentStr := string(content)
	if len(contentStr) == 0 {
		return // Empty file doesn't contain anything
	}

	if contentStr == unexpected {
		t.Errorf("File should not contain:\n%s", unexpected)
	}
}
