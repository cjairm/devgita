package paths_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/paths"
)

var cmd = commands.NewBaseCommand()

func getHome(t *testing.T) string {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal("could not get home dir:", err)
	}
	return home
}

func TestConfigDir(t *testing.T) {
	home := getHome(t)

	t.Run("no subdirs", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		got := paths.GetConfigDir()
		want := filepath.Join(home, ".config")
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("one subdir", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		got := paths.GetConfigDir(constants.AppName)
		want := filepath.Join(home, ".config", constants.AppName)
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("multiple subdirs", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		got := paths.GetConfigDir(constants.AppName, "nvim")
		want := filepath.Join(home, ".config", constants.AppName, "nvim")
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("XDG_CONFIG_HOME override", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "/tmp/xdg-config")
		got := paths.GetConfigDir(constants.AppName)
		want := filepath.Join("/tmp/xdg-config", constants.AppName)
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})
}

func TestDataDir(t *testing.T) {
	home := getHome(t)

	t.Run("default location", func(t *testing.T) {
		t.Setenv("XDG_DATA_HOME", "")
		got := paths.GetDataDir(constants.AppName)
		want := filepath.Join(home, ".local", "share", constants.AppName)
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("no subdir", func(t *testing.T) {
		t.Setenv("XDG_DATA_HOME", "")
		got := paths.GetDataDir()
		want := filepath.Join(home, ".local", "share")
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("override", func(t *testing.T) {
		t.Setenv("XDG_DATA_HOME", "/tmp/xdg-data")
		got := paths.GetDataDir("app")
		want := filepath.Join("/tmp/xdg-data", "app")
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})
}

func TestCacheDir(t *testing.T) {
	home := getHome(t)

	t.Run("default location", func(t *testing.T) {
		t.Setenv("XDG_CACHE_HOME", "")
		got := paths.GetCacheDir(constants.AppName)
		want := filepath.Join(home, ".cache", constants.AppName)
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("multiple subdirs", func(t *testing.T) {
		t.Setenv("XDG_CACHE_HOME", "")
		got := paths.GetCacheDir(constants.AppName, "nvim")
		want := filepath.Join(home, ".cache", constants.AppName, "nvim")
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("override", func(t *testing.T) {
		t.Setenv("XDG_CACHE_HOME", "/tmp/xdg-cache")
		got := paths.GetCacheDir("logs")
		want := filepath.Join("/tmp/xdg-cache", "logs")
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})
}

func TestAppDir(t *testing.T) {
	home := getHome(t)

	t.Run("returns app dir", func(t *testing.T) {
		got := paths.GetAppDir("logs")
		want := filepath.Join(home, ".local", "share", constants.AppName, "logs")
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})
}
