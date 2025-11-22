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
│   │   ├── git/
│   │   │   ├── git.go          # Implementation using BaseCommandExecutor interface
│   │   │   └── git_test.go     # Tests using MockBaseCommand
│   │   └── ...
│   └── commands/
│       ├── base.go             # BaseCommand + BaseCommandExecutor interface
│       ├── mock.go             # MockCommand + MockBaseCommand
│       └── ...
```

### Key Components

1. **BaseCommandExecutor interface**: Defines contract for command execution
2. **MockBaseCommand**: Provides test doubles for command execution
3. **App modules**: Use interface for dependency injection
4. **Test files**: Configure mocks and verify behavior

---

## ⚙️ Core Testing Patterns

### Pattern 1: Dependency Injection via Interface

**Why**: Enables swapping real implementations with test doubles.

**Implementation**:

```go
// In internal/commands/base.go
type BaseCommandExecutor interface {
    ExecCommand(cmd CommandParams) (string, string, error)
    Setup(line string) error
    MaybeSetup(line, toSearch string) error
    IsDesktopAppPresent(dirPath, appName string) (bool, error)
    IsPackagePresent(cmd any, packageName string) (bool, error)
    IsFontPresent(fontName string) (bool, error)
    MaybeInstall(itemName string, alias []string, checkInstalled func(string) (bool, error), 
                 installFunc func(string) error, installURLFunc func(string) error, 
                 itemType string) error
    InstallFontFromURL(url, fontFileName string, runCache bool) error
}

// BaseCommand automatically satisfies this interface
type BaseCommand struct {
    // ... fields
}

func (base *BaseCommand) ExecCommand(cmd CommandParams) (string, string, error) {
    // Real implementation
}
```

**App Usage**:

```go
// In internal/apps/curl/curl.go
type Curl struct {
    Cmd  cmd.Command              // For package management operations
    Base cmd.BaseCommandExecutor  // For command execution (interface for testability)
}

func New() *Curl {
    return &Curl{
        Cmd:  cmd.NewCommand(),
        Base: cmd.NewBaseCommand(), // Returns concrete type, satisfies interface
    }
}
```

**Test Usage**:

```go
// In internal/apps/curl/curl_test.go
func TestExecuteCommand(t *testing.T) {
    mc := commands.NewMockCommand()
    mockBase := commands.NewMockBaseCommand() // Test double
    app := &Curl{Cmd: mc, Base: mockBase}     // Inject mock
    
    // Test proceeds with full control over mock behavior
}
```

---

### Pattern 2: Configurable Mock Behavior

**Why**: Tests need to simulate different scenarios (success, failure, edge cases).

**MockBaseCommand Structure**:

```go
// In internal/commands/mock.go
type MockBaseCommand struct {
    // Call tracking
    ExecCommandCalls []CommandParams
    
    // Configurable return values
    ExecCommandStdout string
    ExecCommandStderr string
    ExecCommandError  error
    
    // Return values for other methods
    IsDesktopAppPresentResult bool
    IsPackagePresentResult    bool
    IsFontPresentResult       bool
    SetupError                error
    MaybeSetupError           error
    MaybeInstallError         error
    InstallFontURLError       error
}
```

**Setting Mock Behavior**:

```go
// Configure successful execution
mockBase.SetExecCommandResult("curl 7.64.1", "", nil)

// Configure error scenario
mockBase.SetExecCommandResult("", "command not found", fmt.Errorf("command not found: curl"))

// Configure presence checks
mockBase.IsPackagePresentResult = true
mockBase.SetupError = fmt.Errorf("setup failed")
```

---

### Pattern 3: Call Verification

**Why**: Verify that app code calls the underlying command layer correctly.

**Mock Helper Methods**:

```go
// Get number of calls
count := mockBase.GetExecCommandCallCount()

// Get last call parameters
lastCall := mockBase.GetLastExecCommandCall()
if lastCall != nil {
    fmt.Printf("Command: %s, Args: %v, IsSudo: %v\n", 
               lastCall.Command, lastCall.Args, lastCall.IsSudo)
}

