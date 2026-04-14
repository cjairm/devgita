---
description: Create PR from branch context. Generate title/body from analysis.
temperature: 0.1
tools:
  write: false
  edit: false
  bash:
    "which gh": allow
    "gh auth status": allow
    "git *": allow
    "cat .github/*": allow
    "test -f *": allow
---

Create PR describing overall branch impact. Commits are granular steps; PR shows the complete picture.

## Process

### 1. Gather Context

```bash
git branch --show-current
git log main..HEAD --format="%s%n%b%n---"
git diff main... --stat
```

If stat insufficient for understanding, get full diff: `git diff main...`

### 2. Check PR Template

```bash
test -f .github/PULL_REQUEST_TEMPLATE.md && cat .github/PULL_REQUEST_TEMPLATE.md
```

### 3. Analyze Changes

From commits + diff, understand:
- What changed (files, functions, logic)
- Why it matters (business/technical impact)
- How it works (implementation approach)
- What could fail (risks, edge cases)

### 4. Generate PR Content

**Title (imperative, max 72 chars):**
- Complete sentence as command: "Add X", "Fix Y", not "Adding X"
- Standalone: understandable without description
- Include ticket ID from $ARGUMENTS if provided
- Good: `[JIRA-123] Fix login redirect on OAuth timeout`
- Bad: "Fix bug", "Update code", "Phase 1"

**Body - use template if exists, otherwise:**

1. **Summary**: What changed and why (readable without diving into code)
2. **Problem**: What was broken/limiting, why it matters, relevant links
3. **Solution**: Approach taken, key decisions not visible in code, tradeoffs
4. **Risk**: What could fail, dependencies, breaking changes, what's NOT included
5. **Testing**: How validated, benchmarks if perf-related, note if screenshots needed

**Content guidance:**
- Why over what (code shows what; explain why)
- Include decision rationale not evident in code
- Write for developers years later without your context
- Note areas needing extra review
- Suggest labels (feature/bugfix/refactor)

### 5. Create PR

Check gh CLI and create:

```bash
which gh
gh auth status
git push -u origin $(git branch --show-current)
gh pr create --title "Generated Title" --body "Generated Description"
```

If gh unavailable, output formatted markdown for manual creation.

## Rules

- Always generate title/body from analysis (never use `--fill`)
- Never modify commits
- Template is source of truth when exists
- PRs address single concern; note if part of larger work
- First line standalone for git log/history

#### References
- https://google.github.io/eng-practices/review/developer/
- https://gist.github.com/hcastro/52c5824a747b901c289261518504effb
