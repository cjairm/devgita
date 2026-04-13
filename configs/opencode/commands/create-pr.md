---
description: Create PR from overall branch context. Generate title/body from analysis. Never modify commits.
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
    "gh pr create --title * --body *": allow
    "git remote get-url origin": allow
---

# Create Pull Request

## Execution Steps

### 1. Gather Context (Token-Efficient)

**Start with commits for understanding:**

```bash
git status
git log main..HEAD --format="%s%n%b%n---"
```

**Get diff context** (commits alone don't show the full picture):

```bash
git diff main... --stat
```

**If --stat insufficient for quality understanding**, get full diff:

```bash
git diff main...
```

**Critical**: Never compromise understanding. Get full diff when needed for quality.

### 2. Analyze Changes

Review all gathered context (commits + diff) to understand the **overall branch impact**:

- **What changed**: Files, functions, logic modified across all commits
- **Why it matters**: Business/technical impact of the complete change
- **How it works**: Implementation approach
- **What could fail**: Risks and edge cases

### 3. Check PR Template

```bash
[ -f .github/PULL_REQUEST_TEMPLATE.md ] && cat .github/PULL_REQUEST_TEMPLATE.md || echo "No template"
```

### 4. Generate PR Content

**If template exists:** Fill all sections completely using context.

**If no template:** Use enhanced 5-section format:

1. **Summary (First Paragraph)**: 
   - **What** changed specifically (summarize major changes)
   - **Why** this change is being made (context, problem being solved)
   - Keep it short but informative enough to understand the PR without reading code

2. **Problem/Background**: 
   - What was broken/limiting
   - Include bug numbers, relevant links, or design doc references
   - Explain why this matters (business/technical impact)

3. **Solution**: 
   - What changed and why this approach
   - Mention key implementation decisions not reflected in code
   - If there are shortcomings/tradeoffs, acknowledge them

4. **Risk & Mitigation**: 
   - What could fail and how mitigated
   - Any dependencies or breaking changes
   - What's explicitly NOT included (if relevant)

5. **Testing/Validation**: 
   - How validated
   - Include benchmark results if performance-related
   - UI changes: mention screenshots needed (let user provide)

**Title format:** Imperative mood + ticket ID (if provided via `$ARGUMENTS`)

- **Complete sentence as an order** (imperative)
- **Short, focused, and standalone** - should be understandable without reading the full description
- Example: `[JIRA-123] Fix login redirect error on OAuth timeout`
- Example: `Add user avatar upload with resize and compression`
- Example: `Remove size limit on RPC server message freelist`

**Avoid vague titles:**
- ❌ "Fix bug", "Fix build", "Add patch"
- ❌ "Moving code from A to B", "Phase 1"
- ✅ "Fix race condition in user session cleanup job"
- ✅ "Refactor authentication middleware for better testability"

**Content enhancements:**

- **First paragraph**: Should summarize what + why, allowing readers to understand the PR without diving into details
- **Context is key**: Include decision rationale not visible in code (why this approach vs alternatives)
- **External links**: Add context inline - links may become inaccessible due to permissions/retention
- **UI changes**: Include screenshots/GIFs (mention need, let user provide)
- **Uncertainty**: Call out areas needing extra review
- **Labels**: Suggest relevant labels (feature/bugfix/refactor, frontend/backend)
- **Scope**: Note what's explicitly NOT included if relevant
- **Small PRs deserve attention**: Even small changes need context for future developers
- **Searchability**: Write descriptions that help future developers find this PR years later

### 5. Create PR

**Check if gh CLI is available:**

```bash
which gh
```

**If gh available**, push and create PR with your generated content:

```bash
git push -u origin $(git branch --show-current)
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

- **PR scope**: Describe the overall branch impact, not individual commits. Commits are granular steps; PR shows the complete picture.
- **Always generate**: Always create title/body from your analysis. Never use `--fill`.
- **Never modify commits**: Commits stay as-is. Generate PR content separately.
- **Template**: Source of truth when it exists.
- **Quality**: Get whatever context needed (full diff if necessary).
- **Tokens**: Start cheap (commits, --stat), escalate only when needed.
- **Focus**: PRs should address single concern. Note if this is part of larger work.
- **Ticket refs**: Use `$ARGUMENTS` for ticket IDs (e.g., JIRA-123, #456).
- **First line standalone**: Title must be understandable without reading description (for git log/history).
- **Why over what**: Code shows what changed; description explains why and context.
- **Future-proof**: Write for developers who will read this years later without your context.
- **Decision rationale**: Include "why" decisions not evident in code (alternatives considered, tradeoffs made).
- **Complete sentences**: Use imperative mood for title ("Add feature" not "Adding feature").
- **Avoid vagueness**: Never use "Fix bug", "Update code", "Phase 1" without specifics.

#### References
- https://google.github.io/eng-practices/review/developer/
- https://gist.github.com/hcastro/52c5824a747b901c289261518504effb
