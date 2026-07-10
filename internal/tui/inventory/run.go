package tuiinventory

import (
	tea "charm.land/bubbletea/v2"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/internal/inventory"
)

// Run starts the shared inventory dashboard. dg list calls this with
// Options{} (unfiltered); dg validate calls it with Options{ProblemsOnly: true}.
func Run(gc *config.GlobalConfig, opts Options) error {
	c := &inventory.Collector{Cmd: commands.NewCommand(), Base: commands.NewBaseCommand()}
	items := c.Collect(gc)
	m := newModel(items, opts)
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}
