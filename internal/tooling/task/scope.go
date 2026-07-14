package task

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// reviewScopeFetchTimeout bounds the network fetch in ReviewScope so a hung
// remote can't block a caller expecting a fast, single-call response.
const reviewScopeFetchTimeout = 10 * time.Second

// fileChange is one file's row from a merged --numstat / --name-status pair.
// Status is left empty until merged with --name-status output. Binary files
// report Added/Removed as 0 (numstat uses "-" for both, per its format).
type fileChange struct {
	Path    string
	Status  string
	Added   int
	Removed int
	Binary  bool
}

// scopeData is the orchestration result handed to formatReviewScope.
type scopeData struct {
	OnDefaultBranch bool
	Detached        bool
	DetachedSHA     string

	CurrentBranch string
	DefaultBranch string
	FetchFailed   bool

	Behind int
	Ahead  int

	Commits []string

	Files    []fileChange
	Excluded []fileChange
}

// ReviewScope fetches origin (best-effort, bounded), resolves the comparison
// base against the repo's default branch, and returns a compact orientation
// report: ahead/behind, commit subjects, and a per-file stat table with
// lockfile-style noise pulled into a separate excluded-files note.
func (tm *TaskManager) ReviewScope() (string, error) {
	fetchFailed := tm.Git.FetchOriginTimeout(reviewScopeFetchTimeout) != nil

	currentBranch, err := tm.Git.CurrentBranch()
	if err != nil {
		return "", fmt.Errorf("review-scope: %w", err)
	}

	if currentBranch == "" {
		sha, err := tm.Git.ShortHead()
		if err != nil {
			return "", fmt.Errorf("review-scope: %w", err)
		}
		return formatReviewScope(scopeData{Detached: true, DetachedSHA: sha}), nil
	}

	defaultBranch := tm.Git.DefaultBranch()
	if currentBranch == defaultBranch {
		return formatReviewScope(scopeData{
			OnDefaultBranch: true,
			DefaultBranch:   defaultBranch,
		}), nil
	}

	base, err := tm.mergeBase(defaultBranch)
	if err != nil {
		return "", fmt.Errorf("review-scope: %w", err)
	}

	behind, ahead, err := tm.aheadBehind(defaultBranch)
	if err != nil {
		return "", fmt.Errorf("review-scope: %w", err)
	}

	commits, err := tm.commitSubjects(base)
	if err != nil {
		return "", fmt.Errorf("review-scope: %w", err)
	}

	files, err := tm.fileChanges(base + "...HEAD")
	if err != nil {
		return "", fmt.Errorf("review-scope: %w", err)
	}
	reviewable, excluded := partitionExcluded(files)

	return formatReviewScope(scopeData{
		CurrentBranch: currentBranch,
		DefaultBranch: defaultBranch,
		FetchFailed:   fetchFailed,
		Behind:        behind,
		Ahead:         ahead,
		Commits:       commits,
		Files:         reviewable,
		Excluded:      excluded,
	}), nil
}

// mergeBase resolves the merge-base between origin/<defaultBranch> and HEAD —
// the comparison base reused for ahead/behind, commit subjects, and the diff.
func (tm *TaskManager) mergeBase(defaultBranch string) (string, error) {
	out, err := tm.Git.RunCapture("merge-base", "origin/"+defaultBranch, "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// aheadBehind returns how many commits HEAD is ahead/behind origin/<defaultBranch>.
func (tm *TaskManager) aheadBehind(defaultBranch string) (behind, ahead int, err error) {
	out, err := tm.Git.RunCapture(
		"rev-list", "--left-right", "--count", "origin/"+defaultBranch+"...HEAD",
	)
	if err != nil {
		return 0, 0, err
	}
	return parseAheadBehind(out)
}

// parseAheadBehind parses `git rev-list --left-right --count A...B` output
// ("<only-in-A>\t<only-in-B>") into (behind, ahead) when A is the upstream
// default branch and B is HEAD.
func parseAheadBehind(raw string) (behind, ahead int, err error) {
	fields := strings.Fields(strings.TrimSpace(raw))
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("unexpected rev-list --left-right --count output: %q", raw)
	}
	behind, err = strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, fmt.Errorf("unexpected rev-list count %q: %w", fields[0], err)
	}
	ahead, err = strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("unexpected rev-list count %q: %w", fields[1], err)
	}
	return behind, ahead, nil
}

