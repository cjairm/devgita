package tmux_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/apps/tmux"
	"github.com/cjairm/devgita/internal/testutil"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/paths"
)

func init() {
	// Initialize logger for tests
	testutil.InitLogger()
}

func TestNew(t *testing.T) {
	t.Helper()

	app := tmux.New()
	if app == nil {
		t.Error("Expected New() to return a non-nil Tmux instance")
	}
}

func TestNameAndKind(t *testing.T) {
	app := &tmux.Tmux{}
	if app.Name() != constants.Tmux {
		t.Errorf("expected Name() %q, got %q", constants.Tmux, app.Name())
	}
	if app.Kind() != apps.KindTerminal {
		t.Errorf("expected Kind() KindTerminal, got %v", app.Kind())
	}
}

func TestInstall(t *testing.T) {
	t.Helper()

	tests := []struct {
		name        string
		shouldError bool
		installErr  error
	}{
		{
			name:        "successful installation",
			shouldError: false,
			installErr:  nil,
		},
		{
			name:        "installation failure",
			shouldError: true,
			installErr:  errors.New("installation failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			mockApp := testutil.NewMockApp()
			mockApp.Cmd.InstallError = tt.installErr

			app := &tmux.Tmux{
				Cmd: mockApp.Cmd,
			}

			err := app.Install()

			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Verify the correct package was passed
			if mockApp.Cmd.InstalledPkg != "tmux" {
				t.Errorf(
					"Expected package 'tmux', got '%s'",
					mockApp.Cmd.InstalledPkg,
				)
			}

			testutil.VerifyNoRealCommands(t, mockApp.Base)
		})
	}
}