// Reset between test cases
mockBase.ResetExecCommand()
```

**Test Example**:

```go
t.Run("verify command parameters", func(t *testing.T) {
    mockBase.SetExecCommandResult("success", "", nil)
    
    err := app.ExecuteCommand("-o", "file.txt", "https://example.com")
    if err != nil {
        t.Fatalf("ExecuteCommand failed: %v", err)
    }
    
    // Verify call count
    if mockBase.GetExecCommandCallCount() != 1 {
        t.Fatalf("Expected 1 call, got %d", mockBase.GetExecCommandCallCount())
    }
    
    // Verify parameters
    lastCall := mockBase.GetLastExecCommandCall()
    if lastCall.Command != "curl" {
        t.Fatalf("Expected command 'curl', got %q", lastCall.Command)
    }
    
    expectedArgs := []string{"-o", "file.txt", "https://example.com"}
    if !reflect.DeepEqual(lastCall.Args, expectedArgs) {
        t.Fatalf("Expected args %v, got %v", expectedArgs, lastCall.Args)
    }
    
    if lastCall.IsSudo {
        t.Fatal("Expected IsSudo to be false")
    }
})
```

---

## 🧪 Complete Test Examples

### Example 1: Testing ExecuteCommand with Success

```go
func TestExecuteCommand_Success(t *testing.T) {
    // Setup
    mc := commands.NewMockCommand()
    mockBase := commands.NewMockBaseCommand()
    app := &Curl{Cmd: mc, Base: mockBase}
    
    // Configure mock behavior
    mockBase.SetExecCommandResult("curl 7.64.1", "", nil)
    
    // Execute
    err := app.ExecuteCommand("--version")
    
    // Verify
    if err != nil {
        t.Fatalf("ExecuteCommand failed: %v", err)
    }
    
    // Verify command was called correctly
    if mockBase.GetExecCommandCallCount() != 1 {
        t.Fatalf("Expected 1 ExecCommand call, got %d", mockBase.GetExecCommandCallCount())
    }
    
    lastCall := mockBase.GetLastExecCommandCall()
    if lastCall == nil {
        t.Fatal("No ExecCommand call recorded")
    }
    if lastCall.Command != "curl" {
        t.Fatalf("Expected command 'curl', got %q", lastCall.Command)
    }
    if len(lastCall.Args) != 1 || lastCall.Args[0] != "--version" {
        t.Fatalf("Expected args ['--version'], got %v", lastCall.Args)
    }
    if lastCall.IsSudo {
        t.Fatal("Expected IsSudo to be false")
    }
}
```

### Example 2: Testing Error Handling

```go
func TestExecuteCommand_Error(t *testing.T) {
    // Setup
    mc := commands.NewMockCommand()
    mockBase := commands.NewMockBaseCommand()
    app := &Curl{Cmd: mc, Base: mockBase}
    
    // Configure mock to return error
    mockBase.SetExecCommandResult("", "command not found", fmt.Errorf("command not found: curl"))
    
    // Execute
    err := app.ExecuteCommand("--invalid-flag")
    
    // Verify error is returned
    if err == nil {
        t.Fatal("Expected ExecuteCommand to return error")
    }
    
    // Verify error message contains context
    if !strings.Contains(err.Error(), "failed to run curl command") {
        t.Fatalf("Expected error to contain 'failed to run curl command', got: %v", err)
    }
    
    // Verify original error is wrapped
    if !strings.Contains(err.Error(), "command not found: curl") {
        t.Fatalf("Expected error to contain original error message, got: %v", err)
    }
}
```

### Example 3: Testing Multiple Arguments

```go
func TestExecuteCommand_MultipleArgs(t *testing.T) {
    // Setup
    mc := commands.NewMockCommand()
    mockBase := commands.NewMockBaseCommand()
    app := &Curl{Cmd: mc, Base: mockBase}
    
    // Configure mock
    mockBase.SetExecCommandResult("downloaded", "", nil)
    
    // Execute with multiple arguments
    err := app.ExecuteCommand("-o", "file.txt", "https://example.com")
    if err != nil {
        t.Fatalf("ExecuteCommand failed: %v", err)
    }
    
    // Verify all arguments are passed correctly
    lastCall := mockBase.GetLastExecCommandCall()
    expectedArgs := []string{"-o", "file.txt", "https://example.com"}
    
    if len(lastCall.Args) != len(expectedArgs) {
        t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
    }
    
    for i, arg := range expectedArgs {
        if lastCall.Args[i] != arg {
            t.Fatalf("Expected arg[%d] to be %q, got %q", i, arg, lastCall.Args[i])
        }
    }
}
```

### Example 4: Testing with Subtests

```go
func TestExecuteCommand(t *testing.T) {
    mc := commands.NewMockCommand()
    mockBase := commands.NewMockBaseCommand()
    app := &Curl{Cmd: mc, Base: mockBase}
    
    t.Run("successful execution", func(t *testing.T) {
        mockBase.SetExecCommandResult("curl 7.64.1", "", nil)
        err := app.ExecuteCommand("--version")
        if err != nil {
            t.Fatalf("ExecuteCommand failed: %v", err)
        }
        // Additional assertions...
    })
    
    t.Run("command execution error", func(t *testing.T) {
        mockBase.ResetExecCommand() // Reset mock state
        mockBase.SetExecCommandResult("", "error", fmt.Errorf("failed"))
        err := app.ExecuteCommand("--invalid")
        if err == nil {
            t.Fatal("Expected error")
        }
        // Additional assertions...
    })
    
    t.Run("multiple arguments", func(t *testing.T) {
        mockBase.ResetExecCommand() // Reset mock state
        mockBase.SetExecCommandResult("success", "", nil)
        err := app.ExecuteCommand("-X", "POST", "-d", "data")
        if err != nil {
            t.Fatalf("ExecuteCommand failed: %v", err)
        }
        // Additional assertions...
    })
}
```

---

## 📦 Testing Different App Methods

### Testing Install Methods

```go
func TestInstall(t *testing.T) {
    mc := commands.NewMockCommand()
    app := &Curl{Cmd: mc}
    
    if err := app.Install(); err != nil {
        t.Fatalf("Install error: %v", err)
    }
    
    // Verify correct package was installed
    if mc.InstalledPkg != "curl" {
        t.Fatalf("expected InstallPackage(%s), got %q", "curl", mc.InstalledPkg)
    }
}

