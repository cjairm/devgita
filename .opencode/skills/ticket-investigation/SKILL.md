---
name: ticket-investigation
description: Investigate a Jira ticket, analyze the relevant codebase architecture and data flow, and produce production-ready fix proposals with safety audits. Use this skill when the user provides a Jira ticket URL and wants to understand what's going on, plan a fix, or get a technical analysis of an issue. Also use when someone says "investigate this ticket", "look into this bug", "analyze this issue", "plan a fix for this", or provides a Jira/issue-tracker link and asks for help.
---

# Ticket Investigation & Fix Planning

You are a Technical Lead investigating a bug or issue ticket in an unfamiliar part of a codebase. Your job: understand the issue deeply, map the relevant architecture, and propose safe, production-ready fixes.

The user provides a ticket URL (Jira or similar issue tracker).

---

## Step 1: Fetch the Ticket

Extract the ticket key from the URL (e.g., `PROJ-123` from `https://company.atlassian.net/browse/PROJ-123`).

Fetch the full ticket via the Jira MCP tools (`getJiraIssue`) or equivalent. If the ticket is unreachable or you lack permissions, report the specific error and stop.

**Save verbatim.** Create `tickets/<TICKET-KEY>/ticket_raw.md` and paste the ticket summary, description, and all relevant fields exactly as they appear. Do not paraphrase or summarize — altered ticket text has caused incorrect fix planning in the past. After saving, re-read the file to confirm nothing was lost.

If the ticket description is vague or lacks reproduction steps, note the gaps explicitly — these become assumptions you'll need to flag in your fix proposals.

## Step 2: Trace the Architecture

The goal is an explanation clear enough that someone new to this part of the codebase can understand the fix.

### How to explore unfamiliar code

Think through this systematically:

1. **Start from symptoms.** What does the ticket describe? An error message, unexpected behavior, a data inconsistency? Search the codebase for those exact strings — error messages, status codes, field names mentioned in the ticket.

2. **Find the entry point.** Which API endpoint, UI action, background job, or event triggers the affected behavior? Search for route definitions, controller methods, or queue consumers that match.

3. **Trace the data flow.** From the entry point, follow the code path: what services, repositories, or modules does it call? What database tables or external APIs does it touch? Read each file along the path.

4. **Identify boundaries and dependencies.** What shared libraries, middleware, configuration, or other services interact with this flow? Are there event listeners, webhooks, or cron jobs that operate on the same data?

5. **Check for related tests.** Look for existing test files covering the affected area — they reveal expected behavior and edge cases the original author considered.

### Self-check before writing

After your initial trace, pause and verify:

- Does the flow you traced actually explain the ticket's symptoms?
- Are there components you found references to but didn't investigate?
- Could the issue originate upstream or downstream of where you looked?

## Step 3: Document the Architecture

Write `tickets/<TICKET-KEY>/architecture_analysis.md`:

```markdown
# Architecture Analysis: <TICKET-KEY>

## Affected Components

List each file, module, service, or table involved, with a one-line description of its role in the flow.

## Data Flow

Describe the end-to-end path: entry point -> service calls -> data access -> response. Include the specific function/method names you traced through.

## Dependencies & Interactions

Other systems, shared libraries, middleware, or background jobs that touch this area.

## Key Observations

Anything surprising: race conditions, implicit assumptions, missing validations, technical debt, or patterns that affect how the fix should be approached.
```

Verify technical accuracy before saving.

## Step 4: Research and Propose Fixes

Develop at least two production-ready options (e.g., a minimal targeted fix vs. a more thorough refactor) so the team can make an informed tradeoff decision.

For each option, think through:

- **Exact changes** — which files, which functions, what modifications
- **Breaking changes** — does this alter API contracts, database schemas, or public interfaces?
- **Performance** — could this introduce high-CPU loops, N+1 queries, or lock contention?
- **Regression risk** — could this break existing behavior? Are there downstream consumers or dependent features?
- **Blast radius** — what other features, services, or teams could be affected?

If any of these checks surface concerns, investigate further before writing the proposal. A confident-sounding but wrong safety assessment is worse than an honest "I'm uncertain about X."

Write `tickets/<TICKET-KEY>/fix_proposals.md`:

```markdown
# Fix Proposals: <TICKET-KEY>

## Option 1: <Short Title>

### Description

What the fix does and why it addresses the root cause.

### Changes Required

Specific files and modifications, with enough detail for a developer to implement.

### Safety Audit

- **Breaking changes:** [assessment]
- **Performance impact:** [assessment]
- **Regression risk:** [assessment]
- **Blast radius:** [assessment]

### Tradeoffs

Pros and cons of this approach.

---

## Option 2: <Short Title>

[Same structure]

---

## Recommendation

Which option to pursue and why, including any assumptions that should be validated first.
```

Flag uncertainties explicitly rather than glossing over them.
