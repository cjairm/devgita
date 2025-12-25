package opencode

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/testutil"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/paths"
)

func init() {
	// Initialize logger for tests
	testutil.InitLogger()
}

func TestNew(t *testing.T) {
	app := New()

	if app == nil {
		t.Fatal("New() returned nil")
	}

	if app.Cmd == nil {
		t.Error("Expected Cmd to be initialized")
	}

	if app.Base == nil {
		t.Error("Expected Base to be initialized")
	}
}

func TestInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &OpenCode{Cmd: mockApp.Cmd, Base: mockApp.Base}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}

	if mockApp.Cmd.InstalledPkg != constants.OpenCode {
		t.Fatalf(
			"expected InstallPackage(%s), got %q",
			constants.OpenCode,
			mockApp.Cmd.InstalledPkg,
		)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestSoftInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &OpenCode{Cmd: mockApp.Cmd, Base: mockApp.Base}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}

	if mockApp.Cmd.MaybeInstalled != constants.OpenCode {
		t.Fatalf(
			"expected MaybeInstallPackage(%s), got %q",
			constants.OpenCode,
			mockApp.Cmd.MaybeInstalled,
		)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestForceInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &OpenCode{Cmd: mockApp.Cmd, Base: mockApp.Base}

	// ForceInstall calls Uninstall first, which returns error
	err := app.ForceInstall()
	if err == nil {
		t.Fatal("Expected ForceInstall to return error (Uninstall not supported)")
	}

	if !strings.Contains(err.Error(), "uninstall") {
		t.Fatalf("Expected error about uninstall, got: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestUninstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &OpenCode{Cmd: mockApp.Cmd, Base: mockApp.Base}

	err := app.Uninstall()
	if err == nil {
		t.Fatal("Expected Uninstall to return error")
	}

	expectedMsg := "opencode uninstall not supported through devgita"
	if !strings.Contains(err.Error(), "uninstall not supported") {
		t.Fatalf("Expected error message %q, got: %v", expectedMsg, err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestUpdate(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &OpenCode{Cmd: mockApp.Cmd, Base: mockApp.Base}

	err := app.Update()
	if err == nil {
		t.Fatal("Expected Update to return error")
	}

	if !strings.Contains(err.Error(), "not implemented") {
		t.Fatalf("Expected error about not implemented, got: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestForceConfigure(t *testing.T) {
	t.Run("ConfigureWithDefaultTheme", func(t *testing.T) {
		tc := testutil.SetupCompleteTest(t)
		defer tc.Cleanup()

		// Setup OpenCode-specific paths
		appConfigDir := filepath.Join(tc.AppDir, "configs", "opencode")
		userConfigDir := filepath.Join(tc.ConfigDir, "opencode")

		// Create source directories
		if err := os.MkdirAll(appConfigDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(appConfigDir, "themes"), 0755); err != nil {
			t.Fatal(err)
		}

		// Create template file
		templateContent := `{
  "version": "1.0.0",
  "theme": "{{ .Theme }}",
  "settings": {
    "fontSize": 14,
    "fontFamily": "JetBrains Mono"
  }
}`
		templatePath := filepath.Join(appConfigDir, "opencode.json.tmpl")
		if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Create default theme file
		themeContent := `{
  "name": "Devgita Gruvbox",
  "type": "dark",
  "colors": {
    "background": "#282828",
    "foreground": "#ebdbb2"
  }
}`
		themeSourcePath := filepath.Join(appConfigDir, "themes", "default.json")
		if err := os.WriteFile(themeSourcePath, []byte(themeContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Override paths
		paths.Paths.App.Configs.OpenCode = appConfigDir
		paths.Paths.Config.OpenCode = userConfigDir

		app := &OpenCode{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

		// Test with default theme
		options := ConfigureOptions{Theme: ""}
		if err := app.ForceConfigure(options); err != nil {
			t.Fatalf("ForceConfigure error: %v", err)
		}

		// Verify config file was generated
		configPath := filepath.Join(userConfigDir, "opencode.json")
		if _, err := os.Stat(configPath); err != nil {
			t.Fatalf("Expected config file at %s: %v", configPath, err)
		}

		// Verify config content
		content, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("Failed to read config: %v", err)
		}

		configStr := string(content)
		if !strings.Contains(configStr, `"theme": "default"`) {
			t.Errorf("Expected theme to be 'default', got: %s", configStr)
		}
		if !strings.Contains(configStr, "JetBrains Mono") {
			t.Error("Expected config to contain font family")
		}

		// Verify theme file was copied
		themePath := filepath.Join(userConfigDir, "themes", "default.json")
		if _, err := os.Stat(themePath); err != nil {
			t.Fatalf("Expected theme file at %s: %v", themePath, err)
		}

		// Verify theme content
		themeContentRead, err := os.ReadFile(themePath)
		if err != nil {
			t.Fatalf("Failed to read theme: %v", err)
		}
		if !strings.Contains(string(themeContentRead), "Devgita Gruvbox") {
			t.Error("Expected theme file to contain Gruvbox theme")
		}

		testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
		testutil.VerifyNoRealConfigChanges(t)
	})

	t.Run("ConfigureWithCustomTheme", func(t *testing.T) {
		tc := testutil.SetupCompleteTest(t)
		defer tc.Cleanup()

		// Setup OpenCode-specific paths
		appConfigDir := filepath.Join(tc.AppDir, "configs", "opencode")
		userConfigDir := filepath.Join(tc.ConfigDir, "opencode")

		// Create source directories
		if err := os.MkdirAll(appConfigDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create template file
		templateContent := `{
  "version": "1.0.0",
  "theme": "{{ .Theme }}",
  "settings": {
    "fontSize": 14
  }
}`
		templatePath := filepath.Join(appConfigDir, "opencode.json.tmpl")
		if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Override paths
		paths.Paths.App.Configs.OpenCode = appConfigDir
		paths.Paths.Config.OpenCode = userConfigDir

		app := &OpenCode{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

		// Test with custom theme
		options := ConfigureOptions{Theme: "solarized-dark"}
		if err := app.ForceConfigure(options); err != nil {
			t.Fatalf("ForceConfigure error: %v", err)
		}

		// Verify config file was generated
		configPath := filepath.Join(userConfigDir, "opencode.json")
		if _, err := os.Stat(configPath); err != nil {
			t.Fatalf("Expected config file at %s: %v", configPath, err)
		}

		// Verify config content has custom theme
		content, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("Failed to read config: %v", err)
		}

		configStr := string(content)
		if !strings.Contains(configStr, `"theme": "solarized-dark"`) {
			t.Errorf("Expected theme to be 'solarized-dark', got: %s", configStr)
		}

		// Verify theme file was NOT copied (custom theme)
		themePath := filepath.Join(userConfigDir, "themes", "default.json")
		if _, err := os.Stat(themePath); err == nil {
			t.Error("Expected NO theme file for custom theme")
		}

		testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
		testutil.VerifyNoRealConfigChanges(t)
	})

	t.Run("ConfigureWithDefaultKeyword", func(t *testing.T) {
		tc := testutil.SetupCompleteTest(t)
		defer tc.Cleanup()

		// Setup OpenCode-specific paths
		appConfigDir := filepath.Join(tc.AppDir, "configs", "opencode")
		userConfigDir := filepath.Join(tc.ConfigDir, "opencode")

		// Create source directories
		if err := os.MkdirAll(appConfigDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(appConfigDir, "themes"), 0755); err != nil {
			t.Fatal(err)
		}

		// Create template file
		templateContent := `{"theme": "{{ .Theme }}"}`
		templatePath := filepath.Join(appConfigDir, "opencode.json.tmpl")
		if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Create theme file
		themeSourcePath := filepath.Join(appConfigDir, "themes", "default.json")
		if err := os.WriteFile(themeSourcePath, []byte(`{"name": "test"}`), 0644); err != nil {
			t.Fatal(err)
		}

		// Override paths
		paths.Paths.App.Configs.OpenCode = appConfigDir
		paths.Paths.Config.OpenCode = userConfigDir

		app := &OpenCode{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

		// Test with explicit "default" theme
		options := ConfigureOptions{Theme: "default"}
		if err := app.ForceConfigure(options); err != nil {
			t.Fatalf("ForceConfigure error: %v", err)
		}

		// Verify theme file was copied (explicit default)
		themePath := filepath.Join(userConfigDir, "themes", "default.json")
		if _, err := os.Stat(themePath); err != nil {
			t.Fatalf("Expected theme file at %s: %v", themePath, err)
		}

		testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
	})

	t.Run("RemovesExistingConfigDirectory", func(t *testing.T) {
		tc := testutil.SetupCompleteTest(t)
		defer tc.Cleanup()

		// Setup OpenCode-specific paths
		appConfigDir := filepath.Join(tc.AppDir, "configs", "opencode")
		userConfigDir := filepath.Join(tc.ConfigDir, "opencode")

		// Create source directories
		if err := os.MkdirAll(appConfigDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create template file
		templatePath := filepath.Join(appConfigDir, "opencode.json.tmpl")
		if err := os.WriteFile(templatePath, []byte(`{"theme": "{{ .Theme }}"}`), 0644); err != nil {
			t.Fatal(err)
		}

		// Create existing config with old file
		if err := os.MkdirAll(userConfigDir, 0755); err != nil {
			t.Fatal(err)
		}
		oldFilePath := filepath.Join(userConfigDir, "old-file.json")
		if err := os.WriteFile(oldFilePath, []byte("old content"), 0644); err != nil {
			t.Fatal(err)
		}

		// Override paths
		paths.Paths.App.Configs.OpenCode = appConfigDir
		paths.Paths.Config.OpenCode = userConfigDir

		app := &OpenCode{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

		options := ConfigureOptions{Theme: "custom"}
		if err := app.ForceConfigure(options); err != nil {
			t.Fatalf("ForceConfigure error: %v", err)
		}

		// Verify old file was removed
		if _, err := os.Stat(oldFilePath); err == nil {
			t.Error("Expected old file to be removed")
		}

		// Verify new config exists
		configPath := filepath.Join(userConfigDir, "opencode.json")
		if _, err := os.Stat(configPath); err != nil {
			t.Fatalf("Expected new config file: %v", err)
		}

		testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
	})
}

func TestSoftConfigure(t *testing.T) {
	t.Run("SkipWhenAlreadyConfigured", func(t *testing.T) {
		tc := testutil.SetupCompleteTest(t)
		defer tc.Cleanup()

		// Setup OpenCode-specific paths
		userConfigDir := filepath.Join(tc.ConfigDir, "opencode")

		// Create existing config file (marker)
		if err := os.MkdirAll(userConfigDir, 0755); err != nil {
			t.Fatal(err)
		}
		markerPath := filepath.Join(userConfigDir, "opencode.json")
		if err := os.WriteFile(markerPath, []byte(`{"theme": "existing"}`), 0644); err != nil {
			t.Fatal(err)
		}

		// Override paths
		paths.Paths.Config.OpenCode = userConfigDir

		app := &OpenCode{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

		options := ConfigureOptions{Theme: "new-theme"}
		if err := app.SoftConfigure(options); err != nil {
			t.Fatalf("SoftConfigure error: %v", err)
		}

		// Verify config was NOT changed
		content, err := os.ReadFile(markerPath)
		if err != nil {
			t.Fatalf("Failed to read config: %v", err)
		}

		if !strings.Contains(string(content), "existing") {
			t.Error("Expected existing config to be preserved")
		}
		if strings.Contains(string(content), "new-theme") {
			t.Error("Expected config NOT to be overwritten")
		}

		testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
	})

	t.Run("ConfigureWhenNotConfigured", func(t *testing.T) {
		tc := testutil.SetupCompleteTest(t)
		defer tc.Cleanup()

		// Setup OpenCode-specific paths
		appConfigDir := filepath.Join(tc.AppDir, "configs", "opencode")
		userConfigDir := filepath.Join(tc.ConfigDir, "opencode")

		// Create source directories
		if err := os.MkdirAll(appConfigDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create template file
		templatePath := filepath.Join(appConfigDir, "opencode.json.tmpl")
		if err := os.WriteFile(templatePath, []byte(`{"theme": "{{ .Theme }}"}`), 0644); err != nil {
			t.Fatal(err)
		}

		// Override paths
		paths.Paths.App.Configs.OpenCode = appConfigDir
		paths.Paths.Config.OpenCode = userConfigDir

		app := &OpenCode{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

		// SoftConfigure when no existing config
		options := ConfigureOptions{Theme: "fresh"}
		if err := app.SoftConfigure(options); err != nil {
			t.Fatalf("SoftConfigure error: %v", err)
		}

		// Verify config was created
		configPath := filepath.Join(userConfigDir, "opencode.json")
		if _, err := os.Stat(configPath); err != nil {
			t.Fatalf("Expected config file to be created: %v", err)
		}

		content, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("Failed to read config: %v", err)
		}

		if !strings.Contains(string(content), `"theme": "fresh"`) {
			t.Error("Expected config to contain new theme")
		}

		testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
	})

	t.Run("SkipWhenAlreadyInstalledByDevgita", func(t *testing.T) {
		tc := testutil.SetupCompleteTest(t)
		defer tc.Cleanup()

		// Create global config with opencode already installed
		globalConfigContent := `app_path: ""
config_path: ""
installed:
  packages:
    - opencode
shell:
  mise: false
`
		globalConfigPath := filepath.Join(tc.ConfigDir, constants.App.Name, constants.App.File.GlobalConfig)
		if err := os.WriteFile(globalConfigPath, []byte(globalConfigContent), 0644); err != nil {
			t.Fatal(err)
		}

		userConfigDir := filepath.Join(tc.ConfigDir, "opencode")

		// Override paths
		paths.Paths.Config.OpenCode = userConfigDir

		app := &OpenCode{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

		options := ConfigureOptions{Theme: "test"}
		if err := app.SoftConfigure(options); err != nil {
			t.Fatalf("SoftConfigure error: %v", err)
		}

		// Verify NO config was created
		configPath := filepath.Join(userConfigDir, "opencode.json")
		if _, err := os.Stat(configPath); err == nil {
			t.Error("Expected NO config file when already installed")
		}

		testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
	})

	t.Run("SkipWhenAlreadyInstalled", func(t *testing.T) {
		tc := testutil.SetupCompleteTest(t)
		defer tc.Cleanup()

		// Create global config with opencode in already_installed
		globalConfigContent := `app_path: ""
config_path: ""
already_installed:
  packages:
    - opencode
shell:
  mise: false
`
		globalConfigPath := filepath.Join(tc.ConfigDir, constants.App.Name, constants.App.File.GlobalConfig)
		if err := os.WriteFile(globalConfigPath, []byte(globalConfigContent), 0644); err != nil {
			t.Fatal(err)
		}

		userConfigDir := filepath.Join(tc.ConfigDir, "opencode")

		// Override paths
		paths.Paths.Config.OpenCode = userConfigDir

		app := &OpenCode{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

		options := ConfigureOptions{Theme: "test"}
		if err := app.SoftConfigure(options); err != nil {
			t.Fatalf("SoftConfigure error: %v", err)
		}

		// Verify NO config was created
		configPath := filepath.Join(userConfigDir, "opencode.json")
		if _, err := os.Stat(configPath); err == nil {
			t.Error("Expected NO config file when already_installed")
		}

		testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
	})
}

func TestExecuteCommand(t *testing.T) {
	t.Run("SuccessfulExecution", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		app := &OpenCode{Cmd: mockApp.Cmd, Base: mockApp.Base}

		mockApp.Base.SetExecCommandResult("OpenCode 1.0.0", "", nil)

		err := app.ExecuteCommand("--version")
		if err != nil {
			t.Fatalf("ExecuteCommand failed: %v", err)
		}

		// Verify command was called
		if mockApp.Base.GetExecCommandCallCount() != 1 {
			t.Fatalf("Expected 1 call, got %d", mockApp.Base.GetExecCommandCallCount())
		}

		// Verify parameters
		lastCall := mockApp.Base.GetLastExecCommandCall()
		if lastCall.Command != constants.OpenCode {
			t.Fatalf("Expected command '%s', got %q", constants.OpenCode, lastCall.Command)
		}
		if len(lastCall.Args) != 1 || lastCall.Args[0] != "--version" {
			t.Fatalf("Expected args ['--version'], got %v", lastCall.Args)
		}
	})

	t.Run("MultipleArguments", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		app := &OpenCode{Cmd: mockApp.Cmd, Base: mockApp.Base}

		mockApp.Base.SetExecCommandResult("", "", nil)

		err := app.ExecuteCommand("file1.go", "file2.go")
		if err != nil {
			t.Fatalf("ExecuteCommand failed: %v", err)
		}

		lastCall := mockApp.Base.GetLastExecCommandCall()
		if len(lastCall.Args) != 2 {
			t.Fatalf("Expected 2 args, got %d", len(lastCall.Args))
		}
		if lastCall.Args[0] != "file1.go" || lastCall.Args[1] != "file2.go" {
			t.Fatalf("Expected args ['file1.go', 'file2.go'], got %v", lastCall.Args)
		}
	})

	t.Run("NoArguments", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		app := &OpenCode{Cmd: mockApp.Cmd, Base: mockApp.Base}

		mockApp.Base.SetExecCommandResult("", "", nil)

		err := app.ExecuteCommand()
		if err != nil {
			t.Fatalf("ExecuteCommand failed: %v", err)
		}

		lastCall := mockApp.Base.GetLastExecCommandCall()
		if len(lastCall.Args) != 0 {
			t.Fatalf("Expected 0 args, got %d", len(lastCall.Args))
		}
	})

	t.Run("CommandError", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		app := &OpenCode{Cmd: mockApp.Cmd, Base: mockApp.Base}

		expectedErr := fmt.Errorf("command failed")
		mockApp.Base.SetExecCommandResult("", "error output", expectedErr)

		err := app.ExecuteCommand("--invalid")
		if err == nil {
			t.Fatal("Expected error from ExecuteCommand")
		}

		if !strings.Contains(err.Error(), "opencode command execution failed") {
			t.Fatalf("Expected wrapped error, got: %v", err)
		}
	})
}

func TestConfigureOptions(t *testing.T) {
	t.Run("EmptyThemeDefaultsToDefault", func(t *testing.T) {
		options := ConfigureOptions{Theme: ""}

		// Empty theme should be treated as default
		expectedTheme := DEFAULT_THEME_NAME
		actualTheme := options.Theme
		if actualTheme == "" {
			actualTheme = DEFAULT_THEME_NAME
		}

		if actualTheme != expectedTheme {
			t.Errorf("Expected empty theme to default to %q, got %q", expectedTheme, actualTheme)
		}
	})

	t.Run("CustomThemePreserved", func(t *testing.T) {
		customTheme := "my-custom-theme"
		options := ConfigureOptions{Theme: customTheme}

		if options.Theme != customTheme {
			t.Errorf("Expected theme %q, got %q", customTheme, options.Theme)
		}
	})
}

func TestDefaultThemeName(t *testing.T) {
	if DEFAULT_THEME_NAME != "default" {
		t.Errorf("Expected DEFAULT_THEME_NAME to be 'default', got %q", DEFAULT_THEME_NAME)
	}
}
