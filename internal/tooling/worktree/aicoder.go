package worktree

import (
	"fmt"
	"os/exec"
	"strings"
)

// AICoder represents an AI coding assistant that can be launched in a worktree window
type AICoder interface {
	Name() string
	Command() string
	EnsureInstalled() error
}

// OpenCodeCoder implements AICoder for OpenCode
type OpenCodeCoder struct{}

func (o *OpenCodeCoder) Name() string    { return "opencode" }
func (o *OpenCodeCoder) Command() string { return "opencode" }
func (o *OpenCodeCoder) EnsureInstalled() error {
	if _, err := exec.LookPath("opencode"); err != nil {
		return fmt.Errorf("opencode is not installed. Install it with: dg install --only terminal")
	}
	return nil
}

// ClaudeCoder implements AICoder for Claude Code
type ClaudeCoder struct{}

func (c *ClaudeCoder) Name() string    { return "claude" }
func (c *ClaudeCoder) Command() string { return "CLAUDE_CODE_NO_FLICKER=1 claude" }
func (c *ClaudeCoder) EnsureInstalled() error {
	if _, err := exec.LookPath("claude"); err != nil {
		return fmt.Errorf("claude is not installed. Install it with: dg install --only terminal")
	}
	return nil
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
		return nil, fmt.Errorf("unknown AI coder alias %q. Valid aliases: opencode, oc, claude, cc, claudecode", alias)
	}
}
