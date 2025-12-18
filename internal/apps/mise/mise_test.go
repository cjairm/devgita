package mise

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/paths"
)

func init() {
	// Initialize logger for tests
	logger.Init(false)
}

func TestNew(t *testing.T) {
	app := New()

	if app == nil {
		t.Fatal("New() returned nil")
	}
}

func TestInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Mise{Cmd: mc}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mc.InstalledPkg != constants.Mise {
		t.Fatalf("expected InstallPackage(%s), got %q", constants.Mise, mc.InstalledPkg)
	}
}

// SKIP: ForceInstall test as per guidelines
// ForceInstall calls Uninstall (which now modifies config) before Install
// Testing this creates complex state management in tests

func TestSoftInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Mise{Cmd: mc}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mc.MaybeInstalled != constants.Mise {
		t.Fatalf("expected MaybeInstallPackage(%s), got %q", constants.Mise, mc.MaybeInstalled)
	}
}

func TestForceConfigure(t *testing.T) {
	// Create temp directories
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	templatesDir := filepath.Join(tempDir, "templates")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("Failed to create templates dir: %v", err)
	}

	// Create global config file
	globalConfigPath := filepath.Join(configDir, constants.AppName, constants.GlobalConfigFile)
	if err := os.MkdirAll(filepath.Dir(globalConfigPath), 0755); err != nil {
		t.Fatalf("Failed to create global config dir: %v", err)
	}

	initialConfig := `app_path: ""
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
`
	if err := os.WriteFile(globalConfigPath, []byte(initialConfig), 0644); err != nil {
		t.Fatalf("Failed to create global config: %v", err)
	}

	// Create simple template
	templatePath := filepath.Join(templatesDir, constants.DevgitaShellTemplate)
	templateContent := `# Test template
{{if .Mise}}
if command -v mise &> /dev/null; then
  eval "$(mise activate zsh)"
fi
{{end}}
`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// Override paths
	oldConfigDir := paths.ConfigDir
	oldAppDir := paths.AppDir
	oldTemplatesAppDir := paths.TemplatesAppDir
	t.Cleanup(func() {
		paths.ConfigDir = oldConfigDir
		paths.AppDir = oldAppDir
		paths.TemplatesAppDir = oldTemplatesAppDir
	})

	paths.ConfigDir = configDir
	paths.AppDir = tempDir
	paths.TemplatesAppDir = templatesDir

	mc := commands.NewMockCommand()
	app := &Mise{Cmd: mc}

	if err := app.ForceConfigure(); err != nil {
		t.Fatalf("ForceConfigure error: %v", err)
	}

	// Verify generated file contains mise activation
	outputPath := filepath.Join(tempDir, "devgita.zsh")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	expectedContent := `eval "$(mise activate zsh)"`
	if !strings.Contains(string(content), expectedContent) {
		t.Fatalf(
			"Expected file to contain %q, but it didn't. Content: %s",
			expectedContent,
			string(content),
		)
	}

	// Verify global config was updated
	updatedConfig, err := os.ReadFile(globalConfigPath)
	if err != nil {
		t.Fatalf("Failed to read updated config: %v", err)
	}
	if !strings.Contains(string(updatedConfig), "mise: true") {
		t.Fatalf(
			"Expected config to have mise enabled. Config: %s",
			string(updatedConfig),
		)
	}
}

