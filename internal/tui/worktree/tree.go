package tuiworktree

import (
	"sort"
	"strings"

	"github.com/cjairm/devgita/internal/tooling/worktree"
)

type rowKind int

const (
	rowRepo rowKind = iota
	rowWorktree
)

type row struct {
	kind   rowKind
	repo   string
	status worktree.WorktreeStatus
}

// buildRows groups statuses by repo (alpha-sorted), applies filter, respects collapsed map.
func buildRows(statuses []worktree.WorktreeStatus, collapsed map[string]bool, filter string) []row {
	// Group by repo
	groups := map[string][]worktree.WorktreeStatus{}
	for _, s := range statuses {
		groups[s.Repo] = append(groups[s.Repo], s)
	}

	// Sort repos
	repos := make([]string, 0, len(groups))
	for r := range groups {
		repos = append(repos, r)
	}
	sort.Strings(repos)

	filter = strings.ToLower(filter)
	var rows []row
	for _, repo := range repos {
		children := groups[repo]
		// Filter: keep only children that match
		var visible []worktree.WorktreeStatus
		for _, s := range children {
			if filter == "" || strings.Contains(strings.ToLower(repo+"/"+s.Name), filter) {
				visible = append(visible, s)
			}
		}
		if len(visible) == 0 {
			continue
		}
		rows = append(rows, row{kind: rowRepo, repo: repo})
		if !collapsed[repo] {
			for _, s := range visible {
				rows = append(rows, row{kind: rowWorktree, repo: repo, status: s})
			}
		}
	}
	return rows
}

// worktreeIndices returns indices into rows that are rowWorktree kind.
func worktreeIndices(rows []row) []int {
	var out []int
	for i, r := range rows {
		if r.kind == rowWorktree {
			out = append(out, i)
		}
	}
	return out
}
