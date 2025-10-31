package fonts

import (
	"testing"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/logger"
)

func init() {
	// Initialize logger for tests
	logger.Init(false)
}

func TestNew(t *testing.T) {
	fonts := New()
	if fonts == nil {
		t.Fatal("New() returned nil")
	}
	if fonts.Cmd == nil {
		t.Fatal("New() did not initialize command")
	}
}

func TestInstallFont(t *testing.T) {
	mockCmd := cmd.NewMockCommand()
	fonts := &Fonts{Cmd: mockCmd}

	fontName := "font-hack-nerd-font"
	err := fonts.Install(fontName)
	if err != nil {
		t.Fatalf("InstallFont() failed: %v", err)
	}

	if mockCmd.InstalledDesktopApp != fontName {
		t.Errorf("Expected font %s, got %s", fontName, mockCmd.InstalledDesktopApp)
	}
}

func TestSoftInstallFont(t *testing.T) {
	mockCmd := cmd.NewMockCommand()
	fonts := &Fonts{Cmd: mockCmd}

	fontName := "font-meslo-lg-nerd-font"
	err := fonts.SoftInstall(fontName)
	if err != nil {
		t.Fatalf("SoftInstallFont() failed: %v", err)
	}

	if mockCmd.FontName != fontName {
		t.Errorf("Expected font %s, got %s", fontName, mockCmd.FontName)
	}
	if mockCmd.FontURL != "" {
		t.Errorf("Expected empty URL, got %s", mockCmd.FontURL)
	}
}

func TestAvailable(t *testing.T) {
	mockCmd := cmd.NewMockCommand()
	fonts := &Fonts{Cmd: mockCmd}

	available := fonts.Available()
	expectedFonts := []string{
		"font-hack-nerd-font",
		"font-meslo-lg-nerd-font",
		"font-caskaydia-mono-nerd-font",
		"font-fira-mono",
		"font-jetbrains-mono-nerd-font",
	}

	if len(available) != len(expectedFonts) {
		t.Fatalf("Expected %d fonts, got %d", len(expectedFonts), len(available))
	}

	for i, expected := range expectedFonts {
		if available[i] != expected {
			t.Errorf("Expected font %s at index %d, got %s", expected, i, available[i])
		}
	}
}