func TestSoftConfigure(t *testing.T) {
	t.Run("ConfigureWhenNotEnabled", func(t *testing.T) {
		// Create temp directories
		tempDir := t.TempDir()
		configDir := filepath.Join(tempDir, "config")
		templatesDir := filepath.Join(tempDir, "templates")

		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatalf("Failed to create config dir: %v", err)
		}
		if err := os.MkdirAll(templatesDir, 0755); err != nil {
			t.Fatalf("Failed to create templates dir: %v", err)
		}

		// Create global config with feature disabled
		globalConfigPath := filepath.Join(configDir, constants.AppName, constants.GlobalConfigFile)
		if err := os.MkdirAll(filepath.Dir(globalConfigPath), 0755); err != nil {
			t.Fatalf("Failed to create global config dir: %v", err)
		}

		initialConfig := `shell:
  mise: false
  zoxide: false
  zsh_autosuggestions: false
  zsh_syntax_highlighting: false
  powerlevel10k: false
`
		if err := os.WriteFile(globalConfigPath, []byte(initialConfig), 0644); err != nil {
			t.Fatalf("Failed to create global config: %v", err)
		}

		// Create template
		templatePath := filepath.Join(templatesDir, constants.DevgitaShellTemplate)
		templateContent := `{{if .Mise}}mise-enabled{{end}}`
		if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
			t.Fatalf("Failed to create template: %v", err)
		}

		// Override paths
		oldConfigDir := paths.ConfigDir
		oldAppDir := paths.AppDir
		oldTemplatesAppDir := paths.TemplatesAppDir
		t.Cleanup(func() {
			paths.ConfigDir = oldConfigDir
			paths.AppDir = oldAppDir
			paths.TemplatesAppDir = oldTemplatesAppDir
		})

		paths.ConfigDir = configDir
		paths.AppDir = tempDir
		paths.TemplatesAppDir = templatesDir

		mc := commands.NewMockCommand()
		app := &Mise{Cmd: mc}

		if err := app.SoftConfigure(); err != nil {
			t.Fatalf("SoftConfigure error: %v", err)
		}

		// Verify file was generated
		outputPath := filepath.Join(tempDir, "devgita.zsh")
		content, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("Failed to read generated file: %v", err)
		}

		if !strings.Contains(string(content), "mise-enabled") {
			t.Fatal("Expected template to be rendered with feature enabled")
		}
	})

	t.Run("SkipWhenAlreadyEnabled", func(t *testing.T) {
		// Create temp directories
		tempDir := t.TempDir()
		configDir := filepath.Join(tempDir, "config")

		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatalf("Failed to create config dir: %v", err)
		}

		// Create global config with feature ALREADY enabled
		globalConfigPath := filepath.Join(configDir, constants.AppName, constants.GlobalConfigFile)
		if err := os.MkdirAll(filepath.Dir(globalConfigPath), 0755); err != nil {
			t.Fatalf("Failed to create global config dir: %v", err)
		}

		configWithFeatureEnabled := `shell:
  mise: true
  zoxide: false
  zsh_autosuggestions: false
  zsh_syntax_highlighting: false
  powerlevel10k: false
`
		if err := os.WriteFile(globalConfigPath, []byte(configWithFeatureEnabled), 0644); err != nil {
			t.Fatalf("Failed to create global config: %v", err)
		}

		// Override paths
		oldConfigDir := paths.ConfigDir
		oldAppDir := paths.AppDir
		t.Cleanup(func() {
			paths.ConfigDir = oldConfigDir
			paths.AppDir = oldAppDir
		})

		paths.ConfigDir = configDir
		paths.AppDir = tempDir

		mc := commands.NewMockCommand()
		app := &Mise{Cmd: mc}

		if err := app.SoftConfigure(); err != nil {
			t.Fatalf("SoftConfigure error: %v", err)
		}

		// Verify NO file was generated (should skip when already enabled)
		outputPath := filepath.Join(tempDir, "devgita.zsh")
		if _, err := os.Stat(outputPath); err == nil {
			t.Fatal("Expected no file to be generated when feature already enabled")
		}
	})
}

func TestUninstall(t *testing.T) {
	// Create temp directories
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	templatesDir := filepath.Join(tempDir, "templates")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("Failed to create templates dir: %v", err)
	}

	// Create global config with feature enabled
	globalConfigPath := filepath.Join(configDir, constants.AppName, constants.GlobalConfigFile)
	if err := os.MkdirAll(filepath.Dir(globalConfigPath), 0755); err != nil {
		t.Fatalf("Failed to create global config dir: %v", err)
	}

	initialConfig := `shell:
  mise: true
  zoxide: false
  zsh_autosuggestions: false
  zsh_syntax_highlighting: false
  powerlevel10k: false
`
	if err := os.WriteFile(globalConfigPath, []byte(initialConfig), 0644); err != nil {
		t.Fatalf("Failed to create global config: %v", err)
	}

	// Create template
	templatePath := filepath.Join(templatesDir, constants.DevgitaShellTemplate)
	templateContent := `{{if .Mise}}enabled{{else}}disabled{{end}}`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// Override paths
	oldConfigDir := paths.ConfigDir
	oldAppDir := paths.AppDir
	oldTemplatesAppDir := paths.TemplatesAppDir
	t.Cleanup(func() {
		paths.ConfigDir = oldConfigDir
		paths.AppDir = oldAppDir
		paths.TemplatesAppDir = oldTemplatesAppDir
	})

	paths.ConfigDir = configDir
	paths.AppDir = tempDir
	paths.TemplatesAppDir = templatesDir

	mc := commands.NewMockCommand()
	app := &Mise{Cmd: mc}

	if err := app.Uninstall(); err != nil {
		t.Fatalf("Uninstall error: %v", err)
	}

	// Verify feature was disabled in config
	updatedConfig, err := os.ReadFile(globalConfigPath)
	if err != nil {
		t.Fatalf("Failed to read updated config: %v", err)
	}
	if !strings.Contains(string(updatedConfig), "mise: false") {
		t.Fatalf(
			"Expected config to have mise disabled. Config: %s",
			string(updatedConfig),
		)
	}

	// Verify generated file reflects disabled state
	outputPath := filepath.Join(tempDir, "devgita.zsh")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}
	if !strings.Contains(string(content), "disabled") {
		t.Fatalf("Expected template to show disabled state. Content: %s", string(content))
	}
}

