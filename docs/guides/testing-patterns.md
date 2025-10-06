# Testing Patterns in Devgita

This guide documents the established testing patterns and conventions used throughout the devgita codebase to ensure consistent, reliable, and maintainable tests.

## Mock Interfaces and Command Testing

### MockCommand Usage

The project provides a reusable `MockCommand` in `internal/commands/mock.go` for testing all app modules:

```go
// Use the centralized mock instead of creating custom mocks
mc := commands.NewMockCommand()
app := &YourApp{Cmd: mc}

// Test installation tracking
err := app.Install()
assert.Equal(t, "expected-package", mc.InstalledPkg)

// Test error scenarios
mc.SetError("install", errors.New("installation failed"))
err := app.Install()
assert.Error(t, err)
```

### Mock Features

- **State tracking**: Records all method calls and parameters
- **Error simulation**: Configure different error scenarios with `SetError()`
- **Reset capability**: Clean state between tests with `Reset()`
- **Sensible defaults**: Pre-configured for common test scenarios

## Filesystem Testing Patterns

### Using t.TempDir()

Always use `t.TempDir()` for filesystem tests to ensure clean, isolated test environments:

```go
func TestConfigureCopiesFiles(t *testing.T) {
    src := t.TempDir()
    dst := t.TempDir()
    
    // Override global paths for test duration
    oldAppDir, oldLocalDir := paths.AppConfigDir, paths.LocalConfigDir
    paths.AppConfigDir, paths.LocalConfigDir = src, dst
    t.Cleanup(func() {
        paths.AppConfigDir, paths.LocalConfigDir = oldAppDir, oldLocalDir
    })
    
    // Test logic here...
}
```

### Path Override Pattern

For apps that use global path variables:

1. **Save original paths** before test
2. **Override with temporary directories**
3. **Restore in t.Cleanup()** to prevent test pollution
4. **Test both creation and modification** scenarios

## Helper Function Conventions

### t.Helper() Usage

Mark test helper functions with `t.Helper()` to improve test failure reporting:

```go
func setupTestEnvironment(t *testing.T) (string, string) {
    t.Helper()
    
    src := t.TempDir()
    dst := t.TempDir()
    
    // Setup logic
    return src, dst
}

func assertFileContent(t *testing.T, filePath, expected string) {
    t.Helper()
    
    content, err := os.ReadFile(filePath)
    if err != nil {
        t.Fatalf("failed to read file %s: %v", filePath, err)
    }
    
    if string(content) != expected {
        t.Fatalf("content mismatch: expected %q, got %q", expected, string(content))
    }
}
```

## Platform-Specific Test Strategies

### Mock Platform Detection

Use the mock command to simulate different platforms:

```go
func TestPlatformSpecificInstall(t *testing.T) {
    tests := []struct {
        name           string
        platform       string
        expectedPkg    string
    }{
        {"macOS", "darwin", "brew-package"},
        {"Linux", "linux", "apt-package"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mc := commands.NewMockCommand()
            // Configure mock for specific platform
            app := &YourApp{Cmd: mc}
            
            err := app.Install()
            assert.NoError(t, err)
            assert.Equal(t, tt.expectedPkg, mc.InstalledPkg)
        })
    }
}
```

### Cross-Platform Path Testing

Test path resolution across different platforms:

```go
func TestPathResolution(t *testing.T) {
    // Test uses actual path resolution but with temporary directories
    tempHome := t.TempDir()
    oldHome := os.Getenv("HOME")
    os.Setenv("HOME", tempHome)
    t.Cleanup(func() { os.Setenv("HOME", oldHome) })
    
    path := paths.GetConfigDir("test-app")
    expected := filepath.Join(tempHome, ".config", "test-app")
    assert.Equal(t, expected, path)
}
```

## Integration Test Patterns

### Configuration Testing

Test both forced and soft configuration patterns:

```go
func TestForceConfigure(t *testing.T) {
    // Test that ForceConfigure overwrites existing files
    setupConfigTest(t)
    
    // Create initial config
    err := app.ForceConfigure()
    assertFileContent(t, configPath, originalContent)
    
    // Modify config
    writeFile(t, configPath, modifiedContent)
    
    // Force configure should overwrite
    err = app.ForceConfigure()
    assertFileContent(t, configPath, originalContent)
}

func TestSoftConfigure(t *testing.T) {
    // Test that SoftConfigure preserves existing files
    setupConfigTest(t)
    
    // Initial configure
    err := app.SoftConfigure()
    assertFileContent(t, configPath, originalContent)
    
    // Modify config
    writeFile(t, configPath, modifiedContent)
    
    // Soft configure should NOT overwrite
    err = app.SoftConfigure()
    assertFileContent(t, configPath, modifiedContent)
}
```

### Command Execution Testing

For apps with command execution, test both success and failure scenarios:

```go
func TestExecuteCommand(t *testing.T) {
    mc := commands.NewMockCommand()
    app := &YourApp{Cmd: mc}
    
    // Test command execution (may fail in test environment)
    err := app.ExecuteCommand("--version")
    
    // Don't assert success/failure, just ensure no panic
    if err == nil {
        t.Log("Command succeeded unexpectedly (binary available)")
    } else {
        t.Logf("Command failed as expected: %v", err)
    }
}
```

## Test Organization Best Practices

### Test Structure

- **One test per method**: Each public method should have dedicated tests
- **Separate positive/negative cases**: Use subtests for different scenarios
- **Table-driven tests**: For multiple similar test cases

### Test Naming

- **TestMethodName**: For basic functionality tests
- **TestMethodName_Scenario**: For specific scenarios
- **TestMethodName_Error**: For error condition tests

### Error Testing

Always test error conditions and edge cases:

```go
func TestInstall_PackageManagerMissing(t *testing.T) {
    mc := commands.NewMockCommand()
    mc.PackageManagerInstalled = false
    
    app := &YourApp{Cmd: mc}
    err := app.Install()
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "package manager")
}
```

## Common Anti-Patterns to Avoid

### ❌ Don't

- Create custom mocks when `MockCommand` exists
- Use absolute paths in tests
- Forget to clean up temporary files/directories
- Test implementation details instead of behavior
- Use `os.RemoveAll()` with hardcoded paths

### ✅ Do

- Use `commands.NewMockCommand()` for command interface testing
- Use `t.TempDir()` for all filesystem operations
- Test both success and error scenarios
- Focus on testing public interface behavior
- Use `t.Cleanup()` for resource management

## Test Utilities

### Available Test Helpers

```go
// Mock command with state tracking
mc := commands.NewMockCommand()
mc.Reset()                    // Clear state between tests
mc.SetError("install", err)   // Configure error scenarios

// Temporary directory creation
tempDir := t.TempDir()        // Automatically cleaned up

// Path override pattern
oldPath := paths.GlobalPath
paths.GlobalPath = tempPath
t.Cleanup(func() { paths.GlobalPath = oldPath })
```

This testing approach ensures reliable, maintainable tests that accurately reflect the behavior of the devgita application modules while remaining fast and deterministic.