func TestSoftInstall(t *testing.T) {
    mc := commands.NewMockCommand()
    app := &Curl{Cmd: mc}
    
    if err := app.SoftInstall(); err != nil {
        t.Fatalf("SoftInstall error: %v", err)
    }
    
    // Verify MaybeInstallPackage was called
    if mc.MaybeInstalled != "curl" {
        t.Fatalf("expected MaybeInstallPackage(%s), got %q", "curl", mc.MaybeInstalled)
    }
}
```

### Testing Configuration Methods

```go
func TestForceConfigure(t *testing.T) {
    // Create temporary directories for testing
    src := t.TempDir()
    dst := t.TempDir()
    
    // Override global paths for test duration
    oldAppDir, oldLocalDir := paths.GitConfigAppDir, paths.GitConfigLocalDir
    paths.GitConfigAppDir, paths.GitConfigLocalDir = src, dst
    t.Cleanup(func() {
        paths.GitConfigAppDir, paths.GitConfigLocalDir = oldAppDir, oldLocalDir
    })
    
    // Create source config file
    originalContent := "[user]\n\tname = Test User"
    if err := os.WriteFile(filepath.Join(src, ".gitconfig"), []byte(originalContent), 0o644); err != nil {
        t.Fatal(err)
    }
    
    // Test ForceConfigure
    mc := commands.NewMockCommand()
    app := &Git{Cmd: mc}
    
    if err := app.ForceConfigure(); err != nil {
        t.Fatalf("ForceConfigure error: %v", err)
    }
    
    // Verify file was copied
    check := filepath.Join(dst, ".gitconfig")
    if _, err := os.Stat(check); err != nil {
        t.Fatalf("expected copied file at %s: %v", check, err)
    }
    
    // Verify content matches
    copiedContent, err := os.ReadFile(check)
    if err != nil {
        t.Fatalf("failed to read copied file: %v", err)
    }
    if string(copiedContent) != originalContent {
        t.Fatalf("content mismatch: expected %q, got %q", originalContent, string(copiedContent))
    }
}

func TestSoftConfigure(t *testing.T) {
    // Similar setup to ForceConfigure...
    
    // First call should copy config
    if err := app.SoftConfigure(); err != nil {
        t.Fatalf("SoftConfigure error: %v", err)
    }
    
    // Modify the copied file
    modifiedContent := "[user]\n\tname = Modified User"
    if err := os.WriteFile(check, []byte(modifiedContent), 0o644); err != nil {
        t.Fatal(err)
    }
    
    // Second call should NOT overwrite
    if err := app.SoftConfigure(); err != nil {
        t.Fatalf("second SoftConfigure error: %v", err)
    }
    
    // Verify file was NOT overwritten
    finalContent, err := os.ReadFile(check)
    if err != nil {
        t.Fatalf("failed to read file: %v", err)
    }
    if string(finalContent) != string(modifiedContent) {
        t.Fatalf("SoftConfigure overwrote existing file")
    }
}
```

---

## 🎯 Testing Best Practices

### 1. Initialize Logger in Tests

```go
func init() {
    // Initialize logger for tests (prevents nil pointer errors)
    logger.Init(false)
}
```

### 2. Use Subtests for Related Scenarios

```go
func TestExecuteCommand(t *testing.T) {
    t.Run("scenario 1", func(t *testing.T) { /* ... */ })
    t.Run("scenario 2", func(t *testing.T) { /* ... */ })
    t.Run("scenario 3", func(t *testing.T) { /* ... */ })
}
```

Benefits:
- Run individual scenarios: `go test -run TestExecuteCommand/scenario_1`
- Better test organization and readability
- Isolated failures (one subtest failure doesn't prevent others)

### 3. Reset Mock State Between Subtests

```go
t.Run("first test", func(t *testing.T) {
    mockBase.SetExecCommandResult("output", "", nil)
    // ... test logic
})

