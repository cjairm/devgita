# Cycle: Worktree UX Improvements - Short Aliases + Interactive Selection

**Date:** 2026-04-22
**Estimated Duration:** ~2 hours
**Status:** Draft

---

## 1. Domain Context

The worktree feature (`dg worktree`) enables parallel AI-assisted development across multiple branches. Current commands are verbose:

- `dg worktree create <name>`
- `dg worktree list`
- `dg worktree remove <name>`

Users want shorter commands following Cobra's CLI-as-UI philosophy and interactive fzf selection for remove/jump operations.

**Reference:** `docs/plans/worktree-management.md`

---

## 2. Engineer Context

**Relevant files:**

- `cmd/worktree.go` - Cobra command definitions (lines 16-124)
- `internal/tooling/worktree/worktree.go` - WorktreeManager coordinator
- `internal/apps/tmux/tmux.go` - Tmux operations (line 122-141 for window ops)
- `configs/templates/devgita.zsh.tmpl` - Shell aliases using fzf-tmux (line 114)

**Key patterns:**

- Cobra `Aliases` field for command aliases
- `fzf-tmux -p 50% --reverse` pattern for popup selection (from devgita.zsh.tmpl)
- `BaseCommandExecutor.ExecCommand()` for command execution with stdin

**Commands to test:**

```bash
go build -o devgita main.go && ./devgita wt --help
go test ./internal/tooling/worktree/... ./cmd/...
```

---

## 3. Objective

Add short command aliases (`dg wt`, `wt c/l/rm`) and interactive fzf selection for `remove` and new `jump` commands, enabling faster worktree workflows.

---

## 4. Scope Boundary

### In Scope

- [x] Add `Aliases: []string{"wt"}` to worktreeCmd
- [x] Add aliases for subcommands: create (c, new), list (l, ls), remove (rm, r)
- [x] Add `dg wt jump` (aliases: j) - fzf select and switch to worktree window
- [x] Make `remove` use fzf selection when no argument provided
- [x] Update help text to show aliases

### Explicitly Out of Scope

- `--prompt` flag for initial OpenCode prompt (future cycle)
- `--from <branch>` flag (future cycle)
- `dg wt cleanup` bulk removal (future cycle)
- Uncommitted changes warnings (future cycle)

**Scope is locked.** If you discover something out of scope is needed, document it for a future cycle.

---

## 5. Implementation Plan

### File Changes

| Action | File Path                               | Description                                                        |
| ------ | --------------------------------------- | ------------------------------------------------------------------ |
| Modify | `cmd/worktree.go:16`                    | Add `Aliases: []string{"wt"}` to worktreeCmd                       |
| Modify | `cmd/worktree.go:33`                    | Add `Aliases: []string{"c", "new"}` to createCmd                   |
| Modify | `cmd/worktree.go:61`                    | Add `Aliases: []string{"l", "ls"}` to listCmd                      |
| Modify | `cmd/worktree.go:95`                    | Add `Aliases: []string{"rm", "r"}` to removeCmd, make arg optional |
| Create | `cmd/worktree.go:~125`                  | Add new `jumpCmd` with fzf selection                               |
| Modify | `internal/tooling/worktree/worktree.go` | Add `SelectWorktreeInteractively()` method                         |
| Modify | `internal/apps/tmux/tmux.go:~142`       | Add `SelectWindow(name)` method                                    |
| Modify | `internal/apps/fzf/fzf.go`              | Add `SelectFromList(items []string)` method                        |

### Step-by-Step

#### Step 1: Add aliases to worktreeCmd (parent command)

- Open `cmd/worktree.go`
- Line 16-31: Add `Aliases: []string{"wt"}` to worktreeCmd
- Update Long description to mention alias
- **Verify:** `go build && ./devgita wt --help` shows aliases

#### Step 2: Add aliases to create subcommand

- Line 33-59: Add `Aliases: []string{"c", "new"}` to worktreeCreateCmd
- **Verify:** `./devgita wt c --help` works

#### Step 3: Add aliases to list subcommand

- Line 61-93: Add `Aliases: []string{"l", "ls"}` to worktreeListCmd
- **Verify:** `./devgita wt l` works

#### Step 4: Add aliases to remove subcommand

- Line 95-117: Add `Aliases: []string{"rm", "r"}` to worktreeRemoveCmd
- Change `Args: cobra.ExactArgs(1)` to `Args: cobra.MaximumNArgs(1)`
- **Verify:** `./devgita wt rm --help` works

#### Step 5: Add SelectFromList to fzf app

- Open `internal/tooling/terminal/dev_tools/fzf/fzf.go`
- Add method to run fzf with stdin input and return selected line
- **Verify:** Unit test passes

#### Step 6: Add interactive selection to WorktreeManager

- Open `internal/tooling/worktree/worktree.go`
- Add `SelectWorktreeInteractively()` method that:
  1. Gets list of worktrees via `List()`
  2. Pipes names to fzf
  3. Returns selected name
- **Verify:** Unit test passes

#### Step 7: Update remove command to use interactive selection

- Modify removeCmd Run function:
  - If no args, call `SelectWorktreeInteractively()`
  - Use selected name for removal
- **Verify:** `./devgita wt rm` (no args) shows fzf picker

#### Step 8: Add SelectWindow to tmux app

- Open `internal/apps/tmux/tmux.go`
- Add `SelectWindow(name string) error` method
- **Verify:** Unit test passes

#### Step 9: Add jump command

- Add new jumpCmd to `cmd/worktree.go`:
  - `Use: "jump"`, `Aliases: []string{"j"}`
  - Call `SelectWorktreeInteractively()`
  - Switch to selected window with `SelectWindow()`
- Register in init()
- **Verify:** `./devgita wt j` shows fzf and switches

#### Step 10: Update help text

- Update Long descriptions to mention available aliases
- **Verify:** `./devgita wt --help` shows all aliases clearly

---

## 6. Verification Plan

### Automated Verification

```bash
go build -o devgita main.go
go test ./internal/tooling/worktree/...
go test ./cmd/...
go vet ./...
```

### Manual Verification

1. `./devgita wt --help` - Shows "Aliases: wt" and all subcommand aliases
2. `./devgita wt c test-feature` - Creates worktree (short form)
3. `./devgita wt l` - Lists worktrees (short form)
4. `./devgita wt j` - Opens fzf, selecting switches window
5. `./devgita wt rm` (no args) - Opens fzf, selecting removes worktree
6. `./devgita wt rm test-feature` - Still works with explicit name

### Regression Check

- Existing `dg worktree create/list/remove` commands still work
- fzf popup appears in tmux environment
- Non-tmux usage still works (graceful fallback)

---

## 7. Risks & Trade-offs

| Risk                             | Likelihood | Mitigation                                  |
| -------------------------------- | ---------- | ------------------------------------------- |
| fzf not installed                | Medium     | Check for fzf, show error with install hint |
| Not in tmux session (jump fails) | Medium     | Detect tmux env, show helpful error         |
| No worktrees to select           | Low        | Handle empty list gracefully                |

### Trade-offs Made

- Using fzf directly vs promptui: fzf provides better UX with preview/search, matches existing shell patterns
- Aliases vs separate commands: Aliases keep help organized while providing shortcuts

---

## 8. Cross-Model Review Notes

Space for review feedback when handing off between sessions or models:

- [ ] Root cause confirmed? (UX friction from verbose commands)
- [ ] All affected files identified?
- [ ] Verification steps are executable?
- [ ] Scope is appropriately bounded?

**Reviewer notes:**
(Fill in during review)
