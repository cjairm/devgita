---
description: Generate a PR title and WHY-focused description from current changes
agent: plan
subtask: true
---

# Goal

Generate a crisp PR title and a WHY-driven description for the current branch.
Focus on: rationale, problem addressed, scope touched, and any migration/risks.
Keep it actionable and reviewer-friendly.

$ARGUMENTS may be the context if present

---

## Repo & Branch

Current branch:
!`git rev-parse --abbrev-ref HEAD || true`

Base branch (detected):
!`bash -lc 'BASE=$(git symbolic-ref -q --short refs/remotes/origin/HEAD 2>/dev/null | sed "s#origin/##"); \
  [[ -z "$BASE" ]] && { git rev-parse --verify main >/dev/null 2>&1 && BASE=main || BASE=master; }; echo "$BASE"'`

## Summary of Work in Progress

Git status:
!`git status --porcelain=v1 || true`

Name-status since base (committed changes):
!`bash -lc 'BASE=$(git symbolic-ref -q --short refs/remotes/origin/HEAD 2>/dev/null | sed "s#origin/##"); \
  [[ -z "$BASE" ]] && { git rev-parse --verify main >/dev/null 2>&1 && BASE=main || BASE=master; }; \
  git diff --name-status "${BASE}..HEAD" || true'`

Staged vs unstaged (uncommitted):
!`git diff --cached --name-status || true`
!`git diff --name-status || true`

Recently touched files (top-level areas):
!`bash -lc '{ git diff --name-only "${BASE:-main}..HEAD" 2>/dev/null || true; \
  git diff --name-only --cached || true; git diff --name-only || true; } \
  | awk -F/ "NF>1{print \$1\"/\"\$2}" | sort -u'`

Recent commits on this branch:
!`bash -lc 'git log --pretty=format:"%h %s" "${BASE:-main}..HEAD" 2>/dev/null || true'`

## Optional Context

Provided via $ARGUMENTS

---

## Produce Output (PR Title + WHY)

Using the information above, write:

1. **PR Title** — one line, imperative mood, ≤ 80 chars.
2. **PR Description** — emphasize the WHY (problem, motivation, constraints), then WHAT (key changes), and HOW (approach).
   - Mention affected areas (paths/subsystems), risks, and rollout/backout notes.
   - If changes are app-scoped (e.g., `$ARGUMENTS` looks like an app name), call that out.

**Output exactly in this format:**

**PR Title:** <your concise title here>

**PR Description:**

- **Why:** <rationale / problem / constraints>
- **What:** <summary of changes>
- **How:** <approach & key decisions>
- **Scope:** <notable files/areas>
- **Risks/Migration:** <breakages, flags, rollbacks>
- **Testing:** <how validated; mention staged/unstaged coverage>
