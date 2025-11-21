---
description: Create a complete new app module following the standardized app structure, conventions, and documentation patterns.
agent: build
---

# Target

Directory: `internal/apps/$ARGUMENTS/`
File: `internal/apps/$ARGUMENTS/$ARGUMENTS.go`
Docs: `docs/apps/$ARGUMENTS.md`
Reference docs: @docs/project-overview.md

# Task

When generating a new app module:

- Produce the full expected file layout.
- Create the main app file using the required function structure and naming conventions.
- Add placeholders or stubs only — **no implementation logic beyond structure, signatures, and error-handling patterns.**
- Include template documentation in `docs/apps/$ARGUMENTS.md`.
- Create a basic test file containing the required test cases with mocks only.
- Ensure all naming follows `$ARGUMENTS` as the app name.

# Guidelines

- DO NOT reference or modify files outside the new app directory.
- DO NOT assume or add constants or paths unless explicitly mentioned — simply document where they must exist.
- Generate only standardized function skeletons:
  - `New()`
  - `Install()`
  - `ForceInstall()`
  - `SoftInstall()`
  - `ForceConfigure()`
  - `SoftConfigure()`
  - `Uninstall()`
  - `ExecuteCommand()`
  - `Update()`

- Follow the same error-handling guidance as the audit prompt (never ignore errors).
- Keep all logic conceptual (pseudo-code or high-level placeholders).
- Include documentation headers in all files (summary, references, purpose).

# Expected Structure

The generated module must include:

1. **Main App File Skeleton**
   - Correct package name
   - Required imports (only references, not full code)
   - Type definition
   - All expected functions in the standardized order
   - Only pseudo-logic (no implementation)

2. **Test File**
   - Tests for:
     - `New`
     - `Install`
     - `SoftInstall`
     - `ForceConfigure`
     - `SoftConfigure`
     - `ExecuteCommand`

   - Tests must be skipped for:
     - `ForceInstall`
     - `Uninstall`
     - `Update`

   - Use structure-only mocks; no real logic
   - Follow same expectation patterns as audit prompt

3. **Documentation File**
   - App overview
   - Lifecycle: Install → Configure → Execute
   - Table of exported functions
   - Description of expected config paths (symbolic only)
   - References to shared docs

# Docs File Details (`docs/apps/$ARGUMENTS.md`)

If missing:

- Create a complete new doc including:
  - App purpose (inferred)
  - Summary of functions
  - Expected interactions
  - Notes on install/configure/execute flow

If existing:

- Update sections to align with the standard structure
- Ensure naming, order, and expectations match the new layout

# Output Format

When invoked, produce:

1. **Summary of what will be created**
2. **Skeleton file for the new app**
3. **Skeleton test file**
4. **Documentation draft (Markdown)**
5. **Notes or assumptions**
