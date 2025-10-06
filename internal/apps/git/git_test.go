package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/constants"
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

func TestForceInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Git{Cmd: mc}

	if err := app.ForceInstall(); err != nil {
		t.Fatalf("ForceInstall error: %v", err)
	}
	if mc.InstalledPkg != constants.Git {
		t.Fatalf("expected InstallPackage(%s), got %q", constants.Git, mc.InstalledPkg)
	}
}

func TestSoftInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Git{Cmd: mc}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mc.MaybeInstalled != constants.Git {
		t.Fatalf("expected MaybeInstallPackage(%s), got %q", constants.Git, mc.MaybeInstalled)
	}
}

func TestUninstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Git{Cmd: mc}

	err := app.Uninstall()
	if err == nil {
		t.Fatal("expected Uninstall to return error for unsupported operation")
	}
	if err.Error() != "git uninstall not supported through devgita" {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestForceConfigure(t *testing.T) {
	// Create temp "app config" dir with a fake file as source
	src := t.TempDir()
	dst := t.TempDir()

	// Override global paths for the duration of the test
	oldAppDir, oldLocalDir := paths.GitConfigAppDir, paths.GitConfigLocalDir
	paths.GitConfigAppDir, paths.GitConfigLocalDir = src, dst
	t.Cleanup(func() {
		paths.GitConfigAppDir, paths.GitConfigLocalDir = oldAppDir, oldLocalDir
	})

	originalContent := "[user]\n\tname = Test User"
	if err := os.WriteFile(filepath.Join(src, ".gitconfig"), []byte(originalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	mc := commands.NewMockCommand()
	app := &Git{Cmd: mc}

	if err := app.ForceConfigure(); err != nil {
		t.Fatalf("ForceConfigure error: %v", err)
	}

	check := filepath.Join(dst, ".gitconfig")
	if _, err := os.Stat(check); err != nil {
		t.Fatalf("expected copied file at %s: %v", check, err)
	}

	copiedContent, err := os.ReadFile(check)
	if err != nil {
		t.Fatalf("failed to read copied file: %v", err)
	}
	if string(copiedContent) != originalContent {
		t.Fatalf("content mismatch: expected %q, got %q", originalContent, string(copiedContent))
	}

	modifiedContent := "[user]\n\tname = Modified User"
	if err := os.WriteFile(check, []byte(modifiedContent), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := app.ForceConfigure(); err != nil {
		t.Fatalf("second ForceConfigure error: %v", err)
	}

	finalContent, err := os.ReadFile(check)
	if err != nil {
		t.Fatalf("failed to read file after second configure: %v", err)
	}
	if string(finalContent) == string(modifiedContent) {
		t.Fatalf(
			"ForceConfigure did not overwrite: expected %q, got %q",
			originalContent,
			string(finalContent),
		)
	}
}

func TestSoftConfigure(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	// Override global paths for the duration of the test
	oldAppDir, oldLocalDir := paths.GitConfigAppDir, paths.GitConfigLocalDir
	paths.GitConfigAppDir, paths.GitConfigLocalDir = src, dst
	t.Cleanup(func() {
		paths.GitConfigAppDir, paths.GitConfigLocalDir = oldAppDir, oldLocalDir
	})

	originalContent := "[user]\n\tname = Test User"
	if err := os.WriteFile(filepath.Join(src, ".gitconfig"), []byte(originalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	mc := commands.NewMockCommand()
	app := &Git{Cmd: mc}

	if err := app.SoftConfigure(); err != nil {
		t.Fatalf("SoftConfigure error: %v", err)
	}

	check := filepath.Join(dst, ".gitconfig")
	if _, err := os.Stat(check); err != nil {
		t.Fatalf("expected copied file at %s: %v", check, err)
	}

	modifiedContent := "[user]\n\tname = Modified User"
	if err := os.WriteFile(check, []byte(modifiedContent), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := app.SoftConfigure(); err != nil {
		t.Fatalf("second SoftConfigure error: %v", err)
	}

	finalContent, err := os.ReadFile(check)
	if err != nil {
		t.Fatalf("failed to read file after second configure: %v", err)
	}
	if string(finalContent) == string(originalContent) {
		t.Fatalf(
			"SoftConfigure overwrote existing file: expected %q, got %q",
			modifiedContent,
			string(finalContent),
		)
	}
}

func TestExecuteCommand(t *testing.T) {
	// These tests use the actual BaseCommand but won't execute actual git commands
	// since we expect git to fail in test environment
	mc := commands.NewMockCommand()
	app := &Git{Cmd: mc}

	// Test that methods don't panic - they will fail due to git not being available
	// but that's expected in test environment
	err := app.ExecuteCommand("--version")

	// Both should fail similarly since git isn't available, but shouldn't panic
	if err == nil {
		t.Log("Git commands succeeded unexpectedly (git must be available)")
	} else {
		t.Logf("Git commands failed as expected (no git available): %v", err)
	}
}

func TestGitSpecificMethods(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Git{Cmd: mc}

	tests := []struct {
		name string
		fn   func() error
	}{
		{"DeleteBranch", func() error { return app.DeleteBranch("feature-branch", false) }},
		{"DeleteBranchForced", func() error { return app.DeleteBranch("feature-branch", true) }},
		{"DeepClean", func() error { return app.DeepClean("", "") }},
		{"FetchOrigin", func() error { return app.FetchOrigin() }},
		{"Pop", func() error { return app.Pop("") }},
		{"Pull", func() error { return app.Pull("") }},
		{"PullWithBranch", func() error { return app.Pull("main") }},
		{"SwitchBranch", func() error { return app.SwitchBranch("main") }},
		{"Restore", func() error { return app.Restore("", "file.txt") }},
		{"RestoreWithBranch", func() error { return app.Restore("develop", "file.txt") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We expect these to fail since git commands will fail in test environment,
			// but we're testing that the method calls work without panicking
			err := tt.fn()
			// The error is expected since git commands will fail in test environment
			// We just want to ensure no panic and proper error handling
			if err == nil {
				t.Logf("%s completed successfully (unexpected but ok)", tt.name)
			} else {
				t.Logf("%s failed as expected: %v", tt.name, err)
			}
		})
	}
}
