# Cycle: `dg wt ui` polish — empty state, cursor editing, existing-branch create, PR title

**Date:** 2026-07-16
**Estimated Duration:** Steps 1–2 shipped; remaining Steps 3–4 ~3 hours (existing-branch
create ~1.5h, PR title ~1.5h). Split as two follow-up commits so each stays within the
~3h cycle-slice guidance.
**Status:** In Progress

---

## 1. Domain Context

`dg wt ui` is the worktree dashboard (list + attach/destroy/repair/filter/diff pane, plus
in-TUI create from cycle `2026-07-15-wt-ui-create-flow.md`). Four rough edges surfaced in
use:

1. **Empty dashboard shows `(loading...)` forever.** With no worktrees, no row is ever
   selected, so the diff pane never computes and sits on its placeholder — it reads as a
   hang, not "nothing here yet."
2. **Text inputs can't edit mid-string.** The name prompt, filter, and repo-picker query
   were append-and-backspace only. Pasting a path and then fixing a character in the
   middle was impossible — you had to backspace back to it and retype.
3. **The diff pane doesn't say what branch/PR you're looking at beyond `base ← branch`.**
   Per-file dividers already exist (`rewriteFileHeaders`), but there's no PR title, so the
   pane doesn't connect to the change being reviewed.

4. **Creating a worktree for an existing local branch fails.** Naming a worktree after a
   branch that already exists locally and is checked out in the main clone dead-ends with
   `git: Preparing worktree (checking out '<branch>')` — git refuses a branch already held
   by another worktree, and nothing frees it first. This ordinary "I made a branch, now
   give it a worktree" flow should just work, from both `dg wt create` and `dg wt ui`.

Related: prior cycles `2026-06-07-worktree-v2-tui-dashboard.md`,
`2026-07-15-wt-ui-create-flow.md`; ROADMAP Worktree Enhancements.

---

## 2. Engineer Context

- **Relevant files:**
  - `internal/tui/worktree/model.go` — dashboard model. Diff computed only when a worktree
    row is selected (`computeDiffCmd`, dispatched from `statusesMsg`/`j`/`k`);
    `renderDiffContent` shows `(loading...)` when `diffContent == ""`; right-pane header in
    `renderRight` (`base ← branch  ±f +a -r`).
  - `internal/tui/worktree/diffview.go` — `rewriteFileHeaders` already renders a styled
    `─── path +a -r` divider per file with a blank line between files; `[`/`]` jump between
    them.
  - `internal/tui/components/textinput.go` — `TextInput` shared single-line editor (value +
    caret; left/right/home/end/backspace/delete/insert; `RenderPlain`).
  - `internal/tui/components/filter.go`, `cmdpalette.go`, `fuzzypicker.go` — the three
    fields that embed `TextInput`.
  - `internal/tooling/worktree/worktree.go` — `Create` / `CreateAt` → shared `create`,
    which calls `git.CreateWorktreeIn(repoRoot, wtPath, name)` (`worktree.go:216`).
  - `internal/apps/git/git.go` — `CreateWorktreeIn(repoDir, path, branch)` (`git.go:388`)
    runs `git worktree add <path> <branch>` when the branch exists; this is where the
    checked-out-in-main-clone case must be handled. Helpers already present:
    `ListWorktreesAt`, `IsWorktreeDirty`, `DefaultBranchIn`, `SwitchBranch`.
  - `internal/tooling/task/branchdiff.go` — `BranchDiffAt` produces the diff + base
    metadata; a natural home for an optional PR-title lookup.
- **Testing:** [docs/guides/testing-patterns.md](../../guides/testing-patterns.md) —
  `testutil.MockApp`, `VerifyNoRealCommands`, isolate `paths.Paths.*` with `t.Cleanup`,
  `func init() { testutil.InitLogger() }`. TUI tests drive `Update` and assert on `View()`
  / seam calls.
- **Test commands:**

  ```bash
  go test ./internal/tui/... ./internal/tooling/worktree/ ./internal/tooling/task/ ./internal/apps/git/
  go test ./...
  make lint
  ```

---

## 3. Objective

An empty dashboard tells the user what to do; every text field supports mid-string
editing; naming a worktree after an existing local branch adopts that branch instead of
failing (from both `dg wt create` and the TUI); and the diff pane shows the PR title when
the branch has a PR.

---

## 4. Scope Boundary

### In Scope

- [x] Empty-state: once the first `List()` returns zero worktrees, the diff pane shows
      "No worktrees yet — press n to create one" instead of a permanent `(loading...)`.
      `(loading...)` stays only for a genuinely pending diff.
- [x] Shared `TextInput` editor (value + caret, left/right/home/end, mid-string
      insert/delete), adopted by the name prompt, filter, and repo-picker query. Removes
      the old append-only `TrimLastRune` path.
- [ ] `create` converts an existing local branch (see Step 3) — when the worktree name
      matches a branch already checked out in the main clone, free that branch (switch the
      clone to its default branch) and create the worktree on it, instead of failing.
      Covers both `dg wt create` and the TUI create flow; refuses on an uncommitted-work
      tree rather than move it.
