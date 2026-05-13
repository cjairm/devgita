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
	"github.com/cjairm/devgita/pkg/constants"
)

// factories maps app names to lazy constructors so apps are only instantiated when needed.
var factories = map[string]func() apps.App{
	constants.Aerospace:  func() apps.App { return aerospace.New() },
	constants.Alacritty:  func() apps.App { return alacritty.New() },
	constants.Brave:      func() apps.App { return brave.New() },
	constants.Claude:     func() apps.App { return claude.New() },
	constants.DevgitaApp: func() apps.App { return devgita.New() },
	constants.Docker:     func() apps.App { return docker.New() },
	constants.Fastfetch:  func() apps.App { return fastfetch.New() },
	constants.Flameshot:  func() apps.App { return flameshot.New() },
	constants.Gimp:       func() apps.App { return gimp.New() },
	constants.Git:        func() apps.App { return git.New() },
	constants.I3:         func() apps.App { return i3.New() },
	constants.LazyDocker: func() apps.App { return lazydocker.New() },
	constants.LazyGit:    func() apps.App { return lazygit.New() },
	constants.Mise:       func() apps.App { return mise.New() },
	constants.Neovim:     func() apps.App { return neovim.New() },
	constants.OpenCode:   func() apps.App { return opencode.New() },
	constants.Raycast:    func() apps.App { return raycast.New() },
	constants.Tmux:       func() apps.App { return tmux.New() },
	constants.Ulauncher:  func() apps.App { return ulauncher.New() },
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

// GetAppsByKind returns sorted names of all registered apps matching the given kind.
func GetAppsByKind(kind apps.AppKind) []string {
	var names []string
	for name, factory := range factories {
		if factory().Kind() == kind {
			names = append(names, name)
		}
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
