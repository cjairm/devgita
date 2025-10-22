---
description: Generate only the Pull Request title and description text (AVOID create or submit the PR).
agent: plan
context (if any): $ARGUMENTS
---

## PR Generation Plan

1. **Identify the Current Git Branch:**
   - Use `git rev-parse --abbrev-ref HEAD` to determine the current branch name.

2. **Fetch the Git Diff:**
   - Retrieve the changes between the current branch and the main branch using `git diff origin/main...HEAD`.

3. **Analyze the Changes:**
   - Examine the diff to classify the change type (e.g., feature, bug fix, refactor, dependency update, documentation/chore).

4. **Generate the PR Title:**
   - Create a concise, imperative-style title summarizing the intent of the change.
     - Example: `Fix login redirect after token refresh`.
     - If context with Jira link: `[JIRA-1298] Fix login redirect after token refresh`.

5. **Draft the PR Description:**
   - Focus primarily on **why** the change was made — the problem solved, or context behind the decision.
   - Optionally include:
     - How the issue was identified
     - How it was solved or implemented
     - Any relevant implications, tradeoffs, or testing notes
   - Avoid restating code-level details already visible in the diff.

6. **Include Relevant Metadata:**
   - Add issue links if the changes address specific issues.
   - Include a checklist for reviewers (e.g., tests added, documentation updated, all tests passing).

7. **Format the PR Output:**
   - Markdown formatted for GitHub or GitLab.
   - Output should include:
     - `## Title:`
     - `## Description:`

8. **Review and Finalize:**
   - Present the **PR Title** and **Description** TEXT ONLY for me to review.
   - **AVOID** create, push, or submit the PR automatically.
