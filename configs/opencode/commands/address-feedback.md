---
description: Address PR review feedback - implement suggestions, respond to comments, and resolve them
temperature: 0.2
permission:
  write: allow
  edit: allow
  bash:
    "*": deny
    "gh auth status": allow
    "gh auth status *": allow
    "gh pr view*": allow
    "gh pr view *": allow
    "gh pr diff*": allow
    "gh pr diff *": allow
    "gh api*": allow
    "gh api *": allow
    "jq *": allow
    "cat *": allow
    "dge fetch-pr-comments *": allow
    "git add *": allow
    "git commit *": allow
    "git push*": allow
    "git status": allow
    "git diff*": allow
    "go test *": allow
    "go build *": allow
    "go vet *": allow
    "go fmt *": allow
    "npm test*": allow
    "npm run *": allow
    "make *": allow
---

Address PR review feedback by implementing suggestions, making code changes, responding to comments, and resolving them.

## Usage

```
/address-feedback [owner/repo] [PR_NUMBER]
```

If no arguments provided, uses current branch's PR.

Examples:

- `/address-feedback` - Address feedback on current PR
- `/address-feedback 62` - Address feedback on PR #62
- `/address-feedback cjairm/devgita 62` - Address feedback on specific repo PR

## Process

### 1. Validate Environment

```bash
gh auth status
```

If not authenticated, stop and tell user to run `gh auth login`.

### 2. Get PR Context

If arguments provided:

```bash
gh pr view $PR_NUMBER --repo $OWNER/$REPO --json number,headRefOid,title,url,state,baseRefName,headRefName
```

Otherwise (current branch):

```bash
gh pr view --json number,headRefOid,title,url,state,baseRefName,headRefName,headRepository
```

Extract:

- `PR_NUM`: The PR number
- `COMMIT_SHA`: The headRefOid (HEAD commit SHA)
- `OWNER`: From headRepository.owner.login
- `REPO`: From headRepository.name
- `HEAD_BRANCH`: Branch name for pushing

### 3. Fetch PR Comments

```bash
dge fetch-pr-comments $OWNER/$REPO $PR_NUM /tmp/pr_feedback_$PR_NUM.json
```

Read the comments:

```bash
cat /tmp/pr_feedback_$PR_NUM.json
```

Returns:

```json
[
  {
    "id": "PRRT_kwDOL...",
    "path": "src/client.ts",
    "line": 42,
    "body": "Missing error handling here...",
    "author": "reviewer",
    "is_resolved": false
  }
]
```

### 4. Filter Unresolved Comments

Only process comments where `is_resolved: false`.

If ALL comments are resolved:

```
## All Feedback Addressed!

All review comments on PR #$PR_NUM are already resolved.
Nothing to do.

PR URL: $PR_URL
```

### 5. Categorize Each Comment

For each unresolved comment, determine the action needed:

| Category            | Description                            | Action                         |
| ------------------- | -------------------------------------- | ------------------------------ |
| CODE_CHANGE         | Specific code suggestion or fix needed | Implement the change           |
| QUESTION            | Reviewer asking for clarification      | Reply with explanation         |
| DISCUSSION          | Opinion/design discussion              | Reply with reasoning           |
| ALREADY_DONE        | Change was already made                | Reply with reference + resolve |
| WONT_FIX            | Valid but won't address (with reason)  | Reply with reasoning           |
| NEEDS_CLARIFICATION | Comment unclear, need reviewer input   | Ask for clarification          |

### 6. Process Each Comment

#### For CODE_CHANGE Comments

1. **Read the file** at `path`
2. **Understand the suggestion** from comment body
3. **Implement the change** using Edit tool
4. **Track the change** for commit message

Example:

```
Comment: "Missing error handling for null response"
File: src/client.ts:42

Action: Add try/catch or null check at line 42
```

#### For QUESTION/DISCUSSION Comments

1. **Analyze the question**
2. **Formulate response** with context from code
3. **Prepare reply** (posted in step 8)

#### For ALREADY_DONE Comments

1. **Verify the change** exists in current code
2. **Prepare reply** referencing what was done
3. **Mark for resolution**

### 7. Run Verification (if code changes made)

If any code changes were made, run appropriate verification:

```bash
# For Go projects
go build ./... && go test ./... && go vet ./...

# For Node projects
npm test

# For Make-based projects
make test
```

If tests fail:

- Show the failure
- Attempt to fix
- Re-run verification
- If still failing, report and ask user for guidance

### 8. Commit Changes (if any)

If code changes were made:

```bash
git add -A
git commit -m "Address PR review feedback

Changes:
- [list of changes made]

Addresses review comments on PR #$PR_NUM"
```

Then push:

```bash
git push origin $HEAD_BRANCH
```

