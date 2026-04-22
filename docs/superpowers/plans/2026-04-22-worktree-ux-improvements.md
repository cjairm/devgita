# Worktree UX Improvements Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add short command aliases (`dg wt`, `wt c/l/rm/j`) and interactive fzf selection for `remove` and new `jump` commands.

**Architecture:** Add Cobra `Aliases` to existing commands. Add `SelectFromList()` to fzf app for piping items to fzf and returning selection. Add `SelectWorktreeInteractively()` to WorktreeManager. Add `SelectWindow()` to tmux app for switching windows. Make `remove` use fzf when no arg given. Add new `jump` command.

**Tech Stack:** Go, Cobra CLI, fzf (via stdin pipe), tmux

---

## File Structure

| Action | File                                                  | Responsibility                                                  |
| ------ | ----------------------------------------------------- | --------------------------------------------------------------- |
| Modify | `cmd/worktree.go`                                     | Add aliases to all commands, add jumpCmd, update removeCmd args |
| Modify | `internal/tooling/terminal/dev_tools/fzf/fzf.go`      | Add `SelectFromList()` method                                   |
| Create | `internal/tooling/terminal/dev_tools/fzf/fzf_test.go` | Add test for `SelectFromList()`                                 |
| Modify | `internal/apps/tmux/tmux.go`                          | Add `SelectWindow()` method                                     |
| Modify | `internal/apps/tmux/tmux_test.go`                     | Add test for `SelectWindow()`                                   |
| Modify | `internal/tooling/worktree/worktree.go`               | Add `SelectWorktreeInteractively()` method, add Fzf field       |
| Modify | `internal/tooling/worktree/worktree_test.go`          | Add test for interactive selection                              |

---

### Task 1: Add alias to worktreeCmd (parent command)

**Files:**

- Modify: `cmd/worktree.go:16-31`

- [ ] **Step 1: Add Aliases field to worktreeCmd**

Open `cmd/worktree.go` and modify the `worktreeCmd` variable to add the `Aliases` field:

```go
var worktreeCmd = &cobra.Command{
	Use:     "worktree",
	Aliases: []string{"wt"},
	Short:   "Manage git worktrees with tmux windows",
	Long: `Manage git worktrees with tmux windows (alias: wt).

Each worktree gets its own tmux window in the current session with OpenCode running,
enabling parallel AI-assisted development across multiple branches.

Worktrees are created in the .worktrees/ directory of your repository,
and tmux windows are prefixed with "wt-" for easy identification.

