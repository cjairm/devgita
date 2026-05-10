package registry

import (
	"sort"
	"testing"

	"github.com/cjairm/devgita/internal/testutil"
)

func init() { testutil.InitLogger() }

var expectedApps = []string{
	"aerospace", "alacritty", "brave", "claude", "devgita", "docker",
	"fastfetch", "flameshot", "gimp", "git", "i3", "lazydocker",
	"lazygit", "mise", "neovim", "opencode", "raycast", "tmux", "ulauncher",
}

func TestGetApp_KnownApp(t *testing.T) {
	app, err := GetApp("git")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if app == nil {
		t.Fatal("expected non-nil app")
	}
	if app.Name() != "git" {
		t.Errorf("expected Name() == %q, got %q", "git", app.Name())
	}
}

func TestGetApp_UnknownApp(t *testing.T) {
	_, err := GetApp("notanapp")
	if err == nil {
		t.Fatal("expected error for unknown app, got nil")
	}
}

func TestGetApp_AllRegisteredApps(t *testing.T) {
	for _, name := range expectedApps {
		t.Run(name, func(t *testing.T) {
			app, err := GetApp(name)
			if err != nil {
				t.Fatalf("GetApp(%q) returned error: %v", name, err)
			}
			if app == nil {
				t.Fatalf("GetApp(%q) returned nil app", name)
			}
			if app.Name() != name {
				t.Errorf("GetApp(%q).Name() = %q, want %q", name, app.Name(), name)
			}
		})
	}
}

func TestNames_ContainsAllApps(t *testing.T) {
	names := Names()
	if len(names) != len(expectedApps) {
		t.Errorf("Names() returned %d names, want %d", len(names), len(expectedApps))
	}
	if !sort.StringsAreSorted(names) {
		t.Error("Names() is not sorted")
	}
	nameSet := make(map[string]bool, len(names))
	for _, n := range names {
		nameSet[n] = true
	}
	for _, expected := range expectedApps {
		if !nameSet[expected] {
			t.Errorf("Names() missing %q", expected)
		}
	}
}
