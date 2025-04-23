package commands_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/cjairm/devgita/internal/commands"
)

type FakePlatform struct {
	Linux bool
	Mac   bool
}

func (f FakePlatform) IsLinux() bool { return f.Linux }
func (f FakePlatform) IsMac() bool   { return f.Mac }

func createFile(t *testing.T, dir, name string) {
	t.Helper()
	err := os.WriteFile(filepath.Join(dir, name), []byte{}, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file %q: %v", name, err)
	}
}

func fakeCmdWithOutput(output string) *exec.Cmd {
	return exec.Command("bash", "-c", "echo -e \""+output+"\"")
}

func TestIsDesktopAppPresent(t *testing.T) {
	t.Run("Linux with matching .desktop file", func(t *testing.T) {
		tmpDir := t.TempDir()
		createFile(t, tmpDir, "myapp.desktop")

		b := commands.NewBaseCommandCustom(FakePlatform{Linux: true})
		found, err := b.IsDesktopAppPresent(tmpDir, "myapp")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !found {
			t.Errorf("Expected to find the desktop app, but did not")
		}
	})

	t.Run("Mac with matching file", func(t *testing.T) {
		tmpDir := t.TempDir()
		createFile(t, tmpDir, "myapp")

		b := commands.NewBaseCommandCustom(FakePlatform{Mac: true})
		found, err := b.IsDesktopAppPresent(tmpDir, "myapp")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !found {
			t.Errorf("Expected to find the app file, but did not")
		}
	})

	t.Run("Linux: no match with wrong extension", func(t *testing.T) {
		tmpDir := t.TempDir()
		createFile(t, tmpDir, "myapp.txt")

		b := commands.NewBaseCommandCustom(FakePlatform{Linux: true})
		found, err := b.IsDesktopAppPresent(tmpDir, "myapp")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if found {
			t.Errorf("Expected not to find a desktop app, but did")
		}
	})

	t.Run("Mac: partial match not found", func(t *testing.T) {
		tmpDir := t.TempDir()
		createFile(t, tmpDir, "unrelated")

		b := commands.NewBaseCommandCustom(FakePlatform{Mac: true})
		found, err := b.IsDesktopAppPresent(tmpDir, "myapp")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if found {
			t.Errorf("Expected not to find the app file, but did")
		}
	})

	t.Run("Linux: case-insensitive match with .desktop", func(t *testing.T) {
		tmpDir := t.TempDir()
		createFile(t, tmpDir, "MyApp.Desktop")

		b := commands.NewBaseCommandCustom(FakePlatform{Linux: true})
		found, err := b.IsDesktopAppPresent(tmpDir, "myapp")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !found {
			t.Errorf("Expected to find desktop app with case-insensitive match")
		}
	})

	t.Run("Linux: directory read error", func(t *testing.T) {
		b := commands.NewBaseCommandCustom(FakePlatform{Linux: true})
		found, err := b.IsDesktopAppPresent("/nonexistent/path", "app")
		if err == nil {
			t.Fatalf("Expected error reading nonexistent directory, got nil")
		}
		if found {
			t.Errorf("Expected not to find the app file, but did")
		}
	})
}

func TestIsPackagePresent_Mac(t *testing.T) {
	b := commands.NewBaseCommandCustom(FakePlatform{Mac: true})

	t.Run("Exact match in brew output", func(t *testing.T) {
		cmd := fakeCmdWithOutput("pkg1\nmytool\npkg2")
		found, err := b.IsPackagePresent(cmd, "mytool")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !found {
			t.Errorf("Expected to find package 'mytool', but did not")
		}
	})

	t.Run("No match in brew output", func(t *testing.T) {
		cmd := fakeCmdWithOutput("pkg1\npkg2")
		found, err := b.IsPackagePresent(cmd, "missing")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if found {
			t.Errorf("Expected not to find package 'missing', but did")
		}
	})
}

func TestIsPackagePresent_Linux(t *testing.T) {
	b := commands.NewBaseCommandCustom(FakePlatform{Linux: true})

	t.Run("Match in dpkg output", func(t *testing.T) {
		// Simulate `dpkg -l` format: status name version
		cmd := fakeCmdWithOutput("ii  mytool  1.0.0\nrc  oldtool  0.9.0")
		found, err := b.IsPackagePresent(cmd, "mytool")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !found {
			t.Errorf("Expected to find package 'mytool', but did not")
		}
	})

	t.Run("No match in dpkg output", func(t *testing.T) {
		cmd := fakeCmdWithOutput("ii  someother  1.0.0")
		found, err := b.IsPackagePresent(cmd, "missing")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if found {
			t.Errorf("Expected not to find package 'missing', but did")
		}
	})
}

func TestIsPackagePresent_CommandError(t *testing.T) {
	b := commands.NewBaseCommandCustom(FakePlatform{Linux: true})

	// Using an invalid command to trigger an error
	cmd := exec.Command("false") // This always exits with non-zero status
	_, err := b.IsPackagePresent(cmd, "anything")
	if err == nil {
		t.Fatalf("Expected error from failed command, got nil")
	}
}
