---
description: Migrate shell configuration apps to the template-based GlobalConfig system
agent: build
---

# Scope

Target app is defined implicitly by `$ARGUMENTS`.

Relevant paths:
- App code: `internal/tooling/terminal/dev_tools/$ARGUMENTS/`
- App entry file: `internal/tooling/terminal/dev_tools/$ARGUMENTS/$ARGUMENTS.go`
- App documentation: `docs/tooling/terminal/dev_tools/$ARGUMENTS.md` or `docs/apps/$ARGUMENTS.md`

Reference materials:
- Project overview: `@docs/project-overview.md`
- Reference implementation: `@internal/tooling/terminal/dev_tools/autosuggestions/autosuggestions.go`

---

# Objective

Migrate the app’s shell configuration behavior from direct file manipulation to the centralized, template-based **GlobalConfig** system.

The app must integrate cleanly with global shell configuration generation and rely exclusively on GlobalConfig as the source of truth.

---

# Migration Principles

- Shell configuration is **state-driven**, not file-driven
- No app may directly read from or write to shell config files
- GlobalConfig represents persisted disk state, not in-memory app state
- Configuration updates must be transactional and idempotent
- Templates are the only mechanism for shell output

---

# Required Behavioral Changes

The app must:

- Stop modifying shell files directly
- Stop searching, appending, or prepending strings in config files
- Use GlobalConfig to enable, disable, and query shell features
- Trigger shell regeneration via the template system
- Treat configuration as load → update → regenerate → persist

Installation logic remains unchanged and outside the scope of this migration.

---

# Constraints & Rules

- Only files within the current app directory may be modified
- App identity must always be derived from shared constants (never literals)
- Each configuration action must operate on a fresh GlobalConfig instance
- GlobalConfig must never be cached or stored in app structs
- Errors must be surfaced explicitly and never ignored
- Legacy behavior must be fully removed, not layered on top

---

# Expected Method Responsibilities

The following behaviors must be supported:

- **Force configuration**
  - Guarantees the feature is enabled in GlobalConfig
  - Regenerates shell configuration via templates
  - Persists updated state

- **Soft configuration**
  - Checks GlobalConfig state
  - Performs no action if already enabled
  - Delegates to force configuration otherwise

- **Uninstall**
  - Disables the feature in GlobalConfig
  - Regenerates shell configuration via templates
  - Persists updated state
  - If needed add "	// TODO: We still uninstall the app or remove downloaded doc - see `Install`" since it may require extra work, but not part of this command.

All other app behaviors are preserved as-is.

---

# Documentation Alignment

The app documentation must be updated to reflect:

- Use of template-based GlobalConfig management
- Elimination of direct file manipulation
- Support for enabling and disabling via GlobalConfig
- Conceptual explanation of template integration
- The load–modify–regenerate–save lifecycle

Any references to legacy file-based behavior must be removed.

---

# Validation Expectations

A successful migration ensures:

- Tests pass for the app in isolation
  - If test no there, add it
- Shell configuration is generated exclusively via templates
- No legacy file manipulation code remains
- GlobalConfig is the sole authority for shell feature state
- Documentation accurately reflects runtime behavior

---

# Output Expectations

The agent should provide:

1. A concise summary of the migration outcome
2. A high-level description of behavior changes (before vs after)
3. Confirmation that tests and builds succeed
4. A summary of documentation updates
5. Notes on any deviations or edge cases encountered

---

# Architectural Notes

- This migration enforces consistency across all shell-enabled apps
- The pattern is shared and repeatable across the codebase
- Correctness depends on treating GlobalConfig as the single source of truth
- Statelessness at the app level is mandatory
