package alacritty

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/testutil"
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
	app := &Alacritty{Cmd: mc}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mc.InstalledDesktopApp != constants.Alacritty {
		t.Fatalf(
			"expected InstallDesktopApp(%s), got %q",
			constants.Alacritty,
			mc.InstalledDesktopApp,
		)
	}
}

// SKIP: ForceInstall test as per guidelines
// func TestForceInstall(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &Alacritty{Cmd: mc}
//
// 	if err := app.ForceInstall(); err != nil {
// 		t.Fatalf("ForceInstall error: %v", err)
// 	}
// 	// ForceInstall should call Install() which uses InstallDesktopApp
// 	if mc.InstalledDesktopApp != constants.Alacritty {
// 		t.Fatalf("expected InstallDesktopApp(%s), got %q", constants.Alacritty, mc.InstalledDesktopApp)
// 	}
// }

func TestSoftInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Alacritty{Cmd: mc}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mc.MaybeInstalledDesktop != constants.Alacritty {
		t.Fatalf(
			"expected MaybeInstallDesktopApp(%s), got %q",
			constants.Alacritty,
			mc.MaybeInstalledDesktop,
		)
	}
}

// SKIP: Uninstall test as per guidelines
// func TestUninstall(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &Alacritty{Cmd: mc}
//
// 	err := app.Uninstall()
// 	if err == nil {
// 		t.Fatal("expected Uninstall to return error for unsupported operation")
// 	}
// 	if err.Error() != "uninstall not implemented for alacritty" {
// 		t.Fatalf("unexpected error message: %v", err)
// 	}
// }

func TestForceConfigure(t *testing.T) {
	// Setup isolated test environment with templates and config
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	// Create alacritty template
	tmplDir := filepath.Join(tc.AppDir, "alacritty")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatal(err)
	}

	tmplContent := `[env]
TERM = "xterm-256color"

[window]
opacity = 0.8

{{if eq .Font "default"}}
[font]
size = 13
{{end}}

{{if eq .Theme "default"}}
[colors.primary]
background = "0x282828"
{{end}}
`
	tmplPath := filepath.Join(tmplDir, "alacritty.toml.tmpl")
	if err := os.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create starter.sh
	starterContent := "#!/bin/bash\nzsh"
	starterPath := filepath.Join(tmplDir, "starter.sh")
	if err := os.WriteFile(starterPath, []byte(starterContent), 0755); err != nil {
		t.Fatal(err)
	}

	// Create destination directory
	destDir := filepath.Join(tc.ConfigDir, "alacritty")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Override paths
	oldAppConfig := paths.Paths.App.Configs.Alacritty
	oldLocalConfig := paths.Paths.Config.Alacritty
	oldConfigRoot := paths.Paths.Config.Root

	paths.Paths.App.Configs.Alacritty = tmplDir
	paths.Paths.Config.Alacritty = destDir
	paths.Paths.Config.Root = tc.ConfigDir

	t.Cleanup(func() {
		paths.Paths.App.Configs.Alacritty = oldAppConfig
		paths.Paths.Config.Alacritty = oldLocalConfig
		paths.Paths.Config.Root = oldConfigRoot
	})

	app := &Alacritty{Cmd: tc.MockApp.Cmd}

	if err := app.ForceConfigure(ConfigureOptions{}); err != nil {
		t.Fatalf("ForceConfigure error: %v", err)
	}

	// Check that config file was generated
	configPath := filepath.Join(destDir, "alacritty.toml")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("expected generated file at %s: %v", configPath, err)
	}

	// Check that starter.sh was copied
	starterDest := filepath.Join(destDir, "starter.sh")
	if _, err := os.Stat(starterDest); err != nil {
		t.Fatalf("expected copied file at %s: %v", starterDest, err)
	}

	// Verify generated content contains expected sections
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read generated config: %v", err)
	}
	configStr := string(content)
	if !strings.Contains(configStr, "TERM") {
		t.Error("expected generated config to contain TERM setting")
	}
	if !strings.Contains(configStr, "opacity = 0.8") {
		t.Error("expected generated config to contain opacity setting")
	}
}

func TestSoftConfigure(t *testing.T) {
	// Setup isolated test environment with templates and config
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	// Create alacritty template
	tmplDir := filepath.Join(tc.AppDir, "alacritty")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatal(err)
	}

	tmplContent := `[env]
TERM = "xterm-256color"

[window]
opacity = 0.9

{{if eq .Font "default"}}
[font]
size = 14
{{end}}

{{if eq .Theme "default"}}
[colors.primary]
background = "0x1e1e1e"
{{end}}
`
	tmplPath := filepath.Join(tmplDir, "alacritty.toml.tmpl")
	if err := os.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create starter.sh
	starterContent := "#!/bin/bash\nzsh"
	starterPath := filepath.Join(tmplDir, "starter.sh")
	if err := os.WriteFile(starterPath, []byte(starterContent), 0755); err != nil {
		t.Fatal(err)
	}

	// Create destination directory
	destDir := filepath.Join(tc.ConfigDir, "alacritty")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Override paths
	oldAppConfig := paths.Paths.App.Configs.Alacritty
	oldLocalConfig := paths.Paths.Config.Alacritty
	oldConfigRoot := paths.Paths.Config.Root

	paths.Paths.App.Configs.Alacritty = tmplDir
	paths.Paths.Config.Alacritty = destDir
	paths.Paths.Config.Root = tc.ConfigDir

	t.Cleanup(func() {
		paths.Paths.App.Configs.Alacritty = oldAppConfig
		paths.Paths.Config.Alacritty = oldLocalConfig
		paths.Paths.Config.Root = oldConfigRoot
	})

	app := &Alacritty{Cmd: tc.MockApp.Cmd}

	// Test 1: Should configure when no existing config
	if err := app.SoftConfigure(ConfigureOptions{}); err != nil {
		t.Fatalf("SoftConfigure error: %v", err)
	}

	// Check that config file was generated
	configPath := filepath.Join(destDir, "alacritty.toml")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("expected generated file at %s: %v", configPath, err)
	}

	// Check that starter.sh was copied
	starterDest := filepath.Join(destDir, "starter.sh")
	if _, err := os.Stat(starterDest); err != nil {
		t.Fatalf("expected copied file at %s: %v", starterDest, err)
	}

	// Verify generated content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read generated config: %v", err)
	}
	configStr := string(content)
	if !strings.Contains(configStr, "opacity = 0.9") {
		t.Error("expected generated config to contain opacity setting")
	}

	// Test 2: Should not overwrite when config already exists
	// Modify the existing config
	modifiedContent := `[window]
opacity = 0.5
option_as_alt = "none"
`
	if err := os.WriteFile(configPath, []byte(modifiedContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Run SoftConfigure again
	if err := app.SoftConfigure(ConfigureOptions{}); err != nil {
		t.Fatalf("second SoftConfigure error: %v", err)
	}

	// Check that the file was not overwritten
	finalContent, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}
	if string(finalContent) != modifiedContent {
		t.Fatalf(
			"SoftConfigure should not overwrite existing config: expected %q, got %q",
			modifiedContent,
			string(finalContent),
		)
	}
}