- [ ] PR title in the diff-pane header (see Step 4) — best-effort, cached, shown only when
      a PR exists; runs off the refresh path with a bounded timeout.
- [ ] Docs updated for the existing-branch behavior (required, not optional):
      `cmd/worktree.go` create help and the worktree section of `docs/spec.md`.

### Explicitly Out of Scope

- Making the per-file dividers "louder" — they already exist; revisit only if the PR-title
  header lands and the file headers then look weak by comparison.
- A dedicated command or TUI picker for existing branches — the behavior is folded into
  the normal create path (name matches a branch → adopt it); a separate entry point is a
  later cycle if ever needed.
- Creating/opening a PR from the TUI — read-only title lookup only.

**Scope is locked.** New needs → document for a future cycle.

---

## 5. Implementation Plan

### Step 1 — Empty-state fix (DONE)

- `Model.loaded bool`, set `true` on the first `statusesMsg`.
- `renderRight`: when `loaded && len(statuses) == 0`, render the guidance line and skip the
  stat header; otherwise unchanged. `renderDiffContent` keeps `(loading...)` for the
  transient selected-but-not-yet-loaded case.
- Test: `TestEmptyDashboardShowsGuidance` (loading before first load → guidance after an
  empty load).

### Step 2 — Shared cursor-aware text editor (DONE)

- `internal/tui/components/textinput.go`: `TextInput` with caret; `HandleKey` returns
  `(handled, changed)`; `RenderPlain` draws a reverse-video block caret (mid-string or at
  end). `renderCaret` also used by `CommandPalette`.
- `FilterField` embeds it (`Value()` accessor replacing the old `Text` field);
  `FuzzyPicker` query and the name prompt (`Model.createInput`) use it. `TrimLastRune`
  removed.
- Tests: `TestTextInputEditing`, `TestTextInputRenderCaret`, updated filter/create/picker
  tests; call sites in the inventory TUI migrated to `Value()`.

### Step 3 — `create` converts an existing local branch (PLANNED)

**The bug.** Giving a new worktree the same name as a branch that already exists locally and
is checked out in the main clone fails:

```
create failed: failed to create worktree: git: Preparing worktree (checking out 'feature/gh-546-...')
```

Root cause: both `dg wt create` and the `dg wt ui` create flow funnel through
`git.CreateWorktreeIn` (`git.go:388`, reached via `worktree.go:216`). When the branch exists
locally it runs `git worktree add <path> <branch>`. Git refuses to check out a branch that
is already checked out in another worktree — here the user's main clone, where they created
it — and aborts with `fatal: '<branch>' is already used by worktree at '<path>'`. Nothing
frees the branch first, so an ordinary "I made a branch, now give it a worktree" flow
dead-ends.

**The fix (one place, both surfaces).** In `CreateWorktreeIn`'s `localExists` branch, before
`git worktree add`, detect whether the branch is already checked out in an existing worktree
of the same repo and, when it is, free it first:

- Enumerate the repo's worktrees with `ListWorktreesAt(repoDir)` (already parses
  `git worktree list --porcelain` into `{Path, Branch}`, with `refs/heads/` stripped). Find
  the entry whose `Branch` equals the target branch (exact match).
- No entry holds it → unchanged; `worktree add` proceeds exactly as today.
- The holder's path equals the new worktree `path` → cannot happen here: the upstream state
  check in `worktree.go`'s `create` already rejects an existing worktree before git is
  called, so no special case is needed.
- Otherwise the holder is the source checkout (the main clone). Free the branch by switching
  that checkout to the repo's default branch (`DefaultBranchIn(repoDir)`, then
  `git -C <holderPath> checkout <default>`), then run `worktree add`.

**Dirty-tree guard (never move uncommitted work).** Before switching the holder off the
branch, check `IsWorktreeDirty(holderPath)`. If it's dirty, do nothing and return a plain
error: the branch has uncommitted changes in `<holderPath>`; commit or stash them, then
retry. `worktree add` only carries committed history, and switching a dirty checkout could
carry changes across or fail mid-way — so refuse up front rather than risk it.

**Default-branch edge:** if the target branch _is_ the repo default, there is nothing to
switch the holder to; leave the existing git error to surface (adopting the default branch
into a side worktree isn't a real workflow and isn't worth special-casing).

**What the user sees:** create succeeds; the branch now lives in the managed worktree +
window, and the main clone is left on its default branch. Print a one-line note that the
source checkout was moved to `<default>`, so the switch isn't a surprise.

**Required docs update (not optional):** `create`'s long help (`cmd/worktree.go:51`) says it
"Creates a new branch with the same name". Extend it to note that when a branch of that name
already exists, create adopts it into the worktree (moving the source checkout to the
default branch). Update the worktree section of `docs/spec.md` in the same change.

- Tests (in `internal/apps/git/` with `testutil.MockApp` + `VerifyNoRealCommands`): branch
  checked out in the main clone with a clean tree → holder switched to default, `worktree
add` runs, success; dirty holder → error, no switch and no add; branch checked out
  nowhere → unchanged existing path; branch absent locally and remotely → still creates a
  new branch as before.

