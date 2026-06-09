package tuicomponents

import "fmt"

// DiffStat renders a colored "+added -removed" string.
func (p *Palette) DiffStat(added, removed int) string {
	return p.DiffAdded.Render(
		fmt.Sprintf("+%d", added),
	) + " " + p.DiffRemoved.Render(
		fmt.Sprintf("-%d", removed),
	)
}

// DirtyCount renders a colored "±N" file-count string.
func (p *Palette) DirtyCount(count int) string {
	return p.DiffFiles.Render(fmt.Sprintf("±%d", count))
}

// DiffStatLine renders "±files +added -removed", omitting the "±0" prefix when files == 0.
func (p *Palette) DiffStatLine(files, added, removed int) string {
	ds := p.DiffStat(added, removed)
	if files <= 0 {
		return ds
	}
	return p.DirtyCount(files) + " " + ds
}
