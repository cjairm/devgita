package tuiworktree

import (
	lip "charm.land/lipgloss/v2"
)

type Styles struct {
	RepoHeader         lip.Style
	WorktreeRow        lip.Style
	SelectedRow        lip.Style
	ArmedRow           lip.Style
	ActiveGlyph        lip.Style
	InactiveGlyph      lip.Style
	DirtyAdded         lip.Style
	DirtyRemoved       lip.Style
	ActiveTab          lip.Style
	InactiveTab        lip.Style
	HintBar            lip.Style
	StatusMsg          lip.Style
	Divider            lip.Style
	OfflinePlaceholder lip.Style
	HelpKey            lip.Style
	HelpBorder         lip.Style
}

func newStyles() Styles {
	return Styles{
		RepoHeader:  lip.NewStyle().Bold(true).Foreground(lip.Color("6")),
		WorktreeRow: lip.NewStyle(),
		SelectedRow: lip.NewStyle().
			Background(lip.Color("4")).
			Foreground(lip.Color("15")).
			Bold(true),
		ArmedRow:           lip.NewStyle().Background(lip.Color("1")).Foreground(lip.Color("15")),
		ActiveGlyph:        lip.NewStyle().Foreground(lip.Color("2")),
		InactiveGlyph:      lip.NewStyle().Foreground(lip.Color("8")),
		DirtyAdded:         lip.NewStyle().Foreground(lip.Color("2")),
		DirtyRemoved:       lip.NewStyle().Foreground(lip.Color("1")),
		ActiveTab:          lip.NewStyle().Bold(true).Underline(true).Foreground(lip.Color("15")),
		InactiveTab:        lip.NewStyle().Foreground(lip.Color("8")),
		HintBar:            lip.NewStyle().Foreground(lip.Color("8")),
		StatusMsg:          lip.NewStyle().Foreground(lip.Color("3")),
		Divider:            lip.NewStyle().Foreground(lip.Color("8")),
		OfflinePlaceholder: lip.NewStyle().Foreground(lip.Color("8")).Italic(true),
		HelpKey:            lip.NewStyle().Bold(true).Foreground(lip.Color("6")),
		HelpBorder:         lip.NewStyle().Foreground(lip.Color("8")),
	}
}
