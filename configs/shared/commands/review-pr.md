---
description: Post review feedback to PR as a single review with inline comments
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
    "gh api*": allow
    "gh api *": allow
    "jq *": allow
    "dge fetch-pr-comments *": allow
---

Post review feedback from conversation context to a PR as a **single cohesive review** with inline comments attached.

**Prerequisite**: You should already have review feedback in the conversation (from document-reviewer, code-reviewer agents, or prior analysis). This command posts that feedback to GitHub.

## Process

### 1. Validate Environment

```bash
gh auth status
```

If not authenticated, stop and tell user to run `gh auth login`.

### 2. Get PR Context

```bash
gh pr view --json number,headRefOid,title,url,headRepository
```

If $ARGUMENTS contains a PR number, use: `gh pr view $ARGUMENTS --json number,headRefOid,title,url,headRepository`

Extract:
- `PR_NUM`: The PR number
- `COMMIT_SHA`: The headRefOid (HEAD commit SHA)
- `OWNER`: From headRepository.owner.login
- `REPO`: From headRepository.name

### 3. Download Existing PR Comments (Avoid Duplications)

Before posting new comments, download existing review comments to avoid duplicating feedback:

```bash
dge fetch-pr-comments "$OWNER/$REPO" $PR_NUM existing_comments.json
```

This fetches all existing review thread comments and saves them to `existing_comments.json`. The file contains:
```json
[
  {
    "path": "docs/plan.md",
    "line": 15,
    "body": "Comment text...",
    "author": "username",
    "is_resolved": false
  }
]
```

**Check for duplicates** by comparing:
- Same `path` and `line` combination
- Similar comment `body` content (fuzzy match for substantial overlap)

**Actions:**
- Skip posting comments that already exist (same location + similar content)
- Include in review body a note like: "Note: Skipped 2 comments already posted by previous reviews"
- Only post NEW comments that don't duplicate existing feedback

### 4. Get the PR Diff

To post inline comments, we need to know which lines are in the diff. Get the diff:

```bash
gh pr diff $PR_NUM
```

Parse this to understand:
- Which files are changed
- Which line numbers are part of the diff (new lines on RIGHT side)
- The "position" in the diff for each line (needed for API)

**Important**: The `line` parameter in the API refers to the line number in the file, but only lines that are part of the diff can receive comments.

### 5. Parse Review Feedback from Context

Look at the conversation history for review feedback. Extract:

**For inline comments**, identify items with specific file paths and line numbers:
- `path/to/file.ext:123` - Issue description
- Items under "Concerns / Gaps" or "CRITICAL/IMPORTANT/MINOR" sections that reference specific locations

**Cross-check with existing comments** (from step 3):
- For each inline comment, check if similar feedback already exists at same location
- Skip duplicate comments
- Track skipped count for summary

**For general summary**, use:
- Summary section
- Strengths section  
- Overall recommendation/risk rating
- Questions for author

### 6. Determine Review Event Type

Based on review feedback:
- **APPROVE**: Low risk, no critical/blocking issues
- **REQUEST_CHANGES**: High risk, critical issues, or blocking problems
- **COMMENT**: Medium risk, suggestions but no blockers

### 7. Post Review with Inline Comments (Single API Call)

Use `gh api` to post the review with all **NEW** inline comments in one atomic request (excluding duplicates found in step 3):

```bash
gh api \
  --method POST \
  -H "Accept: application/vnd.github+json" \
  "/repos/OWNER/REPO/pulls/PR_NUM/reviews" \
  -f commit_id="COMMIT_SHA" \
  -f event="REQUEST_CHANGES" \
  -f body="## Review Summary

### Overview
[Summary from review]

### Strengths
[Strengths from review]

### General Concerns
[Concerns without specific line references]

### Questions
[Questions for author]

### Risk Assessment
[Risk rating and reasoning]

---
*Review includes X inline comments on specific locations*
*Note: Skipped Y duplicate comments from existing reviews*" \
  --input - <<'EOF'
{
  "comments": [
    {
      "path": "docs/plan.md",
      "line": 15,
      "body": "**[CRITICAL]** Missing rollback strategy\n\nThis section defines deployment steps but doesn't address what happens if deployment fails.\n\n**Suggestion:** Add a rollback section with specific steps."
    },
    {
      "path": "docs/plan.md", 
      "line": 42,
      "body": "**[IMPORTANT]** Success criteria unclear\n\nThe success criteria should be measurable and specific."
    }
  ]
}
EOF
```

**Comment Parameters:**
- `path`: File path relative to repo root
- `line`: Line number in the file (must be part of the diff)
- `body`: Comment text with severity and suggestion
- `side`: Optional, defaults to "RIGHT" (new/modified lines)
- `start_line`: Optional, for multi-line comments

**For line ranges** (e.g., `file.ext:40-45`), add `start_line`:
```json
{
  "path": "docs/plan.md",
  "start_line": 40,
  "line": 45,
  "body": "Comment spanning multiple lines"
}
```

