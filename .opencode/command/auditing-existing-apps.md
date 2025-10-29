---
description: Audit and refactor existing apps to match the updated conventions and structure.
agent: build
---

# Target

Directory: `internal/apps/$ARGUMENTS/`
File: `internal/apps/$ARGUMENTS/$ARGUMENTS.go`
Docs: `docs/apps/$ARGUMENTS.md`
Reference docs: @docs/project-overview.md

# Task

- Audit existing functions in the app file.
- Refactor and standardize structure to match the conventions below.
- Generate or update a documentation file at `docs/apps/$ARGUMENTS.md` describing:
  - App purpose and overview (inferred from @docs/project-overview.md or app code context).
  - Function responsibilities and high-level behaviors (pseudo-descriptions only).
  - Expected interactions between functions.
  - Notes on constants/paths relevant to this app.

# Guidelines

- AVOID modify files outside the current app directory.
- AVOID add implementation logic or code specifics — only describe function signatures, expected outputs, and logical purpose.
- Follow pseudo-code / function-skeleton conventions only.
- Add or rename functions as specified below.
- Include short documentation headers at the top of each file (with relevant references if applicable).

# Expected structure

// Relevant docs (possible look on internet) for basic commands (not only those used in this app)

package $ARGUMENTS

imports

- New() // Same logic
- Install() // Calls regular Install - Logic should exist already
- ForceInstall() // Must call Uninstall() first, then Install() - Logic can be inferred
- SoftInstall() // Checks if app exists first; equivalent to MaybeInstall()
- ForceConfigure() // Previously Setup(); overwrites configuration
- SoftConfigure() // Previously MaybeSetup(); configures if not already set
- Uninstall() // New placeholder; do not implement logic
- ExecuteCommand() // Previously Run(); executes primary command
- Update() // New placeholder; do not implement logic

# Additional notes

- Do not remove or alter any existing constants or paths under pkg/constants or pkg/paths.
- If references such as `$ARGUMENTSConfigAppDir` or `$ARGUMENTSConfigLocalDir` already exist, leave them unchanged.
- All renames (e.g., Setup → ForceConfigure) and new additions (e.g., ForceInstall, Uninstall) should occur **only within the app file**.
- Make sure to handle errors (Never skip errors)
  AVOID:
  ```
  _ = app.X()
  return app.Y()
  ```
  INSTEAD do:
  ```
  err = = app.X()
  if err != nil {
    return fmt.Errorf("failed to X app: %w", err)
  }
  return app.Y()
  ```

# Tests

For each app under `internal/apps/$ARGUMENTS/`, create a corresponding test file.

Tests must exist for:

- New
- Install
- SoftInstall
- ForceConfigure
- SoftConfigure
- ExecuteCommand

Tests to be SKIPPED for:

- ForceInstall (SKIP test)
- Uninstall (SKIP test)
- Update (No test)

Use existing mocks located at: `@internal/commands/mock.go`
Follow existing test patterns where possible, focusing on structure and expectations, not logic details.

# Docs File (`docs/apps/$ARGUMENTS.md`)

If docs/apps/$ARGUMENTS.md does not exist:

- Create it with:
  - App overview (inferred purpose) - Lifecycle summary (Install, Configure, Execute) - Table of exported functions and their high-level purpose
- If it exists:
  - Update or append missing sections, ensuring consistency with the new function layout.

# Output

Provide for each app:

1. Summary of the audit (existing vs expected functions).
2. Adjusted pseudo-structure following the standardized format.
3. List of test cases (function names and output expectations).
4. Documentation diff or new doc draft (in Markdown).
5. Notes on anything missing, unclear, or inferred.

```

```
