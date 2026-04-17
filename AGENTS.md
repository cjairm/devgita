# Agent Guidelines for devgita

## Required Workflow

**CRITICAL:** Before implementing ANY code changes (bug fixes, features, refactoring), you MUST:

1. **Use the cycle-doc-planner skill** to create a cycle document
2. **Save the cycle doc** to `docs/plans/cycles/YYYY-MM-DD-<cycle-name>.md`
3. **Get user approval** before implementing
4. **Track progress** by checking off steps as completed

See: `.opencode/skills/cycle-doc-planner/SKILL.md`

## Quick Reference

| Task | Command |
|------|---------|
| Build | `go build -o devgita main.go` |
| Test all | `go test ./...` |
| Test single | `go test -run TestName ./pkg/package` |
| Lint | `go vet ./...` |
| Format | `go fmt ./...` |

## Documentation Index

| Topic | Location | Description |
|-------|----------|-------------|
| **Project Overview** | `docs/project-overview.md` | Architecture, installation flow, commands |
| **Testing Patterns** | `docs/guides/testing-patterns.md` | Mocking, dependency injection, test isolation |
| **Error Handling** | `docs/guides/error-handling.md` | Error patterns, `MaybeExitWithError()` |
| **CLI Patterns** | `docs/guides/cli-patterns.md` | Cobra usage, flag handling |
| **Cross-Platform** | `docs/architecture/cross-platform-installation.md` | Strategy pattern, package mappings, Debian strategies |
| **Languages** | `docs/tooling/languages.md` | Language coordinator (Mise integration) |
| **Databases** | `docs/tooling/databases.md` | Database coordinator |
| **Releasing** | `docs/guides/releasing.md` | GitHub releases workflow |

## Code Style (Essential)

- Standard Go formatting (`go fmt`)
- Import groups: stdlib, third-party, internal
- Error handling: always check and handle explicitly
- Logger: `logger.Init(verbose)` before use
- Tests: `t.Helper()` in helpers, `_test.go` suffix

## Key Architecture Patterns

### Cross-Platform Installation
See `docs/architecture/cross-platform-installation.md` for full details.

**Package Mappings:** `pkg/constants/package_mappings.go`
- Translates Homebrew names → apt names (e.g., `gdbm` → `libgdbm-dev`)

**Installation Strategies:** `internal/commands/debian_strategies.go`
- `AptStrategy` - Standard apt install with name translation
- `PPAStrategy` - PPA with GPG key configuration
- `LaunchpadPPAStrategy` - Launchpad PPA via add-apt-repository
- `InstallScriptStrategy` - curl | sh installations
- `NerdFontStrategy` - GitHub release font downloads
- `GitCloneStrategy` - Git repository cloning

### App Interface Pattern
All apps in `internal/apps/` implement:
- `Install()` / `SoftInstall()` / `ForceInstall()`
- `ForceConfigure()` / `SoftConfigure()`
- `ExecuteCommand(args...)`
- `Uninstall()` / `Update()`

See individual app docs in `docs/apps/` for specifics.

### Testing Pattern
```go
func init() { testutil.InitLogger() }

func TestFeature(t *testing.T) {
    mockApp := testutil.NewMockApp()
    app := &MyApp{Cmd: mockApp.Cmd, Base: mockApp.Base}
    // ... test logic
    testutil.VerifyNoRealCommands(t, mockApp.Base)
}
```

## Active Technologies
- Go 1.21+ with Cobra CLI, gopkg.in/yaml.v3, Go `embed`, Go `text/template`
- YAML state: `~/.config/devgita/global_config.yaml`
- Strategy pattern for cross-platform installation
- `BaseCommandExecutor` interface with `IsMac()` for platform detection

## Recent Changes
- 001-binary-dist-audit: Go embed, text/template for config generation
- 002-debian-package-fixes: Strategy pattern, package mappings, exponential backoff downloads
