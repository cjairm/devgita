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

