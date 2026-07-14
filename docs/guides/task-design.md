# Task Design — AI-first, token-wise `dg task` subcommands

`dg task` subcommands are consumed by two audiences: **an LLM agent first** (Claude
Code / OpenCode invoking `devgita task …` from the shared configs), **a human second**
(via the `dge` shell wrapper). When those audiences conflict, design for the agent —
humans tolerate terse output far better than agents tolerate noisy output.

Origin: the output principles were decided in the
[2026-07-14 token-aware git tasks cycle](../plans/cycles/2026-07-14-token-aware-git-tasks.md)
(§5) and generalized here. This guide is their canonical home going forward.

---

## When to build a task (and when not to)

Build a `dg task` subcommand only when it adds **semantic value** an agent can't get
from raw commands. The three justifications, strongest first:

1. **Collapse round-trips.** Each agent tool call costs ~100–200 tokens of framing
   plus a reasoning turn. If a flow needs N commands that always run together
   (fetch → detect default branch → merge-base → stat), one task that runs them
   in-process wins even if the output were identical. Example: `review-scope`
   replaced a 6-call orientation dance.
2. **Enforce policy and contracts.** Things an agent _should_ always do but might
   skip: dedup against resolved threads, keep the merge-base stable across a review
   session, never silently hide excluded files. Encoding the policy in the task makes
   it non-optional. Example: `submit-review` posts body + inline comments atomically.
3. **Render for the reader.** Reshape verbose structured output into compact text
   (see principles below). Example: `review-threads` turns a GraphQL payload into
   one markdown block per thread.

**Don't build a task when:**

- The saving is negligible — measure first (see below). `git status` → `--porcelain`
  saves ~350 bytes; that's prompt guidance, not a subcommand.
- Only generic compression is needed, with no orchestration or policy. That is a
  general-purpose tool's job, not bespoke Go code — see [rtk](#future-rtk) below.
- The output feeds a human only. Streaming raw output (`Stream: true`) is fine there.

## Output principles

1. **Labeled plain text, not markdown scaffolding.** Line-oriented `key: value`
   labels, `- ` lists, aligned stat lines. No headers, tables, bold, or emoji — that
   is rendering decoration an LLM pays for without needing. Markdown syntax only
   where structure earns its tokens (e.g. ` ```diff ` fences), per the
   `reviewThreadsFilter` precedent in `internal/tooling/terminal/dev_tools/jq/jq.go`.
2. **Payload only.** Never wrap output in prose ("Here is the scope:", "Done! ✓").
   The first byte of output is data.
3. **Mutations confirm with one line: verb + target** — e.g. `Resolved thread
PRRT_abc`. Never bare `ok` (the echoed target lets an agent verify it acted on the
   right object and reuse the id without re-fetching) and never more than one line.
   Success/failure for scripts lives in the exit code, not the text.
4. **Stable sentinels.** Fixed phrases agents match verbatim (`No unresolved review
threads.`, `On <default> — no branch to compare…`) are contracts with the shared
   configs in `configs/shared/` — changing one is a breaking change and needs the
   consuming config updated in the same commit. Empty results always get a sentinel;
   empty stdout is ambiguous to an agent (success? nothing found? crash?).
5. **Lossy is allowed only with a receipt.** Filtering noise (lockfiles, generated
   files) is encouraged, but every omission must be announced with a one-line note
   and an escape hatch (`branch-diff` prints `excluded: go.sum (+40/-12)` and offers
   `--file`). A reviewer agent's verification bar requires reading actual code —
   silently compressed payloads produce confidently wrong reviews.

## Architecture: orchestrate, then format

Mirror `PRManager` (`internal/tooling/task/pr.go`): the manager method
**orchestrates** raw fetches, then hands raw output to a **pure formatter** that
renders the final text.

- Input is JSON (a `gh` payload) → a `jq` filter is the formatter.
- Input is line-oriented text (git plumbing) → a pure Go function is the formatter;
  a jq subprocess would add a hop just to split strings.

The split is what makes testing cheap: formatters get golden-fixture unit tests with
zero mocking; orchestration tests only assert which commands ran and the error paths
(`testutil.MockApp` + `VerifyNoRealCommands`, per
[testing-patterns.md](testing-patterns.md)).

## Measure before and after

Token estimates: `bytes / 4` is close enough for prose; code runs ~3.5 chars/token.
Before building, measure the raw output the agent ingests today (`… | wc -c`); after,
measure the task's output on the same input. If the delta plus saved round-trips
isn't clearly worth a Go file and its tests, don't build it. Reference baselines from
the 2026-07-14 cycle: a 25-file branch diff was 192 KB (~50k tokens) raw vs 3.3 KB
(~950 tokens) as a stat table; npm-repo lockfile diffs alone run 10–50k tokens.

---

## Future: rtk

[rtk](https://github.com/rtk-ai/rtk) ("Rust Token Killer", Apache-2.0) is a CLI proxy
that compresses the output of 100+ common commands (`git`, test runners, `docker`,
`cat`, `grep`, …) before it reaches the LLM — claiming 60–90% reduction — and hooks
into Claude Code by auto-rewriting Bash calls (`rtk init -g`). Verified 2026-07-14:
~71k stars, single Rust binary, official Homebrew formula, very active; young
(created 2026-01) with a large open-issue count.

**Relationship to `dg task`: complementary, not competing.** rtk is _generic lossy
compression_ of whatever command runs; `dg task` is _semantic orchestration + policy_
(one-call flows, dedup contracts, stable sentinels, atomic reviews) that no generic
proxy can provide. rtk would cover the long tail this guide says **not** to build
tasks for: test-runner output, `docker ps`, `cat`/`grep` noise.

Adoption stance (tracked in [ROADMAP.md](../../ROADMAP.md) under AI & Development
Tools):

- Candidate app installer for the planned `ai-tools` category — brew on macOS,
  install-script or GitHub-release strategy on Debian.
- **Install the binary; make the hook opt-in.** `rtk init -g` rewrites _every_ agent
  Bash call, including inside our carefully designed flows; per CLAUDE.md's security
  non-negotiables the install script and hook behavior must be reviewed before we
  automate them, and lossy compression must never silently apply to a reviewer's
  diff payload (principle 5).
- Revisit once the project stabilizes (it is six months old and moving fast).
