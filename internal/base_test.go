package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	commands "github.com/cjairm/devgita/internal"
	"github.com/cjairm/devgita/pkg/utils"
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
		got, err := cmd.ConfigDir()
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(home, ".config")
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("one subdir", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		got, err := cmd.ConfigDir(utils.APP_NAME)
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(home, ".config", utils.APP_NAME)
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("multiple subdirs", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		got, err := cmd.ConfigDir(utils.APP_NAME, "nvim")
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(home, ".config", utils.APP_NAME, "nvim")
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("XDG_CONFIG_HOME override", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "/tmp/xdg-config")
		got, err := cmd.ConfigDir(utils.APP_NAME)
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join("/tmp/xdg-config", utils.APP_NAME)
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})
}

func TestDataDir(t *testing.T) {
	home := getHome(t)

	t.Run("default location", func(t *testing.T) {
		t.Setenv("XDG_DATA_HOME", "")
		got, err := cmd.DataDir(utils.APP_NAME)
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(home, ".local", "share", utils.APP_NAME)
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("no subdir", func(t *testing.T) {
		t.Setenv("XDG_DATA_HOME", "")
		got, err := cmd.DataDir()
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(home, ".local", "share")
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("override", func(t *testing.T) {
		t.Setenv("XDG_DATA_HOME", "/tmp/xdg-data")
		got, err := cmd.DataDir("app")
		if err != nil {
			t.Fatal(err)
		}
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
		got, err := cmd.CacheDir(utils.APP_NAME)
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(home, ".cache", utils.APP_NAME)
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("multiple subdirs", func(t *testing.T) {
		t.Setenv("XDG_CACHE_HOME", "")
		got, err := cmd.CacheDir(utils.APP_NAME, "nvim")
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(home, ".cache", utils.APP_NAME, "nvim")
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("override", func(t *testing.T) {
		t.Setenv("XDG_CACHE_HOME", "/tmp/xdg-cache")
		got, err := cmd.CacheDir("logs")
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join("/tmp/xdg-cache", "logs")
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})
}

func TestAppDir(t *testing.T) {
	home := getHome(t)

	t.Run("returns app dir", func(t *testing.T) {
		got, err := cmd.AppDir("logs")
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(home, ".local", "share", utils.APP_NAME, "logs")
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})
}
