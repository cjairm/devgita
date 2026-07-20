package main

import (
	"io/fs"
	"strings"
	"testing"
)

// tmux-resurrect restores panes with `cat <saved-contents>; exec <default-command>`.
// If default-command is a compound command (contains `;` or `&&`), the splice
// truncates at the first separator, the pane execs the wrong word and dies at
// birth — restored sessions collapse and the tmux server exits (terminal
// window auto-closes). This guard keeps default-command a single simple
// command so that class of breakage cannot be reintroduced.
func TestTmuxDefaultCommandStaysResurrectSafe(t *testing.T) {
	data, err := fs.ReadFile(ConfigsFS, "configs/tmux/tmux.conf")
	if err != nil {
		t.Fatalf("failed to read embedded tmux.conf: %v", err)
	}

	var defaultCommandLine string
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "set -g default-command") {
			defaultCommandLine = trimmed
			break
		}
	}
	if defaultCommandLine == "" {
		// Not setting default-command at all is safe (tmux spawns login shells).
		return
	}

	value := strings.TrimSpace(strings.TrimPrefix(defaultCommandLine, "set -g default-command"))
	value = strings.Trim(value, `'"`)
	for _, sep := range []string{";", "&&", "||"} {
		if strings.Contains(value, sep) {
			t.Errorf(
				"default-command %q contains %q: it must stay a single simple command — tmux-resurrect splices it after `exec` when restoring panes, and a compound value kills every restored pane at birth",
				value,
				sep,
			)
		}
	}
}

// configs/zsh/zshenv.zsh is sourced from ~/.zshenv on every zsh startup
// (login or not), before /etc/zshrc runs. zsh startup semantics — not a
// convention we control — impose the guard shape this test checks: without
// the PATH probe the script would eval path_helper unconditionally on every
// shell (including already-healthy ones), and without the executable check
// it would blow up on Linux, where path_helper doesn't exist.
func TestZshenvStaysPathRepairSafe(t *testing.T) {
	data, err := fs.ReadFile(ConfigsFS, "configs/zsh/zshenv.zsh")
	if err != nil {
		t.Fatalf("failed to read embedded zshenv.zsh: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, `":$PATH:" != *":/usr/bin:"*`) {
		t.Error(
			"zshenv.zsh must guard on PATH missing /usr/bin — without it, the script would " +
				"re-run path_helper on every shell startup instead of only broken ones",
		)
	}
	if !strings.Contains(content, "/usr/libexec/path_helper") {
		t.Error(
			"zshenv.zsh must check for /usr/libexec/path_helper — without it, the script " +
				"would fail on Linux, where path_helper doesn't exist",
		)
	}
}
