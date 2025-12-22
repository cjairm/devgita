package devgita

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/apps/git"
	"github.com/cjairm/devgita/internal/testutil"
	"github.com/cjairm/devgita/pkg/paths"
)

func init() {
	testutil.InitLogger()
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

	mockApp := testutil.NewMockApp()
	gitInstance := &git.Git{
		Cmd:  mockApp.Cmd,
		Base: mockApp.Base,
	}
	dg := &Devgita{
		Git:  *gitInstance,
		Base: mockApp.Base,
	}

	// Mock successful git clone
	mockApp.Base.SetExecCommandResult("Cloning into...", "", nil)

	err := dg.SoftInstall()
	if err != nil {
		t.Fatalf("SoftInstall() failed: %v", err)
	}

	// Verify git clone was called
	if mockApp.Base.GetExecCommandCallCount() != 1 {
		t.Fatalf("Expected 1 git command call, got %d", mockApp.Base.GetExecCommandCallCount())
	}

	lastCall := mockApp.Base.GetLastExecCommandCall()
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

	mockApp := testutil.NewMockApp()
	gitInstance := &git.Git{
		Cmd:  mockApp.Cmd,
		Base: mockApp.Base,
	}
	dg := &Devgita{
		Git:  *gitInstance,
		Base: mockApp.Base,
	}

	err := dg.SoftInstall()
	if err != nil {
		t.Fatalf("SoftInstall() failed: %v", err)
	}

	// Verify git clone was NOT called (directory already exists)
	if mockApp.Base.GetExecCommandCallCount() != 0 {
		t.Fatalf(
			"Expected 0 git command calls (directory exists), got %d",
			mockApp.Base.GetExecCommandCallCount(),
		)
	}

	// Verify file was preserved
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

	// Override paths
	oldAppDir := paths.Paths.App.Root
	oldConfigDir := paths.Paths.Config.Root
	oldConfigDirPath := configDirPath
	oldGlobalConfigPath := globalConfigPath
	oldZshConfigPath := zshConfigPath

	configDir := filepath.Join(tempDir, "config")
	devgitaConfigDir := filepath.Join(configDir, "devgita")

	paths.Paths.App.Root = emptyDir
	paths.Paths.Config.Root = configDir
	configDirPath = devgitaConfigDir
	globalConfigPath = filepath.Join(devgitaConfigDir, "global_config.yaml")
	zshConfigPath = filepath.Join(devgitaConfigDir, "devgita.zsh")

	t.Cleanup(func() {
		paths.Paths.App.Root = oldAppDir
		paths.Paths.Config.Root = oldConfigDir
		configDirPath = oldConfigDirPath
		globalConfigPath = oldGlobalConfigPath
		zshConfigPath = oldZshConfigPath
	})

	mockApp := testutil.NewMockApp()
	gitInstance := &git.Git{
		Cmd:  mockApp.Cmd,
		Base: mockApp.Base,
	}
	dg := &Devgita{
		Git:  *gitInstance,
		Base: mockApp.Base,
	}

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

	mockApp := testutil.NewMockApp()
	gitInstance := &git.Git{
		Cmd:  mockApp.Cmd,
		Base: mockApp.Base,
	}
	dg := &Devgita{
		Git:  *gitInstance,
		Base: mockApp.Base,
	}

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

	// Override paths
	oldAppDir := paths.Paths.App.Root
	oldConfigDir := paths.Paths.Config.Root
	oldConfigDirPath := configDirPath
	oldGlobalConfigPath := globalConfigPath
	oldZshConfigPath := zshConfigPath

	configDir := filepath.Join(tempDir, "config")
	devgitaConfigDir := filepath.Join(configDir, "devgita")

	paths.Paths.App.Root = appDir
	paths.Paths.Config.Root = configDir
	configDirPath = devgitaConfigDir
	globalConfigPath = filepath.Join(devgitaConfigDir, "global_config.yaml")
	zshConfigPath = filepath.Join(devgitaConfigDir, "devgita.zsh")

	t.Cleanup(func() {
		paths.Paths.App.Root = oldAppDir
		paths.Paths.Config.Root = oldConfigDir
		configDirPath = oldConfigDirPath
		globalConfigPath = oldGlobalConfigPath
		zshConfigPath = oldZshConfigPath
	})

	mockApp := testutil.NewMockApp()
	gitInstance := &git.Git{
		Cmd:  mockApp.Cmd,
		Base: mockApp.Base,
	}
	dg := &Devgita{
		Git:  *gitInstance,
		Base: mockApp.Base,
	}

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

	// Override paths
	oldAppDir := paths.Paths.App.Root
	oldConfigDir := paths.Paths.Config.Root
	oldConfigDirPath := configDirPath
	oldGlobalConfigPath := globalConfigPath
	oldZshConfigPath := zshConfigPath

	configDir := filepath.Join(tempDir, "config")
	devgitaConfigDir := filepath.Join(configDir, "devgita")

	paths.Paths.App.Root = tempDir
	paths.Paths.Config.Root = configDir
	configDirPath = devgitaConfigDir
	globalConfigPath = filepath.Join(devgitaConfigDir, "global_config.yaml")
	zshConfigPath = filepath.Join(devgitaConfigDir, "devgita.zsh")

	t.Cleanup(func() {
		paths.Paths.App.Root = oldAppDir
		paths.Paths.Config.Root = oldConfigDir
		configDirPath = oldConfigDirPath
		globalConfigPath = oldGlobalConfigPath
		zshConfigPath = oldZshConfigPath
	})

	mockApp := testutil.NewMockApp()
	gitInstance := &git.Git{
		Cmd:  mockApp.Cmd,
		Base: mockApp.Base,
	}
	dg := &Devgita{
		Git:  *gitInstance,
		Base: mockApp.Base,
	}

	// Mock successful git clone
	mockApp.Base.SetExecCommandResult("Cloning into...", "", nil)

	err := dg.ForceInstall()
	if err != nil {
		t.Fatalf("ForceInstall() failed: %v", err)
	}

	// Verify old file was removed by Uninstall
	if _, err := os.Stat(existingFile); !os.IsNotExist(err) {
		t.Fatal("Expected existing file to be removed by Uninstall")
	}

	// Verify git clone was called
	if mockApp.Base.GetExecCommandCallCount() != 1 {
		t.Fatalf("Expected 1 git clone call, got %d", mockApp.Base.GetExecCommandCallCount())
	}
}

