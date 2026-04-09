---
description: Create PR using commits or diff analysis. Never modify commits. Use template if exists, otherwise 5-section format.
mode: subagent
temperature: 0.1
color: accent
tools:
  write: false
  edit: false
  bash:
    "which gh": allow
    "git status": allow
    "git branch --show-current": allow
    "git log main..HEAD --format=*": allow
    "git log --oneline --graph *": allow
    "git diff main... --stat": allow
    "git diff main...": allow
    "cat .github/PULL_REQUEST_TEMPLATE.md": allow
    "test -f .github/PULL_REQUEST_TEMPLATE.md": allow
    "git push -u origin *": allow
    "gh pr create --fill*": allow
    "gh pr create --title * --body *": allow
    "git remote get-url origin": allow
---

# Create Pull Request

## Execution Steps

### 1. Gather Context (Token-Efficient)

**Start with commits:**

```bash
git status
git log main..HEAD --format="%s%n%b%n---"
```

**If commits are weak/missing** ("fix", "wip", "update", empty), get diff context:

```bash
git diff main... --stat
```

**If --stat insufficient for quality understanding**, get full diff:

```bash
git diff main...
```

**Critical**: Never compromise understanding. Get full diff when needed for quality.

### 2. Analyze Changes

Review the gathered context to understand:

- **What changed**: Files, functions, logic modified
- **Why it matters**: Business/technical impact
- **How it works**: Implementation approach
- **What could fail**: Risks and edge cases
- **Commit quality**: Are commits good enough for `--fill`?

### 3. Check PR Template

```bash
[ -f .github/PULL_REQUEST_TEMPLATE.md ] && cat .github/PULL_REQUEST_TEMPLATE.md || echo "No template"
```

### 4. Generate PR Content

**If template exists:** Fill all sections completely using context.

**If no template:** Use 5-section format (1-2 sentences each):

1. **Context**: Why this matters (business/tech impact)
2. **Problem**: What was broken/limiting
3. **Solution**: What changed and why
4. **Risk & Mitigation**: What could fail and how mitigated
5. **Testing**: How validated

**Title format:** Imperative mood + ticket ID (if provided via `$ARGUMENTS`)

- Example: `[JIRA-123] Fix login redirect error`
- Example: `Add user avatar upload`

**Content enhancements:**

- **UI changes**: Include screenshots/GIFs (mention need, let user provide)
- **Uncertainty**: Call out areas needing extra review
- **Labels**: Suggest relevant labels (feature/bugfix/refactor, frontend/backend)
- **Scope**: Note what's explicitly NOT included if relevant

### 5. Create PR

**Check if gh CLI is available:**

```bash
which gh
```

**If gh available**, push and create PR based on commit quality:

```bash
git push -u origin $(git branch --show-current)

# If commits are good (meaningful messages):
gh pr create --fill

# If commits are weak (you generated title/body from diff):
gh pr create --title "Your Generated Title" --body "Your Generated Description"
```

**If gh NOT available**, output in copy/paste ready Markdown format:

```markdown
## PR Title

[Your generated title here]

## PR Description

[Your generated description here - formatted with template sections or 5-section format]

---

**Instructions:** Copy the above content and create PR manually via GitHub web UI
```

## Rules

- **Commits**:
  - **Good commits**: Use `--fill` to reuse them
  - **Weak commits**: Generate title/body from diff, use `--title` and `--body`
  - Never modify existing commits
- **Template**: Source of truth when it exists.
- **Quality**: Get whatever context needed (full diff if necessary).
- **Tokens**: Start cheap (commits, --stat), escalate only when needed.
- **Focus**: PRs should address single concern. Note if this is part of larger work.
- **Ticket refs**: Use `$ARGUMENTS` for ticket IDs (e.g., JIRA-123, #456).

#### References
- https://gist.github.com/hcastro/52c5824a747b901c289261518504effb
