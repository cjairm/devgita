<!--
  Sync Impact Report
  ==================
  Version change: 0.0.0 (initial template) → 1.0.0
  Bump rationale: MAJOR — first ratification, all principles newly defined.

  Modified principles:
    - [PRINCIPLE_1..5] (template placeholders) → 7 concrete principles

  Added sections:
    - Core Principles (7 principles replacing 5 placeholders)
    - Technology Stack (replacing SECTION_2)
    - Development Workflow (replacing SECTION_3)
    - Governance (filled)

  Removed sections:
    - None (all template sections filled or replaced)

  Templates requiring updates:
    - .specify/templates/plan-template.md
      Constitution Check section references "[Gates determined based on
      constitution file]" — generic, no update needed. ✅
    - .specify/templates/spec-template.md
      Uses MUST/SHOULD language consistent with constitution. ✅
    - .specify/templates/tasks-template.md
      Phase structure and parallel markers compatible. No principle-
      specific task types required beyond existing structure. ✅
    - .specify/templates/checklist-template.md
      Generic checklist, no constitution-specific references. ✅
    - .specify/templates/agent-file-template.md
      No constitution references. ✅

  Deferred items:
    - None
-->

# Devgita Constitution

## Core Principles

### I. Zero-Dependency Distribution

The user MUST NOT need any pre-installed tools to run devgita. Pre-built
binaries are provided per platform and architecture. Config templates are
embedded in the binary via Go `embed`. The binary handles everything end
to end, including package manager bootstrap (Homebrew on macOS, apt on
Debian/Ubuntu).

Binaries are built locally and uploaded to GitHub Releases. Supported
targets: darwin-arm64, darwin-amd64, linux-amd64.

### II. Platform Parity with Isolation

macOS logic and Debian/Ubuntu logic MUST NOT leak into each other.
Platform-specific code lives behind the factory pattern (`MacOSCommand`,
`DebianCommand`).

Every app MUST have a macOS equivalent and a Debian/Ubuntu equivalent
where applicable (e.g., Raycast on macOS / Ulauncher on Debian,
Aerospace on macOS / equivalent on Debian). If no equivalent exists for
a platform, the app MUST be excluded from that platform's installation
flow entirely.

Configs for the same app MUST be functionally identical across platforms
unless the app is platform-exclusive.

### III. Idempotent and Safe

Every operation MUST be safe to re-run without side effects. Smart
detection distinguishes pre-existing packages from devgita-installed
ones. Devgita MUST NOT uninstall anything it did not install. State is
tracked in `~/.config/devgita/global_config.yaml`.

### IV. Simplicity Over Verbosity

Code, output, scripts, and documentation go straight to the point. No
unnecessary abstractions. No over-engineering. Simple does not mean
skipping audits or quality — it means no waste. If something can be said
in one line, do not use three.

### V. Testability

All platform commands MUST be behind interfaces (`BaseCommandExecutor`,
`Command`). `MockBaseCommand` is used for isolated unit tests. No real
system commands MUST execute during tests. Three isolation levels are
defined in the testing-patterns guide and MUST be followed:

- Level 1: Simple mock (no configuration)
- Level 2: Isolated paths (configuration without shell)
- Level 3: Complete test environment (shell configuration)

### VI. Configuration as Templates

Every app that applies configuration MUST have its config stored as a
Go-embedded template under `configs/`. Templates MUST exist for both
platforms where the app is available. On first run, configs are extracted
from the binary to `~/.config/devgita/`. Template-based shell config
generation is driven by `GlobalConfig` state.

### VII. Audit Before Shipping

Every app module MUST be audited before a feature is considered complete.
The audit covers:

- Template existence in `configs/` for the app
- Platform equivalence (macOS and Debian implementations)
- Config consistency across platforms
- Proper error handling per the error-handling guide
- Test coverage per the testing-patterns guide

No feature ships without passing audit.

## Technology Stack

- **Language**: Go
- **CLI framework**: Cobra
- **Configuration**: YAML (`global_config.yaml`), Go templates
- **Package managers**: Homebrew (macOS), apt/snap (Debian/Ubuntu)
- **Runtime management**: Mise (for Node.js, Go, Python, Rust)
- **Testing**: `go test` with mock interfaces
- **Distribution**: Pre-built binaries via GitHub Releases
- **Embedding**: Go `embed` for config templates

## Development Workflow

1. **Build**: `go build -o devgita main.go`
2. **Test**: `go test ./...`
3. **Lint**: `go vet ./...`
4. **Format**: `go fmt ./...`
5. **Release**: Build binaries locally per target, upload to GitHub
   Releases

Code review MUST verify compliance with all seven principles. Complexity
MUST be justified in writing if it deviates from Principle IV. Runtime
development guidance lives in `AGENTS.md` and `docs/`.

## Governance

This constitution supersedes all other development practices for devgita.
Amendments require:

1. Documentation of the proposed change
2. Rationale explaining why the change is needed
3. Impact assessment on existing code and templates
4. Version bump following semantic versioning:
   - MAJOR: Principle removal or backward-incompatible redefinition
   - MINOR: New principle or materially expanded guidance
   - PATCH: Clarifications, wording, non-semantic refinements

All pull requests and code reviews MUST verify compliance with this
constitution. Violations MUST be resolved before merge.

**Version**: 1.0.0 | **Ratified**: 2026-03-29 | **Last Amended**: 2026-03-29
