---
description: Check if PR review feedback has been addressed, auto-approve if complete, or trigger focused re-review
temperature: 0.1
permission:
  write: deny
  edit: deny
  bash:
    "*": deny
    "gh auth status": allow
    "gh auth status *": allow
    "gh pr view*": allow
    "gh pr view *": allow
    "gh pr diff*": allow
    "gh pr diff *": allow
    "gh pr review*": allow
    "gh pr review *": allow
    "gh api*": allow
    "gh api *": allow
    "jq *": allow
    "cat *": allow
    "dge fetch-pr-comments *": allow
---

Verify that your previously-left PR review comments have been addressed. Auto-approve if all feedback is resolved, or trigger a focused re-review on pending items.

## Usage

```
/approve-pr owner/repo PR_NUMBER
```

Example: `/approve-pr lever/foundational-ai-client 62`

## Process

### 1. Validate Arguments and Environment

Parse `$ARGUMENTS` to extract:
- `OWNER/REPO`: Repository identifier
- `PR_NUMBER`: Pull request number

```bash
gh auth status
```

If not authenticated, stop and tell user to run `gh auth login`.

### 2. Get PR Context

```bash
gh pr view $PR_NUMBER --repo $OWNER/$REPO --json number,headRefOid,title,url,state,files
```

Extract:
- `COMMIT_SHA`: The headRefOid (HEAD commit SHA)
- `STATE`: Ensure PR is open
- `FILES`: List of changed files (for PR type detection)

### 3. Fetch Existing Review Comments

```bash
dge fetch-pr-comments $OWNER/$REPO $PR_NUMBER /tmp/pr_comments_$PR_NUMBER.json
```

Read and parse the comments file:
```bash
cat /tmp/pr_comments_$PR_NUMBER.json
```

This returns JSON array with:
```json
[
  {
    "path": "src/client.ts",
    "line": 42,
    "body": "Missing error handling here...",
    "author": "username",
    "is_resolved": false
  }
]
```

### 4. Analyze Each Unresolved Comment

For each comment where `is_resolved: false`:

1. **Read the current file** at the specified `path`
2. **Check the line/region** mentioned in the comment
3. **Determine status:**

| Condition | Status |
|-----------|--------|
| Code changed and issue fixed | ADDRESSED |
| Code unchanged, valid reply from author | ADDRESSED |
| Code unchanged, no response | PENDING |
| Code changed but issue still exists | PENDING |

**Skip comments where `is_resolved: true`** - GitHub already marked these resolved.

### 5. Generate Assessment Table

Create a summary:

```
| Status | Location | Issue Summary | Resolution |
|--------|----------|---------------|------------|
| DONE | src/client.ts:42 | Missing error handling | Added try/catch at line 40-45 |
| DONE | README.md:15 | Outdated example | Updated with new API |
| PENDING | src/client.ts:88 | Race condition | Code unchanged |
```

### 6. Decision Point

**IF all comments ADDRESSED:**

```bash
# Post approval
gh pr review $PR_NUMBER --repo $OWNER/$REPO --approve --body "$(cat <<'EOF'
All review feedback has been addressed:

[SUMMARY_TABLE]

LGTM!
EOF
)"
```

Report to user:
```
## PR #$PR_NUMBER Approved

All X review comments addressed.

[Summary table]

Action: Approved via `gh pr review --approve`
```

**IF any comments PENDING:**

1. Report current status
2. Detect PR type
3. Trigger focused re-review (see step 7)

### 7. Focused Re-Review (if pending items exist)

**Detect PR type:**

```bash
gh pr diff $PR_NUMBER --repo $OWNER/$REPO --name-only
```

- Majority `.md` files -> document-reviewer
- Otherwise -> code-reviewer

**Dispatch focused re-review:**

Use Task tool with appropriate subagent type (`document-reviewer` or `code-reviewer`).

Prompt template:
```
You are re-reviewing specific PENDING issues from a previous PR review.

## Context
Repository: $OWNER/$REPO
PR: #$PR_NUMBER

## PENDING Issues to Verify

[List ONLY the pending items with file:line and original comment]

## Your Task

For EACH pending issue:
1. Read the current file at the path
2. Check if the issue is now fixed
3. Report: FIXED (what changed) or STILL_PENDING (why)

DO NOT:
- Do a full PR review
- Comment on already-addressed items
- Add unrelated nitpicks
- Duplicate existing feedback

## Output

| Issue | Location | Status | Finding |
|-------|----------|--------|---------|
| ... | ... | FIXED/STILL_PENDING | ... |

Recommendation: Ready to approve? Yes/No
If No: What specifically remains to be done?
```

### 8. After Re-Review

If re-review finds all issues now fixed:
- Auto-approve the PR
- Report success

If issues still pending:
- Report what remains
- Do NOT approve
- User can fix and run `/approve-pr` again

## Output Format

### All Addressed
```
## PR #62 - Review Status

All 4 review comments addressed!

| Status | Location | Issue | Resolution |
|--------|----------|-------|------------|
| DONE | src/client.ts:42 | Missing error handling | Added try/catch |
| DONE | src/client.ts:88 | Unclear naming | Renamed to `connectionTimeout` |
| DONE | README.md:15 | Outdated example | Updated |
| DONE | tests/client.test.ts:30 | Missing edge case | Added null test |

Action: PR approved via `gh pr review --approve`
PR URL: https://github.com/owner/repo/pull/62
```

### Pending Items
```
## PR #62 - Review Status

2/4 review comments addressed, 2 PENDING

| Status | Location | Issue | Resolution |
|--------|----------|-------|------------|
| DONE | src/client.ts:42 | Missing error handling | Added try/catch |
| DONE | README.md:15 | Outdated example | Updated |
| PENDING | src/client.ts:88 | Race condition | Code unchanged |
| PENDING | tests/client.test.ts:30 | Missing edge case | Not addressed |

Triggering focused re-review on pending items...

[Re-review results]
```

## Error Handling

**No arguments provided:**
```
Usage: /approve-pr owner/repo PR_NUMBER
Example: /approve-pr lever/foundational-ai-client 62
```

**PR not found:**
```
Error: PR #62 not found in owner/repo
Check the PR number and repository name.
```

**No review comments found:**
```
No review comments found on PR #62.
Nothing to verify. Would you like to run a fresh review instead?
- Use code-reviewer subagent for code
- Use document-reviewer subagent for docs
```

**Approval fails:**
```
Error: Could not approve PR. GitHub response: [error]

Manual approval command:
gh pr review 62 --repo owner/repo --approve
```

## Integration Notes

Uses existing agents for re-review:
- `agents/code-reviewer.md` - For code-focused PRs (built-in `dge fetch-pr-comments` duplicate avoidance)
- `agents/document-reviewer.md` - For doc-focused PRs (also has duplicate avoidance)

Both agents automatically fetch existing PR comments and avoid duplicating them.
