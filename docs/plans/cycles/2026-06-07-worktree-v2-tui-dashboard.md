# Cycle: Worktree v2 — TUI Dashboard (Phase 1: core two-pane dashboard)

**Date:** 2026-06-07
**Estimated Duration:** ~3 hours
**Status:** Draft

---

## 1. Domain Context

`dg worktree` (`wt`) pairs a git worktree with a tmux window running an AI coder
(Claude Code focus). Today, navigation is an **fzf popup** (`dg wt j`, bound to
`prefix + u` in `configs/tmux/tmux.conf:113`) that lists worktrees and supports
jump (Enter), delete (ctrl-d), repair (ctrl-r).

This is the feature the maintainer uses most. **v2** turns the fzf popup into a
persistent full-screen **TUI dashboard** inspired by claude-squad / vibe-kanban /
bernstein, but deliberately minimal: it should feel like **NERDTree** (familiar
keys, no new mental model) and stay **tmux-backed** (each worktree's Claude keeps
running in its own tmux window; the TUI is a dashboard over those windows).

The full v2 is larger than one cycle, so it is **phased**:

- **Phase 1 (this cycle):** Core two-pane dashboard — left NERDTree-style worktree
  tree grouped by repo; right underline tabs Preview (`tmux capture-pane`) + Diff
  (`git diff`); navigate, attach, filter, quit.
- **Phase 2 (future):** Inline actions — new / delete / repair worktree, open in nvim.
- **Phase 3 (future):** Completion notifications via Claude Code Stop hook writing a
  marker file the TUI watches, plus optional desktop notification.

Design reference: the wireframe set the maintainer approved selects **Layout A
(classic two-pane)**, **tree grouping T1 (by repo)**, **right pane R1 (underline
tabs)**, and **keybinding model K1 (persistent hint bar)**.

---

## 2. Engineer Context

**Relevant existing files:**

| File                                    | Role                                                                  |
| --------------------------------------- | --------------------------------------------------------------------- |
| `cmd/worktree.go`                       | Cobra subcommands; add `dg wt ui` here, register in `init()`          |
| `internal/tooling/worktree/worktree.go` | `WorktreeManager`; `List()` returns `[]WorktreeStatus`                |
| `internal/apps/tmux/tmux.go`            | tmux primitives; `WindowSession`, `SwitchToWindow` exist; add capture |
| `internal/apps/git/git.go`              | git helpers; `IsWorktreeDirty` exists; add `Diff`/`DiffStat`          |
| `internal/commands/base.go:191`         | `BaseCommand.ExecCommand(CommandParams) (stdout, stderr, error)`      |
| `internal/commands/mock.go:204`         | `MockBaseCommand.ExecCommand` for tests                               |
| `configs/tmux/tmux.conf:113`            | `prefix + u` popup binding (currently `devgita wt j`)                 |
| `go.mod`                                | Add Bubble Tea + Lipgloss deps                                        |

**Key existing types/functions to reuse (do NOT reimplement):**

- `worktree.WorktreeStatus{ Name, Path, Branch, TmuxWindow, WindowActive, Repo }`
- `(*WorktreeManager).List() ([]WorktreeStatus, error)` — walks the centralized
  base path, two levels deep, across all repos.
- `worktree.GetWindowName(name) string` — `wt-<flat-name>`.
- `(*Tmux).WindowSession(window) (session string, ok bool)` — finds the session
  hosting a window across all sessions.
- `(*Tmux).SwitchToWindow(session, window) error` — switch-client + select-window.
- `(*Git).IsWorktreeDirty(path) (bool, error)` — `git -C path status --porcelain`.

**Testing patterns (CLAUDE.md §6, docs/guides/testing-patterns.md):**

- `func init() { testutil.InitLogger() }` at top of each test file.
- Use `testutil`/`MockBaseCommand`; **never execute real commands** in tests.
- `testutil.VerifyNoRealCommands(t, mockApp.Base)` where applicable.
- For the TUI model, unit-test the pure `Update`/state transitions with synthetic
  `tea.KeyMsg` and injected data — **do not** spawn a real `tea.Program` or tmux.

**Commands:**

