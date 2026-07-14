---
description: Review a PR and post one cohesive review with inline comments — apply findings already in context (from a code/doc reviewer agent or another model) or review directly, dedup against existing threads, and submit a single verdict. Use for "review this PR", code review, or doc review.
temperature: 0.1
permission:
  write: allow
  edit: deny
  bash:
    "*": deny
    "devgita task *": allow
    "git diff*": allow
    "git fetch*": allow
    "git log*": allow
    "git branch*": allow
---

Post review feedback to a PR as **one cohesive review**. Findings often already sit in the conversation — produced by a `code-reviewer`/`document-reviewer` agent or another model (gpt, qwen, kimi, …). Use those; if context is thin, review directly with the lens below. The repo is the current working directory.

## Usage

```
/review-pr [PR_NUMBER]
```

The PR is resolved from the current branch unless you pass a number.

## Process

### 1. Find the PR

If a `PR_NUMBER` was given, use it (pass `--pr PR_NUMBER` below). Otherwise:

```bash
devgita task current-pr
```

If it prints "No pull request found for the current branch.", stop and tell the user this branch has no PR.

### 2. Load context

```bash
devgita task pr-view          # add --pr PR_NUMBER if you have one
```

Read the PR's purpose first — the description and linked ticket — before any code. Gather the findings already in the conversation. If there are none, review the change yourself with the lens in step 4. For a locally checked-out branch, run `devgita task review-scope` for the orientation (branch, ahead/behind, commits, per-file stats), then `devgita task branch-diff` (or `--file <path>` for one file) for the full noise-filtered diff.

### 3. Fetch existing threads and dedup — never repeat addressed feedback

```bash
devgita task review-threads --state all
```

This returns **resolved and unresolved** threads. Before posting, drop any finding that:

- targets the same `path:line` as an existing thread making substantially the same point, or
- is already **resolved** — a resolved thread is handled; re-raising it is noise.

Keep a count of what you skipped for the summary.

### 4. The review lens (high-leverage first — order matters)

Governing principle: **approve when the PR leaves the codebase healthier than without it**, not when it's "perfect". The question is "is the codebase better merged than not?", not "would I have written it this way?". If the PR is huge, the first finding is to split it — small PRs get genuinely reviewed; large ones get rubber-stamped.

Work the passes in this order so design problems surface before you nitpick code that shouldn't exist:

1. **Design** — does it belong here, fit existing patterns, sit at the right abstraction? Flag over-engineering (generality/features not needed now).
2. **Functionality** — does it do what it claims on the unhappy paths too? Edge cases, nulls, empty inputs, boundaries, downstream failures, concurrency.
3. **Complexity** — too complex = can't be understood quickly, or invites bugs on the next edit. If you can't follow it, others won't either.
4. **Tests** — real coverage of the new logic and edge cases; would they fail if the logic broke?
5. **Naming / comments / docs** — names convey intent; comments explain _why_; update READMEs when behavior changes.
6. **Style** — last and lightest. Prefix optional style points with `Nit:`; never block on personal preference.

Across every pass, a **security lens** for anything touching data, auth, or external input: input validation, authz, injection (SQL/XSS), committed secrets, and the safety of new dependencies.

For a **doc review**, swap passes 2–4 for: accuracy, completeness, structure, and clarity.

Severity tags drive the verdict: `[CRITICAL]` (data loss, security, correctness — a blocker), `[IMPORTANT]`, `[MINOR]`/`[Nit]`. Anchor disagreement on engineering principle, not authority — and call out what's genuinely good, especially well-addressed prior feedback.

### 5. Compose the review — a summary body plus inline comments

Findings that point at a specific line become **inline comments** anchored to the diff; everything else goes in the summary **body**.

**Write plainly.** Everything posted must be understandable by any engineer, including a junior one: everyday words, short sentences, no fancy vocabulary or filler. Each comment says what's wrong, why it matters, and the fix — nothing more.

**Body** — GitHub-Flavored Markdown, written to a scratch file (`/tmp/review.md`); pass it with `--body-file` so backticks and apostrophes survive:

```markdown
## Summary

<!-- 1–2 lines: does this improve code health? -->

## Strengths

<!-- What's done well. Don't skip this. -->

## General notes

<!-- Findings with no single line to anchor to, and any cross-cutting concern. -->

## Questions

<!-- Anything you need the author to clarify. -->

---

<!-- footer when applicable: "Skipped N finding(s) already covered by existing threads." -->
```

**Inline comments** — write a JSON array to a scratch file (`/tmp/comments.json`). Each entry anchors to a diff line; only lines present in the diff can carry one. Lead the body with the severity tag:

```json
[
  {
    "path": "internal/client.go",
    "line": 42,
    "body": "**[CRITICAL]** Missing error handling — a nil response here panics. Guard before dereferencing."
  },
  {
    "path": "internal/client.go",
    "start_line": 60,
    "line": 65,
    "body": "**[Nit]** This block reads more clearly as an early return."
  }
]
```

`line` is the line in the file (right side of the diff); add `start_line` for a multi-line range. Drop any finding already covered by an existing thread (step 3).

### 6. Submit one review

Post the body and the inline comments together as a single review, choosing the verdict:

| Verdict         | When                                       | `--event`         |
| --------------- | ------------------------------------------ | ----------------- |
| Request changes | Any `[CRITICAL]` / blocking issue          | `request-changes` |
| Approve         | No blockers; leaves the codebase healthier | `approve`         |
| Comment         | Suggestions only, nothing blocking         | `comment`         |

```bash
devgita task submit-review \
  --event request-changes \
  --body-file /tmp/review.md \
  --comments-file /tmp/comments.json      # omit when there are no inline findings
```

Add `--pr PR_NUMBER` when you resolved a number in step 1. The review posts atomically — one notification, all inline comments grouped under it.

If you have **nothing new** to add (everything is already covered or addressed), don't post an empty review — post one short `comment-pr` saying so. Likewise, if existing **unresolved** threads remain unaddressed and that's the main issue, flag it in a brief comment rather than re-listing each thread.

## Output

Return a terse summary to the user:

```
## Review posted to PR #<num> — <request changes | approve | comment>

- findings: <N posted>
- skipped: <M already covered by existing threads>

<PR URL>
```

## Notes

- This command never edits code. It reads, then posts exactly one review.
- **Dedup is mandatory**: never duplicate a finding already raised, and treat a resolved thread as handled.
- A line that isn't part of the diff can't take an inline comment — move that finding to the body's "General notes" instead.

## References

- Google Engineering Practices — the standard of code review & what to look for: https://google.github.io/eng-practices/review/reviewer/
- "Software Engineering at Google", Code Review chapter: https://abseil.io/resources/swe-book/html/ch09.html