// commitSubjects returns one-line commit subjects for base..HEAD, oldest first.
func (tm *TaskManager) commitSubjects(base string) ([]string, error) {
	out, err := tm.Git.RunCapture("log", "--format=%s", "--reverse", base+"..HEAD")
	if err != nil {
		return nil, err
	}
	out = strings.TrimSpace(out)
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

// fileChanges runs numstat + name-status over rangeSpec (both --no-renames,
// so a rename deterministically renders as a D + A pair) and merges them by
// path into one fileChange list.
func (tm *TaskManager) fileChanges(rangeSpec string) ([]fileChange, error) {
	numstatOut, err := tm.Git.RunCapture("diff", "--numstat", "--no-renames", rangeSpec)
	if err != nil {
		return nil, err
	}
	nameStatusOut, err := tm.Git.RunCapture("diff", "--name-status", "--no-renames", rangeSpec)
	if err != nil {
		return nil, err
	}

	changes, err := parseNumstat(numstatOut)
	if err != nil {
		return nil, err
	}
	statuses, err := parseNameStatus(nameStatusOut)
	if err != nil {
		return nil, err
	}
	for i := range changes {
		if s, ok := statuses[changes[i].Path]; ok {
			changes[i].Status = s
		} else {
			changes[i].Status = "?"
		}
	}
	return changes, nil
}

// parseNumstat parses `git diff --numstat` output into fileChange entries.
// Status is left unset here; callers merge in --name-status separately.
func parseNumstat(raw string) ([]fileChange, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	var out []fileChange
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("malformed numstat line: %q", line)
		}
		fc := fileChange{Path: parts[2]}
		if parts[0] == "-" && parts[1] == "-" {
			fc.Binary = true
		} else {
			added, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("malformed numstat added count in %q: %w", line, err)
			}
			removed, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("malformed numstat removed count in %q: %w", line, err)
			}
			fc.Added, fc.Removed = added, removed
		}
		out = append(out, fc)
	}
	return out, nil
}

// parseNameStatus parses `git diff --name-status` output into a path->status
// map. A status wider than one letter (e.g. a similarity score on a copy) is
// truncated to its leading letter.
func parseNameStatus(raw string) (map[string]string, error) {
	raw = strings.TrimSpace(raw)
	statuses := map[string]string{}
	if raw == "" {
		return statuses, nil
	}
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed name-status line: %q", line)
		}
		status := parts[0]
		if len(status) > 1 {
			status = status[:1]
		}
		statuses[parts[1]] = status
	}
	return statuses, nil
}

// partitionExcluded splits changes into reviewable and excluded (default
// lockfile/generated-asset patterns), preserving each list's relative order.
func partitionExcluded(changes []fileChange) (reviewable, excluded []fileChange) {
	for _, c := range changes {
		if isExcludedPath(c.Path) {
			excluded = append(excluded, c)
		} else {
			reviewable = append(reviewable, c)
		}
	}
	return reviewable, excluded
}

// formatReviewScope renders scopeData as the compact, LLM-oriented report
// (or one of the two stable edge sentinels). Output is payload-only: no
// leading prose, no markdown headers/tables — see the cycle doc's output
// design principles.
func formatReviewScope(s scopeData) string {
	if s.Detached {
		return fmt.Sprintf(
			"Detached HEAD at %s — no branch to compare. Check out a branch or name a target.",
			s.DetachedSHA,
		)
	}
	if s.OnDefaultBranch {
		return fmt.Sprintf(
			"On %s — no branch to compare. Review uncommitted changes or name a target.",
			s.DefaultBranch,
		)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "branch: %s -> %s (default)  [ahead %d, behind %d]",
		s.CurrentBranch, s.DefaultBranch, s.Ahead, s.Behind)
	if s.FetchFailed {
		b.WriteString("  (fetch failed — comparing against local refs)")
	}
	b.WriteString("\n")

	b.WriteString("commits:\n")
	if len(s.Commits) == 0 {
		b.WriteString("(none)\n")
	} else {
		for _, c := range s.Commits {
			fmt.Fprintf(&b, "- %s\n", c)
		}
	}

	fmt.Fprintf(&b, "files (%d):\n", len(s.Files))
	var totalAdded, totalRemoved int
	for _, f := range s.Files {
		if f.Binary {
			fmt.Fprintf(&b, "%-2s %s  binary\n", f.Status, f.Path)
			continue
		}
		fmt.Fprintf(&b, "%-2s %s  +%d/-%d\n", f.Status, f.Path, f.Added, f.Removed)
		totalAdded += f.Added
		totalRemoved += f.Removed
	}
	fmt.Fprintf(&b, "total: +%d/-%d", totalAdded, totalRemoved)

	if len(s.Excluded) > 0 {
		b.WriteString("\n")
		notes := make([]string, len(s.Excluded))
		for i, f := range s.Excluded {
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