func TestForceConfigure_CreatesConfig(t *testing.T) {
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	mockApp := tc.MockApp
	gitInstance := &git.Git{
		Cmd:  mockApp.Cmd,
		Base: mockApp.Base,
	}
	dg := &Devgita{
		Git:  *gitInstance,
		Base: mockApp.Base,
	}

	err := dg.ForceConfigure()
	if err != nil {
		t.Fatalf("ForceConfigure() failed: %v", err)
	}

	// Verify config file was created
	if _, err := os.Stat(tc.ConfigPath); os.IsNotExist(err) {
		t.Fatal("Expected config file to be created")
	}

	// Verify config content includes AppPath and ConfigPath
	content, err := os.ReadFile(tc.ConfigPath)
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
	if !strings.Contains(contentStr, tc.AppDir) {
		t.Fatalf("Expected config to contain AppDir path %q", tc.AppDir)
	}

	// Verify shell config file (devgita.zsh) was created by RegenerateShellConfig
	if _, err := os.Stat(tc.ZshConfigPath); os.IsNotExist(err) {
		t.Fatalf("Expected shell config file (devgita.zsh) to be created at %s", tc.ZshConfigPath)
	}

	// Verify zsh config content exists
	zshContent, err := os.ReadFile(tc.ZshConfigPath)
	if err != nil {
		t.Fatalf("Failed to read zsh config file: %v", err)
	}
	if len(zshContent) == 0 {
		t.Fatal("Expected zsh config file to have content")
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
	testutil.VerifyNoRealConfigChanges(t)
}

func TestForceConfigure_OverwritesExisting(t *testing.T) {
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	// Create existing config with old values
	oldContent := "app_path: /old/path\nconfig_path: /old/config\n"
	if err := os.WriteFile(tc.ConfigPath, []byte(oldContent), 0o644); err != nil {
		t.Fatalf("Failed to create existing config: %v", err)
	}

	mockApp := tc.MockApp
	gitInstance := &git.Git{
		Cmd:  mockApp.Cmd,
		Base: mockApp.Base,
	}
	dg := &Devgita{
		Git:  *gitInstance,
		Base: mockApp.Base,
	}

	err := dg.ForceConfigure()
	if err != nil {
		t.Fatalf("ForceConfigure() failed: %v", err)
	}

	content, err := os.ReadFile(tc.ConfigPath)
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
	if !strings.Contains(contentStr, tc.AppDir) {
		t.Fatalf("Expected config to contain new AppDir path %q", tc.AppDir)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
	testutil.VerifyNoRealConfigChanges(t)
}

func TestSoftConfigure_PreservesExistingConfig(t *testing.T) {
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	// Compute the actual zsh config path used by devgita
	// devgita uses configDirPath (Config.Root/devgita) not appDir
	devgitaConfigDir := filepath.Join(tc.ConfigDir, "devgita")
	actualZshConfigPath := filepath.Join(devgitaConfigDir, "devgita.zsh")

	// Override package-level variables that are computed at init time
	oldConfigDirPath := configDirPath
	oldGlobalConfigPath := globalConfigPath
	oldZshConfigPath := zshConfigPath

	configDirPath = devgitaConfigDir
	globalConfigPath = tc.ConfigPath
	zshConfigPath = actualZshConfigPath

	t.Cleanup(func() {
		configDirPath = oldConfigDirPath
		globalConfigPath = oldGlobalConfigPath
		zshConfigPath = oldZshConfigPath
	})

	// Create existing config with custom values
	customContent := "app_path: /custom/path\nconfig_path: /custom/config\ncustom_field: custom_value\n"
	if err := os.WriteFile(tc.ConfigPath, []byte(customContent), 0o644); err != nil {
		t.Fatalf("Failed to create existing config: %v", err)
	}

	// Create zsh config file so SoftConfigure checks pass
	// Must create in the same location that devgita.go expects
	if err := os.MkdirAll(devgitaConfigDir, 0o755); err != nil {
		t.Fatalf("Failed to create devgita config dir: %v", err)
	}
	if err := os.WriteFile(actualZshConfigPath, []byte("# zsh config"), 0o644); err != nil {
		t.Fatalf("Failed to create zsh config: %v", err)
	}

	mockApp := tc.MockApp
	gitInstance := &git.Git{
		Cmd:  mockApp.Cmd,
		Base: mockApp.Base,
	}
	dg := &Devgita{
		Git:  *gitInstance,
		Base: mockApp.Base,
	}

	err := dg.SoftConfigure()
	if err != nil {
		t.Fatalf("SoftConfigure() failed: %v", err)
	}

	content, err := os.ReadFile(tc.ConfigPath)
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

	testutil.VerifyNoRealCommands(t, mockApp.Base)
	testutil.VerifyNoRealConfigChanges(t)
}