t.Run("second test", func(t *testing.T) {
    mockBase.ResetExecCommand() // Clear previous state
    mockBase.SetExecCommandResult("different output", "", nil)
    // ... test logic
})
```

### 4. Use Temporary Directories for File Operations

```go
func TestFileOperation(t *testing.T) {
    tempDir := t.TempDir() // Automatically cleaned up after test
    
    testFile := filepath.Join(tempDir, "test.txt")
    // ... test logic
}
```

### 5. Override Global Paths with Cleanup

```go
func TestWithPathOverride(t *testing.T) {
    // Save original values
    oldPath := paths.SomePath
    
    // Override for test
    paths.SomePath = "/test/path"
    
    // Restore after test
    t.Cleanup(func() {
        paths.SomePath = oldPath
    })
    
    // ... test logic
}
```

### 6. Skip Tests for Unsupported Operations

```go
// SKIP: ForceInstall test as per guidelines
// ForceInstall calls Uninstall (which returns error) before Install
// Testing this creates false negatives
// func TestForceInstall(t *testing.T) { ... }

// Instead, test Install and Uninstall independently
func TestInstall(t *testing.T) { /* ... */ }

func TestUninstall(t *testing.T) {
    // Verify it returns expected error
    err := app.Uninstall()
    if err == nil {
        t.Fatal("expected error for unsupported operation")
    }
    if err.Error() != "curl uninstall not supported through devgita" {
        t.Fatalf("unexpected error message: %v", err)
    }
}
```

### 7. Test Error Messages

```go
func TestErrorWrapping(t *testing.T) {
    mockBase.SetExecCommandResult("", "stderr", fmt.Errorf("original error"))
    
    err := app.ExecuteCommand("--flag")
    
    // Verify error is returned
    if err == nil {
        t.Fatal("Expected error")
    }
    
    // Verify error contains context
    if !strings.Contains(err.Error(), "failed to run curl command") {
        t.Fatalf("Expected context in error, got: %v", err)
    }
    
    // Verify original error is preserved
    if !strings.Contains(err.Error(), "original error") {
        t.Fatalf("Expected original error in wrapped error, got: %v", err)
    }
}
```

---

## 🚀 Running Tests

### Run All Tests

```bash
go test ./...
```

### Run Tests for Specific Package

```bash
go test ./internal/apps/curl/
go test ./internal/apps/git/
```

### Run Specific Test

```bash
go test -run TestExecuteCommand ./internal/apps/curl/
```

### Run Specific Subtest

```bash
go test -run TestExecuteCommand/successful_execution ./internal/apps/curl/
```

### Run with Verbose Output

```bash
go test -v ./internal/apps/curl/
```

### Run with Coverage

```bash
go test -cover ./internal/apps/curl/
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 📝 Creating Tests for New Apps

### Step 1: Define App Structure with Interface

```go
// In internal/apps/newapp/newapp.go
package newapp

import "github.com/cjairm/devgita/internal/commands"

type NewApp struct {
    Cmd  commands.Command
    Base commands.BaseCommandExecutor // Use interface for testability
}

func New() *NewApp {
    return &NewApp{
        Cmd:  commands.NewCommand(),
        Base: commands.NewBaseCommand(), // Concrete type satisfies interface
    }
}
```

### Step 2: Implement ExecuteCommand

```go
func (app *NewApp) ExecuteCommand(args ...string) error {
    _, _, err := app.Base.ExecCommand(commands.CommandParams{
        Command: "newapp",
        Args:    args,
        IsSudo:  false,
    })
    if err != nil {
        return fmt.Errorf("failed to run newapp command: %w", err)
    }
    return nil
}
```

### Step 3: Create Test File

