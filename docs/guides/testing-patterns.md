# Testing Patterns in Devgita

## 🧭 Why

Ensure consistent, maintainable, and reliable tests across the devgita codebase by using dependency injection and mocking patterns. This guide documents the testing architecture that enables:

- **Isolated unit tests**: Test app logic without running actual system commands
- **Fast test execution**: No external dependencies or slow operations
- **Predictable results**: Control test outcomes with configurable mocks
- **Easy debugging**: Clear separation between test setup and execution
- **Consistent patterns**: Standardized approach across all app modules

---

## 📋 Testing Architecture Overview

### Component Structure

```
devgita/
├── internal/
│   ├── apps/                    # App implementations
│   │   ├── curl/
│   │   │   ├── curl.go         # Implementation using BaseCommandExecutor interface
│   │   │   └── curl_test.go    # Tests using MockBaseCommand
│   │   └── ...
│   ├── commands/
│   │   ├── base.go             # BaseCommand + BaseCommandExecutor interface
│   │   └── mock.go             # MockCommand + MockBaseCommand
│   └── testutil/
│       └── testutil.go         # Test isolation helpers
```

### Key Components

1. **BaseCommandExecutor interface**: Defines contract for command execution (in `commands/`)
2. **MockCommand/MockBaseCommand**: Provides test doubles (in `commands/mock.go`)
3. **testutil package**: Provides test isolation and orchestration (in `internal/testutil/`)
4. **App modules**: Use interfaces for dependency injection
5. **Test files**: Configure mocks and verify behavior

---

## 🔒 Test Isolation Principles

**Problem**: Without proper isolation, tests can:
- Execute real system commands (brew, apt, git, etc.)
- Modify your actual `.zshrc` file
- Write to real configuration directories (`~/.config/devgita/`)
- Create temporary files that get sourced in your shell

**Solution**: Use `testutil` package for complete isolation:
- Override global paths to temp directories
- Provide mock instances that track calls without executing real commands
- Verify tests don't touch real system

---

## 🎯 Three Levels of Test Isolation

### Level 1: Simple Mock (No Configuration)

For apps that don't touch configuration files:

```go
func TestInstall(t *testing.T) {
    mockApp := testutil.NewMockApp()
    app := &MyApp{Cmd: mockApp.Cmd, Base: mockApp.Base}
    
    err := app.Install()
    // ... assertions
    
    testutil.VerifyNoRealCommands(t, mockApp.Base)
}
```

### Level 2: Isolated Paths (Configuration without Shell)

For apps that write configuration files:

```go
func TestConfigure(t *testing.T) {
    cleanup := testutil.SetupIsolatedPaths(t)
    defer cleanup()
    
    appDir, configDir, templatesDir, _ := testutil.SetupTestDirs(t)
    // ... test logic
}
```

### Level 3: Complete Test Environment (Shell Configuration)

For apps that modify shell configuration:

```go
func TestShellFeature(t *testing.T) {
    tc := testutil.SetupCompleteTest(t)
    defer tc.Cleanup()
    
    app := &MyApp{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}
    // Everything is ready: paths, templates, config, mocks
}
```

---

## ⚙️ Core Testing Patterns

### Pattern 1: Dependency Injection via Interface

```go
// In internal/commands/base.go
type BaseCommandExecutor interface {
    ExecCommand(cmd CommandParams) (string, string, error)
    // ... other methods
}

// In internal/apps/curl/curl.go
type Curl struct {
    Cmd  cmd.Command
    Base cmd.BaseCommandExecutor  // Interface for testability
}

// In tests
func TestExecuteCommand(t *testing.T) {
    mockBase := commands.NewMockBaseCommand()
    app := &Curl{Base: mockBase}  // Inject mock
}
```

### Pattern 2: Configurable Mock Behavior

```go
// Configure success
mockBase.SetExecCommandResult("output", "", nil)

// Configure error
mockBase.SetExecCommandResult("", "error", fmt.Errorf("failed"))

// Configure presence checks
mockBase.IsPackagePresentResult = true
```

### Pattern 3: Call Verification

```go
// Get number of calls
count := mockBase.GetExecCommandCallCount()

// Get last call parameters
lastCall := mockBase.GetLastExecCommandCall()

// Reset between test cases
mockBase.ResetExecCommand()
```

---

## 📦 testutil Package Reference

### Path Management

- `SetupIsolatedPaths(t)` - Override global paths, returns cleanup function
- `SetupTestDirs(t)` - Create isolated directory structure
- `SetupTestEnvironment(t)` - Combines both, recommended for config tests
- `SetupCompleteTest(t)` - Full setup with templates, config, and mocks

### Template & Config Creation

- `CreateGlobalConfigTemplate(t, templatesDir)` - Create config template
- `CreateShellConfigTemplate(t, templatesDir, content)` - Create shell template
- `CreateGlobalConfigFile(t, configDir, content)` - Create actual config file

### Mock Management

- `NewMockApp()` - Creates `MockApp` with both `Cmd` and `Base` mocks
- `MockApp.Reset()` - Clear all mock state

