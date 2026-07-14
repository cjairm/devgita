package task

import (
	"fmt"
	"strings"
)

// BranchDiff returns the merge-base diff against the default branch, with
// lockfile-style noise excluded by default (see exclusions.go) and called
// out in a trailing notes line so nothing is silently hidden.
//
// It does not fetch: review-scope is the orient-and-fetch step, and
// branch-diff is follow-up retrieval within the same review session.
// Re-fetching per file pull would be wasteful and could shift the
// merge-base mid-session if origin moved between calls.
//
// file, when non-empty, bypasses exclusions and returns only that file's
// diff — an explicit request wins over the default noise filter.
func (tm *TaskManager) BranchDiff(file string) (string, error) {
	defaultBranch := tm.Git.DefaultBranch()
	base, err := tm.mergeBase(defaultBranch)
	if err != nil {
		return "", fmt.Errorf("branch-diff: %w", err)
	}
	rangeSpec := base + "...HEAD"

	if file != "" {
		return tm.branchDiffFile(rangeSpec, file)
	}
	return tm.branchDiffAll(rangeSpec)
}

// branchDiffFile returns file's diff over rangeSpec, without exclusions.
// file is passed as its own argv element (never shell-interpolated), so it
// needs no escaping.
func (tm *TaskManager) branchDiffFile(rangeSpec, file string) (string, error) {
	diff, err := tm.Git.RunCapture("diff", rangeSpec, "--", file)
	if err != nil {
		return "", fmt.Errorf("branch-diff: %w", err)
	}
	if strings.TrimSpace(diff) == "" {
		return fmt.Sprintf("No changes for %s in %s.", file, rangeSpec), nil
	}
	return diff, nil
}

// branchDiffAll returns the full diff over rangeSpec with the default
// exclusion patterns applied in a single `git diff` invocation, plus a
// trailing note for every excluded file that actually changed.
func (tm *TaskManager) branchDiffAll(rangeSpec string) (string, error) {
	args := append([]string{"diff", rangeSpec, "--", "."}, exclusionPathspecs()...)
	diff, err := tm.Git.RunCapture(args...)
	if err != nil {
		return "", fmt.Errorf("branch-diff: %w", err)
	}

	numstatOut, err := tm.Git.RunCapture("diff", "--numstat", "--no-renames", rangeSpec)
	if err != nil {
		return "", fmt.Errorf("branch-diff: %w", err)
	}
	changes, err := parseNumstat(numstatOut)
	if err != nil {
		return "", fmt.Errorf("branch-diff: %w", err)
	}
	_, excluded := partitionExcluded(changes)

	return formatBranchDiff(rangeSpec, diff, excluded), nil
}

// formatBranchDiff renders the diff payload plus an exclusion-notes line.
// When the filtered diff is empty but some files were excluded, a sentinel
// takes the diff's place so the payload is never empty; when nothing
// changed at all (nothing to exclude either), a distinct sentinel says so.
func formatBranchDiff(rangeSpec, diff string, excluded []fileChange) string {
	trimmed := strings.TrimSpace(diff)

	var b strings.Builder
	switch {
	case trimmed != "":
		b.WriteString(diff)
	case len(excluded) > 0:
		fmt.Fprintf(
			&b,
			"No reviewable changes in %s (all changes excluded — see notes below).",
			rangeSpec,
		)
	default:
		fmt.Fprintf(&b, "No changes in %s.", rangeSpec)
	}

	if len(excluded) > 0 {
		b.WriteString("\n")
		notes := make([]string, len(excluded))
		for i, f := range excluded {
			if f.Binary {
				notes[i] = fmt.Sprintf("%s (binary)", f.Path)
			} else {
				notes[i] = fmt.Sprintf("%s (+%d/-%d)", f.Path, f.Added, f.Removed)
			}
		}
		fmt.Fprintf(&b, "excluded (see `dg task branch-diff --file <path>` to inspect): %s",
			strings.Join(notes, ", "))
	}

	return b.String()
}