- Build: `go build -o devgita main.go`
- Test: `go test ./...`
- Lint: `make lint` (`go fmt ./...` + `go vet ./...`)

---

## 3. Objective

Ship `dg wt ui`: a full-screen, tmux-backed Bubble Tea dashboard with a
NERDTree-style worktree tree (grouped by repo) on the left and Preview/Diff tabs
on the right, supporting `j/k` move, `h/l`+`z` fold, `Enter` attach, `Tab` switch
tab, `/` filter, and `q` quit — replacing the `prefix + u` fzf popup while leaving
all existing `dg wt` subcommands untouched.

---

## 4. Scope Boundary

### In Scope

- [ ] Add `charmbracelet/bubbletea` + `charmbracelet/lipgloss` to `go.mod`/`go.sum`.
- [ ] New `internal/tui/worktree/` package: Bubble Tea model, tree rendering,
      styles, and a `Run()` entrypoint.
- [ ] Left pane: worktree tree grouped by repo (T1), status glyphs
      (`●` running/has-window, `○` no session) and dirty `±N` from `IsWorktreeDirty`.
- [ ] Tree nav: `j/k` move, `h/l` collapse/expand current repo node, `z` toggle
      collapse-all, wrap-around within visible rows.
- [ ] Right pane: underline tabs `Preview` and `Diff`; `Tab` switches; Preview shows
      `tmux capture-pane` of the selected worktree's window; Diff shows `git diff`.
- [ ] `Enter` = attach: when inside tmux, `SwitchToWindow` to the worktree's window
      and quit the TUI; when window missing or not in tmux, show an inline status
      message (no crash). Repair is deferred to Phase 2.
- [ ] `/` filter: type to filter worktrees by `repo/name` substring; `Esc` clears.
- [ ] Persistent hint bar (K1) at the bottom showing the active keys.
- [ ] Selected-row refresh: re-capture Preview on a timer tick (~1.5s) and recompute
      Diff on selection change.
- [ ] `tmux.CapturePane(session, window) (string, error)` helper + test.
- [ ] `git.Diff(path) (string, error)` and `git.DiffStat(path) (files, added, removed int, error)` helpers + tests.
- [ ] `dg wt ui` Cobra subcommand (aliases `dash`, `dashboard`) registered in `cmd/worktree.go`.
- [ ] Update `configs/tmux/tmux.conf:113` to launch `devgita wt ui` in the popup.
- [ ] Unit tests for the model's `Update` transitions (nav, fold, tab switch, filter).

### Explicitly Out of Scope

- Inline **new / delete / repair** worktree actions (Phase 2).
- **Open in nvim** (Phase 2).
- **Stop-hook completion notifications** + desktop notify + `◆ needs review` glyph
  (Phase 3).
- Alternate layouts/groupings (B–E, T2/T3, R2/R3), command palette, which-key popup.
- Syntax-highlighted / scrollable diff with per-hunk navigation (Phase 1 renders
  plain `git diff` text in a scrollable viewport only).
- PR/merge status via `gh` (stretch, later).
- Removing or changing existing `dg wt` subcommands or `dg wt j` behavior.

**Scope is locked.** Anything discovered beyond this is logged for a later phase.

---

## 5. Implementation Plan

### File Changes

| Action | File Path                                      | Description                                                          |
| ------ | ---------------------------------------------- | -------------------------------------------------------------------- |
| Modify | `go.mod` / `go.sum`                            | Add bubbletea + lipgloss (`go get`)                                  |
| Create | `internal/tui/worktree/model.go`               | Bubble Tea `Model`, `Init`/`Update`/`View`, key handling, tick       |
| Create | `internal/tui/worktree/tree.go`                | Tree node type, build-from-statuses, flatten-to-visible, render rows |
| Create | `internal/tui/worktree/styles.go`              | Lipgloss styles (panes, tabs, selection, status glyphs, hint bar)    |
| Create | `internal/tui/worktree/run.go`                 | `Run()` — wires `WorktreeManager`/`Tmux`/`Git`, runs `tea.Program`   |
| Create | `internal/tui/worktree/model_test.go`          | Unit tests for `Update` transitions with synthetic msgs              |
| Modify | `internal/apps/tmux/tmux.go` (after ~line 244) | Add `CapturePane(session, window) (string, error)`                   |
| Modify | `internal/apps/tmux/tmux_test.go`              | Test `CapturePane` with mock                                         |
| Modify | `internal/apps/git/git.go` (after ~line 356)   | Add `Diff(path)` and `DiffStat(path)`                                |
| Modify | `internal/apps/git/git_test.go`                | Test `Diff`/`DiffStat` with mock                                     |
| Modify | `cmd/worktree.go` (~line 282 vars, ~287 init)  | Add `worktreeUICmd`, register it                                     |
| Modify | `configs/tmux/tmux.conf:113`                   | `devgita wt j` → `devgita wt ui`                                     |

