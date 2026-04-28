---
name: cycle-doc-planner
description: Use before any implementation work to create a bounded cycle document (~3 hours scope). Locks scope, lists all file changes, defines verification steps. Required for bug fixes, features, and refactoring.
---

# Cycle Document Planning

You are creating a **cycle document** - a bounded work package that locks scope and provides everything needed for implementation. Cycles are ~3 hours of focused work with clear deliverables.

**Announce at start:** "I'm using the cycle-doc-planner skill to create a cycle document."

**Save cycle docs to:** `docs/plans/cycles/YYYY-MM-DD-<cycle-name>.md`

---

## When to Use This Skill

- Before implementing any bug fix
- Before adding any new feature
- Before any refactoring work
- When picking up work from a previous session

**Do NOT skip this step.** Cycle docs prevent scope creep, ensure verification, and enable clean handoffs between sessions.

---

## Pre-Planning Investigation

Before writing the cycle doc, you MUST:

1. **Identify root cause** (for bugs) - Don't guess. Read the code. Trace the execution path.
2. **Map affected files** - List every file that will be created, modified, or deleted
3. **Check existing patterns** - How does similar code work in this codebase?
4. **Identify dependencies** - What other code depends on what you're changing?
5. **Find verification method** - How will you prove the fix works?

---

## Cycle Document Template

```markdown
# Cycle: [Short Descriptive Name]

**Date:** YYYY-MM-DD
**Estimated Duration:** ~X hours
**Status:** Draft | In Progress | Complete | Blocked

---

## 1. Domain Context

Brief background for someone unfamiliar with this area:
- What system/feature does this relate to?
- What problem exists or what's being added?
- Link to any relevant docs, tickets, or prior discussions

---

## 2. Engineer Context

What the implementer needs to know:
- Relevant files and their purposes
- Key functions/types involved
- Testing patterns used in this area
- Commands to run tests

---

## 3. Objective

One clear sentence: what does "done" look like?

Example: "Fix the fd-find package name so `dg install` succeeds on macOS without Homebrew warnings."

---

## 4. Scope Boundary

### In Scope
- [ ] Specific deliverable 1
- [ ] Specific deliverable 2
- [ ] Specific deliverable 3

### Explicitly Out of Scope
- Thing that might seem related but we're NOT doing
- Another thing we're deferring

**Scope is locked.** If you discover something out of scope is needed, document it for a future cycle.

---

## 5. Implementation Plan

### File Changes

| Action | File Path | Description |
|--------|-----------|-------------|
| Modify | `path/to/file.go:45` | Change X to Y |
| Create | `path/to/new_file.go` | New component for Z |
| Delete | `path/to/old_file.go` | No longer needed |

### Step-by-Step

Each step should be completable in 5-15 minutes:

#### Step 1: [Action]
- What to do (be specific)
- Expected outcome
- How to verify this step worked

#### Step 2: [Action]
- What to do
- Expected outcome
- How to verify

[Continue for all steps...]

---

## 6. Verification Plan

### Automated Verification
```bash
# Commands that must pass
go test ./path/to/package/...
go build -o devgita main.go
go vet ./...
```

### Manual Verification
1. Step to manually test the change
2. Expected behavior to observe
3. How to confirm fix worked

### Regression Check
- What existing functionality could break?
- How to verify it still works?

---

## 7. Risks & Trade-offs

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Description of risk | Low/Med/High | How to handle |

### Trade-offs Made
- Why this approach vs alternatives
- What we're accepting

---

## 8. Cross-Model Review Notes

Space for review feedback when handing off between sessions or models:

- [ ] Root cause confirmed?
- [ ] All affected files identified?
- [ ] Verification steps are executable?
- [ ] Scope is appropriately bounded?

**Reviewer notes:**
(Fill in during review)
```

---

## Cycle Doc Quality Checklist

Before finalizing, verify:

- [ ] **Root cause identified** - For bugs, the actual cause is documented, not just symptoms
- [ ] **All files listed** - Every file to create/modify/delete is in the table
- [ ] **Line numbers included** - For modifications, include approximate line numbers
- [ ] **Steps are atomic** - Each step is one action, completable in 5-15 minutes
- [ ] **Verification is executable** - Commands and manual steps are specific and runnable
- [ ] **Scope is bounded** - Clear what's in and out of scope
- [ ] **No placeholders** - No "TBD", "TODO", or vague language

---

## After Creating the Cycle Doc

1. **Save the document** to `docs/plans/cycles/YYYY-MM-DD-<cycle-name>.md`
2. **Present to user** for review and approval
3. **Wait for approval** before implementing
4. **Track progress** by checking off steps as completed

---

## Example: Bug Fix Cycle Doc

```markdown
# Cycle: Fix fd-find Package Name on macOS

**Date:** 2025-04-17
**Estimated Duration:** ~30 minutes
**Status:** Draft

## 1. Domain Context

devgita uses `MaybeInstallPackage` to install CLI tools. On macOS, Homebrew 
package names sometimes differ from the display/alias names. The `fd` tool 
is being installed with alias "fd-find" but Homebrew expects just "fd".

## 2. Engineer Context

- `internal/tooling/terminal/dev_tools/fdfind/fdfind.go` - fd-find app module
- `pkg/constants/constants.go` - Package name constants
- Pattern: `MaybeInstallPackage(constant, displayAlias)` - displayAlias is optional

## 3. Objective

Make `dg install` install `fd` on macOS without warnings about package name mismatch.

## 4. Scope Boundary

### In Scope
- [ ] Fix the MaybeInstallPackage call in fdfind.go
- [ ] Verify fix works on macOS

### Explicitly Out of Scope
- Debian package name mapping (separate concern)
- Other package name issues (separate cycles)

## 5. Implementation Plan

### File Changes

| Action | File Path | Description |
|--------|-----------|-------------|
| Modify | `internal/tooling/terminal/dev_tools/fdfind/fdfind.go:46` | Remove "fd-find" alias parameter |

### Step-by-Step

#### Step 1: Modify SoftInstall method
- Open `internal/tooling/terminal/dev_tools/fdfind/fdfind.go`
- Line 46: Change `f.Cmd.MaybeInstallPackage(constants.FdFind, "fd-find")` 
  to `f.Cmd.MaybeInstallPackage(constants.FdFind)`
- The constant `FdFind` already equals `"fd"` which is correct for Homebrew

## 6. Verification Plan

### Automated Verification
```bash
go build -o devgita main.go
go test ./internal/tooling/terminal/dev_tools/fdfind/...
```

### Manual Verification
1. Run `./devgita install --only terminal` on macOS
2. Observe no warning about "fd-find" package name
3. Verify `fd --version` works after installation

## 7. Risks & Trade-offs

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Breaking Debian install | Low | FdFind constant handles both platforms |

## 8. Cross-Model Review Notes

- [x] Root cause confirmed: Alias param unnecessary, constant is correct
- [x] All affected files identified: Just one file
- [x] Verification steps are executable
- [x] Scope is appropriately bounded
```

---

## Handling Multiple Related Issues

If investigating reveals multiple issues (like 4 bugs found during testing):

1. **Create separate cycles** for truly independent fixes
2. **Batch related fixes** into one cycle if they're in the same area and < 3 hours total
3. **Document dependencies** if cycles must be done in order

For batched fixes, the cycle doc should have separate sections in the Implementation Plan for each fix, but shared Verification Plan.