### 9. Reply to Comments and Resolve

For each processed comment, post a reply and resolve:

#### Reply to Comment Thread

```bash
gh api graphql -f query='
mutation($threadId: ID!, $body: String!) {
  addPullRequestReviewThreadReply(input: {
    pullRequestReviewThreadId: $threadId
    body: $body
  }) {
    comment {
      id
    }
  }
}' -f threadId="$THREAD_ID" -f body="$REPLY_BODY"
```

#### Resolve the Thread

```bash
gh api graphql -f query='
mutation($threadId: ID!) {
  resolvePullRequestReviewThread(input: {
    threadId: $threadId
  }) {
    thread {
      isResolved
    }
  }
}' -f threadId="$THREAD_ID"
```

### Reply Templates

**For CODE_CHANGE (implemented):**

```
Done! Fixed in commit [SHORT_SHA].

[Brief description of what was changed]
```

**For QUESTION:**

```
[Direct answer to the question]

[Additional context if helpful]
```

**For DISCUSSION:**

```
Good point! [Response to the discussion]

[Reasoning or decision made]
```

**For ALREADY_DONE:**

```
This was already addressed in [commit/line reference].

[Point to specific code if helpful]
```

**For WONT_FIX:**

```
Acknowledged. Keeping current implementation because:

[Reasoning]

Happy to discuss further if needed.
```

**For NEEDS_CLARIFICATION:**

```
Could you clarify what you mean?

[Specific question about the feedback]
```

**Note:** Don't resolve NEEDS_CLARIFICATION comments - leave them open for reviewer response.

## Output

### Summary Report

```
## PR #$PR_NUM Feedback Addressed

### Changes Made
| File | Line | Change | Commit |
|------|------|--------|--------|
| src/client.ts | 42 | Added error handling | abc1234 |
| src/client.ts | 88 | Renamed variable | abc1234 |

### Comments Resolved
| Location | Category | Response |
|----------|----------|----------|
| src/client.ts:42 | CODE_CHANGE | Fixed + resolved |
| README.md:15 | QUESTION | Replied + resolved |
| src/types.ts:30 | ALREADY_DONE | Referenced existing fix + resolved |

### Comments Left Open
| Location | Category | Reason |
|----------|----------|--------|
| src/api.ts:55 | NEEDS_CLARIFICATION | Asked reviewer for more details |

### Verification
- Tests: PASSED
- Build: PASSED
- Lint: PASSED

### Actions Taken
- Committed: abc1234 "Address PR review feedback"
- Pushed to: feature-branch
- Resolved: 4/5 comment threads
- Left open: 1 (needs clarification)

PR URL: https://github.com/owner/repo/pull/62
```

## Decision Flow

```
For each unresolved comment:

  Is it a specific code suggestion?
    YES -> Implement change -> Reply with commit -> Resolve
    NO  -> Continue

  Is it a question?
    YES -> Formulate answer -> Reply -> Resolve
    NO  -> Continue

  Is it already addressed in code?
    YES -> Find evidence -> Reply with reference -> Resolve
    NO  -> Continue

  Is it unclear what's being asked?
    YES -> Ask for clarification -> DO NOT resolve
    NO  -> Continue

  Is it a valid point you disagree with?
    YES -> Reply with reasoning -> Resolve (unless blocking)
    NO  -> Flag for manual review
```

## Error Handling

### No PR found

```
Error: No PR associated with current branch.

Either:
1. Specify PR: /address-feedback owner/repo PR_NUMBER
2. Create PR first: /create-pr
```

### Push fails

```
Error: Could not push changes.

Your changes are committed locally. To push manually:
git push origin $BRANCH
```

### Thread resolution fails

```
Warning: Could not resolve thread for $PATH:$LINE

Manual resolution: Visit $PR_URL and resolve the thread manually.
Continuing with remaining comments...
```

### API rate limit

```
Error: GitHub API rate limit reached.

Wait a few minutes and try again, or resolve comments manually at:
$PR_URL
```

## Best Practices

### DO

- Read the full comment context before acting
- Run tests after code changes
- Provide meaningful replies (not just "Done")
- Reference specific commits or lines in replies
- Ask for clarification when feedback is ambiguous

### DON'T

- Resolve comments without addressing them
- Make unrelated changes
- Ignore test failures
- Resolve NEEDS_CLARIFICATION comments
- Reply with just "Fixed" without context

## Integration

This command pairs well with:

- `/review-pr` - Receive feedback first
- `/approve-pr` - Verify feedback addressed and get approval
- `/create-pr` - Initial PR creation

Typical workflow:

```
1. Create PR: /create-pr
2. Receive review feedback (from reviewer)
3. Address feedback: /address-feedback
4. Request re-review or: /approve-pr
```
