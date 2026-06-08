package worktree

import (
	"os"
	"testing"

	"github.com/cjairm/devgita/internal/config"
)

func TestResolveAICoder(t *testing.T) {
	tests := []struct {
		name     string
		alias    string
		wantName string
		wantErr  bool
	}{
		{"opencode full", "opencode", "opencode", false},
		{"opencode short", "oc", "opencode", false},
		{"claude full", "claude", "claude", false},
		{"claude short cc", "cc", "claude", false},
		{"claudecode", "claudecode", "claude", false},
		{"case insensitive OPENCODE", "OPENCODE", "opencode", false},
		{"case insensitive Claude", "Claude", "claude", false},
		{"case insensitive CC", "CC", "claude", false},
		{"unknown alias", "cursor", "", true},
		{"empty alias", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coder, err := ResolveAICoder(tt.alias)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if coder.Name() != tt.wantName {
				t.Errorf("expected name %q, got %q", tt.wantName, coder.Name())
			}
		})
	}
}

func TestOpenCodeCoderCommand(t *testing.T) {
	coder := &OpenCodeCoder{}
	if coder.Command() != "opencode" {
		t.Errorf("expected command 'opencode', got %q", coder.Command())
	}
}

func TestClaudeCoderCommand(t *testing.T) {
	coder := &ClaudeCoder{}
	if coder.Command() != "CLAUDE_CODE_NO_FLICKER=1 claude" {
		t.Errorf("expected command 'CLAUDE_CODE_NO_FLICKER=1 claude', got %q", coder.Command())
	}
}

func TestResolveAIAlias(t *testing.T) {
	// Ensure a clean env for every sub-test.
	t.Cleanup(func() { os.Unsetenv("DEVGITA_WORKTREE_AI") })

	t.Run("flag takes highest priority", func(t *testing.T) {
		os.Setenv("DEVGITA_WORKTREE_AI", "claude")
		gc := &config.GlobalConfig{}
		gc.Worktree.DefaultAI = "claude"
		got := ResolveAIAlias("opencode", gc)
		if got != "opencode" {
			t.Errorf("expected flag value 'opencode', got %q", got)
		}
	})

	t.Run("env var used when no flag", func(t *testing.T) {
		os.Setenv("DEVGITA_WORKTREE_AI", "claude")
		gc := &config.GlobalConfig{}
		gc.Worktree.DefaultAI = "opencode"
		got := ResolveAIAlias("", gc)
		if got != "claude" {
			t.Errorf("expected env value 'claude', got %q", got)
		}
	})

	t.Run("global config used when no flag or env", func(t *testing.T) {
		os.Unsetenv("DEVGITA_WORKTREE_AI")
		gc := &config.GlobalConfig{}
		gc.Worktree.DefaultAI = "claude"
		got := ResolveAIAlias("", gc)
		if got != "claude" {
			t.Errorf("expected config value 'claude', got %q", got)
		}
	})

	t.Run("defaults to opencode when nothing set", func(t *testing.T) {
		os.Unsetenv("DEVGITA_WORKTREE_AI")
		gc := &config.GlobalConfig{}
		got := ResolveAIAlias("", gc)
		if got != "opencode" {
			t.Errorf("expected default 'opencode', got %q", got)
		}
	})

	t.Run("nil gc falls back to default", func(t *testing.T) {
		os.Unsetenv("DEVGITA_WORKTREE_AI")
		got := ResolveAIAlias("", nil)
		if got != "opencode" {
			t.Errorf("expected default 'opencode' with nil gc, got %q", got)
		}
	})
}
