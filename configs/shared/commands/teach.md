---
description: Teach a work item — ticket, issue, PR, or idea — as if onboarding a new engineer; what it is, why it exists, the workflow, realistic examples, and pitfalls
temperature: 0.1
permission:
  write: deny
  edit: deny
  bash:
    "*": deny
    "gh issue view*": allow
    "gh pr view*": allow
---

Explain a work item so a new engineer could pick it up. Don't just summarize the ticket — teach the workflow. The output must stand on its own: a reader with zero project background should understand what the work is, why it matters, and how the flow behaves.

## Usage

```
/teach [JIRA-KEY | github issue/PR number or URL | topic]
```

## Process

### 1. Get the work item

- **Jira key** (`PROJ-123` shape): fetch it with the Jira tools, including comments and linked issues.
- **GitHub issue/PR** (number or URL): fetch it with `gh issue view` / `gh pr view` (use `--comments`).
- **Anything else in `$ARGUMENTS`**: treat it as a topic or idea and work from the conversation plus the codebase.
- **No arguments**: teach the work item or idea currently being discussed in the conversation.

### 2. Ground yourself before writing

Read the code, configs, or docs the work touches. Examples must come from the real system — real commands, real field names, real failure modes — never generic placeholders. This is also what lets you explain jargon instead of repeating it.

### 3. Write the explanation

Use this structure — these sections, in this order:

```markdown
## What is this?

One or two short paragraphs on what the work is about.

## Why does it exist?

What problem are we solving? Who benefits?

## Workflow

Step-by-step, in simple language:

1. User does...
2. Service receives...
3. We validate...
4. We store...
5. We return...

For a bug, the workflow is instead: what happens today, what should
happen, and where the two diverge.

Don't describe the implementation unless it helps explain the workflow.

## Examples

2-5 realistic scenarios, each as:
Scenario: ...
What happens: ...
Expected result: ...

Include edge cases when they matter.

## Things to watch for

Assumptions, edge cases, common mistakes — and anything important the
ticket does NOT specify, named as an open question rather than guessed.

## References

Tickets, PRs, ADRs, or docs for further reading. Include this section
only when you have at least one real source to cite — otherwise omit
it entirely. The explanation must stand on its own without them.
```

## Rules

- Assume the reader is a new engineer with no project background.
- Use simple English and short paragraphs.
- Avoid project jargon unless you explain it the first time it appears.
- Prefer concrete examples over abstract descriptions.
- Don't invent requirements. When the source is silent on something important, surface it in "Things to watch for" — never fill the gap with a guess.
- Never modify files or run commands that change state — this command only explains.
