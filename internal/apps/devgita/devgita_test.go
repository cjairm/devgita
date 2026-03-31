package devgita

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/embedded"
	"github.com/cjairm/devgita/internal/testutil"
	"github.com/cjairm/devgita/pkg/paths"
)

func init() {
	testutil.InitLogger()
}

// mockExtractor creates a test configs directory with sample files
func mockExtractor(destDir string) error {
	// Create sample config directories
	dirs := []string{
		filepath.Join(destDir, "git"),
		filepath.Join(destDir, "neovim"),
		filepath.Join(destDir, "tmux"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Create sample config files
	files := map[string]string{
		filepath.Join(destDir, "git", ".gitconfig"):  "[user]\n\tname = Test\n",
		filepath.Join(destDir, "neovim", "init.lua"): "-- Test config\n",
		filepath.Join(destDir, "tmux", "tmux.conf"):  "# Test config\n",
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

func TestNew(t *testing.T) {
	// Set temporary default extractor for test
	oldExtractor := embedded.DefaultExtractor
	embedded.DefaultExtractor = mockExtractor
	t.Cleanup(func() {
		embedded.DefaultExtractor = oldExtractor
	})

	dg := New()
	if dg == nil {
		t.Fatal("New() returned nil")
	}
	if dg.ExtractEmbedded == nil {
		t.Fatal("ExtractEmbedded function is nil")
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
	dg := &Devgita{
		Base:            mockApp.Base,
		ExtractEmbedded: mockExtractor,
	}

	err := dg.SoftInstall()
	if err != nil {
		t.Fatalf("SoftInstall() failed: %v", err)
	}

	// Verify configs directory was created
	configsDir := filepath.Join(appDir, "configs")
	if _, err := os.Stat(configsDir); os.IsNotExist(err) {
		t.Fatal("Configs directory was not created")
	}

	// Verify sample files exist
	gitConfig := filepath.Join(configsDir, "git", ".gitconfig")
	if _, err := os.Stat(gitConfig); os.IsNotExist(err) {
		t.Fatal("Git config was not extracted")
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestSoftInstall_DirectoryExistsWithFiles(t *testing.T) {
	tempDir := t.TempDir()
	appDir := filepath.Join(tempDir, "devgita")
	configsDir := filepath.Join(appDir, "configs")

	// Create existing configs directory with a file
	if err := os.MkdirAll(configsDir, 0755); err != nil {
		t.Fatalf("Failed to create configs dir: %v", err)
	}
	testFile := filepath.Join(configsDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	oldAppDir := paths.Paths.App.Root
	paths.Paths.App.Root = appDir
	t.Cleanup(func() {
		paths.Paths.App.Root = oldAppDir
	})

	mockApp := testutil.NewMockApp()
	dg := &Devgita{
		Base:            mockApp.Base,
		ExtractEmbedded: mockExtractor,
	}

	err := dg.SoftInstall()
	if err != nil {
		t.Fatalf("SoftInstall() failed: %v", err)
	}

	// Verify existing file was preserved (extraction not called)
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("Expected existing file to be preserved")
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestUninstall_RemovesConfigsDirectory(t *testing.T) {
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	// Create configs directory in app root
	configsDir := filepath.Join(paths.Paths.App.Root, "configs")
	if err := os.MkdirAll(configsDir, 0755); err != nil {
		t.Fatalf("Failed to create configs dir: %v", err)
	}
	testFile := filepath.Join(configsDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	dg := &Devgita{
		Base:            tc.MockApp.Base,
		ExtractEmbedded: mockExtractor,
	}

	err := dg.Uninstall()
	if err != nil {
		t.Fatalf("Uninstall() failed: %v", err)
	}

	// Verify configs directory was removed
	if _, err := os.Stat(configsDir); !os.IsNotExist(err) {
		t.Fatal("Expected configs directory to be removed")
	}

	testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
}

func TestForceInstall(t *testing.T) {
	tempDir := t.TempDir()
	appDir := filepath.Join(tempDir, "devgita")
	configsDir := filepath.Join(appDir, "configs")

	// Create existing configs directory
	if err := os.MkdirAll(configsDir, 0755); err != nil {
		t.Fatalf("Failed to create configs dir: %v", err)
	}
	oldFile := filepath.Join(configsDir, "old.txt")
	if err := os.WriteFile(oldFile, []byte("old"), 0644); err != nil {
		t.Fatalf("Failed to create old file: %v", err)
	}

	oldAppDir := paths.Paths.App.Root
	paths.Paths.App.Root = appDir
	t.Cleanup(func() {
		paths.Paths.App.Root = oldAppDir
	})

	mockApp := testutil.NewMockApp()
	dg := &Devgita{
		Base:            mockApp.Base,
		ExtractEmbedded: mockExtractor,
	}

	err := dg.ForceInstall()
	if err != nil {
		t.Fatalf("ForceInstall() failed: %v", err)
	}

	// Verify old file was removed and new configs extracted
	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Fatal("Expected old file to be removed")
	}

	gitConfig := filepath.Join(configsDir, "git", ".gitconfig")
	if _, err := os.Stat(gitConfig); os.IsNotExist(err) {
		t.Fatal("Expected new configs to be extracted")
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestForceConfigure_CreatesConfig(t *testing.T) {
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	dg := &Devgita{
		Base:            tc.MockApp.Base,
		ExtractEmbedded: mockExtractor,
	}

	err := dg.ForceConfigure()
	if err != nil {
		t.Fatalf("ForceConfigure() failed: %v", err)
	}

	// Verify global config was created
	if _, err := os.Stat(tc.ConfigPath); os.IsNotExist(err) {
		t.Fatal("Global config was not created")
	}

	// Verify IsMac field was set
	content, err := os.ReadFile(tc.ConfigPath)
	if err != nil {
		t.Fatalf("Failed to read global config: %v", err)
	}

	// Check for is_mac field in YAML
	if !strings.Contains(string(content), "is_mac:") {
		t.Error("Expected is_mac field in global config")
	}

	testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
}

func TestForceConfigure_OverwritesExisting(t *testing.T) {
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	// Create existing config file
	oldContent := "old: config\n"
	if err := os.WriteFile(tc.ConfigPath, []byte(oldContent), 0644); err != nil {
		t.Fatalf("Failed to create old config: %v", err)
	}

	dg := &Devgita{
		Base:            tc.MockApp.Base,
		ExtractEmbedded: mockExtractor,
	}

	err := dg.ForceConfigure()
	if err != nil {
		t.Fatalf("ForceConfigure() failed: %v", err)
	}

	// Verify config was overwritten
	content, err := os.ReadFile(tc.ConfigPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	if strings.Contains(string(content), "old: config") {
		t.Error("Expected old config to be overwritten")
	}

	testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
}

func TestSoftConfigure_PreservesExistingConfig(t *testing.T) {
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	// Create existing config file with custom marker
	existingContent := "# Custom config marker\napp_path: /custom/path\n"
	if err := os.WriteFile(tc.ConfigPath, []byte(existingContent), 0644); err != nil {
		t.Fatalf("Failed to create existing config: %v", err)
	}

	// Create existing zsh config at the correct location
	// The devgita app expects it at config/devgita/.devgita.zsh, not app/devgita.zsh
	actualZshPath := getZshConfigPath()
	existingZsh := "# Custom zsh marker\nsource /custom/path\n"
	if err := os.WriteFile(actualZshPath, []byte(existingZsh), 0644); err != nil {
		t.Fatalf("Failed to create existing zsh config: %v", err)
	}

	// Record modification time
	configInfo, err := os.Stat(tc.ConfigPath)
	if err != nil {
		t.Fatalf("Failed to stat config: %v", err)
	}
	originalModTime := configInfo.ModTime()

	dg := &Devgita{
		Base:            tc.MockApp.Base,
		ExtractEmbedded: mockExtractor,
	}

	err = dg.SoftConfigure()
	if err != nil {
		t.Fatalf("SoftConfigure() failed: %v", err)
	}

	// Verify config was not regenerated (modification time should be same)
	newConfigInfo, err := os.Stat(tc.ConfigPath)
	if err != nil {
		t.Fatalf("Failed to stat config after soft configure: %v", err)
	}

	if !newConfigInfo.ModTime().Equal(originalModTime) {
		t.Error("Expected config file to be preserved (not regenerated)")
	}

	// Verify custom marker still exists
	content, err := os.ReadFile(tc.ConfigPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	if !strings.Contains(string(content), "Custom config marker") {
		t.Error("Expected custom config marker to be preserved")
	}

	testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
}