func TestForceInstall(t *testing.T) {
	t.Helper()

	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	testutil.IsolateXDGDirs(t)

	oldHome := paths.Paths.Home.Root
	t.Cleanup(func() { paths.Paths.Home.Root = oldHome })
	paths.Paths.Home.Root = tc.ConfigDir

	app := &tmux.Tmux{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

	if err := app.ForceInstall(); err != nil {
		t.Fatalf("ForceInstall() error: %v", err)
	}
	if tc.MockApp.Cmd.InstalledPkg != constants.Tmux {
		t.Errorf("expected Install to be called, got %q", tc.MockApp.Cmd.InstalledPkg)
	}

	testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
}

func TestSoftInstall(t *testing.T) {
	t.Helper()

	tests := []struct {
		name        string
		shouldError bool
		installErr  error
	}{
		{
			name:        "successful soft installation",
			shouldError: false,
			installErr:  nil,
		},
		{
			name:        "soft installation failure",
			shouldError: true,
			installErr:  errors.New("soft installation failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			mockApp := testutil.NewMockApp()
			mockApp.Cmd.MaybeInstallError = tt.installErr

			app := &tmux.Tmux{
				Cmd: mockApp.Cmd,
			}

			err := app.SoftInstall()

			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Verify the correct package was passed
			if mockApp.Cmd.MaybeInstalled != "tmux" {
				t.Errorf(
					"Expected package 'tmux', got '%s'",
					mockApp.Cmd.MaybeInstalled,
				)
			}

			testutil.VerifyNoRealCommands(t, mockApp.Base)
		})
	}
}

func TestUninstall(t *testing.T) {
	t.Helper()

	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	testutil.IsolateXDGDirs(t)

	oldHome := paths.Paths.Home.Root
	t.Cleanup(func() { paths.Paths.Home.Root = oldHome })
	paths.Paths.Home.Root = tc.ConfigDir

	app := &tmux.Tmux{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

	if err := app.Uninstall(); err != nil {
		t.Fatalf("Uninstall error: %v", err)
	}
	if tc.MockApp.Cmd.UninstalledPkg != constants.Tmux {
		t.Errorf("expected UninstallPackage(%s), got %q", constants.Tmux, tc.MockApp.Cmd.UninstalledPkg)
	}

	testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
}

func TestUpdate(t *testing.T) {
	t.Helper()

	mockApp := testutil.NewMockApp()
	app := &tmux.Tmux{Cmd: mockApp.Cmd, Base: mockApp.Base}

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
	t.Helper()

	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	// Create source directory with tmux config
	sourceDir := filepath.Join(tc.AppDir, "tmux")
	err := os.MkdirAll(sourceDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create destination directory
	destDir := tc.ConfigDir
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create destination directory: %v", err)
	}

	testutil.IsolateXDGDirs(t)

	oldTmux := paths.Paths.App.Configs.Tmux
	t.Cleanup(func() { paths.Paths.App.Configs.Tmux = oldTmux })
	paths.Paths.App.Configs.Tmux = sourceDir

	oldHome := paths.Paths.Home.Root
	t.Cleanup(func() { paths.Paths.Home.Root = oldHome })
	paths.Paths.Home.Root = destDir

	// Create source tmux.conf file (without leading dot in source)
	sourceConfig := filepath.Join(sourceDir, "tmux.conf")
	configContent := "# Test tmux configuration\nset -g default-terminal \"screen-256color\""
	err = os.WriteFile(sourceConfig, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create source config: %v", err)
	}

	app := tmux.New()

	err = app.ForceConfigure()
	if err != nil {
		t.Errorf("ForceConfigure returned error: %v", err)
	}

	// Verify config was copied to destination
	destConfig := filepath.Join(destDir, ".tmux.conf")
	content, err := os.ReadFile(destConfig)
	if err != nil {
		t.Fatalf("Failed to read destination config: %v", err)
	}

	if string(content) != configContent {
		t.Errorf(
			"Expected config content %q, got %q",
			configContent,
			string(content),
		)
	}

	// Verify shell config was generated
	shellContent, err := os.ReadFile(tc.ZshConfigPath)
	if err != nil {
		t.Fatalf("Failed to read shell config: %v", err)
	}

	if !strings.Contains(string(shellContent), "# Tmux enabled") {
		t.Error("Expected shell config to contain Tmux feature")
	}

	testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
}

func TestSoftConfigure(t *testing.T) {
	t.Helper()

	// Test case 1: Configuration doesn't exist - should configure
	t.Run("ConfigureWhenNotExists", func(t *testing.T) {
		tc := testutil.SetupCompleteTest(t)
		defer tc.Cleanup()

		// Create source directory with tmux config
		sourceDir := filepath.Join(tc.AppDir, "tmux")
		err := os.MkdirAll(sourceDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create source directory: %v", err)
		}

		destDir := tc.ConfigDir

		testutil.IsolateXDGDirs(t)

		oldTmux := paths.Paths.App.Configs.Tmux
		t.Cleanup(func() { paths.Paths.App.Configs.Tmux = oldTmux })
		paths.Paths.App.Configs.Tmux = sourceDir

		oldHome := paths.Paths.Home.Root
		t.Cleanup(func() { paths.Paths.Home.Root = oldHome })
		paths.Paths.Home.Root = destDir

		// Create source tmux.conf file (without leading dot in source)
		sourceConfig := filepath.Join(sourceDir, "tmux.conf")
		configContent := "# Test tmux configuration\nset -g default-terminal \"screen-256color\""
		err = os.WriteFile(sourceConfig, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create source config: %v", err)
		}

		// Mock the HOME environment variable since SoftConfigure uses os.UserHomeDir()
		t.Setenv("HOME", destDir)

		app := tmux.New()

		err = app.SoftConfigure()
		if err != nil {
			t.Errorf("SoftConfigure returned error: %v", err)
		}

		// Verify config was copied to destination
		destConfig := filepath.Join(destDir, ".tmux.conf")
		content, err := os.ReadFile(destConfig)
		if err != nil {
			t.Fatalf("Failed to read destination config: %v", err)
		}

		if string(content) != configContent {
			t.Errorf(
				"Expected config content %q, got %q",
				configContent,
				string(content),
			)
		}

		// Verify shell config was generated
		shellContent, err := os.ReadFile(tc.ZshConfigPath)
		if err != nil {
			t.Fatalf("Failed to read shell config: %v", err)
		}

		if !strings.Contains(string(shellContent), "# Tmux enabled") {
			t.Error("Expected shell config to contain Tmux feature on first call")
		}

		testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
	})

	// Test case 2: Configuration already exists - should skip file copy but enable shell feature
	t.Run("SkipWhenExists", func(t *testing.T) {
		tc := testutil.SetupCompleteTest(t)
		defer tc.Cleanup()

		// Create home directory with existing .tmux.conf
		homeDir := tc.ConfigDir

		// Set Home path before creating the config file
		testutil.IsolateXDGDirs(t)

		oldHome := paths.Paths.Home.Root
		t.Cleanup(func() { paths.Paths.Home.Root = oldHome })
		paths.Paths.Home.Root = homeDir

		existingConfig := filepath.Join(homeDir, ".tmux.conf")
		existingContent := "# Existing tmux configuration\nset -g mouse on"
		err := os.WriteFile(existingConfig, []byte(existingContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create existing config: %v", err)
		}

		// Create source directory (though it shouldn't be used for file copy)
		sourceDir := filepath.Join(tc.AppDir, "tmux")
		err = os.MkdirAll(sourceDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create source directory: %v", err)
		}
		oldTmux2 := paths.Paths.App.Configs.Tmux
		t.Cleanup(func() { paths.Paths.App.Configs.Tmux = oldTmux2 })
		paths.Paths.App.Configs.Tmux = sourceDir

		// Create a different source config to prove it's not copied (without leading dot in source)
		sourceConfig := filepath.Join(sourceDir, "tmux.conf")
		sourceContent := "# Different config that should not be copied"
		err = os.WriteFile(sourceConfig, []byte(sourceContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create source config: %v", err)
		}

		// Mock the UserHomeDir for this test by temporarily setting HOME env var
		t.Setenv("HOME", homeDir)

		app := tmux.New()

		err = app.SoftConfigure()
		if err != nil {
			t.Errorf("SoftConfigure returned error: %v", err)
		}

		// Verify content wasn't changed
		contentAfter, err := os.ReadFile(existingConfig)
		if err != nil {
			t.Fatalf("Failed to read config after test: %v", err)
		}

		if string(contentAfter) != existingContent {
			t.Errorf("Expected config to remain unchanged, but it was modified")
		}

		if string(contentAfter) == sourceContent {
			t.Errorf(
				"Config was overwritten with source content when it should have been preserved",
			)
		}

		// Verify shell config was still generated (feature should be enabled even when file exists)
		shellContent, err := os.ReadFile(tc.ZshConfigPath)
		if err != nil {
			t.Fatalf("Failed to read shell config: %v", err)
		}

		if !strings.Contains(string(shellContent), "# Tmux enabled") {
			t.Error("Expected shell config to contain Tmux feature even when config file exists")
		}

		testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
	})
}

func TestExecuteCommand(t *testing.T) {
	t.Helper()

	tests := []struct {
		name        string
		args        []string
		shouldError bool
		execErr     error
	}{
		{
			name:        "successful execution",
			args:        []string{"--version"},
			shouldError: false,
			execErr:     nil,
		},
		{
			name:        "execution failure",
			args:        []string{"invalid-command"},
			shouldError: true,
			execErr:     errors.New("command failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			mockApp := testutil.NewMockApp()

			if tt.execErr != nil {
				mockApp.Base.SetExecCommandResult("", "error", tt.execErr)
			} else {
				mockApp.Base.SetExecCommandResult("tmux 3.3a", "", nil)
			}

			app := &tmux.Tmux{
				Cmd:  mockApp.Cmd,
				Base: mockApp.Base,
			}

			err := app.ExecuteCommand(tt.args...)

			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Verify the command was called
			if mockApp.Base.GetExecCommandCallCount() != 1 {
				t.Errorf("Expected 1 ExecCommand call, got %d", mockApp.Base.GetExecCommandCallCount())
			}

			lastCall := mockApp.Base.GetLastExecCommandCall()
			if lastCall == nil {
				t.Fatal("No ExecCommand call recorded")
			}
			if lastCall.Command != "tmux" {
				t.Errorf("Expected command 'tmux', got %q", lastCall.Command)
			}
		})
	}
}

func TestCreateSession(t *testing.T) {
	t.Helper()

	tests := []struct {
		name        string
		sessionName string
		workdir     string
		shouldError bool
		execErr     error
	}{
		{
			name:        "successful session creation",
			sessionName: "my-session",
			workdir:     "/path/to/project",
			shouldError: false,
			execErr:     nil,
		},
		{
			name:        "session creation failure",
			sessionName: "duplicate-session",
			workdir:     "/path/to/project",
			shouldError: true,
			execErr:     errors.New("duplicate session: duplicate-session"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			mockApp := testutil.NewMockApp()

			if tt.execErr != nil {
				mockApp.Base.SetExecCommandResult("", "error", tt.execErr)
			} else {
				mockApp.Base.SetExecCommandResult("", "", nil)
			}

			app := &tmux.Tmux{
				Cmd:  mockApp.Cmd,
				Base: mockApp.Base,
			}

			err := app.CreateSession(tt.sessionName, tt.workdir)

			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Verify correct arguments were passed
			lastCall := mockApp.Base.GetLastExecCommandCall()
			if lastCall == nil {
				t.Fatal("No ExecCommand call recorded")
			}

			expectedArgs := []string{"new-session", "-d", "-s", tt.sessionName, "-c", tt.workdir}
			if len(lastCall.Args) != len(expectedArgs) {
				t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
			}
			for i, arg := range expectedArgs {
				if lastCall.Args[i] != arg {
					t.Errorf("Expected arg[%d] to be %q, got %q", i, arg, lastCall.Args[i])
				}
			}
		})
	}
}

func TestKillSession(t *testing.T) {
	t.Helper()

	tests := []struct {
		name        string
		sessionName string
		shouldError bool
		execErr     error
	}{
		{
			name:        "successful session kill",
			sessionName: "my-session",
			shouldError: false,
			execErr:     nil,
		},
		{
			name:        "session not found",
			sessionName: "nonexistent-session",
			shouldError: true,
			execErr:     errors.New("session not found: nonexistent-session"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			mockApp := testutil.NewMockApp()

			if tt.execErr != nil {
				mockApp.Base.SetExecCommandResult("", "error", tt.execErr)
			} else {
				mockApp.Base.SetExecCommandResult("", "", nil)
			}

			app := &tmux.Tmux{
				Cmd:  mockApp.Cmd,
				Base: mockApp.Base,
			}

			err := app.KillSession(tt.sessionName)

			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Verify correct arguments were passed
			lastCall := mockApp.Base.GetLastExecCommandCall()
			if lastCall == nil {
				t.Fatal("No ExecCommand call recorded")
			}

			expectedArgs := []string{"kill-session", "-t", tt.sessionName}
			if len(lastCall.Args) != len(expectedArgs) {
				t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
			}
			for i, arg := range expectedArgs {
				if lastCall.Args[i] != arg {
					t.Errorf("Expected arg[%d] to be %q, got %q", i, arg, lastCall.Args[i])
				}
			}
		})
	}
}

func TestHasSession(t *testing.T) {
	t.Helper()

	tests := []struct {
		name        string
		sessionName string
		exists      bool
		execErr     error
	}{
		{
			name:        "session exists",
			sessionName: "my-session",
			exists:      true,
			execErr:     nil,
		},
		{
			name:        "session does not exist",
			sessionName: "nonexistent-session",
			exists:      false,
			execErr:     errors.New("session not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			mockApp := testutil.NewMockApp()

			if tt.execErr != nil {
				mockApp.Base.SetExecCommandResult("", "error", tt.execErr)
			} else {
				mockApp.Base.SetExecCommandResult("", "", nil)
			}

			app := &tmux.Tmux{
				Cmd:  mockApp.Cmd,
				Base: mockApp.Base,
			}

			result := app.HasSession(tt.sessionName)

			if result != tt.exists {
				t.Errorf("Expected HasSession to return %v, got %v", tt.exists, result)
			}

			// Verify correct arguments were passed
			lastCall := mockApp.Base.GetLastExecCommandCall()
			if lastCall == nil {
				t.Fatal("No ExecCommand call recorded")
			}

			expectedArgs := []string{"has-session", "-t", tt.sessionName}
			if len(lastCall.Args) != len(expectedArgs) {
				t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
			}
			for i, arg := range expectedArgs {
				if lastCall.Args[i] != arg {
					t.Errorf("Expected arg[%d] to be %q, got %q", i, arg, lastCall.Args[i])
				}
			}
		})
	}
}

func TestSendKeys(t *testing.T) {
	t.Helper()

	tests := []struct {
		name        string
		sessionName string
		keys        string
		shouldError bool
		execErr     error
	}{
		{
			name:        "successful send keys",
			sessionName: "my-session",
			keys:        "opencode",
			shouldError: false,
			execErr:     nil,
		},
		{
			name:        "send keys failure",
			sessionName: "nonexistent-session",
			keys:        "opencode",
			shouldError: true,
			execErr:     errors.New("session not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			mockApp := testutil.NewMockApp()

			if tt.execErr != nil {
				mockApp.Base.SetExecCommandResult("", "error", tt.execErr)
			} else {
				mockApp.Base.SetExecCommandResult("", "", nil)
			}

			app := &tmux.Tmux{
				Cmd:  mockApp.Cmd,
				Base: mockApp.Base,
			}

			err := app.SendKeys(tt.sessionName, tt.keys)

			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Verify correct arguments were passed
			lastCall := mockApp.Base.GetLastExecCommandCall()
			if lastCall == nil {
				t.Fatal("No ExecCommand call recorded")
			}

			expectedArgs := []string{"send-keys", "-t", tt.sessionName, tt.keys, "Enter"}
			if len(lastCall.Args) != len(expectedArgs) {
				t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
			}
			for i, arg := range expectedArgs {
				if lastCall.Args[i] != arg {
					t.Errorf("Expected arg[%d] to be %q, got %q", i, arg, lastCall.Args[i])
				}
			}
		})
	}
}

func TestSelectWindow(t *testing.T) {
	t.Helper()

	tests := []struct {
		name        string
		windowName  string
		shouldError bool
		execErr     error
	}{
		{
			name:        "successful window selection",
			windowName:  "wt-feature",
			shouldError: false,
			execErr:     nil,
		},
		{
			name:        "window not found",
			windowName:  "wt-nonexistent",
			shouldError: true,
			execErr:     errors.New("can't find window: wt-nonexistent"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			mockApp := testutil.NewMockApp()

			if tt.execErr != nil {
				mockApp.Base.SetExecCommandResult("", "error", tt.execErr)
			} else {
				mockApp.Base.SetExecCommandResult("", "", nil)
			}

			app := &tmux.Tmux{
				Cmd:  mockApp.Cmd,
				Base: mockApp.Base,
			}

			err := app.SelectWindow(tt.windowName)

			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Verify correct arguments were passed
			lastCall := mockApp.Base.GetLastExecCommandCall()
			if lastCall == nil {
				t.Fatal("No ExecCommand call recorded")
			}

			expectedArgs := []string{"select-window", "-t", tt.windowName}
			if len(lastCall.Args) != len(expectedArgs) {
				t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
			}
			for i, arg := range expectedArgs {
				if lastCall.Args[i] != arg {
					t.Errorf("Expected arg[%d] to be %q, got %q", i, arg, lastCall.Args[i])
				}
			}
		})
	}
}