### Verification Helpers

- `VerifyNoRealCommands(t, mockBase)` - Ensure no real commands executed
- `VerifyNoRealConfigChanges(t)` - Check paths aren't pointing to real config
- `AssertFileContains(t, path, expected)` - Verify file content
- `AssertFileNotContains(t, path, unexpected)` - Verify file doesn't contain content
- `InitLogger()` - Initialize logger for tests (call in `init()`)

---

## 🧪 Test Examples

### Example 1: Simple Installation Test

```go
func init() {
    testutil.InitLogger()
}

func TestInstall(t *testing.T) {
    mockApp := testutil.NewMockApp()
    app := &Curl{Cmd: mockApp.Cmd}
    
    if err := app.Install(); err != nil {
        t.Fatalf("Install error: %v", err)
    }
    
    if mockApp.Cmd.InstalledPkg != "curl" {
        t.Errorf("Expected package curl, got %s", mockApp.Cmd.InstalledPkg)
    }
    
    testutil.VerifyNoRealCommands(t, mockApp.Base)
}
```

### Example 2: Command Execution Test

```go
func TestExecuteCommand(t *testing.T) {
    mockApp := testutil.NewMockApp()
    app := &Curl{Cmd: mockApp.Cmd, Base: mockApp.Base}
    
    t.Run("successful execution", func(t *testing.T) {
        mockApp.Base.SetExecCommandResult("curl 7.64.1", "", nil)
        
        err := app.ExecuteCommand("--version")
        if err != nil {
            t.Fatalf("ExecuteCommand failed: %v", err)
        }
        
        if mockApp.Base.GetExecCommandCallCount() != 1 {
            t.Errorf("Expected 1 call, got %d", mockApp.Base.GetExecCommandCallCount())
        }
        
        lastCall := mockApp.Base.GetLastExecCommandCall()
        if lastCall.Command != "curl" {
            t.Errorf("Expected curl command, got %s", lastCall.Command)
        }
    })
    
    t.Run("error handling", func(t *testing.T) {
        mockApp.Reset()
        mockApp.Base.SetExecCommandResult("", "error", fmt.Errorf("failed"))
        
        err := app.ExecuteCommand("--invalid")
        if err == nil {
            t.Fatal("Expected error")
        }
    })
}
```

### Example 3: Shell Feature Configuration Test

```go
func TestForceConfigure(t *testing.T) {
    tc := testutil.SetupCompleteTest(t)
    defer tc.Cleanup()
    
    // Create custom shell template
    template := `{{if .Mise}}eval "$(mise activate zsh)"{{end}}`
    testutil.CreateShellConfigTemplate(t, tc.TemplatesDir, template)
    
    app := &Mise{Cmd: tc.MockApp.Cmd}
    
    if err := app.ForceConfigure(); err != nil {
        t.Fatalf("ForceConfigure failed: %v", err)
    }
    
    // Verify shell config was generated
    content, _ := os.ReadFile(tc.ZshConfigPath)
    if !strings.Contains(string(content), "mise activate") {
        t.Error("Shell config missing mise activation")
    }
    
    // Verify global config was updated
    configContent, _ := os.ReadFile(tc.ConfigPath)
    if !strings.Contains(string(configContent), "mise: true") {
        t.Error("Expected mise to be enabled in config")
    }
    
    testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
    testutil.VerifyNoRealConfigChanges(t)
}
```

### Example 4: Configuration File Copy Test

```go
func TestForceConfigure(t *testing.T) {
    cleanup := testutil.SetupIsolatedPaths(t)
    defer cleanup()
    
    appDir, configDir, _, _ := testutil.SetupTestDirs(t)
    
    // Create source config
    gitConfigAppDir := filepath.Join(appDir, "git")
    os.MkdirAll(gitConfigAppDir, 0755)
    
    sourceConfig := filepath.Join(gitConfigAppDir, ".gitconfig")
    sourceContent := "[user]\n\tname = Test User"
    os.WriteFile(sourceConfig, []byte(sourceContent), 0644)
    
    // Override paths
    paths.GitConfigAppDir = gitConfigAppDir
    paths.GitConfigLocalDir = filepath.Join(configDir, "git")
    
    // Test
    app := &Git{Cmd: commands.NewMockCommand()}
    if err := app.ForceConfigure(); err != nil {
        t.Fatalf("ForceConfigure failed: %v", err)
    }
    
    // Verify
    dstConfig := filepath.Join(paths.GitConfigLocalDir, ".gitconfig")
    content, _ := os.ReadFile(dstConfig)
    if string(content) != sourceContent {
        t.Errorf("Content mismatch")
    }
    
    testutil.VerifyNoRealConfigChanges(t)
}
```

---

## 🎯 Testing Best Practices

### 1. Always Initialize Logger

```go
func init() {
    testutil.InitLogger()
}
```

### 2. Use Subtests for Related Scenarios

```go
func TestFeature(t *testing.T) {
    t.Run("scenario 1", func(t *testing.T) { /* ... */ })
    t.Run("scenario 2", func(t *testing.T) { /* ... */ })
}
```