### Data flow

```
Run() → WorktreeManager.List() ──► []WorktreeStatus
                                     │
        build tree (group by repo) ─┘
                                     │ on selection / tick
        selected worktree ──► Tmux.CapturePane(session, window)  → Preview text
                          └─► Git.Diff(path) / Git.DiffStat(path) → Diff text + ±N
        Enter ──► Tmux.WindowSession(window) → SwitchToWindow(session, window) → quit
```

### Step-by-Step

#### Step 1: Add dependencies

- Run `go get github.com/charmbracelet/bubbletea@latest github.com/charmbracelet/lipgloss@latest`.
- Run `go mod tidy`.
- Verify: `go build ./...` succeeds; `go.mod` lists both modules.

#### Step 2: Add `tmux.CapturePane`

- In `internal/apps/tmux/tmux.go`, add:
  ```go
  // CapturePane returns the visible content of a window's active pane, including
  // ANSI color escapes (-e). target is "session:window".
  func (t *Tmux) CapturePane(session, window string) (string, error) {
      target := session + ":" + window
      execCommand := cmd.CommandParams{
          Command: constants.Tmux,
          Args:    []string{"capture-pane", "-p", "-e", "-t", target},
      }
      stdout, stderr, err := t.Base.ExecCommand(execCommand)
      if err != nil {
          if stderr != "" {
              return "", fmt.Errorf("capture-pane: %s", stderr)
          }
          return "", fmt.Errorf("failed to capture pane %s: %w", target, err)
      }
      return stdout, nil
  }
  ```
- Verify: compiles; mirrors `WindowSession`'s stdout/stderr handling.

#### Step 3: Add `git.Diff` and `git.DiffStat`

- In `internal/apps/git/git.go`, add `Diff(path)` running `git -C path diff` and
  returning stdout (reuse the `ExecCommand` + stderr-surfacing pattern from
  `IsWorktreeDirty`).
- Add `DiffStat(path) (files, added, removed int, err error)` running
  `git -C path diff --numstat` and summing columns (numstat lines:
  `<added>\t<removed>\t<file>`; binary files report `-`/`-` → count as 0 lines but
  still increment `files`).
- Verify: compiles.

#### Step 4: Tree model (`tree.go`)

- Define a flat, render-friendly representation:

  ```go
  type rowKind int
  const ( rowRepo rowKind = iota; rowWorktree )

  type row struct {
      kind     rowKind
      repo     string
      status   worktree.WorktreeStatus // zero for repo rows
      depth    int
  }
  ```

- `buildRows(statuses []worktree.WorktreeStatus, collapsed map[string]bool, filter string) []row`:
  group by `Repo` (stable alpha sort), emit a repo header row then its worktree
  rows unless `collapsed[repo]`; apply case-insensitive `repo/name` substring filter
  (a repo row is kept if any child matches).
- `glyph(s worktree.WorktreeStatus) string`: `●` when `WindowActive`, else `○`.
- Verify: pure functions, covered by Step 9 tests.

#### Step 5: Styles (`styles.go`)

- Lipgloss styles: left/right pane borders, active/inactive tab, selected row
  (reverse/﻿highlight), repo header (bold), dirty `±N` (yellow), added/removed
  counts (green/red), dim hint bar. Keep a single `Styles` struct constructed once.
- Verify: compiles.

#### Step 6: Model core (`model.go`)

