package tuicomponents_test

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	tuicomponents "github.com/cjairm/devgita/internal/tui/components"
)

func TestTabBarActiveContainsLabel(t *testing.T) {
	p := tuicomponents.NewPalette()
	tabs := []tuicomponents.Tab{{Label: "Agent"}, {Label: "Diff"}}
	got := ansi.Strip(p.TabBar(tabs, 0))
	if !strings.Contains(got, "Agent") {
		t.Errorf("TabBar output (stripped) %q missing active label Agent", got)
	}
}

func TestTabBarInactiveContainsLabel(t *testing.T) {
	p := tuicomponents.NewPalette()
	tabs := []tuicomponents.Tab{{Label: "Agent"}, {Label: "Diff"}}
	got := ansi.Strip(p.TabBar(tabs, 0))
	if !strings.Contains(got, "Diff") {
		t.Errorf("TabBar output (stripped) %q missing inactive label Diff", got)
	}
}

func TestTabBarOutOfRangeNoPanic(t *testing.T) {
	p := tuicomponents.NewPalette()
	tabs := []tuicomponents.Tab{{Label: "A"}, {Label: "B"}}
	_ = p.TabBar(tabs, -1)
	_ = p.TabBar(tabs, 99)
}

func TestTabBarEmptyReturnsEmpty(t *testing.T) {
	p := tuicomponents.NewPalette()
	got := p.TabBar(nil, 0)
	if got != "" {
		t.Errorf("TabBar(nil, 0) = %q, want empty string", got)
	}
}
