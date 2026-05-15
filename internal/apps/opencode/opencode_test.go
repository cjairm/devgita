package opencode

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/testutil"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/paths"
)

func init() {
	testutil.InitLogger()
}

func setupSharedDir(t *testing.T, baseDir string) {
	t.Helper()
	sharedDir := filepath.Join(baseDir, "configs", "shared")
	for _, sub := range []string{"skills", "commands", "agents"} {
		if err := os.MkdirAll(filepath.Join(sharedDir, sub), 0755); err != nil {
			t.Fatal(err)
		}
	}
	oldShared := paths.Paths.App.Configs.Shared
	t.Cleanup(func() { paths.Paths.App.Configs.Shared = oldShared })
	paths.Paths.App.Configs.Shared = sharedDir
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

func TestNameAndKind(t *testing.T) {
	o := &OpenCode{}
	if o.Name() != constants.OpenCode {
		t.Errorf("expected Name() %q, got %q", constants.OpenCode, o.Name())
	}
	if o.Kind() != apps.KindTerminal {
		t.Errorf("expected Kind() KindTerminal, got %v", o.Kind())
	}
}

func TestInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &OpenCode{Cmd: mockApp.Cmd, Base: mockApp.Base}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mockApp.Cmd.InstalledPkg != constants.OpenCode {
		t.Fatalf("expected InstallPackage(%s), got %q", constants.OpenCode, mockApp.Cmd.InstalledPkg)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestForceInstall(t *testing.T) {
	testutil.IsolateXDGDirs(t)
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	userConfigDir := filepath.Join(tc.ConfigDir, "opencode")
	oldOpenCodeDir := paths.Paths.Config.OpenCode
	t.Cleanup(func() { paths.Paths.Config.OpenCode = oldOpenCodeDir })
	paths.Paths.Config.OpenCode = userConfigDir

	app := &OpenCode{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

	if err := app.ForceInstall(); err != nil {
		t.Fatalf("ForceInstall() error: %v", err)
	}
	if tc.MockApp.Cmd.InstalledPkg != constants.OpenCode {
		t.Errorf("expected Install to be called, got %q", tc.MockApp.Cmd.InstalledPkg)
	}

	testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
}

func TestSoftInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &OpenCode{Cmd: mockApp.Cmd, Base: mockApp.Base}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mockApp.Cmd.MaybeInstalled != constants.OpenCode {
		t.Fatalf("expected MaybeInstallPackage(%s), got %q", constants.OpenCode, mockApp.Cmd.MaybeInstalled)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestUninstall(t *testing.T) {
	testutil.IsolateXDGDirs(t)
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	userConfigDir := filepath.Join(tc.ConfigDir, "opencode")
	if err := os.MkdirAll(userConfigDir, 0755); err != nil {
		t.Fatal(err)
	}
	oldOpenCodeDir := paths.Paths.Config.OpenCode
	t.Cleanup(func() { paths.Paths.Config.OpenCode = oldOpenCodeDir })
	paths.Paths.Config.OpenCode = userConfigDir

	app := &OpenCode{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

	if err := app.Uninstall(); err != nil {
		t.Fatalf("Uninstall error: %v", err)
	}
	if tc.MockApp.Cmd.UninstalledPkg != constants.OpenCode {
		t.Errorf("expected UninstallPackage(%s), got %q", constants.OpenCode, tc.MockApp.Cmd.UninstalledPkg)
	}

	testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
}

func TestUpdate(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &OpenCode{Cmd: mockApp.Cmd, Base: mockApp.Base}

	err := app.Update()
	if err == nil {
		t.Fatal("expected Update to return error")
	}
	if !errors.Is(err, apps.ErrUpdateNotSupported) {
		t.Errorf("expected ErrUpdateNotSupported, got: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestForceConfigure(t *testing.T) {
	t.Run("ConfigureWithDefaultTheme", func(t *testing.T) {
		testutil.IsolateXDGDirs(t)
		tc := testutil.SetupCompleteTest(t)
		defer tc.Cleanup()

		appConfigDir := filepath.Join(tc.AppDir, "configs", "opencode")
		userConfigDir := filepath.Join(tc.ConfigDir, "opencode")

		if err := os.MkdirAll(appConfigDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(appConfigDir, "themes"), 0755); err != nil {
			t.Fatal(err)
		}

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

		themeContent := `{"name": "Devgita Gruvbox", "type": "dark"}`
		themeSourcePath := filepath.Join(appConfigDir, "themes", "default.json")
		if err := os.WriteFile(themeSourcePath, []byte(themeContent), 0644); err != nil {
			t.Fatal(err)
		}

		setupSharedDir(t, tc.AppDir)

		oldAppConfigs := paths.Paths.App.Configs.OpenCode
		t.Cleanup(func() { paths.Paths.App.Configs.OpenCode = oldAppConfigs })
		paths.Paths.App.Configs.OpenCode = appConfigDir

		oldConfigOpenCode := paths.Paths.Config.OpenCode
		t.Cleanup(func() { paths.Paths.Config.OpenCode = oldConfigOpenCode })
		paths.Paths.Config.OpenCode = userConfigDir

		app := &OpenCode{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

		if err := app.ForceConfigure(); err != nil {
			t.Fatalf("ForceConfigure error: %v", err)
		}

		configPath := filepath.Join(userConfigDir, "opencode.json")
		if _, err := os.Stat(configPath); err != nil {
			t.Fatalf("Expected config file at %s: %v", configPath, err)
		}

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

		themePath := filepath.Join(userConfigDir, "themes", "default.json")
		if _, err := os.Stat(themePath); err != nil {
			t.Fatalf("Expected theme file at %s: %v", themePath, err)
		}

		themeContentRead, err := os.ReadFile(themePath)
		if err != nil {
			t.Fatalf("Failed to read theme: %v", err)
		}
		if !strings.Contains(string(themeContentRead), "Devgita Gruvbox") {
			t.Error("Expected theme file to contain Gruvbox theme")
		}

		testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
	})

	t.Run("RemovesExistingConfigDirectory", func(t *testing.T) {
		testutil.IsolateXDGDirs(t)
		tc := testutil.SetupCompleteTest(t)
		defer tc.Cleanup()

		appConfigDir := filepath.Join(tc.AppDir, "configs", "opencode")
		userConfigDir := filepath.Join(tc.ConfigDir, "opencode")

		if err := os.MkdirAll(appConfigDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(appConfigDir, "themes"), 0755); err != nil {
			t.Fatal(err)
		}

		templatePath := filepath.Join(appConfigDir, "opencode.json.tmpl")
		if err := os.WriteFile(templatePath, []byte(`{"theme": "{{ .Theme }}"}`), 0644); err != nil {
			t.Fatal(err)
		}
		themeSourcePath := filepath.Join(appConfigDir, "themes", "default.json")
		if err := os.WriteFile(themeSourcePath, []byte(`{"name": "test"}`), 0644); err != nil {
			t.Fatal(err)
		}

		if err := os.MkdirAll(userConfigDir, 0755); err != nil {
			t.Fatal(err)
		}
		oldFilePath := filepath.Join(userConfigDir, "old-file.json")
		if err := os.WriteFile(oldFilePath, []byte("old content"), 0644); err != nil {
			t.Fatal(err)
		}

		setupSharedDir(t, tc.AppDir)

		oldAppConfigs := paths.Paths.App.Configs.OpenCode
		t.Cleanup(func() { paths.Paths.App.Configs.OpenCode = oldAppConfigs })
		paths.Paths.App.Configs.OpenCode = appConfigDir

		oldConfigOpenCode := paths.Paths.Config.OpenCode
		t.Cleanup(func() { paths.Paths.Config.OpenCode = oldConfigOpenCode })
		paths.Paths.Config.OpenCode = userConfigDir

		app := &OpenCode{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

		if err := app.ForceConfigure(); err != nil {
			t.Fatalf("ForceConfigure error: %v", err)
		}

		if _, err := os.Stat(oldFilePath); err == nil {
			t.Error("Expected old file to be removed")
		}

		configPath := filepath.Join(userConfigDir, "opencode.json")
		if _, err := os.Stat(configPath); err != nil {
			t.Fatalf("Expected new config file: %v", err)
		}

		testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
	})
}

func TestSoftConfigure(t *testing.T) {
	t.Run("SkipWhenAlreadyConfigured", func(t *testing.T) {
		testutil.IsolateXDGDirs(t)
		tc := testutil.SetupCompleteTest(t)
		defer tc.Cleanup()

		userConfigDir := filepath.Join(tc.ConfigDir, "opencode")
		if err := os.MkdirAll(userConfigDir, 0755); err != nil {
			t.Fatal(err)
		}
		markerPath := filepath.Join(userConfigDir, "opencode.json")
		if err := os.WriteFile(markerPath, []byte(`{"theme": "existing"}`), 0644); err != nil {
			t.Fatal(err)
		}

		oldConfigOpenCode := paths.Paths.Config.OpenCode
		t.Cleanup(func() { paths.Paths.Config.OpenCode = oldConfigOpenCode })
		paths.Paths.Config.OpenCode = userConfigDir

		app := &OpenCode{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

		if err := app.SoftConfigure(); err != nil {
			t.Fatalf("SoftConfigure error: %v", err)
		}

		content, err := os.ReadFile(markerPath)
		if err != nil {
			t.Fatalf("Failed to read config: %v", err)
		}
		if !strings.Contains(string(content), "existing") {
			t.Error("Expected existing config to be preserved")
		}

		testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
	})

	t.Run("ConfigureWhenNotConfigured", func(t *testing.T) {
		testutil.IsolateXDGDirs(t)
		tc := testutil.SetupCompleteTest(t)
		defer tc.Cleanup()

		appConfigDir := filepath.Join(tc.AppDir, "configs", "opencode")
		userConfigDir := filepath.Join(tc.ConfigDir, "opencode")

		if err := os.MkdirAll(appConfigDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(appConfigDir, "themes"), 0755); err != nil {
			t.Fatal(err)
		}

		templatePath := filepath.Join(appConfigDir, "opencode.json.tmpl")
		if err := os.WriteFile(templatePath, []byte(`{"theme": "{{ .Theme }}"}`), 0644); err != nil {
			t.Fatal(err)
		}
		themeSourcePath := filepath.Join(appConfigDir, "themes", "default.json")
		if err := os.WriteFile(themeSourcePath, []byte(`{"name": "test"}`), 0644); err != nil {
			t.Fatal(err)
		}

		setupSharedDir(t, tc.AppDir)

		oldAppConfigs := paths.Paths.App.Configs.OpenCode
		t.Cleanup(func() { paths.Paths.App.Configs.OpenCode = oldAppConfigs })
		paths.Paths.App.Configs.OpenCode = appConfigDir

		oldConfigOpenCode := paths.Paths.Config.OpenCode
		t.Cleanup(func() { paths.Paths.Config.OpenCode = oldConfigOpenCode })
		paths.Paths.Config.OpenCode = userConfigDir

		app := &OpenCode{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

		if err := app.SoftConfigure(); err != nil {
			t.Fatalf("SoftConfigure error: %v", err)
		}

		configPath := filepath.Join(userConfigDir, "opencode.json")
		if _, err := os.Stat(configPath); err != nil {
			t.Fatalf("Expected config file to be created: %v", err)
		}

		testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
	})

	t.Run("ConfigureWhenAlreadyInstalledButNotConfigured", func(t *testing.T) {
		testutil.IsolateXDGDirs(t)
		tc := testutil.SetupCompleteTest(t)
		defer tc.Cleanup()

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

		appConfigDir := filepath.Join(tc.AppDir, "configs", "opencode")
		userConfigDir := filepath.Join(tc.ConfigDir, "opencode")

		if err := os.MkdirAll(appConfigDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(appConfigDir, "themes"), 0755); err != nil {
			t.Fatal(err)
		}
		templatePath := filepath.Join(appConfigDir, "opencode.json.tmpl")
		if err := os.WriteFile(templatePath, []byte(`{"theme": "{{ .Theme }}"}`), 0644); err != nil {
			t.Fatal(err)
		}
		themeSourcePath := filepath.Join(appConfigDir, "themes", "default.json")
		if err := os.WriteFile(themeSourcePath, []byte(`{"name": "test"}`), 0644); err != nil {
			t.Fatal(err)
		}

		setupSharedDir(t, tc.AppDir)

		oldAppConfigs := paths.Paths.App.Configs.OpenCode
		t.Cleanup(func() { paths.Paths.App.Configs.OpenCode = oldAppConfigs })
		paths.Paths.App.Configs.OpenCode = appConfigDir

		oldConfigOpenCode := paths.Paths.Config.OpenCode
		t.Cleanup(func() { paths.Paths.Config.OpenCode = oldConfigOpenCode })
		paths.Paths.Config.OpenCode = userConfigDir

		app := &OpenCode{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

		if err := app.SoftConfigure(); err != nil {
			t.Fatalf("SoftConfigure error: %v", err)
		}

		configPath := filepath.Join(userConfigDir, "opencode.json")
		if _, err := os.Stat(configPath); err != nil {
			t.Fatalf("Expected config file to be created even when already installed: %v", err)
		}

		testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
	})
}

func TestExecuteCommand(t *testing.T) {
	t.Run("SuccessfulExecution", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		app := &OpenCode{Cmd: mockApp.Cmd, Base: mockApp.Base}

		mockApp.Base.SetExecCommandResult("OpenCode 1.0.0", "", nil)

		if err := app.ExecuteCommand("--version"); err != nil {
			t.Fatalf("ExecuteCommand failed: %v", err)
		}

		if mockApp.Base.GetExecCommandCallCount() != 1 {
			t.Fatalf("Expected 1 call, got %d", mockApp.Base.GetExecCommandCallCount())
		}

		lastCall := mockApp.Base.GetLastExecCommandCall()
		if lastCall.Command != constants.OpenCode {
			t.Fatalf("Expected command '%s', got %q", constants.OpenCode, lastCall.Command)
		}
		if len(lastCall.Args) != 1 || lastCall.Args[0] != "--version" {
			t.Fatalf("Expected args ['--version'], got %v", lastCall.Args)
		}
	})

	t.Run("CommandError", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		app := &OpenCode{Cmd: mockApp.Cmd, Base: mockApp.Base}

		mockApp.Base.SetExecCommandResult("", "error output", fmt.Errorf("command failed"))

		err := app.ExecuteCommand("--invalid")
		if err == nil {
			t.Fatal("Expected error from ExecuteCommand")
		}
		if !strings.Contains(err.Error(), "opencode command execution failed") {
			t.Fatalf("Expected wrapped error, got: %v", err)
		}
	})
}

func TestDefaultThemeName(t *testing.T) {
	if DEFAULT_THEME_NAME != "default" {
		t.Errorf("Expected DEFAULT_THEME_NAME to be 'default', got %q", DEFAULT_THEME_NAME)
	}
}