### Alternative: Using jq to build JSON

If you have many comments, build the JSON dynamically:

```bash
# Build comments array
COMMENTS=$(jq -n '[
  {path: "docs/plan.md", line: 15, body: "**[CRITICAL]** Issue 1"},
  {path: "docs/plan.md", line: 42, body: "**[IMPORTANT]** Issue 2"}
]')

# Post review
gh api \
  --method POST \
  -H "Accept: application/vnd.github+json" \
  "/repos/OWNER/REPO/pulls/PR_NUM/reviews" \
  -f commit_id="$COMMIT_SHA" \
  -f event="REQUEST_CHANGES" \
  -f body="Review summary here" \
  --input <(echo "{\"comments\": $COMMENTS}")
```

## Output

After posting, confirm to user:

```
## Review Posted to PR #123

**Review Status:** REQUEST_CHANGES

**Inline Comments:** 4 new comments attached to review
- docs/plan.md:15 - [CRITICAL] Missing rollback strategy
- docs/plan.md:42 - [IMPORTANT] Unclear success criteria
- docs/plan.md:78 - [MINOR] Consider adding diagram
- docs/plan.md:95 - [MINOR] Typo in section header

**Skipped Duplicates:** 2 comments already exist from previous reviews
- docs/plan.md:30 - Similar feedback already posted
- docs/plan.md:55 - Already addressed in existing comments

**PR URL:** https://github.com/org/repo/pull/123

All new comments posted as a single cohesive review.
```

## Mapping Review Output to PR Comments

| Review Section | Where it Goes |
|----------------|---------------|
| Issues with `file:line` references | `comments` array (inline) |
| Summary | Review `body` |
| Strengths | Review `body` |
| Concerns (no line ref) | Review `body` |
| Questions | Review `body` |
| Risk Rating | Review `body` + determines `event` |

## Error Handling

### Duplicate Comment Detection

Compare new comments against existing comments from step 3:

1. **Exact match**: Same path, line, and substantially similar body text → Skip
2. **Location match**: Same path and line, different body → Consider if it adds new value
3. **Keep statistics**: Track skipped count to report to user

**Similarity threshold**: Consider 70%+ text overlap as duplicate

### Lines Not in Diff

If a line is not in the diff, the entire API call will fail. To handle this:

1. First, check which lines are in the diff
2. Only include comments for lines that are in the diff
3. Move comments for lines NOT in diff to the general review body

```
Note: Could not post inline comment for docs/plan.md:15 (line not in diff)
Including in general review body instead.
```

## Example

**Given this review feedback in context:**
```
## Summary
The plan is well-structured but missing critical rollback strategy.

## Strengths
- Clear problem definition
- Good task breakdown

## Concerns / Gaps
- Missing rollback strategy at `docs/plan.md:15-20`
- Success criteria unclear at `docs/plan.md:42`
- No monitoring plan mentioned

## Risk Rating
Medium - rollback gap is significant
```

**Posts single review with:**

**Review body:**
> ## Review Summary
> 
> ### Overview
> The plan is well-structured but missing critical rollback strategy.
> 
> ### Strengths
> - Clear problem definition
> - Good task breakdown
> 
> ### General Concerns
> - No monitoring plan mentioned
> 
> ### Risk Assessment
> **Medium** - rollback gap is significant
> 
> ---
> *Review includes 2 inline comments on specific locations*

**Inline comments attached:**
1. `docs/plan.md:15` - **[CRITICAL]** Missing rollback strategy
2. `docs/plan.md:42` - **[IMPORTANT]** Success criteria unclear

**Event:** `REQUEST_CHANGES`

## Benefits of Single Review API

1. **One notification** instead of many separate comment notifications
2. **Atomic operation** - all comments posted together or none
3. **Proper review status** - can set APPROVE/REQUEST_CHANGES with comments
4. **Better UX** - comments appear as part of a cohesive review in GitHub UI
5. **Threaded** - all inline comments are grouped under one review
6. **No duplicates** - existing comments checked before posting to avoid redundancy

## Duplicate Detection Logic

### Fetching Existing Comments

The `dge fetch-pr-comments` command uses GitHub GraphQL API to fetch:
- All review thread comments (first 100 threads)
- For each thread: path, line, body, author, resolution status
- Filters out URL-only comments (auto-generated references)

### Comparison Algorithm

For each new comment candidate:

```bash
# Pseudocode for duplicate detection
for new_comment in new_comments:
  for existing in existing_comments:
    if new_comment.path == existing.path and new_comment.line == existing.line:
      similarity = calculate_text_similarity(new_comment.body, existing.body)
      if similarity > 0.70:  # 70% threshold
        mark_as_duplicate(new_comment)
        break
```

**Text similarity**: Use simple word overlap or Levenshtein distance to compare comment bodies.

### Handling Duplicates

- **Skip**: Don't include in `comments` array for API call
- **Report**: Include in review body summary
- **Log**: Show user which comments were skipped and why