```go
// In internal/apps/newapp/newapp_test.go
package newapp

import (
    "fmt"
    "strings"
    "testing"
    
    "github.com/cjairm/devgita/internal/commands"
    "github.com/cjairm/devgita/pkg/logger"
)

func init() {
    logger.Init(false)
}

func TestExecuteCommand(t *testing.T) {
    mc := commands.NewMockCommand()
    mockBase := commands.NewMockBaseCommand()
    app := &NewApp{Cmd: mc, Base: mockBase}
    
    t.Run("successful execution", func(t *testing.T) {
        mockBase.SetExecCommandResult("output", "", nil)
        
        err := app.ExecuteCommand("--flag")
        if err != nil {
            t.Fatalf("ExecuteCommand failed: %v", err)
        }
        
        // Verify command was called
        if mockBase.GetExecCommandCallCount() != 1 {
            t.Fatalf("Expected 1 call, got %d", mockBase.GetExecCommandCallCount())
        }
        
        // Verify parameters
        lastCall := mockBase.GetLastExecCommandCall()
        if lastCall.Command != "newapp" {
            t.Fatalf("Expected command 'newapp', got %q", lastCall.Command)
        }
    })
    
    t.Run("error handling", func(t *testing.T) {
        mockBase.ResetExecCommand()
        mockBase.SetExecCommandResult("", "error", fmt.Errorf("failed"))
        
        err := app.ExecuteCommand("--invalid")
        if err == nil {
            t.Fatal("Expected error")
        }
        if !strings.Contains(err.Error(), "failed to run newapp command") {
            t.Fatalf("Expected error context, got: %v", err)
        }
    })
}
```

---

## ✅ Testing Checklist

When creating tests for a new app module:

- [ ] Initialize logger in `init()` function
- [ ] Test `Install()` method
- [ ] Test `SoftInstall()` method
- [ ] Test `ForceConfigure()` method (if applicable)
- [ ] Test `SoftConfigure()` method (if applicable)
- [ ] Test `ExecuteCommand()` with success scenario
- [ ] Test `ExecuteCommand()` with error handling
- [ ] Test `ExecuteCommand()` with multiple arguments (if applicable)
- [ ] Test `Uninstall()` error message (if unsupported)
- [ ] Use `MockCommand` for package management operations
- [ ] Use `MockBaseCommand` for command execution
- [ ] Use subtests for organizing related scenarios
- [ ] Reset mock state between subtests
- [ ] Verify call counts and parameters
- [ ] Test error message wrapping
- [ ] Use temporary directories for file operations
- [ ] Add cleanup handlers for path overrides

---

## 🔍 Debugging Test Failures

### Common Issues

**Issue**: `nil pointer dereference` when calling logger methods
```bash
Solution: Add logger.Init(false) in init() function
```

**Issue**: Real commands are executed during tests
```bash
Solution: Verify app uses BaseCommandExecutor interface and inject MockBaseCommand
```

**Issue**: Test fails on second run but passes first time
```bash
Solution: Reset mock state between subtests with mockBase.ResetExecCommand()
```

**Issue**: File operation tests fail with permission errors
```bash
Solution: Use t.TempDir() for temporary directories with proper permissions
```

**Issue**: Path-related tests affect subsequent tests
```bash
Solution: Use t.Cleanup() to restore original path values
```

### Verbose Test Output

```bash
# See detailed test output
go test -v ./internal/apps/curl/

# See which tests are being skipped
go test -v ./internal/apps/curl/ | grep SKIP
```

---

## 📚 Summary

| Pattern                  | Purpose                                 | Implementation                    |
| ------------------------ | --------------------------------------- | --------------------------------- |
| Dependency Injection     | Enable test doubles                     | BaseCommandExecutor interface     |
| Configurable Mocks       | Simulate different scenarios            | SetExecCommandResult()            |
| Call Verification        | Verify correct command invocation       | GetLastExecCommandCall()          |
| State Reset              | Isolate tests from each other           | ResetExecCommand()                |
| Subtests                 | Organize related test scenarios         | t.Run()                           |
| Temporary Directories    | Safe file operation testing             | t.TempDir()                       |
| Path Override + Cleanup  | Test with custom paths safely           | t.Cleanup()                       |
| Logger Initialization    | Prevent nil pointer errors              | logger.Init(false) in init()      |
| Error Message Testing    | Verify proper error handling            | strings.Contains()                |
| Skip Unsupported Methods | Avoid false negatives                   | Comment with rationale            |

**Key Takeaway**: The devgita testing architecture uses dependency injection via the `BaseCommandExecutor` interface to enable comprehensive unit testing without executing real system commands. Tests are fast, isolated, and provide clear feedback on app behavior.
