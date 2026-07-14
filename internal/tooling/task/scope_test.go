package task

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/commands"
)

// --- Pure formatter / parser tests (fixtures, no mocks) ---

func TestParseNumstat(t *testing.T) {
	t.Run("parses regular and binary lines", func(t *testing.T) {
		raw := "120\t30\tinternal/tooling/task/task.go\n" +
			"200\t0\tinternal/tooling/task/scope.go\n" +
			"-\t-\tassets/logo.png\n"
		got, err := parseNumstat(raw)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 3 {
			t.Fatalf("expected 3 entries, got %d", len(got))
		}
		if got[0].Path != "internal/tooling/task/task.go" || got[0].Added != 120 ||
			got[0].Removed != 30 {
			t.Fatalf("unexpected entry 0: %+v", got[0])
		}
		if got[2].Path != "assets/logo.png" || !got[2].Binary {
			t.Fatalf("expected entry 2 to be binary: %+v", got[2])
		}
	})

	t.Run("empty input yields nil", func(t *testing.T) {
		got, err := parseNumstat("   \n  ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}
	})

	t.Run("malformed line errors", func(t *testing.T) {
		if _, err := parseNumstat("not-a-numstat-line"); err == nil {
			t.Fatal("expected error for malformed line")
		}
	})

	t.Run("non-numeric count errors", func(t *testing.T) {
		if _, err := parseNumstat("abc\t1\tfile.go"); err == nil {
			t.Fatal("expected error for non-numeric added count")
		}
	})
}

func TestParseNameStatus(t *testing.T) {
	t.Run("parses statuses", func(t *testing.T) {
		raw := "M\tinternal/tooling/task/task.go\nA\tinternal/tooling/task/scope.go\n"
		got, err := parseNameStatus(raw)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got["internal/tooling/task/task.go"] != "M" ||
			got["internal/tooling/task/scope.go"] != "A" {
			t.Fatalf("unexpected statuses: %+v", got)
		}
	})

	t.Run("truncates multi-char status to leading letter", func(t *testing.T) {
		got, err := parseNameStatus("R100\told.go\n")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got["old.go"] != "R" {
			t.Fatalf("expected status 'R', got %q", got["old.go"])
		}
	})

	t.Run("empty input yields empty map", func(t *testing.T) {
		got, err := parseNameStatus("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("expected empty map, got %+v", got)
		}
	})

	t.Run("malformed line errors", func(t *testing.T) {
		if _, err := parseNameStatus("nofield"); err == nil {
			t.Fatal("expected error for malformed line")
		}
	})
}

func TestParseAheadBehind(t *testing.T) {
	t.Run("parses behind and ahead", func(t *testing.T) {
		behind, ahead, err := parseAheadBehind("3\t5\n")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if behind != 3 || ahead != 5 {
			t.Fatalf("expected behind=3 ahead=5, got behind=%d ahead=%d", behind, ahead)
		}
	})

	t.Run("wrong field count errors", func(t *testing.T) {
		if _, _, err := parseAheadBehind("3"); err == nil {
			t.Fatal("expected error for wrong field count")
		}
	})

	t.Run("non-numeric field errors", func(t *testing.T) {
		if _, _, err := parseAheadBehind("x\t5"); err == nil {
			t.Fatal("expected error for non-numeric field")
		}
	})
}

func TestPartitionExcluded(t *testing.T) {
	changes := []fileChange{
		{Path: "internal/tooling/task/task.go", Status: "M", Added: 120, Removed: 30},
		{Path: "go.sum", Status: "M", Added: 40, Removed: 12},
		{Path: "packages/app/package-lock.json", Status: "M", Added: 5, Removed: 1},
	}
	reviewable, excluded := partitionExcluded(changes)
	if len(reviewable) != 1 || reviewable[0].Path != "internal/tooling/task/task.go" {
		t.Fatalf("unexpected reviewable: %+v", reviewable)
	}
	if len(excluded) != 2 {
		t.Fatalf("expected 2 excluded, got %+v", excluded)
	}
}

func TestFormatReviewScope(t *testing.T) {
	t.Run("detached HEAD sentinel", func(t *testing.T) {
		got := formatReviewScope(scopeData{Detached: true, DetachedSHA: "abc1234"})
		want := "Detached HEAD at abc1234 — no branch to compare. Check out a branch or name a target."
		if got != want {
			t.Fatalf("unexpected output: %q", got)
		}
	})

	t.Run("on default branch sentinel", func(t *testing.T) {
		got := formatReviewScope(scopeData{OnDefaultBranch: true, DefaultBranch: "main"})
		want := "On main — no branch to compare. Review uncommitted changes or name a target."
		if got != want {
			t.Fatalf("unexpected output: %q", got)
		}
	})

	t.Run("full report with excluded file", func(t *testing.T) {
		got := formatReviewScope(scopeData{
			CurrentBranch: "feat/x",
			DefaultBranch: "main",
			Ahead:         2,
			Behind:        0,
			Commits: []string{
				"feat(task): add review-scope",
				"test(task): cover offline fetch",
			},
			Files: []fileChange{
				{Path: "internal/tooling/task/task.go", Status: "M", Added: 120, Removed: 30},
				{Path: "internal/tooling/task/scope.go", Status: "A", Added: 200, Removed: 0},
			},
			Excluded: []fileChange{
				{Path: "go.sum", Status: "M", Added: 40, Removed: 12},
			},
		})
		want := "branch: feat/x -> main (default)  [ahead 2, behind 0]\n" +
			"commits:\n" +
			"- feat(task): add review-scope\n" +
			"- test(task): cover offline fetch\n" +
			"files (2):\n" +
			"M  internal/tooling/task/task.go  +120/-30\n" +
			"A  internal/tooling/task/scope.go  +200/-0\n" +
			"total: +320/-30\n" +
			"excluded (see `dg task branch-diff --file <path>` to inspect): go.sum (+40/-12)"
		if got != want {
			t.Fatalf("unexpected output:\n%s\n---want---\n%s", got, want)
		}
	})

	t.Run("marks fetch failure", func(t *testing.T) {
		got := formatReviewScope(scopeData{
			CurrentBranch: "feat/x",
			DefaultBranch: "main",
			FetchFailed:   true,
		})
		if !strings.Contains(got, "(fetch failed — comparing against local refs)") {
			t.Fatalf("expected fetch-failed marker, got: %q", got)
		}
	})

	t.Run("no commits renders (none)", func(t *testing.T) {
		got := formatReviewScope(scopeData{CurrentBranch: "feat/x", DefaultBranch: "main"})
		if !strings.Contains(got, "commits:\n(none)\n") {
			t.Fatalf("expected '(none)' for empty commits, got: %q", got)
		}
	})

	t.Run("binary file renders without counts", func(t *testing.T) {
		got := formatReviewScope(scopeData{
			CurrentBranch: "feat/x",
			DefaultBranch: "main",
			Files:         []fileChange{{Path: "assets/logo.png", Status: "A", Binary: true}},
		})
		if !strings.Contains(got, "A  assets/logo.png  binary") {
			t.Fatalf("expected binary file row, got: %q", got)
		}
	})
}

// --- Orchestration tests (mocked git.Base, no real commands) ---

func TestReviewScope(t *testing.T) {
	t.Run("happy path renders full report", func(t *testing.T) {
		tm, gitBase, _ := newTaskSetup()
		gitBase.SetExecCommandResults(
			commands.ExecCommandResult("", "", nil),              // fetch origin
			commands.ExecCommandResult("feat/x\n", "", nil),      // branch --show-current
			commands.ExecCommandResult("origin/main\n", "", nil), // symbolic-ref (default branch)
			commands.ExecCommandResult("abc123\n", "", nil),      // merge-base
			commands.ExecCommandResult("0\t2\n", "", nil),        // rev-list --left-right --count
			commands.ExecCommandResult(
				"feat(task): add review-scope\ntest(task): cover offline fetch\n", "", nil,
			), // log --format=%s --reverse
			commands.ExecCommandResult(
				"120\t30\tinternal/tooling/task/task.go\n"+
					"200\t0\tinternal/tooling/task/scope.go\n"+
					"40\t12\tgo.sum\n", "", nil,
			), // diff --numstat --no-renames
			commands.ExecCommandResult(
				"M\tinternal/tooling/task/task.go\n"+
					"A\tinternal/tooling/task/scope.go\n"+
					"M\tgo.sum\n", "", nil,
			), // diff --name-status --no-renames
		)

		out, err := tm.ReviewScope()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		wantLines := []string{
			"branch: feat/x -> main (default)  [ahead 2, behind 0]",
			"- feat(task): add review-scope",
			"- test(task): cover offline fetch",
			"files (2):",
			"M  internal/tooling/task/task.go  +120/-30",
			"A  internal/tooling/task/scope.go  +200/-0",
			"total: +320/-30",
			"excluded (see `dg task branch-diff --file <path>` to inspect): go.sum (+40/-12)",
		}
		for _, want := range wantLines {
			if !strings.Contains(out, want) {
				t.Fatalf("expected output to contain %q, got:\n%s", want, out)
			}
		}
	})

	t.Run("fetch failure does not abort, marks output", func(t *testing.T) {
		tm, gitBase, _ := newTaskSetup()
		gitBase.SetExecCommandResults(
			commands.ExecCommandResult("", "fatal: unreachable", fmt.Errorf("dial tcp: timeout")),
			commands.ExecCommandResult("feat/x\n", "", nil),
			commands.ExecCommandResult("origin/main\n", "", nil),
			commands.ExecCommandResult("abc123\n", "", nil),
			commands.ExecCommandResult("0\t0\n", "", nil),
			commands.ExecCommandResult("", "", nil),
			commands.ExecCommandResult("", "", nil),
			commands.ExecCommandResult("", "", nil),
		)

		out, err := tm.ReviewScope()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out, "(fetch failed — comparing against local refs)") {
			t.Fatalf("expected fetch-failed marker, got:\n%s", out)
		}
	})

	t.Run("on default branch returns sentinel without diffing", func(t *testing.T) {
		tm, gitBase, _ := newTaskSetup()
		gitBase.SetExecCommandResults(
			commands.ExecCommandResult("", "", nil),              // fetch origin
			commands.ExecCommandResult("main\n", "", nil),        // branch --show-current
			commands.ExecCommandResult("origin/main\n", "", nil), // symbolic-ref (default branch)
		)

		out, err := tm.ReviewScope()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := "On main — no branch to compare. Review uncommitted changes or name a target."
		if out != want {
			t.Fatalf("unexpected output: %q", out)
		}
		if got := gitBase.GetExecCommandCallCount(); got != 3 {
			t.Fatalf("expected exactly 3 git calls (no diffing), got %d", got)
		}
	})

	t.Run("detached HEAD returns sentinel without diffing", func(t *testing.T) {
		tm, gitBase, _ := newTaskSetup()
		gitBase.SetExecCommandResults(
			commands.ExecCommandResult("", "", nil), // fetch origin
			commands.ExecCommandResult(
				"",
				"",
				nil,
			), // branch --show-current -> empty (detached)
			commands.ExecCommandResult("abc1234\n", "", nil), // rev-parse --short HEAD
		)

		out, err := tm.ReviewScope()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := "Detached HEAD at abc1234 — no branch to compare. Check out a branch or name a target."
		if out != want {
			t.Fatalf("unexpected output: %q", out)
		}
		if got := gitBase.GetExecCommandCallCount(); got != 3 {
			t.Fatalf("expected exactly 3 git calls (no diffing), got %d", got)
		}
	})
}
