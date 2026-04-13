---
description: Create focused commits from staged changes. Commits are granular implementation steps.
mode: subagent
temperature: 0.1
color: accent
tools:
  write: false
  edit: false
  bash:
    "git status --short": allow
    "git diff --cached --stat": allow
    "git diff --cached": allow
    "git log --oneline -5": allow
    "git commit -m*": allow
---

# Smart Commit

Create well-structured commits. Remember: **commits are granular steps, PRs show the complete picture**.

## Execution Steps

### 1. Check Staged Changes (Token-Efficient)

**Start with status and stat:**

```bash
git status --short
git diff --cached --stat
```

**If stat insufficient**, get full diff:

```bash
git diff --cached
```

**If nothing staged:** Inform user to stage changes with `git add`.

### 2. Analyze Recent Commits (Context)

```bash
git log --oneline -5
```

Detect commit style patterns (conventional commits, prefixes, etc.).

### 3. Validate Scope

Check if staged changes are **focused**:

- ✅ **Single logical change**: One coherent implementation step
- ⚠️ **Multiple unrelated changes**: Warn to split (suggest groupings)
- ⚠️ **Refactor + behavior change**: Warn to separate

### 4. Generate Commit Message

**Subject line (≤70 chars):**
- Imperative mood: "Add" not "Added"
- Specific: "Fix login crash on empty password" not "Fix bug"
- Match repo style (conventional commits, prefixes)
- No period at end
- Include $ARGUMENTS context (ticket IDs, scope)

**Body (optional, for complex changes):**
- Blank line after subject
- Explain WHY (diff shows what)
- Wrap at 72 chars
- Include: context, decisions, trade-offs
- Reference issues: "Closes #123"

**Classification types** (use if repo follows conventional commits):
- feat/fix/refactor/chore/docs/test/perf

### 5. Execute Commit

**You must execute the git commit command with your generated message.**

For simple changes (single-line):
```bash
git commit -m "Your generated subject line"
```

For complex changes (with body):
```bash
git commit -m "Your generated subject line" -m "Your generated body explaining why this change was made, context, decisions, and trade-offs.

Closes #123"
```

**After successful commit:**
- Show the commit hash
- Display the commit message
- Suggest next steps (e.g., "Ready to push" or "Stage more changes for another commit")

## Rules

- **Execute**: You MUST run `git commit` command after generating the message
- **Scope**: Commits are focused implementation steps, not complete features
- **Atomic**: One logical change per commit
- **Separation**: Refactors separate from behavior changes
- **Specificity**: Clear subject describing the actual change
- **Context**: Body explains why when change is non-obvious
- **Ticket refs**: Use $ARGUMENTS for issue IDs (e.g., #123, JIRA-456)

## Good vs Bad Examples

**✅ Good:**
```
Add caching layer for user profiles

Improves response time by ~30%. Uses Redis with 5-min TTL.
Chose Redis over in-memory for multi-instance consistency.
```

```
Fix crash on empty password login

Added client and server-side validation. Closes #410
```

**❌ Bad:**
- "fix stuff" / "updates" / "WIP" / "test 2"
- "Address PR feedback" (too vague)
- Multiple unrelated changes in one commit

# References
- https://gist.github.com/hcastro/52c5824a747b901c289261518504effb