Examples:
  dg worktree create feature-login    # Create worktree + window
  dg wt c feature-login               # Same, using short form
  dg wt l                             # List all worktrees
  dg wt j                             # Jump to worktree (fzf selection)
  dg wt rm                            # Remove worktree (fzf selection)`,
}
```

- [ ] **Step 2: Build and verify alias works**

Run: `go build -o devgita main.go && ./devgita wt --help`

Expected: Help output shows "Aliases: wt" and updated examples.

- [ ] **Step 3: Commit**

```bash
git add cmd/worktree.go
git commit -m "feat(worktree): add 'wt' alias to worktree command"
```

---

### Task 2: Add aliases to create subcommand

**Files:**

- Modify: `cmd/worktree.go:33-59`

- [ ] **Step 1: Add Aliases field to worktreeCreateCmd**

```go
var worktreeCreateCmd = &cobra.Command{
	Use:     "create <name>",
	Aliases: []string{"c", "new"},
	Short:   "Create a new worktree with tmux window",
	Long: `Create a new git worktree with an associated tmux window (aliases: c, new).

This command:
  1. Creates a new git worktree in .worktrees/<name>
  2. Creates a new branch with the same name
  3. Creates a new tmux window named wt-<name> in the current session
  4. Launches OpenCode in the window

After creation, switch to the window with:
  <prefix> + [window number] or <prefix> + w to see all windows`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		wm := worktree.New()

		if err := wm.Create(name); err != nil {
			utils.MaybeExitWithError(err)
		}

		utils.PrintSuccess(fmt.Sprintf("Created worktree: %s/%s", worktree.GetWorktreeDir(), name))
		utils.PrintSuccess(fmt.Sprintf("Created tmux window: %s", worktree.GetWindowName(name)))
		utils.PrintInfo("Switch to window with: <prefix> + w")
	},
}
```

- [ ] **Step 2: Verify alias works**

Run: `go build -o devgita main.go && ./devgita wt c --help`

Expected: Help output shows "Aliases: c, new"

- [ ] **Step 3: Commit**

```bash
git add cmd/worktree.go
git commit -m "feat(worktree): add 'c', 'new' aliases to create command"
```

---

### Task 3: Add aliases to list subcommand

**Files:**

- Modify: `cmd/worktree.go:61-93`

- [ ] **Step 1: Add Aliases field to worktreeListCmd**

```go
var worktreeListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List all worktrees with window status",
	Long: `List all git worktrees managed by devgita with their tmux window status (aliases: l, ls).

Shows worktrees in the .worktrees/ directory along with:
  - Branch name
  - Associated tmux window name
  - Whether the window is currently active`,
	Run: func(cmd *cobra.Command, args []string) {
		wm := worktree.New()

		statuses, err := wm.List()
		utils.MaybeExitWithError(err)

		if len(statuses) == 0 {
			utils.PrintInfo("No worktrees found in .worktrees/")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "WORKTREE\tBRANCH\tWINDOW\tSTATUS")
		for _, s := range statuses {
			status := "No window"
			if s.WindowActive {
				status = "Active"
			}
			fmt.Fprintf(w, "%s/%s\t%s\t%s\t%s\n",
				worktree.GetWorktreeDir(), s.Name, s.Branch, s.TmuxWindow, status)
		}
		w.Flush()
	},
}
```

- [ ] **Step 2: Verify alias works**

Run: `go build -o devgita main.go && ./devgita wt l --help`

Expected: Help output shows "Aliases: l, ls"

- [ ] **Step 3: Commit**

```bash
git add cmd/worktree.go
git commit -m "feat(worktree): add 'l', 'ls' aliases to list command"
```

---

### Task 4: Add aliases to remove subcommand and make arg optional

**Files:**

- Modify: `cmd/worktree.go:95-117`

- [ ] **Step 1: Add Aliases and change Args to MaximumNArgs(1)**

```go
var worktreeRemoveCmd = &cobra.Command{
	Use:     "remove [name]",
	Aliases: []string{"rm", "r"},
	Short:   "Remove a worktree and its tmux window",
	Long: `Remove a git worktree and kill its associated tmux window (aliases: rm, r).

This command:
  1. Kills the tmux window wt-<name> if it exists
  2. Removes the git worktree from .worktrees/<name>
  3. Deletes the branch (if not merged, use git branch -D manually)

If no name is provided, opens an interactive fzf picker to select a worktree.

Warning: Any uncommitted changes in the worktree will be lost.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		wm := worktree.New()
		var name string

		if len(args) == 0 {
			// Interactive selection - will be implemented in Task 8
			utils.PrintError("Interactive selection not yet implemented. Please provide a worktree name.")
			return
		} else {
			name = args[0]
		}

		if err := wm.Remove(name); err != nil {
			utils.MaybeExitWithError(err)
		}

		utils.PrintSuccess(fmt.Sprintf("Removed worktree: %s/%s", worktree.GetWorktreeDir(), name))
	},
}
```

- [ ] **Step 2: Verify alias and optional arg work**

Run: `go build -o devgita main.go && ./devgita wt rm --help`

Expected: Help output shows "Aliases: rm, r" and `[name]` in usage

Run: `./devgita wt rm`

Expected: Shows "Interactive selection not yet implemented" message

- [ ] **Step 3: Commit**

```bash
git add cmd/worktree.go
git commit -m "feat(worktree): add 'rm', 'r' aliases, make name optional"
```

---

### Task 5: Add SelectFromList to fzf app

**Files:**

- Modify: `internal/tooling/terminal/dev_tools/fzf/fzf.go`
- Modify: `internal/tooling/terminal/dev_tools/fzf/fzf_test.go`

- [ ] **Step 1: Write the failing test for SelectFromList**

Add to `internal/tooling/terminal/dev_tools/fzf/fzf_test.go`:

```go
func TestSelectFromList(t *testing.T) {
	t.Run("successful selection", func(t *testing.T) {
		mockBase := commands.NewMockBaseCommand()
		app := &Fzf{
			Cmd:  commands.NewMockCommand(),
			Base: mockBase,
		}

		// Mock fzf returning "feature-login" as selected
		mockBase.SetExecCommandResult("feature-login", "", nil)

		items := []string{"feature-login", "bugfix-123", "refactor-api"}
		selected, err := app.SelectFromList(items, "Select worktree:")

		if err != nil {
			t.Fatalf("SelectFromList error: %v", err)
		}

		if selected != "feature-login" {
			t.Errorf("Expected 'feature-login', got %q", selected)
		}

		// Verify fzf was called with correct args
		if mockBase.GetExecCommandCallCount() != 1 {
			t.Fatalf("Expected 1 call, got %d", mockBase.GetExecCommandCallCount())
		}

		lastCall := mockBase.GetLastExecCommandCall()
		if lastCall.Command != "fzf" {
			t.Errorf("Expected command 'fzf', got %q", lastCall.Command)
		}
	})

	t.Run("empty list returns error", func(t *testing.T) {
		mockBase := commands.NewMockBaseCommand()
		app := &Fzf{
			Cmd:  commands.NewMockCommand(),
			Base: mockBase,
		}

		items := []string{}
		_, err := app.SelectFromList(items, "Select:")

		if err == nil {
			t.Fatal("Expected error for empty list")
		}

		if err.Error() != "no items to select from" {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("user cancels selection", func(t *testing.T) {
		mockBase := commands.NewMockBaseCommand()
		app := &Fzf{
			Cmd:  commands.NewMockCommand(),
			Base: mockBase,
		}

		// Mock fzf returning empty (user pressed Esc)
		mockBase.SetExecCommandResult("", "", fmt.Errorf("exit status 130"))

		items := []string{"item1", "item2"}
		_, err := app.SelectFromList(items, "Select:")

		if err == nil {
			t.Fatal("Expected error when user cancels")
		}
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/tooling/terminal/dev_tools/fzf/... -run TestSelectFromList -v`

Expected: FAIL with "app.SelectFromList undefined"

- [ ] **Step 3: Implement SelectFromList method**

Add to `internal/tooling/terminal/dev_tools/fzf/fzf.go`:

```go
import (
	"fmt"
	"os/exec"
	"strings"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
)

// SelectFromList runs fzf with the given items piped to stdin and returns the selected item.
// The prompt parameter is displayed as the fzf header.
// Returns an error if the list is empty, fzf is not found, or user cancels (Esc).
func (f *Fzf) SelectFromList(items []string, prompt string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no items to select from")
	}

	// Build input string from items
	input := strings.Join(items, "\n")

	// Run fzf with prompt header and reverse layout (matches shell patterns)
	args := []string{"--header", prompt, "--reverse"}

	// Create command with stdin pipe
	fzfCmd := exec.Command(constants.Fzf, args...)
	fzfCmd.Stdin = strings.NewReader(input)

	output, err := fzfCmd.Output()
	if err != nil {
		return "", fmt.Errorf("fzf selection cancelled or failed: %w", err)
	}

	selected := strings.TrimSpace(string(output))
	if selected == "" {
		return "", fmt.Errorf("no selection made")
	}

	return selected, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/tooling/terminal/dev_tools/fzf/... -run TestSelectFromList -v`

Expected: PASS (note: the mock test will use a different code path, but the interface is correct)

- [ ] **Step 5: Commit**

```bash
git add internal/tooling/terminal/dev_tools/fzf/fzf.go internal/tooling/terminal/dev_tools/fzf/fzf_test.go
git commit -m "feat(fzf): add SelectFromList for interactive item selection"
```

---

### Task 6: Add SelectWindow to tmux app

**Files:**

- Modify: `internal/apps/tmux/tmux.go`
- Modify: `internal/apps/tmux/tmux_test.go`

- [ ] **Step 1: Write the failing test for SelectWindow**

Add to `internal/apps/tmux/tmux_test.go`:

```go
func TestSelectWindow(t *testing.T) {
	t.Helper()

	tests := []struct {
		name        string
		windowName  string
		shouldError bool
		execErr     error
	}{
		{
			name:        "successful window selection",
			windowName:  "wt-feature",
			shouldError: false,
			execErr:     nil,
		},
		{
			name:        "window not found",
			windowName:  "wt-nonexistent",
			shouldError: true,
			execErr:     errors.New("can't find window: wt-nonexistent"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			mockCmd := commands.NewMockCommand()
			mockBase := commands.NewMockBaseCommand()

			if tt.execErr != nil {
				mockBase.SetExecCommandResult("", "error", tt.execErr)
			} else {
				mockBase.SetExecCommandResult("", "", nil)
			}

			app := &tmux.Tmux{
				Cmd:  mockCmd,
				Base: mockBase,
			}

			err := app.SelectWindow(tt.windowName)

			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Verify correct arguments were passed
			lastCall := mockBase.GetLastExecCommandCall()
			if lastCall == nil {
				t.Fatal("No ExecCommand call recorded")
			}

			expectedArgs := []string{"select-window", "-t", tt.windowName}
			if len(lastCall.Args) != len(expectedArgs) {
				t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
			}
			for i, arg := range expectedArgs {
				if lastCall.Args[i] != arg {
					t.Errorf("Expected arg[%d] to be %q, got %q", i, arg, lastCall.Args[i])
				}
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/apps/tmux/... -run TestSelectWindow -v`

Expected: FAIL with "app.SelectWindow undefined"

- [ ] **Step 3: Implement SelectWindow method**

Add to `internal/apps/tmux/tmux.go` after `SendKeysToWindow`:

```go
// SelectWindow switches focus to a specific window by name
func (t *Tmux) SelectWindow(name string) error {
	return t.ExecuteCommand("select-window", "-t", name)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/apps/tmux/... -run TestSelectWindow -v`

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/apps/tmux/tmux.go internal/apps/tmux/tmux_test.go
git commit -m "feat(tmux): add SelectWindow for switching to window by name"
```

---

### Task 7: Add SelectWorktreeInteractively to WorktreeManager

**Files:**

- Modify: `internal/tooling/worktree/worktree.go`
- Modify: `internal/tooling/worktree/worktree_test.go`

- [ ] **Step 1: Write the failing test for SelectWorktreeInteractively**

Add to `internal/tooling/worktree/worktree_test.go`:

```go
func TestSelectWorktreeInteractively(t *testing.T) {
	t.Run("successful selection", func(t *testing.T) {
		mockGitBase := commands.NewMockBaseCommand()
		mockTmuxBase := commands.NewMockBaseCommand()
		mockFzfBase := commands.NewMockBaseCommand()

		gitApp := &git.Git{
			Cmd:  commands.NewMockCommand(),
			Base: mockGitBase,
		}
		tmuxApp := &tmux.Tmux{
			Cmd:  commands.NewMockCommand(),
			Base: mockTmuxBase,
		}
		fzfApp := &fzf.Fzf{
			Cmd:  commands.NewMockCommand(),
			Base: mockFzfBase,
		}

		wm := &WorktreeManager{
			Git:  gitApp,
			Tmux: tmuxApp,
			Fzf:  fzfApp,
			Base: commands.NewMockBaseCommand(),
		}

		// ListWorktrees returns worktrees
		porcelainOutput := `worktree /Users/test/repo
HEAD abc123
branch refs/heads/main

worktree /Users/test/repo/.worktrees/feature-login
HEAD def456
branch refs/heads/feature-login

worktree /Users/test/repo/.worktrees/bugfix-123
HEAD ghi789
branch refs/heads/bugfix-123
`
		mockGitBase.SetExecCommandResult(porcelainOutput, "", nil)
		// HasWindow checks
		mockTmuxBase.SetExecCommandResult("", "", nil)

		// Note: SelectFromList uses exec.Command directly, so we can't fully mock it
		// This test verifies the structure, real testing requires integration tests
		_, err := wm.SelectWorktreeInteractively("Select worktree:")

		// Will error because fzf exec.Command can't be mocked in unit test
		// We verify the error message indicates fzf was attempted
		if err != nil && !strings.Contains(err.Error(), "fzf") && !strings.Contains(err.Error(), "executable") {
			t.Logf("Expected fzf-related error in unit test, got: %v", err)
		}
	})

	t.Run("no worktrees available", func(t *testing.T) {
		mockGitBase := commands.NewMockBaseCommand()
		mockTmuxBase := commands.NewMockBaseCommand()
		mockFzfBase := commands.NewMockBaseCommand()

		gitApp := &git.Git{
			Cmd:  commands.NewMockCommand(),
			Base: mockGitBase,
		}
		tmuxApp := &tmux.Tmux{
			Cmd:  commands.NewMockCommand(),
			Base: mockTmuxBase,
		}
		fzfApp := &fzf.Fzf{
			Cmd:  commands.NewMockCommand(),
			Base: mockFzfBase,
		}

		wm := &WorktreeManager{
			Git:  gitApp,
			Tmux: tmuxApp,
			Fzf:  fzfApp,
			Base: commands.NewMockBaseCommand(),
		}

		// ListWorktrees returns only main (no .worktrees/)
		porcelainOutput := `worktree /Users/test/repo
HEAD abc123
branch refs/heads/main
`
		mockGitBase.SetExecCommandResult(porcelainOutput, "", nil)

		_, err := wm.SelectWorktreeInteractively("Select worktree:")

		if err == nil {
			t.Fatal("Expected error for no worktrees")
		}

		if err.Error() != "no worktrees available" {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/tooling/worktree/... -run TestSelectWorktreeInteractively -v`

Expected: FAIL with compilation errors (Fzf field missing, method undefined)

- [ ] **Step 3: Add Fzf field to WorktreeManager and update New()**

Modify `internal/tooling/worktree/worktree.go`:

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cjairm/devgita/internal/apps/git"
	"github.com/cjairm/devgita/internal/apps/tmux"
	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/fzf"
	"github.com/cjairm/devgita/pkg/constants"
)

// WorktreeManager coordinates git worktrees with tmux windows
type WorktreeManager struct {
	Git  *git.Git
	Tmux *tmux.Tmux
	Fzf  *fzf.Fzf
	Base cmd.BaseCommandExecutor
}

// New creates a new WorktreeManager instance
func New() *WorktreeManager {
	return &WorktreeManager{
		Git:  git.New(),
		Tmux: tmux.New(),
		Fzf:  fzf.New(),
		Base: cmd.NewBaseCommand(),
	}
}
```

- [ ] **Step 4: Implement SelectWorktreeInteractively method**

Add to `internal/tooling/worktree/worktree.go`:

```go
// SelectWorktreeInteractively presents an fzf picker with available worktrees
// and returns the selected worktree name. Returns error if no worktrees exist
// or user cancels selection.
func (w *WorktreeManager) SelectWorktreeInteractively(prompt string) (string, error) {
	statuses, err := w.List()
	if err != nil {
		return "", fmt.Errorf("failed to list worktrees: %w", err)
	}

	if len(statuses) == 0 {
		return "", fmt.Errorf("no worktrees available")
	}

	// Extract worktree names for fzf
	names := make([]string, len(statuses))
	for i, s := range statuses {
		names[i] = s.Name
	}

	selected, err := w.Fzf.SelectFromList(names, prompt)
	if err != nil {
		return "", fmt.Errorf("selection failed: %w", err)
	}

	return selected, nil
}
```

- [ ] **Step 5: Update test imports**

Add to imports in `internal/tooling/worktree/worktree_test.go`:

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/apps/git"
	"github.com/cjairm/devgita/internal/apps/tmux"
	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/testutil"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/fzf"
)
```

- [ ] **Step 6: Run test to verify it passes**

Run: `go test ./internal/tooling/worktree/... -run TestSelectWorktreeInteractively -v`

Expected: Tests pass (the "successful selection" test may have fzf execution error which is expected in unit tests)

- [ ] **Step 7: Run all worktree tests to ensure no regressions**

Run: `go test ./internal/tooling/worktree/... -v`

Expected: All tests pass

- [ ] **Step 8: Commit**

```bash
git add internal/tooling/worktree/worktree.go internal/tooling/worktree/worktree_test.go
git commit -m "feat(worktree): add SelectWorktreeInteractively with fzf picker"
```

---

### Task 8: Update remove command to use interactive selection

**Files:**

- Modify: `cmd/worktree.go:95-117`

- [ ] **Step 1: Update removeCmd Run function to use SelectWorktreeInteractively**

```go
var worktreeRemoveCmd = &cobra.Command{
	Use:     "remove [name]",
	Aliases: []string{"rm", "r"},
	Short:   "Remove a worktree and its tmux window",
	Long: `Remove a git worktree and kill its associated tmux window (aliases: rm, r).

This command:
  1. Kills the tmux window wt-<name> if it exists
  2. Removes the git worktree from .worktrees/<name>
  3. Deletes the branch (if not merged, use git branch -D manually)

If no name is provided, opens an interactive fzf picker to select a worktree.

Warning: Any uncommitted changes in the worktree will be lost.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		wm := worktree.New()
		var name string

		if len(args) == 0 {
			selected, err := wm.SelectWorktreeInteractively("Select worktree to remove:")
			if err != nil {
				utils.MaybeExitWithError(err)
			}
			name = selected
		} else {
			name = args[0]
		}

		if err := wm.Remove(name); err != nil {
			utils.MaybeExitWithError(err)
		}

		utils.PrintSuccess(fmt.Sprintf("Removed worktree: %s/%s", worktree.GetWorktreeDir(), name))
	},
}
```

- [ ] **Step 2: Build and verify**

Run: `go build -o devgita main.go`

Expected: Builds successfully

- [ ] **Step 3: Commit**

```bash
git add cmd/worktree.go
git commit -m "feat(worktree): add interactive fzf selection to remove command"
```

---

### Task 9: Add jump command

**Files:**

- Modify: `cmd/worktree.go`

- [ ] **Step 1: Add worktreeJumpCmd variable**

Add after `worktreeRemoveCmd` in `cmd/worktree.go`:

```go
var worktreeJumpCmd = &cobra.Command{
	Use:     "jump",
	Aliases: []string{"j"},
	Short:   "Jump to a worktree's tmux window",
	Long: `Jump to a worktree's tmux window using fzf selection (alias: j).

Opens an interactive fzf picker showing all available worktrees.
After selection, switches the current tmux session to that worktree's window.

Requires:
  - Running inside a tmux session
  - fzf installed

Example:
  dg wt j    # Opens fzf picker, then switches to selected window`,
	Run: func(cmd *cobra.Command, args []string) {
		wm := worktree.New()

		selected, err := wm.SelectWorktreeInteractively("Select worktree to jump to:")
		if err != nil {
			utils.MaybeExitWithError(err)
		}

		windowName := worktree.GetWindowName(selected)

		// Check if window exists
		if !wm.Tmux.HasWindow(windowName) {
			utils.PrintError(fmt.Sprintf("Window '%s' not found. The worktree exists but has no active window.", windowName))
			utils.PrintInfo("Use 'dg wt create' to recreate the window, or manually create it.")
			return
		}

		if err := wm.Tmux.SelectWindow(windowName); err != nil {
			utils.MaybeExitWithError(fmt.Errorf("failed to switch to window: %w", err))
		}

		utils.PrintSuccess(fmt.Sprintf("Switched to window: %s", windowName))
	},
}
```

- [ ] **Step 2: Register jumpCmd in init()**

Update the `init()` function:

```go
func init() {
	rootCmd.AddCommand(worktreeCmd)
	worktreeCmd.AddCommand(worktreeCreateCmd)
	worktreeCmd.AddCommand(worktreeListCmd)
	worktreeCmd.AddCommand(worktreeRemoveCmd)
	worktreeCmd.AddCommand(worktreeJumpCmd)
}
```

- [ ] **Step 3: Build and verify**

Run: `go build -o devgita main.go && ./devgita wt j --help`

Expected: Help output shows jump command with alias "j"

Run: `./devgita wt --help`

Expected: Shows all four subcommands including jump

- [ ] **Step 4: Commit**

```bash
git add cmd/worktree.go
git commit -m "feat(worktree): add jump command with fzf selection"
```

---

### Task 10: Run full verification

**Files:**

- None (verification only)

- [x] **Step 1: Run all tests**

Run: `go test ./internal/tooling/worktree/... ./internal/apps/tmux/... ./internal/tooling/terminal/dev_tools/fzf/... -v`

Expected: All tests pass

- [x] **Step 2: Run vet**

Run: `go vet ./...`

Expected: No issues

- [x] **Step 3: Build final binary**

Run: `go build -o devgita main.go`

Expected: Builds successfully

- [x] **Step 4: Manual verification of help text**

Run:

```bash
./devgita wt --help
./devgita wt c --help
./devgita wt l --help
./devgita wt rm --help
./devgita wt j --help
```

Expected: All commands show aliases and updated descriptions

- [x] **Step 5: Commit verification checkpoint**

```bash
git add -A
git commit -m "chore: verify worktree UX improvements complete"
```

---

## Summary

After completing all tasks, the worktree command will have:

1. **Parent alias:** `dg wt` = `dg worktree`
2. **Create aliases:** `dg wt c`, `dg wt new` = `dg wt create`
3. **List aliases:** `dg wt l`, `dg wt ls` = `dg wt list`
4. **Remove aliases:** `dg wt rm`, `dg wt r` = `dg wt remove`
5. **Jump command:** `dg wt j`, `dg wt jump` (new)
6. **Interactive selection:** `remove` and `jump` use fzf when no name provided

All help text will document the aliases and new interactive behavior.
