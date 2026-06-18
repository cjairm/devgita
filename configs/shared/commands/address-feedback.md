---
description: Address PR review feedback — implement suggestions, reply to reviewers, and resolve threads on the current branch's PR
temperature: 0.2
permission:
  write: allow
  edit: allow
  bash:
    "*": deny
    "devgita task *": allow
    "git add *": allow
    "git commit *": allow
    "git push*": allow
    "git status": allow
    "git diff*": allow
    "go *": allow
    "npm *": allow
    "npx *": allow
    "pnpm *": allow
    "yarn *": allow
    "make *": allow
---

Address PR review feedback: implement the requested changes, reply to each reviewer thread, and resolve the ones you handled.

## Usage

```
/address-feedback [PR_NUMBER]
```

The PR is resolved from the current branch unless you pass a number. The repo is the current working directory.

## Process

### 1. Find the PR

If a `PR_NUMBER` was given, use it and pass `--pr PR_NUMBER` to every command below.

Otherwise resolve it from the branch:

```bash
devgita task current-pr
```

If that prints "No pull request found for the current branch.", stop and tell the user this branch has no PR — they can open one with `/create-pr`. Do nothing else.

### 2. Read the unresolved threads

```bash
devgita task review-threads --state unresolved   # add --pr PR_NUMBER if you have one
```

`--state unresolved` is the default, but pass it explicitly — these are the only threads to act on; resolved ones are already handled. The output is compact markdown. Each block is headed `## path:line (thread <id>)`, may include a diff hunk, then `**author** (comment-id): body`. The thread `<id>` is what you reply to and resolve.

Read the whole list before changing anything — a later comment can change how you handle an earlier one.

### 3. Triage each thread

Sort every thread into one of these, then work in logical batches (not one commit per comment):

| Bucket        | What it is                                  | Action                                              |
| ------------- | ------------------------------------------- | --------------------------------------------------- |
| Fix           | A concrete change you agree with            | Implement it, reply, resolve                        |
| Already done  | Code already satisfies the comment          | Reply pointing to the code/commit, resolve          |
| Discuss       | Disagreement or a tradeoff worth explaining | Reply with reasoning, resolve unless it's blocking  |
| Out of scope  | Valid but not for this PR                   | Reply that you'll track it separately, resolve      |
| Needs clarity | You genuinely don't understand the ask      | Reply with a specific question — **do not resolve** |

Before implementing, look at how the surrounding code already solves similar problems so your change matches the codebase's existing patterns.

### 4. Implement, verify, and commit

Make the changes, then verify and commit **following the repo's own conventions** — read `AGENTS.md`, `CLAUDE.md`, or `INSTRUCTIONS.md` for how this project builds, tests, and commits. Don't assume a language or toolchain.

Group related fixes into meaningful commits and push so reviewers can see what changed since their pass.

### 5. Reply and resolve

For each thread you handled, reply then resolve:

```bash
devgita task reply-thread <thread-id> "Done in <short-sha> — <what changed>"
devgita task resolve-thread <thread-id>
```

Use `--body-file <path>` instead of an inline body when the reply needs multi-line Markdown.

Reply to every thread, even trivial ones, so the reviewer never has to guess whether you saw it. Reference the commit when a fix landed there. Leave "needs clarity" threads open.

## Output

Return only this — no preamble, no narration:

```
## PR #<num> — feedback addressed

### Changes
| File | Change | Commit |
|------|--------|--------|
| ... | ... | <sha> |

### Threads
| Location | Bucket | Outcome |
|----------|--------|---------|
| path:line | Fix | replied + resolved |
| path:line | Needs clarity | asked, left open |

### Verification
<result of the repo's build/test command, or "n/a — no code changes">
```

## Notes

- Resolve a thread only after you've actually addressed it; never resolve to clear noise.
- Don't make changes unrelated to the feedback.
- If a push or a resolve fails, report it and continue with the rest — your local commits are safe.
