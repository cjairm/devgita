package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveFromInstalled(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*GlobalConfig)
		removeItem  string
		removeType  string
		checkField  func(*GlobalConfig) []string
		wantRemains []string
	}{
		{
			name: "removes existing package",
			setup: func(gc *GlobalConfig) {
				gc.Installed.Packages = []string{"git", "tmux", "neovim"}
			},
			removeItem:  "tmux",
			removeType:  "package",
			checkField:  func(gc *GlobalConfig) []string { return gc.Installed.Packages },
			wantRemains: []string{"git", "neovim"},
		},
		{
			name: "no-op when package absent",
			setup: func(gc *GlobalConfig) {
				gc.Installed.Packages = []string{"git"}
			},
			removeItem:  "tmux",
			removeType:  "package",
			checkField:  func(gc *GlobalConfig) []string { return gc.Installed.Packages },
			wantRemains: []string{"git"},
		},
		{
			name: "removes desktop_app",
			setup: func(gc *GlobalConfig) {
				gc.Installed.DesktopApps = []string{"brave", "alacritty", "raycast"}
			},
			removeItem:  "brave",
			removeType:  "desktop_app",
			checkField:  func(gc *GlobalConfig) []string { return gc.Installed.DesktopApps },
			wantRemains: []string{"alacritty", "raycast"},
		},
		{
			name: "does not affect already_installed",
			setup: func(gc *GlobalConfig) {
				gc.Installed.Packages = []string{"git"}
				gc.AlreadyInstalled.Packages = []string{"git"}
			},
			removeItem:  "git",
			removeType:  "package",
			checkField:  func(gc *GlobalConfig) []string { return gc.AlreadyInstalled.Packages },
			wantRemains: []string{"git"},
		},
		{
			name: "unknown item type is no-op",
			setup: func(gc *GlobalConfig) {
				gc.Installed.Packages = []string{"git"}
			},
			removeItem:  "git",
			removeType:  "unknown_type",
			checkField:  func(gc *GlobalConfig) []string { return gc.Installed.Packages },
			wantRemains: []string{"git"},
		},
		{
			name: "removes terminal_tool",
			setup: func(gc *GlobalConfig) {
				gc.Installed.TerminalTools = []string{"tool1", "tool2"}
			},
			removeItem:  "tool1",
			removeType:  "terminal_tool",
			checkField:  func(gc *GlobalConfig) []string { return gc.Installed.TerminalTools },
			wantRemains: []string{"tool2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := &GlobalConfig{}
			tt.setup(gc)
			gc.RemoveFromInstalled(tt.removeItem, tt.removeType)
			assert.Equal(t, tt.wantRemains, tt.checkField(gc))
		})
	}
}
