package worktree

import (
	"os"
	"strings"
	"testing"
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

// OpenCodeCoder/ClaudeCoder.EnsureInstalled route through the shared
// ensureToolInstalled helper, which itself goes through the swappable
// commands.LookPathFn (see setLookPathFn in repo_candidates_test.go) rather
// than calling exec.LookPath directly - so both the success and failure
// paths are exercisable here without executing a real command.

func TestOpenCodeCoderEnsureInstalledOK(t *testing.T) {
	setLookPathFn(t, func(string) (string, error) {
		return "/usr/bin/opencode", nil
	})

	if err := (&OpenCodeCoder{}).EnsureInstalled(); err != nil {
		t.Fatalf("unexpected error when opencode is on PATH: %v", err)
	}
}

func TestOpenCodeCoderEnsureInstalledMissing(t *testing.T) {
	setLookPathFn(t, func(string) (string, error) {
		return "", os.ErrNotExist
	})

	err := (&OpenCodeCoder{}).EnsureInstalled()
	if err == nil {
		t.Fatal("expected error when opencode is not on PATH, got nil")
	}
	if got := err.Error(); !strings.Contains(got, "opencode") {
		t.Errorf("expected error to mention opencode, got %q", got)
	}
}

func TestClaudeCoderEnsureInstalledOK(t *testing.T) {
	setLookPathFn(t, func(string) (string, error) {
		return "/usr/bin/claude", nil
	})

	if err := (&ClaudeCoder{}).EnsureInstalled(); err != nil {
		t.Fatalf("unexpected error when claude is on PATH: %v", err)
	}
}

func TestClaudeCoderEnsureInstalledMissing(t *testing.T) {
	setLookPathFn(t, func(string) (string, error) {
		return "", os.ErrNotExist
	})

	err := (&ClaudeCoder{}).EnsureInstalled()
	if err == nil {
		t.Fatal("expected error when claude is not on PATH, got nil")
	}
	if got := err.Error(); !strings.Contains(got, "claude") {
		t.Errorf("expected error to mention claude, got %q", got)
	}
}