- `Model` fields: `mgr *worktree.WorktreeManager`, `tmux *tmux.Tmux`, `git *git.Git`,
  `statuses []WorktreeStatus`, `rows []row`, `cursor int`, `collapsed map[string]bool`,
  `allCollapsed bool`, `activeTab` (preview|diff), `preview string`, `diff string`,
  `diffAdded/Removed/Files int`, `filtering bool`, `filter string`, `status string`
  (inline message), `width/height int`, `styles Styles`.
- `Init()` returns a batch: initial `loadCmd` + `tickCmd` (1.5s).
- Messages: `statusesMsg`, `previewMsg`, `diffMsg`, `tickMsg`, plus `tea.KeyMsg`,
  `tea.WindowSizeMsg`.
- `Update` handling:
  - `WindowSizeMsg`: store dims, recompute pane widths.
  - `tickMsg`: re-issue `capturePreviewCmd(selected)` + reschedule tick.
  - filtering mode on: printable runes append to `filter`, Backspace edits,
    `Esc` exits filtering (keep/clear filter), `Enter` exits filtering keeping filter.
  - normal mode keys: `j/k` move cursor over **worktree** rows (skip repo headers or
    land on them harmlessly), `h`/`l` collapse/expand the repo of the current row,
    `z` toggle `allCollapsed`, `Tab` toggle `activeTab`, `/` enter filtering,
    `Enter` attach (see Step 7), `q`/`ctrl+c` quit.
  - on cursor move or filter change → recompute `rows`, issue `capturePreviewCmd`
    and `computeDiffCmd` for the newly selected worktree.
- `View`: join left tree + right tabbed pane with Lipgloss; render hint bar:
  `↵ attach · j/k move · h/l fold · ⇥ tab · / filter · q quit` (and a footer line
  for the inline `status` message when set).
- Verify: compiles; `go vet` clean.

#### Step 7: Attach behavior

- `attachCmd(selected)`:
  - resolve `window := worktree.GetWindowName(selected.Name)`.
  - if `os.Getenv("TMUX") == ""` → return a `statusMsg` "not inside tmux; run `dg wt ui` from a tmux session".
  - `session, ok := tmux.WindowSession(window)`; if `!ok` → `statusMsg`
    "no live window for <name> (repair coming in Phase 2)".
  - else `tmux.SwitchToWindow(session, window)` then return `tea.Quit`.
- Verify: covered by model test (mock tmux via injected funcs — see note below).

#### Step 8: `Run()` entrypoint + command wiring

- `internal/tui/worktree/run.go`: `func Run() error` builds `worktree.New()`,
  `tmux.New()`, `git.New()`, constructs the `Model`, and runs
  `tea.NewProgram(m, tea.WithAltScreen()).Run()`.
- `cmd/worktree.go`: add
  ```go
  var worktreeUICmd = &cobra.Command{
      Use:     "ui",
      Aliases: []string{"dash", "dashboard"},
      Short:   "Open the worktree dashboard (TUI)",
      Args:    cobra.NoArgs,
      RunE:    func(cmd *cobra.Command, args []string) error { return tuiworktree.Run() },
  }
  ```
  and `worktreeCmd.AddCommand(worktreeUICmd)` in `init()`.
- Verify: `./devgita wt ui` launches (manual, Step 11).

#### Step 9: Model unit tests (`model_test.go`)

- `func init() { testutil.InitLogger() }`.
- Construct a `Model` with injected fake data and **injected I/O functions** so no
  real commands run: give the model function fields (e.g. `captureFn`, `diffFn`,
  `attachFn`) defaulting to the real tmux/git calls in `Run()` but overridable in
  tests. (Adjust Step 6 fields to expose these seams.)
- Tests: build-rows grouping/order; `j/k` cursor movement and bounds; `h/l`/`z`
  fold behavior changes visible rows; `Tab` toggles `activeTab`; `/` then runes
  filter `rows`; attach with no `TMUX` sets the inline status message and does not quit.
- Verify: `go test ./internal/tui/worktree/...` passes.

#### Step 10: Helper tests

- `tmux_test.go`: `CapturePane` issues `capture-pane -p -e -t s:w` and returns stdout
  (assert via `MockBaseCommand` + `GetLastExecCommandCall`); error path surfaces stderr.