### Step 4 — PR title in the diff header (PLANNED)

The lookup is **its own async message, not part of `diffFn`/`computeDiffCmd`.** The
dashboard reloads every 3s (`refreshInterval`, `model.go:25`) and re-derives the selected
worktree's diff each time; folding a shell-out to `gh` into that path would make a slow or
hung `gh` stall every perceived diff update. So the PR title runs on a separate, cached,
bounded command instead.

- **Model state:** `prTitles map[string]string` keyed by branch, plus a
  `prTitlePending map[string]bool` (or a small set) so the same branch isn't looked up
  twice concurrently. Session-only cache: a branch is looked up once; the entry is never
  invalidated during the session (a stale title is a far smaller cost than re-shelling on
  every tick). `computeDiffCmd`/`diffFn` are unchanged.
- **Trigger:** when the selected worktree changes (`j`/`k`/`statusesMsg`) and its branch
  has no cache entry and none is pending, dispatch `prTitleCmd(branch, path)`. The 3s tick
  never triggers a lookup for an already-cached or pending branch, so steady state does
  zero `gh` calls.
- **The command:** best-effort `gh pr view --json title -q .title` run in the worktree dir,
  **wrapped in a 2s `context.WithTimeout` via `exec.CommandContext`** so a hung `gh` can't
  outlive one refresh interval. Returns empty (never an error surfaced to the user) when
  `gh` is absent from PATH, unauthenticated, times out, or the branch has no PR. Result
  flows back as `prTitleMsg{branch, title}`; Update stores it (empty title still cached, so
  a no-PR branch isn't retried every selection).
- **Render:** `renderRight` shows `PRTitle` as the first header line (title style) with the
  existing `base ← branch  ±f +a -r` line beneath it; nothing extra when the title is empty.
- **Injectable seam:** add a `prTitleFn func(branch, path string) string` alongside the
  other seams (`diffFn`, `attachFn`, …) so tests never shell out.
- Tests: header renders the title when present and omits the line when empty; a cached or
  pending branch is not re-looked-up on tick; `prTitleFn` mocked with
  `VerifyNoRealCommands`; absent-`gh`/timeout path yields empty title without error.

---

## 6. Verification Plan

### Automated

```bash
go test ./internal/tui/... ./internal/tooling/worktree/ ./internal/tooling/task/ ./internal/apps/git/
go test ./... -cover
make lint
```

### Manual

1. [x] `dg wt ui` with no worktrees → right pane shows the create guidance, not `(loading...)`.
2. [x] `n` → name prompt → paste a path → move the caret with ←/→ and fix a character mid-string.
3. [ ] In a repo, create a branch and stay on it. In `dg wt ui` pick that repo and name the worktree the same as the branch → worktree + window created on it; the main clone is left on its default branch.
4. [ ] Same as 3 but with uncommitted changes in the main clone → clear error saying the branch has uncommitted changes; nothing moved.
5. [ ] A worktree whose branch has an open PR → diff header shows the PR title; a branch with no PR → no title line, no delay.

### Regression Check

- [x] Attach/destroy/repair/filter/diff-pane and the `n` create flow unchanged.
- [ ] `dg wt new`, `rm`, `ls`, `repair`, `prune` unchanged; `create` unchanged except that a name matching an existing local branch now adopts it instead of erroring.

---

## 7. Risks & Trade-offs

| Risk                                                       | Likelihood | Mitigation                                                                                 |
| ---------------------------------------------------------- | ---------- | ------------------------------------------------------------------------------------------ |
| Switching the main clone off the branch surprises the user | Med        | Only done to free a branch for its worktree; print a one-line note; refuse on a dirty tree |
| Uncommitted work in the main clone could be disturbed      | Med        | Guard with `IsWorktreeDirty`; refuse and do nothing when dirty                             |
| `gh pr view` adds latency to the diff pane                 | Med        | Only when `gh` present; cache per branch; never block the diff on it                       |
| Reverse-video caret renders oddly on some terminals        | Low        | Standard SGR reverse; block caret at end preserves the old visible cursor                  |

### Trade-offs

- **Fold the behavior into `create` rather than a new command** — the checkout logic already
  lives in `CreateWorktreeIn`; the only gap is freeing a branch that's checked out in the
  main clone. Fixing it there covers both the CLI and the TUI in one place.
- **PR title via `gh` rather than a git API client** — `gh` is already the project's GitHub
  surface elsewhere; no new dependency, and absence degrades to "no title" cleanly.

---

## 8. Notes for Implementers

- Steps 1–2 are shipped. Steps 3–4 need approval before implementation (this cycle doc is
  the request for that).
- Commit after each step once its verify check passes.
- The Step 3 fix belongs in `git.CreateWorktreeIn` (the single funnel both surfaces use),
  not in the CLI or TUI layers. If it needs to grow beyond the dirty-tree guard, stop and
  document a follow-up rather than widening scope.
