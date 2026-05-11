# CLAUDE.md — Devgita Development Guide

⚠️ **DOCUMENTATION CURRENCY:** This file is the source of truth for development practices. Keep it up to date as patterns and decisions evolve. Stale documentation is worse than no documentation—it misleads contributors. If you change how we do something, update this file and/or linked files in the same PR.

---

## 1. What this is

Devgita is a cross-platform development environment manager that automates installation and configuration of terminal tools, programming language runtimes, database systems, and desktop applications on macOS and Debian/Ubuntu.

**Core functionality:**
- **Installation automation** — `dg install` with interactive category selection
- **Category-based setup** — terminal tools, languages, databases, desktop apps
- **Cross-platform support** — Single command syntax works on macOS and Linux
- **Smart state tracking** — `global_config.yaml` tracks what was installed by devgita
- **Configuration templates** — Embedded configs applied consistently across machines
- **Idempotent operations** — Safe to re-run; detects existing packages

For planned features and roadmap, see [ROADMAP.md](ROADMAP.md)

---

## 2. Source of truth

Read these **in order** before starting work:

| File                    | Governs                                                                  |
| ----------------------- | ------------------------------------------------------------------------ |
| `docs/spec.md`          | What features exist, how they work, edge cases, testing strategy         |
| `CLAUDE.md` (this file) | Development practices, tech stack, architecture patterns, code standards |
| `docs/decisions/`       | Individual architectural decisions and their rationale                   |
| `docs/plans/cycles/`    | Current cycle scope and priorities (always check if a cycle is active)   |
| `CONTRIBUTING.md`       | Setup, build, test, and release workflows                                |
| `ROADMAP.md`            | Planned features, future commands, open questions                        |

---

## 3. Product principles

1. **Zero-dependency installer** — Installation happens via shell script alone; no pre-installed tools required beyond bash/curl
2. **Idempotent operations** — Running `dg install` twice produces the same result; safe to re-run
3. **Cross-platform consistency** — Same command syntax and behavior on macOS and Linux; platform differences are transparent to users
4. **Configuration persistence** — User edits to configs are never overwritten; new installs preserve existing customizations
5. **Modular architecture** — Each app is independent; failures in one app don't cascade to others
6. **Transparent state** — All installation state tracked in `~/.config/devgita/global_config.yaml`; users can inspect what was installed

---

## 4. Non-negotiable rules

Hard constraints that override all other considerations:

### Security

- Never execute arbitrary downloaded code without verification
- All shell scripts (`install.sh`, embedded configs) must be reviewed before execution
- Credentials and secrets must never be stored in configs or committed to git
- User input must always be validated before use (especially paths, command arguments)

### Data Integrity

- Installation state must be atomic: either complete or fully roll back (no partial installations)
- User home directory must never be assumed writable in global locations; respect XDG Base Directory if needed
- Config files installed by devgita must be distinguishable from user edits (version markers, checksums)

### Platform Support

- macOS 13+ (Ventura or newer) must be supported; don't use features that break on older versions
- Debian 12+ (Bookworm) and Ubuntu 24+ must be supported; test on both
- Only amd64 and arm64 architectures; drop support only with major version bump

---

## 5. Tech stack

| Layer                       | Technology                    | Notes                                                                                               |
| --------------------------- | ----------------------------- | --------------------------------------------------------------------------------------------------- |
| **Language**                | Go 1.23+                      | stdlib, no cgo where possible (cross-compilation)                                                   |
| **Build System**            | Make                          | See Makefile for targets                                                                            |
| **CLI Framework**           | Cobra                         | Used in `cmd/` for command structure                                                                |
| **Config Format**           | YAML (`gopkg.in/yaml.v3`)     | State stored in `~/.config/devgita/global_config.yaml`                                              |
| **Config Generation**       | Go `text/template` + `embed`  | Templates in `configs/` embedded at compile time                                                    |
| **Package Manager**         | Homebrew (macOS), APT (Linux) | With package name translation (see `pkg/constants/package_mappings.go`)                             |
| **Logging**                 | Custom (zap-like logger)      | Initialized with `logger.Init(verbose)`                                                             |
| **Testing**                 | Go `testing` package          | Unit tests in `*_test.go` alongside code                                                            |
| **Installation Strategies** | Strategy pattern              | In `internal/commands/debian_strategies.go` (AptStrategy, PPAStrategy, InstallScriptStrategy, etc.) |
| **CI/CD**                   | GitHub Actions                | `.github/workflows/release.yml` builds multiplatform binaries on git tag                            |

---

## 6. Implementation behavior

### Coding standards