- `git_test.go`: `Diff` returns stdout; `DiffStat` parses `--numstat` (incl. a binary
  `-\t-\tfile` line). `testutil.VerifyNoRealCommands` where applicable.
- Verify: `go test ./internal/apps/tmux/... ./internal/apps/git/...` passes.

#### Step 11: tmux binding + manual smoke

- `configs/tmux/tmux.conf:113`: change the popup command from `devgita wt j` to
  `devgita wt ui`. Keep `display-popup -E -w 80% -h 60%`.
- Verify: see Manual Verification below.

---

## 6. Verification Plan

### Automated Verification

```bash
go build -o devgita main.go
go vet ./...
go test ./...
go test ./internal/tui/worktree/... ./internal/apps/tmux/... ./internal/apps/git/...
```

### Manual Verification

1. In a tmux session, create two worktrees with Claude running:
   `./devgita wt new demo-a --ai claude` and `./devgita wt new demo-b --ai claude`.
2. Run `./devgita wt ui`.
3. Confirm: tree shows the repo with `demo-a`/`demo-b` children, status glyph + any
   `±N`; `j/k` moves; `h/l` and `z` fold/unfold; right pane Preview shows the
   selected worktree's Claude output; `Tab` switches to Diff and shows `git diff`.
4. Press `/`, type part of a name → list filters; `Esc` clears.
5. Select `demo-a`, press `Enter` → TUI exits and tmux is now on the `wt-demo-a`
   window.
6. Re-bind check: `prefix + u` opens the same dashboard in a popup.
7. Run `./devgita wt ui` **outside** tmux → it still renders; `Enter` shows the
   inline "not inside tmux" status instead of crashing.

### Regression Check

- `dg wt j`, `dg wt l`, `dg wt new`, `dg wt rm`, `dg wt repair`, `dg wt prune` still
  behave exactly as before (only added a command + changed one tmux binding).
- `go test ./internal/tooling/worktree/...` still passes (untouched logic).

---

## 7. Risks & Trade-offs

| Risk                                                             | Likelihood | Mitigation                                                                                        |
| ---------------------------------------------------------------- | ---------- | ------------------------------------------------------------------------------------------------- |
| New deps (bubbletea/lipgloss) increase binary size               | High       | Both are pure Go, no cgo → cross-compilation (darwin/linux, amd64/arm64) unaffected. Acceptable.  |
| `capture-pane` only returns the **visible** pane, not scrollback | Medium     | Acceptable for a live preview in Phase 1; document. Scrollback/viewport scroll deferred.          |
| `switch-client` from inside a `display-popup` behaves oddly      | Medium     | Same mechanism the current fzf popup already uses successfully; verified in manual Step 5/6.      |
| Bubble Tea models are awkward to unit-test                       | Medium     | Test the pure `Update`/state transitions with synthetic msgs + injected I/O fns; no real program. |
| Preview refresh tick causes flicker/CPU                          | Low        | 1.5s tick, capture only the selected window; `tea.WithAltScreen` double-buffers.                  |

### Trade-offs Made

- **Bubble Tea over hand-rolled raw-terminal rendering** — matches Go ecosystem
  norms, gives us alt-screen/diffed rendering for free, and Phases 2–3 (filtering,
  popups, notifications) build naturally on it. Cost: two new dependencies.
- **Keep `dg wt j` (fzf) as-is** — no migration risk; the TUI is additive and the
  only behavioral change is the `prefix + u` binding target.
- **NERDTree-familiar keys, no palette/which-key** — honors the maintainer's
  "least thinking / fewest new keybindings" constraint; richer models from the
  wireframes are explicitly deferred.

---

## 8. Cross-Model Review Notes

- [ ] Root cause / goal confirmed? (This is a feature, not a bug — objective in §3.)
- [ ] All affected files identified? (§5 table)
- [ ] Verification steps executable? (§6)
- [ ] Scope appropriately bounded to ~3h and Phase 1 only? (§4)
- [ ] Dependency addition acceptable (bubbletea/lipgloss, pure Go)?
- [ ] TUI testability seam (injected I/O fns) acceptable vs. teatest dependency?

**Reviewer notes:**
(Fill in during review)

```

```
