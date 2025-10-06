---
description: Generate a Pull Request based on the current Git branch and changes.
agent: plan
---

## PR Generation Plan

1. **Identify the Current Git Branch:**
   - Use `git rev-parse --abbrev-ref HEAD` to determine the current branch name.

2. **Fetch the Git Diff:**
   - Retrieve the changes between the current branch and the main branch using `git diff origin/main...HEAD`.

3. **Analyze the Changes:**
   - Examine the diff to classify the changes (e.g., feature, bug fix, refactor, dependency update, documentation/chore).

4. **Generate the PR Title:**
   - Based on the analysis, create a concise and descriptive title for the PR.

5. **Draft the PR Description:**
   - Provide a detailed description including:
     - What was changed
     - Why the change was made
     - How it was implemented
     - Any relevant testing instructions
     - Screenshots or GIFs if applicable

6. **Include Relevant Metadata:**
   - Add issue links if the changes address specific issues.
   - Include a checklist for reviewers (e.g., tests added, documentation updated, all tests passing).

7. **Format the PR:**
   - Ensure the PR is formatted in Markdown, suitable for platforms like GitHub or GitLab.

8. **Review and Finalize:**
   - Present the generated **PR Description** for review before submission.
   - Present the generated **PR Title** — one line, imperative mood, ≤ 80 chars.