- Follow [Effective Go](https://golang.org/doc/effective_go) conventions
- Naming: camelCase for functions/variables, PascalCase for exports
- Run `go fmt` before committing (make lint does this)
- Comments explain WHY, not WHAT (code should be self-documenting)
- Never ignore errors; always handle or return them explicitly

### Logger usage

- Initialize once at startup: `logger.Init(verbose)` in cmd/root.go
- Use logger for all output, not println/fmt.Print
- Log levels: error (always), warn (important), info (user actions), debug (detailed)

### Error handling

**See [docs/guides/error-handling.md](docs/guides/error-handling.md) for detailed patterns.**

Key principles:
- Always check errors: `if err != nil { return err }` or `logger.Fatal(err)`
- Use `MaybeExitWithError()` for user-facing errors
- Provide actionable error messages: tell users what went wrong and how to fix it
- Never expose raw Go errors to users; wrap and clarify

### Testing requirements

**CRITICAL: Always use mocks. Never execute real commands in tests.**

**See [docs/guides/testing-patterns.md](docs/guides/testing-patterns.md) for complete patterns, examples, and reference.**

Accidental real command execution is a common mistake. It can:
- Break tests on different systems (missing tools, platform differences)
- Modify user state (install packages, create files)
- Cause CI failures in shared environments
- Hide bugs (tests pass only if side effects succeed)

**Testing checklist:**
- [ ] All public functionality has tests
- [ ] Use `testutil.MockApp` for command mocking (see `internal/testutil/`)
- [ ] Verify no real commands executed: `testutil.VerifyNoRealCommands(t, mockApp.Base)`
- [ ] Use `t.Helper()` in test helper functions
- [ ] Test both success and failure paths
- [ ] No real Homebrew/apt calls in tests
- [ ] Initialize logger in test file: `func init() { testutil.InitLogger() }`

### App interface pattern

All apps in `internal/apps/` implement the `App` interface defined in `internal/apps/contract.go`. The full contract, sentinel errors, `AppKind` enum, `baseapp.Reinstall`, and constructor patterns are documented in **[docs/guides/app-interface.md](docs/guides/app-interface.md)**.

Quick reference:

- Every app adds `var _ apps.App = (*X)(nil)` for compile-time enforcement
- `ForceInstall` must use `baseapp.Reinstall(a.Install, a.Uninstall)` — never call `Uninstall` directly without handling `ErrUninstallNotSupported`
- Unsupported ops return sentinel errors (`apps.ErrUninstallNotSupported`, `apps.ErrUpdateNotSupported`, …) — **never free-form strings**
- `Fonts` satisfies `FontInstaller`, not `App` — see the guide for details

```go
// Minimal new app pattern
var _ apps.App = (*MyApp)(nil)

func (a *MyApp) Name() string       { return constants.MyApp }
func (a *MyApp) Kind() apps.AppKind { return apps.KindTerminal }

func (a *MyApp) ForceInstall() error { return baseapp.Reinstall(a.Install, a.Uninstall) }
func (a *MyApp) Uninstall() error    { return fmt.Errorf("%w for myapp", apps.ErrUninstallNotSupported) }
func (a *MyApp) Update() error       { return fmt.Errorf("%w for myapp", apps.ErrUpdateNotSupported) }
```

---

## 7. Critical surfaces

Verify these before changing code in these areas:

### Installation State Management

- [ ] Global config file is updated atomically (write to temp, then rename)
- [ ] Duplicate installations are prevented (check global_config.yaml before install)
- [ ] Rollback on failure works (test by simulating install failure partway through)
- [ ] Installation state persists across shell restarts

### Cross-platform Installation (macOS ↔ Debian/Ubuntu)

- [ ] Package names translated correctly (check `pkg/constants/package_mappings.go`)
- [ ] Debian strategies handle all cases (apt, PPA, Launchpad, script, git clone)
- [ ] Platform detection works reliably (use `BaseCommandExecutor.IsMac()`)
- [ ] Both platforms produce equivalent results (same tool versions, configs)

### Shell Integration

- [ ] Shell config files (`devgita.zsh`) are sourced correctly
- [ ] Mise activation works after install (test: `eval "$(mise activate zsh)"`)
- [ ] User shell customizations are not overwritten
- [ ] Aliases and functions don't conflict with user's existing setup

---

## 8. Platform scope

**Supported platforms:**
- macOS 13+ (Ventura or newer) with Homebrew
- Debian 12+ (Bookworm) and Ubuntu 24+ with APT
- Architectures: amd64, arm64 (Apple Silicon)

**Supported categories:**
- Terminal tools (40+): shells, editors, utilities, runtime managers
- Languages: Node.js, Python, Go, Rust, PHP (via Mise or native)
- Databases: PostgreSQL, Redis, MySQL, MongoDB, SQLite
- Desktop apps: Platform-specific GUIs (Docker, browsers, window managers, etc.)

**Single-command installation:**
- `dg install` — interactive full setup
- `dg install --only <category>` — install specific category
- `dg install --skip <category>` — install all except category

See [ROADMAP.md](ROADMAP.md) for planned features and future platform support

---

## 9. Versioning & tagging

Devgita follows [Semantic Versioning](https://semver.org/) strictly: **`vMAJOR.MINOR.PATCH`**

### Which bump to use

| Change type | Bump | Example |
| --- | --- | --- |
| Bug fix, typo, test fix, docs correction | **PATCH** `x.x.^` | `v0.10.2` -> `v0.10.3` |
| New feature, new app installer, new command, new flag | **MINOR** `x.^.x` | `v0.10.3` -> `v0.11.0` |
| Breaking change to CLI interface, config format change, removed platform support | **MAJOR** `^.x.x` | `v0.11.0` -> `v1.0.0` |

**Rules:**
- Tags always start with `v` (e.g., `v0.10.3`, not `0.10.3`)
- PATCH resets to 0 on MINOR bump; MINOR and PATCH reset to 0 on MAJOR bump
- Refactoring with no behavior change = PATCH (conservative)
- Multiple bug fixes in one release = single PATCH bump
- A release mixing features and fixes = MINOR bump (the higher bump wins)
- When in doubt, ask before tagging

### Tagging workflow

```bash
git tag v0.10.3
git push origin v0.10.3
```

GitHub Actions builds and publishes automatically. See [docs/guides/releasing.md](docs/guides/releasing.md) for the full release process.

---

## 10. Change discipline

Things that must never happen silently (always require explicit PR discussion and test):

- Altering command signatures (`dg install --something`) without deprecation plan
- Adding new package categories without updating installer logic and tests
- Changing config file format (`global_config.yaml`) without migration strategy
- Removing support for a platform (macOS or Linux) without major version bump
- Modifying what "terminal tools" category includes — users depend on this being stable
- Changing installation paths or config directories — affects existing installations

---

## 11. Spec-driven development

When to write documentation **before** code:

- **Write a cycle doc** when tackling a feature that spans multiple commands or touches multiple layers (see `docs/plans/TEMPLATE.md`)
- **Write an ADR** when choosing between technologies, patterns, or approaches with lasting impact (see `docs/decisions/TEMPLATE.md`)
- **Skip both** for bug fixes, incremental improvements to existing features, or obvious changes
- **Quick ADR** (one page) for local design decisions; full ADR for platform-level choices

**Required workflow before implementing ANY code changes:**

1. If the change is substantial, create a cycle document in `docs/plans/cycles/YYYY-MM-DD-<name>.md`
2. Get user/team approval before implementing
3. Track progress by checking off steps as you go

---

## 11. Architecture Patterns

### Cross-platform installation

See `docs/guides/cross-platform-installation.md` for full details.

**Package Mappings:** `pkg/constants/package_mappings.go`

- Translates Homebrew package names → APT package names (e.g., `gdbm` → `libgdbm-dev`)

**Installation Strategies:** `internal/commands/debian_strategies.go`

- `AptStrategy` — Standard apt install with automatic name translation
- `PPAStrategy` — Personal Package Archives with GPG key configuration
- `LaunchpadPPAStrategy` — Launchpad PPA via `add-apt-repository`
- `InstallScriptStrategy` — Executable install scripts (`curl | sh`)
- `NerdFontStrategy` — GitHub release downloads for fonts
- `GitCloneStrategy` — Git repository cloning and setup

**Strategy pattern flow:**

1. Each app has a `GetInstallStrategy()` method
2. Strategy is platform-aware (returns different strategy on macOS vs. Debian)
3. Strategies handle error cases and retries (exponential backoff for downloads)
4. No knowledge of specific tools in strategy base—fully generic

### Testing pattern

```go
func init() { testutil.InitLogger() }

func TestFeature(t *testing.T) {
    mockApp := testutil.NewMockApp()
    app := &MyApp{Cmd: mockApp.Cmd, Base: mockApp.Base}

    // Test logic here

    // Always verify no real commands executed
    testutil.VerifyNoRealCommands(t, mockApp.Base)
}
```

---

## 12. Codebase landmarks

Where to find and add code:

| Purpose                   | Location                       | Notes                                                   |
| ------------------------- | ------------------------------ | ------------------------------------------------------- |
| **CLI commands**          | `cmd/`                         | Entry points; register in cmd/root.go                   |
| **App modules**           | `internal/apps/{appname}/`     | 2 files per app: `{appname}.go` + `{appname}_test.go`   |
| **Category coordinators** | `internal/tooling/`            | terminal, languages, databases, worktree                |
| **Platform installers**   | `internal/commands/`           | Strategy implementations for Debian, Darwin             |
| **Configuration logic**   | `internal/config/`             | Global state management                                 |
| **TUI components**        | `internal/tui/`                | Terminal UI helpers (minimal)                           |
| **Shared utilities**      | `pkg/`                         | Logger, paths, file ops, constants, package mappings    |
| **Embedded configs**      | `configs/`                     | Templates and static files (embedded at compile time)   |
| **Tests**                 | `*_test.go` alongside impl     | Use testutil mocks; never execute real commands         |
| **User docs**             | `docs/`                        | Feature docs, architecture, app guides, tooling details |
| **Developer docs**        | `CLAUDE.md`, `CONTRIBUTING.md` | This file and contributor guide                         |

### Adding a new command

**See [docs/guides/cli-patterns.md](docs/guides/cli-patterns.md) for detailed patterns, examples, and best practices.**

1. Read cli-patterns.md to understand command structure and patterns
2. Create handler in `cmd/{command}.go` (or create subdirectory for complex commands with subcommands)
3. Implement command logic using Cobra following patterns from cli-patterns.md
4. Add tests alongside implementation (`*_test.go`) using [docs/guides/testing-patterns.md](docs/guides/testing-patterns.md)
5. Register in `cmd/root.go`
6. Document in README.md and `docs/spec.md` if user-facing
7. If substantial, create a cycle doc first (see section 10)

### Adding a new app installer

1. Create directory `internal/apps/{appname}/`
2. Implement `{appname}.go` with app interface
3. Implement `{appname}_test.go` with tests
4. Add config templates to `configs/{appname}/` if applicable
5. Register in appropriate category in `internal/tooling/{category}/`
6. Document in `docs/apps/{appname}.md`

---

## Quick Reference: Common Commands

| Task        | Command                               | Location                                |
| ----------- | ------------------------------------- | --------------------------------------- |
| Build       | `make build`                          | Current platform                        |
| Build all   | `make all`                            | darwin-arm64, darwin-amd64, linux-amd64 |
| Test all    | `go test ./...`                       | All tests with coverage                 |
| Test single | `go test -run TestName ./pkg/package` | Specific test                           |
| Lint        | `make lint`                           | Format + vet                            |
| Format      | `go fmt ./...`                        | Auto-format code                        |
| Clean       | `make clean`                          | Remove binaries                         |

---

## Documentation Index

Quick reference to where things live:

| Topic                | Location                                           | Description                                           |
| -------------------- | -------------------------------------------------- | ----------------------------------------------------- |
| **Development Guides** | `docs/guides/README.md`                            | Index of all guides with quick-start by task          |
| **Feature Spec**     | `docs/spec.md`                                     | What features exist, architecture, edge cases, testing strategy |
| **Testing Patterns** | `docs/guides/testing-patterns.md`                  | Mocking, dependency injection, test isolation         |
| **Error Handling**   | `docs/guides/error-handling.md`                    | Error patterns, user-facing messages                  |
| **CLI Patterns**     | `docs/guides/cli-patterns.md`                      | Command structure, Cobra patterns, flags, subcommands |
| **Cross-Platform**   | `docs/guides/cross-platform-installation.md`       | Strategy pattern, package mappings, Debian strategies |
| **Releasing**        | `docs/guides/releasing.md`                         | GitHub releases workflow, versioning                  |
| **Roadmap**          | `ROADMAP.md`                                       | Planned commands, future features, open questions     |
| **Decisions**        | `docs/decisions/README.md`                         | Architectural decisions with rationale                |
| **Contributing**     | `CONTRIBUTING.md`                                  | Dev setup, build, test, git workflow, release process |

---

## Recent Changes & Active Work

**Last updated:** 2026-04-28

**Recent specs completed:**

- `specs/001-binary-dist-audit/` — Go embed, text/template for config generation
- `specs/002-debian-package-fixes/` — Strategy pattern, package mappings, exponential backoff downloads

**Active cycles:**

- Check `docs/plans/cycles/` for current work (e.g., worktree UX improvements)

**Known patterns:**

- Strategy pattern for cross-platform installation
- Interface-based testing with mock apps
- YAML state management with global config
- Embedded configs via Go `embed` package

---

## Deprecation Note

**AGENTS.md has been consolidated into this file.** OpenCode recognizes CLAUDE.md as the development guide (with AGENTS.md as fallback). All project-specific development practices are now documented here.

---
