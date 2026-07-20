package worktree

import (
	"fmt"
	"strings"

	"github.com/cjairm/devgita/internal/commands"
)

// AICoder represents an AI coding assistant that can be launched in a worktree window
type AICoder interface {
	Name() string
	Command() string
	EnsureInstalled() error
}

// ensureToolInstalled looks up name via commands.LookPathFn and returns a
// consistent, actionable error if it isn't found. Shared by every
// EnsureInstalled below (opencode, claude) and by layout.go's nvim check
// (nvim has no AICoder wrapper since it isn't an AI coder) - one lookup +
// error-format shape instead of three hand-rolled copies of it.
//
// Going through commands.LookPathFn (rather than calling exec.LookPath
// directly) is what makes all three checks swappable in tests - see
// repo_candidates_test.go's setLookPathFn helper, reused by this package's
// tests.
func ensureToolInstalled(name string) error {
	if _, err := commands.LookPathFn(name); err != nil {
		return fmt.Errorf("%s is not installed. Install it with: dg install --only terminal", name)
	}
	return nil
}

// OpenCodeCoder implements AICoder for OpenCode
type OpenCodeCoder struct{}

func (o *OpenCodeCoder) Name() string    { return "opencode" }
func (o *OpenCodeCoder) Command() string { return "opencode" }
func (o *OpenCodeCoder) EnsureInstalled() error {
	return ensureToolInstalled(o.Name())
}

// ClaudeCoder implements AICoder for Claude Code
type ClaudeCoder struct{}

func (c *ClaudeCoder) Name() string    { return "claude" }
func (c *ClaudeCoder) Command() string { return "CLAUDE_CODE_NO_FLICKER=1 claude" }
func (c *ClaudeCoder) EnsureInstalled() error {
	return ensureToolInstalled(c.Name())
}

// ResolveAICoder resolves an alias to an AICoder implementation
// Valid aliases (case-insensitive):
//   - opencode, oc -> OpenCodeCoder
//   - claude, cc, claudecode -> ClaudeCoder
func ResolveAICoder(alias string) (AICoder, error) {
	switch strings.ToLower(alias) {
	case "opencode", "oc":
		return &OpenCodeCoder{}, nil
	case "claude", "cc", "claudecode":
		return &ClaudeCoder{}, nil
	default:
		return nil, fmt.Errorf(
			"unknown AI coder alias %q. Valid aliases: opencode, oc, claude, cc, claudecode",
			alias,
		)
	}
}
