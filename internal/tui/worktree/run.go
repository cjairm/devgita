package tuiworktree

import (
	tea "charm.land/bubbletea/v2"

	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/internal/tooling/worktree"
)

// Run starts the worktree TUI dashboard.
func Run() error {
	gc := &config.GlobalConfig{}
	// Best-effort: load global config for AI alias resolution; if it fails, defaults apply.
	_ = gc.Load()

	mgr := worktree.New()
	m := newModel(mgr, mgr.Tmux, mgr.Git, gc)

	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}