func TestExecuteCommand(t *testing.T) {
	mc := commands.NewMockCommand()
	mockBase := commands.NewMockBaseCommand()
	app := &Mise{Cmd: mc, Base: mockBase}

	t.Run("successful execution", func(t *testing.T) {
		mockBase.SetExecCommandResult("mise 2024.1.0", "", nil)

		err := app.ExecuteCommand("--version")
		if err != nil {
			t.Fatalf("ExecuteCommand failed: %v", err)
		}

		// Verify command was called
		if mockBase.GetExecCommandCallCount() != 1 {
			t.Fatalf("Expected 1 call, got %d", mockBase.GetExecCommandCallCount())
		}

		// Verify parameters
		lastCall := mockBase.GetLastExecCommandCall()
		if lastCall.Command != constants.Mise {
			t.Fatalf("Expected command '%s', got %q", constants.Mise, lastCall.Command)
		}
		if len(lastCall.Args) != 1 || lastCall.Args[0] != "--version" {
			t.Fatalf("Expected args ['--version'], got %v", lastCall.Args)
		}
	})

	t.Run("multiple arguments", func(t *testing.T) {
		mockBase.ResetExecCommand()
		mockBase.SetExecCommandResult("", "", nil)

		err := app.ExecuteCommand("install", "node@20")
		if err != nil {
			t.Fatalf("ExecuteCommand failed: %v", err)
		}

		lastCall := mockBase.GetLastExecCommandCall()
		if len(lastCall.Args) != 2 {
			t.Fatalf("Expected 2 args, got %d", len(lastCall.Args))
		}
		if lastCall.Args[0] != "install" || lastCall.Args[1] != "node@20" {
			t.Fatalf("Expected args ['install', 'node@20'], got %v", lastCall.Args)
		}
	})
}

func TestUseGlobal(t *testing.T) {
	mc := commands.NewMockCommand()
	mockBase := commands.NewMockBaseCommand()
	app := &Mise{Cmd: mc, Base: mockBase}

	t.Run("successful use global", func(t *testing.T) {
		mockBase.SetExecCommandResult("", "", nil)

		err := app.UseGlobal("node", "20")
		if err != nil {
			t.Fatalf("UseGlobal failed: %v", err)
		}

		lastCall := mockBase.GetLastExecCommandCall()
		expectedArgs := []string{"use", "--global", "node@20"}
		if len(lastCall.Args) != len(expectedArgs) {
			t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
		}
		for i, arg := range expectedArgs {
			if lastCall.Args[i] != arg {
				t.Fatalf("Expected arg[%d] to be %q, got %q", i, arg, lastCall.Args[i])
			}
		}
	})

	t.Run("missing language parameter", func(t *testing.T) {
		mockBase.ResetExecCommand()
		err := app.UseGlobal("", "20")
		if err == nil {
			t.Fatal("Expected error for missing language")
		}
		if !strings.Contains(err.Error(), "language") {
			t.Fatalf("Expected error about language, got: %v", err)
		}
	})

	t.Run("missing version parameter", func(t *testing.T) {
		err := app.UseGlobal("node", "")
		if err == nil {
			t.Fatal("Expected error for missing version")
		}
		if !strings.Contains(err.Error(), "version") {
			t.Fatalf("Expected error about version, got: %v", err)
		}
	})
}
