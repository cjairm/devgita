package registry

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/apps/aerospace"
	"github.com/cjairm/devgita/internal/apps/alacritty"
	"github.com/cjairm/devgita/internal/apps/brave"
	"github.com/cjairm/devgita/internal/apps/claude"
	"github.com/cjairm/devgita/internal/apps/devgita"
	"github.com/cjairm/devgita/internal/apps/docker"
	"github.com/cjairm/devgita/internal/apps/fastfetch"
	"github.com/cjairm/devgita/internal/apps/flameshot"
	"github.com/cjairm/devgita/internal/apps/gimp"
	"github.com/cjairm/devgita/internal/apps/git"
	"github.com/cjairm/devgita/internal/apps/i3"
	"github.com/cjairm/devgita/internal/apps/lazydocker"
	"github.com/cjairm/devgita/internal/apps/lazygit"
	"github.com/cjairm/devgita/internal/apps/mise"
	"github.com/cjairm/devgita/internal/apps/neovim"
	"github.com/cjairm/devgita/internal/apps/opencode"
	"github.com/cjairm/devgita/internal/apps/raycast"
	"github.com/cjairm/devgita/internal/apps/tmux"
	"github.com/cjairm/devgita/internal/apps/ulauncher"
)

// factories maps app names to lazy constructors so apps are only instantiated when needed.
var factories = map[string]func() apps.App{
	"aerospace":  func() apps.App { return aerospace.New() },
	"alacritty":  func() apps.App { return alacritty.New() },
	"brave":      func() apps.App { return brave.New() },
	"claude":     func() apps.App { return claude.New() },
	"devgita":    func() apps.App { return devgita.New() },
	"docker":     func() apps.App { return docker.New() },
	"fastfetch":  func() apps.App { return fastfetch.New() },
	"flameshot":  func() apps.App { return flameshot.New() },
	"gimp":       func() apps.App { return gimp.New() },
	"git":        func() apps.App { return git.New() },
	"i3":         func() apps.App { return i3.New() },
	"lazydocker": func() apps.App { return lazydocker.New() },
	"lazygit":    func() apps.App { return lazygit.New() },
	"mise":       func() apps.App { return mise.New() },
	"neovim":     func() apps.App { return neovim.New() },
	"opencode":   func() apps.App { return opencode.New() },
	"raycast":    func() apps.App { return raycast.New() },
	"tmux":       func() apps.App { return tmux.New() },
	"ulauncher":  func() apps.App { return ulauncher.New() },
}

// GetApp returns the App for the given name, or an error listing all supported names.
func GetApp(name string) (apps.App, error) {
	factory, ok := factories[name]
	if !ok {
		return nil, fmt.Errorf("unknown app %q\n\nSupported apps:\n  %s", name, formatNames())
	}
	return factory(), nil
}

// Names returns a sorted slice of all registered app names.
func Names() []string {
	names := make([]string, 0, len(factories))
	for name := range factories {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// formatNames formats registered names into aligned columns for display.
func formatNames() string {
	names := Names()
	const cols = 7
	var rows []string
	for i := 0; i < len(names); i += cols {
		end := i + cols
		if end > len(names) {
			end = len(names)
		}
		rows = append(rows, strings.Join(names[i:end], "  "))
	}
	return strings.Join(rows, "\n  ")
}
