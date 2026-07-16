# Cycle: `dg wt ui` polish — empty state, cursor editing, adopt, PR title

**Date:** 2026-07-16
**Estimated Duration:** Steps 1–2 shipped; remaining Steps 3–4 ~3 hours (adopt ~1.5h,
PR title ~1.5h). Split as two follow-up commits so each stays within the ~3h cycle-slice
guidance.
**Status:** In Progress

---

## 1. Domain Context

`dg wt ui` is the worktree dashboard (list + attach/destroy/repair/filter/diff pane, plus
in-TUI create from cycle `2026-07-15-wt-ui-create-flow.md`). Three rough edges surfaced in
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

Also raised: a way to take a branch created outside devgita and move it under worktree
management. The building block already exists — `CreateWorktreeIn` (`git.go:388`) checks
out an existing local/remote branch when one is present — but there's no discoverable verb
for it and no handling of the "branch is checked out elsewhere" case.

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
  - `internal/tooling/worktree/worktree.go` — `Create` / `CreateAt` → shared `create`;
    `internal/apps/git/git.go` — `CreateWorktreeIn(repoDir, path, branch)` already adopts an
    existing branch via `git worktree add <path> <branch>`.
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
editing; a branch created outside devgita can be adopted into a managed worktree with one
discoverable command; and the diff pane shows the PR title when the branch has a PR.

---

## 4. Scope Boundary

### In Scope

- [x] Empty-state: once the first `List()` returns zero worktrees, the diff pane shows
      "No worktrees yet — press n to create one" instead of a permanent `(loading...)`.
      `(loading...)` stays only for a genuinely pending diff.
- [x] Shared `TextInput` editor (value + caret, left/right/home/end, mid-string
      insert/delete), adopted by the name prompt, filter, and repo-picker query. Removes
      the old append-only `TrimLastRune` path.
- [ ] `dg wt adopt <branch>` (see Step 3) — adopt an existing branch into a managed
      worktree + window, with the checked-out-elsewhere and uncommitted-work cases handled.
- [ ] PR title in the diff-pane header (see Step 4) — best-effort, cached, shown only when
      a PR exists; runs off the refresh path with a bounded timeout.
- [ ] Docs updated for `adopt` (required, not optional): `cmd/worktree.go` help for both
      `create` and `adopt`, and the worktree subcommand list in `docs/spec.md`.

### Explicitly Out of Scope

- Making the per-file dividers "louder" — they already exist; revisit only if the PR-title
  header lands and the file headers then look weak by comparison.
- Adopting multiple branches at once, or a picker of adoptable branches in the TUI (CLI
  verb first; a TUI entry point is a later cycle).
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

### Step 3 — `dg wt adopt <branch>` (PLANNED)

- New Cobra subcommand `adopt` (alias none for now) under `worktreeCmd`, `--repo` flag like
  `create`. Delegates to a new `WorktreeManager.Adopt(repoPath, branch, coder)` that reuses
  the shared `create` flow but requires the branch to already exist (local or remote) —
  `CreateWorktreeIn` already does the `git worktree add <path> <branch>` checkout in that
  case, so `Adopt` is `create` with an "existing branch required" precondition plus the two
  guards below. No new worktree layout or window logic.
- **Checked-out-elsewhere guard:** git refuses `worktree add` for a branch already checked
  out in another worktree (typically the user's main clone, where they created it). Detect
  this before calling git via `git worktree list --porcelain` with precise matching:
  - Each worktree block ends with `branch refs/heads/<name>`; strip the `refs/heads/`
    prefix and compare against the target branch by exact name.
  - A block with no `branch` line (detached HEAD, or a bare repo entry) is skipped — it
    holds no branch and can't collide.
  - Block **only** when the branch is checked out at a path **different from** this
    adopt's target worktree path. The same-path case is a re-run and is already covered by
    `create`'s existing "worktree already exists" state check, so it must not be reported
    as a checked-out-elsewhere conflict.
  - Error names the conflicting path and tells the user to switch that checkout off the
    branch first. Never `--force` silently — that would detach the other checkout.
- **Uncommitted-work caveat:** `worktree add` brings only committed history; uncommitted
  changes stay in the origin checkout. `adopt` prints a plain-language note that uncommitted
  work does not move, so nobody assumes it did.
- **Required docs update (not optional):** `create`'s help text currently states it
  "Creates a new branch with the same name" (`cmd/worktree.go:51`). Adding `adopt` without
  updating that is contradictory UX. Update `cmd/worktree.go` (create's long help notes the
  existing-branch-adopts behavior; new `adopt` help) and `docs/spec.md`'s worktree
  subcommand list in the same change.
- Tests: `Adopt` happy path (existing local branch → worktree + window, via mocks),
  branch-not-found error, checked-out-elsewhere error (different path), same-path re-run
  falls through to the existing state check, `VerifyNoRealCommands`.
- **Open decision (needs author sign-off before build):** after `adopt` exists, does
  `dg wt create <existing-branch>` keep adopting the existing branch (current behavior,
  non-breaking, but two verbs do the overlapping thing), or does `create` become
  new-branch-only for a cleaner mental model (clearer, but a behavior change subject to
  change discipline — deprecation note + spec update)? Recommendation: **keep `create`'s
  current behavior** (non-breaking; `adopt` is the discoverable explicit verb), and only
  tighten `create` if a future cycle deprecates it deliberately.

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
3. [ ] Create a branch in a normal clone, switch that clone off it, then `dg wt adopt <branch>` → worktree + window on the existing branch.
4. [ ] `dg wt adopt <branch>` while the branch is still checked out in the main clone → clear error naming where it's checked out; nothing changed.
5. [ ] A worktree whose branch has an open PR → diff header shows the PR title; a branch with no PR → no title line, no delay.

### Regression Check

- [x] Attach/destroy/repair/filter/diff-pane and the `n` create flow unchanged.
- [ ] `dg wt create`, `new`, `rm`, `ls`, `repair`, `prune` unchanged.

---

## 7. Risks & Trade-offs

| Risk                                                 | Likelihood | Mitigation                                                                 |
| ---------------------------------------------------- | ---------- | -------------------------------------------------------------------------- |
| `adopt` `--force`-checks out a branch already in use | Med        | Detect checked-out-elsewhere before git runs; return an error, never force |
| Users expect uncommitted work to move with `adopt`   | Med        | Explicit plain-language note; document in help + spec                      |
| `gh pr view` adds latency to the diff pane           | Med        | Only when `gh` present; cache per branch; never block the diff on it       |
| Reverse-video caret renders oddly on some terminals  | Low        | Standard SGR reverse; block caret at end preserves the old visible cursor  |

### Trade-offs

- **`adopt` as a thin verb over `create`'s existing-branch path** rather than a new
  subsystem — the checkout logic already exists; the value is discoverability + the two
  guards, not new plumbing.
- **PR title via `gh` rather than a git API client** — `gh` is already the project's GitHub
  surface elsewhere; no new dependency, and absence degrades to "no title" cleanly.

---

## 8. Notes for Implementers

- Steps 1–2 are shipped. Steps 3–4 need approval before implementation (this cycle doc is
  the request for that).
- Commit after each step once its verify check passes.
- If `adopt` needs to grow beyond the two guards, stop and document a follow-up rather than
  widening scope.