Run specific subtest: `go test -run TestFeature/scenario_1`

### 3. Reset Mock State Between Subtests

```go
t.Run("first", func(t *testing.T) {
    mockApp.Base.SetExecCommandResult("output1", "", nil)
    // ... test
})

t.Run("second", func(t *testing.T) {
    mockApp.Reset()  // Clear previous state
    mockApp.Base.SetExecCommandResult("output2", "", nil)
    // ... test
})
```

### 4. Use Temporary Directories

```go
func TestFileOperation(t *testing.T) {
    tempDir := t.TempDir()  // Automatically cleaned up
    // ... test logic
}
```

### 5. Override Global Paths with Cleanup

```go
func TestWithPathOverride(t *testing.T) {
    oldPath := paths.SomePath
    paths.SomePath = "/test/path"
    
    t.Cleanup(func() {
        paths.SomePath = oldPath
    })
}
```

### 6. Verify Error Messages

```go
if err == nil {
    t.Fatal("Expected error")
}
if !strings.Contains(err.Error(), "expected context") {
    t.Fatalf("Expected error context, got: %v", err)
}
```

---

## ✅ Testing Checklist

- [ ] Initialize logger: `testutil.InitLogger()` in `init()`
- [ ] Use appropriate isolation level (Simple Mock / Isolated Paths / Complete)
- [ ] Test `Install()`, `SoftInstall()`, `ExecuteCommand()`
- [ ] Test error handling scenarios
- [ ] Use subtests for organizing related scenarios
- [ ] Reset mock state between subtests: `mockApp.Reset()`
- [ ] Verify no real commands: `testutil.VerifyNoRealCommands(t, mockApp.Base)`
- [ ] Verify isolation: `testutil.VerifyNoRealConfigChanges(t)`
- [ ] Use `t.TempDir()` for temporary files
- [ ] Add cleanup with `defer` or `t.Cleanup()`

---

## 🚀 Running Tests

```bash
# Run all tests
go test ./...

# Run specific package
go test ./internal/apps/curl/

# Run specific test
go test -run TestExecuteCommand ./internal/apps/curl/

# Run specific subtest
go test -run TestExecuteCommand/successful_execution ./internal/apps/curl/

# Verbose output
go test -v ./internal/apps/curl/

# With coverage
go test -cover ./internal/apps/curl/
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 🔍 Common Issues & Solutions

| Issue | Solution |
|-------|----------|
| `nil pointer dereference` in logger | Add `testutil.InitLogger()` in `init()` |
| Real commands executed during tests | Use `MockBaseCommand` and verify with `VerifyNoRealCommands()` |
| Test fails on second run | Reset mock state: `mockApp.Reset()` |
| Permission errors in file tests | Use `t.TempDir()` |
| Path-related tests affect others | Use `t.Cleanup()` to restore paths |
| Tests modify real `.zshrc` | Use `SetupIsolatedPaths()` or `SetupCompleteTest()` |

---

## 📝 Creating Tests for New Apps

1. **Define app structure with interface**
```go
type NewApp struct {
    Cmd  commands.Command
    Base commands.BaseCommandExecutor  // Interface for testability
}
```

2. **Initialize logger in test file**
```go
func init() {
    testutil.InitLogger()
}
```

3. **Choose isolation level**
   - No config files → Use `testutil.NewMockApp()`
   - Config files → Use `testutil.SetupIsolatedPaths()`
   - Shell config → Use `testutil.SetupCompleteTest()`

4. **Write tests with subtests**
```go
func TestExecuteCommand(t *testing.T) {
    mockApp := testutil.NewMockApp()
    app := &NewApp{Cmd: mockApp.Cmd, Base: mockApp.Base}
    
    t.Run("success", func(t *testing.T) { /* ... */ })
    t.Run("error", func(t *testing.T) { mockApp.Reset(); /* ... */ })
}
```

5. **Verify isolation**
```go
testutil.VerifyNoRealCommands(t, mockApp.Base)
testutil.VerifyNoRealConfigChanges(t)
```

---

## 📚 Summary

| Component | Location | Purpose |
|-----------|----------|---------|
| `MockCommand` | `internal/commands/mock.go` | Mock package management operations |
| `MockBaseCommand` | `internal/commands/mock.go` | Mock command execution |
| `MockGit` | `internal/commands/mock.go` | Mock Git operations |
| `testutil` helpers | `internal/testutil/testutil.go` | Test orchestration & isolation |
| `MockApp` wrapper | `internal/testutil/testutil.go` | Combines Cmd + Base mocks |

**Key Principles**:
1. **Mocks stay in `commands/`** - They implement domain interfaces
2. **Orchestration in `testutil/`** - Setup, cleanup, verification helpers
3. **Complete isolation** - Tests never touch real system files or execute real commands
4. **Three isolation levels** - Choose based on what your app touches

**Architecture Decision**: Keep mock implementations in `internal/commands/mock.go` because they are domain-specific test doubles. The `testutil` package provides higher-level orchestration and convenience wrappers.
