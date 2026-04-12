# Agent Guidelines for devgita

## Build/Test/Lint Commands

- **Build**: `go build -o devgita main.go`
- **Test all**: `go test ./...`
- **Test single**: `go test -run TestName ./pkg/package`
- **Lint**: `go vet ./...`
- **Format**: `go fmt ./...`

## Code Style Guidelines

- Use standard Go formatting (`go fmt`)
- Import groups: stdlib, third-party, internal packages
- Package naming: lowercase, single words (e.g., `commands`, `config`)
- Function naming: camelCase for private, PascalCase for public
- Error handling: always check and handle errors explicitly
- Use `t.Helper()` in test helper functions
- Test files end with `_test.go` and use package `package_test`
- Use descriptive variable names (`tempDir` not `tmp`)
- Comments use `//` format, capitalize first word
- Cobra commands in `cmd/` package, business logic in `internal/`
- Configuration structs use YAML tags for serialization
- Logger initialization with `logger.Init(verbose)` before use

## Testing Guidelines

- Follow testing patterns documented in `docs/guides/testing-patterns.md`
- Use dependency injection via `BaseCommandExecutor` interface for testability
- Initialize logger in test `init()` functions with `logger.Init(false)`
- Use `MockBaseCommand` for testing command execution without running actual commands
- Reset mock state between subtests with `ResetExecCommand()`
- Use `t.TempDir()` for temporary directories in file operation tests
- Organize related test scenarios with subtests using `t.Run()`
- Verify command parameters with `GetLastExecCommandCall()` and `GetExecCommandCallCount()`
- Test error wrapping and message context with `strings.Contains()`
- Skip tests for unsupported methods (e.g., ForceInstall) with rationale comments

## Future App Docs Structure

```
docs/apps/
├── neovim.md       # Your Git integration + config details (create/update when needed)
├── alacritty.md    # Terminal configuration specifics (create/update when needed)
├── tmux.md         # Session management setup (create/update when needed)
└── aerospace.md    # Window manager configuration (create/update when needed)
└── ...
```


## Active Technologies
- Go 1.21+ (existing project, uses `embed` package from Go 1.16+) + Cobra CLI, gopkg.in/yaml.v3, Go `embed`, Go `text/template` (001-binary-dist-audit)
- YAML files on disk (`~/.config/devgita/global_config.yaml`), embedded filesystem via `embed.FS` (001-binary-dist-audit)
- Go 1.21+ standard library: `net/http`, `context`, `os/exec`, `time`, `encoding/json` (002-debian-package-fixes)
- YAML files on disk (~/.config/devgita/global_config.yaml), embedded filesystem via embed.FS (002-debian-package-fixes)
- Strategy pattern for installation methods: AptStrategy, PPAStrategy, GitHubBinaryStrategy, GitCloneStrategy, InstallScriptStrategy, NerdFontStrategy (002-debian-package-fixes)
- Exponential backoff with jitter for download retries: 3 attempts, 1s/2s/4s delays (002-debian-package-fixes)
- Platform-specific package name mapping: macOS Homebrew names → Debian apt names (002-debian-package-fixes)
- `BaseCommandExecutor` interface includes `IsMac() bool` for platform detection in apps (002-debian-package-fixes)

## Recent Changes
- 001-binary-dist-audit: Added Go 1.21+ (existing project, uses `embed` package from Go 1.16+) + Cobra CLI, gopkg.in/yaml.v3, Go `embed`, Go `text/template`
- 002-debian-package-fixes: Added strategy pattern for Debian installations, exponential backoff downloads, platform-specific package name mappings, IsMac() to BaseCommandExecutor interface
