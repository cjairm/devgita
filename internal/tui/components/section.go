package tuicomponents

import "fmt"

// SectionHeader renders a T2/T3 tree section header: "LABEL · N".
// Width <= 0 returns "".
func (p *Palette) SectionHeader(label string, count int, width int) string {
	if width <= 0 {
		return ""
	}
	return p.SectionHead.Render(fmt.Sprintf("%s · %d", label, count))
}
