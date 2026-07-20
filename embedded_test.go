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
