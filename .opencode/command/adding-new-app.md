# Adding New App: $ARGUMENTS

This command creates a complete new app module for **$ARGUMENTS** following devgita's standardized patterns.

**Usage**: `$ARGUMENTS` should be the app name you want to add (e.g., `btop`, `ripgrep`, `fd`)

## What This Command Does

1. **Creates a new app module** in `internal/apps/$ARGUMENTS/`
2. **Implements all standardized methods** (ForceInstall, SoftInstall, ForceConfigure, SoftConfigure)
3. **Adds constants and paths** for the new app
4. **Creates basic unit tests** with simple mocks
5. **Sets up configuration templates** if needed

## Step-by-Step Process

### 1. Create App Module Directory

```bash
mkdir -p internal/apps/$ARGUMENTS
```

### 2. Create Main App File

Create `internal/apps/$ARGUMENTS/$ARGUMENTS.go`:

```go
// -------------------------
// NOTE: Brief description of $ARGUMENTS and its purpose
// - **Documentation**: https://link-to-official-docs
// - **Devgita documentation**: $ARGUMENTS provides [functionality]. Useful commands:
//   - Basic usage: `$ARGUMENTS`
//   - Help: `$ARGUMENTS --help`
// -------------------------

package $ARGUMENTS

import (
	"fmt"
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type $ARGUMENTS struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *$ARGUMENTS {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &$ARGUMENTS{Cmd: osCmd, Base: *baseCmd}
}

func install() error {
	osCmd := cmd.NewCommand()
	return osCmd.InstallPackage("$ARGUMENTS")  // Adjust package name if different
}

func (a *$ARGUMENTS) ForceInstall() error {
	return install()
}

func (a *$ARGUMENTS) SoftInstall() error {
	return a.Cmd.MaybeInstallPackage("$ARGUMENTS")  // Adjust package name if different
}

func configure() error {
	return files.CopyDir(paths.$ARGUMENTSConfigAppDir, paths.$ARGUMENTSConfigLocalDir)
}

func (a *$ARGUMENTS) ForceConfigure() error {
	return configure()
}

func (a *$ARGUMENTS) SoftConfigure() error {
	configFile := filepath.Join(paths.$ARGUMENTSConfigLocalDir, "config-file")  // Adjust config filename
	if files.FileAlreadyExist(configFile) {
		return nil
	}
	return configure()
}

func (a *$ARGUMENTS) Uninstall() error {
	return a.Cmd.UninstallPackage("$ARGUMENTS")  // Adjust package name if different
}

func (a *$ARGUMENTS) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		IsSudo:      false,
		Command:     constants.$ARGUMENTS,
		Args:        args,
	}
	if _, err := a.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run %s command: %w", constants.$ARGUMENTS, err)
	}
	return nil
}
```

### 3. Add Constants

Add to `pkg/constants/constants.go`:

```go
// App names section
$ARGUMENTS = "$ARGUMENTS"
```

### 4. Add Paths  

Add to `pkg/paths/paths.go`:

```go
// In "Configs from Devgita app" section:
$ARGUMENTSConfigAppDir = GetAppDir(constants.ConfigAppDirName, constants.$ARGUMENTS)

// In "Config in local" section:
$ARGUMENTSConfigLocalDir = GetConfigDir(constants.$ARGUMENTS)
```

### 5. Create Configuration (Optional)

If the app needs configuration, create:

```bash
mkdir -p configs/$ARGUMENTS
```

Add default config files in `configs/$ARGUMENTS/` directory.

### 6. Create Unit Tests

Create `internal/apps/$ARGUMENTS/$ARGUMENTS_test.go`:

```go
package $ARGUMENTS

import (
	"testing"
	"path/filepath"
	"os"
)

// Simple mock for testing - no real commands executed
type mockCommand struct {
	installCalled   bool
	uninstallCalled bool
	packageName     string
}

func (m *mockCommand) InstallPackage(pkg string) error {
	m.installCalled = true
	m.packageName = pkg
	return nil
}

func (m *mockCommand) MaybeInstallPackage(pkg string) error {
	return m.InstallPackage(pkg)
}

func (m *mockCommand) UninstallPackage(pkg string) error {
	m.uninstallCalled = true
	m.packageName = pkg
	return nil
}

func TestNew(t *testing.T) {
	app := New()
	if app == nil {
		t.Fatal("New() returned nil")
	}
	if app.Cmd == nil {
		t.Error("Cmd not initialized")
	}
}

func TestForceInstall(t *testing.T) {
	mockCmd := &mockCommand{}
	app := &$ARGUMENTS{Cmd: mockCmd}

	err := app.ForceInstall()
	if err != nil {
		t.Errorf("ForceInstall failed: %v", err)
	}

	if !mockCmd.installCalled {
		t.Error("Expected InstallPackage to be called")
	}

	if mockCmd.packageName != "$ARGUMENTS" {
		t.Errorf("Expected package name '$ARGUMENTS', got %s", mockCmd.packageName)
	}
}

func TestSoftInstall(t *testing.T) {
	mockCmd := &mockCommand{}
	app := &$ARGUMENTS{Cmd: mockCmd}

	err := app.SoftInstall()
	if err != nil {
		t.Errorf("SoftInstall failed: %v", err)
	}

	if !mockCmd.installCalled {
		t.Error("Expected MaybeInstallPackage to be called")
	}
}

func TestUninstall(t *testing.T) {
	mockCmd := &mockCommand{}
	app := &$ARGUMENTS{Cmd: mockCmd}

	err := app.Uninstall()
	if err != nil {
		t.Errorf("Uninstall failed: %v", err)
	}

	if !mockCmd.uninstallCalled {
		t.Error("Expected UninstallPackage to be called")
	}
}

func TestSoftConfigure(t *testing.T) {
	tempDir := t.TempDir()

	app := New()

	// Test when config doesn't exist - should configure
	err := app.SoftConfigure()
	if err != nil {
		t.Errorf("SoftConfigure failed when no config exists: %v", err)
	}

	// Test when config exists - should not overwrite
	configFile := filepath.Join(tempDir, "test-config")
	err = os.WriteFile(configFile, []byte("existing config"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = app.SoftConfigure()
	if err != nil {
		t.Errorf("SoftConfigure failed when config exists: %v", err)
	}
}

func TestForceConfigure(t *testing.T) {
	app := New()

	err := app.ForceConfigure()
	if err != nil {
		t.Errorf("ForceConfigure failed: %v", err)
	}
}
```

### 7. Test the New Module

```bash
# Test the new module
go test ./internal/apps/$ARGUMENTS -v

# Test overall build
go test ./...

# Verify compilation
go build -o devgita main.go
```

## Installation Method Selection

Choose the appropriate installation method for your app:

### Regular Package (Most Common)
```go
func install() error {
	osCmd := cmd.NewCommand()
	return osCmd.InstallPackage("$ARGUMENTS")
}
```

### Desktop Application
```go
func install() error {
	osCmd := cmd.NewCommand()
	return osCmd.InstallDesktopApp("$ARGUMENTS")
}
```

### Platform-Specific Packages
```go
func install() error {
	osCmd := cmd.NewCommand()
	if osCmd.IsMacOS() {
		return osCmd.InstallPackage("macos-package-name")
	}
	return osCmd.InstallPackage("linux-package-name")
}
```

### Multiple Packages
```go
func install() error {
	osCmd := cmd.NewCommand()
	packages := []string{"package1", "package2", "package3"}
	for _, pkg := range packages {
		if err := osCmd.MaybeInstallPackage(pkg); err != nil {
			return fmt.Errorf("failed to install %s: %w", pkg, err)
		}
	}
	return nil
}
```

## Configuration Options

### No Configuration Needed
If the app doesn't need configuration, simplify the methods:

```go
func (a *$ARGUMENTS) ForceConfigure() error {
	return nil  // No configuration needed
}

func (a *$ARGUMENTS) SoftConfigure() error {
	return nil  // No configuration needed
}
```

### Simple Configuration
```go
func configure() error {
	return files.CopyDir(paths.$ARGUMENTSConfigAppDir, paths.$ARGUMENTSConfigLocalDir)
}
```

### Custom Configuration Logic
```go
func configure() error {
	// Custom configuration steps
	if err := files.CopyDir(paths.$ARGUMENTSConfigAppDir, paths.$ARGUMENTSConfigLocalDir); err != nil {
		return err
	}
	
	// Additional setup steps
	return setupCustomFeatures()
}
```

## Testing Guidelines

### Keep Tests Simple
- **Basic mocks only** - Just track method calls
- **No real commands** - All package managers stubbed
- **Use t.TempDir()** - For any file operations
- **Test core functionality** - Installation, configuration, uninstallation

### Mock Structure
```go
type mockCommand struct {
	installCalled   bool
	uninstallCalled bool
	packageName     string
	installError    error  // Optional: test error scenarios
}
```

### Test All Standard Methods
- TestNew() - Constructor
- TestForceInstall() - Force installation
- TestSoftInstall() - Conditional installation  
- TestUninstall() - Package removal
- TestForceConfigure() - Force configuration
- TestSoftConfigure() - Conditional configuration

## Example: Adding ripgrep

Here's a complete example for adding `ripgrep`:

1. **Create the module** (`internal/apps/ripgrep/ripgrep.go`):

```go
// -------------------------
// NOTE: Fast text search tool that recursively searches directories for patterns
// - **Documentation**: https://github.com/BurntSushi/ripgrep
// - **Devgita documentation**: ripgrep is a fast grep alternative. Useful commands:
//   - Search for pattern: `rg "pattern"`
//   - Search in specific files: `rg "pattern" --type js`
//   - Case insensitive: `rg -i "pattern"`
// -------------------------

package ripgrep

import (
	"fmt"
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type Ripgrep struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *Ripgrep {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Ripgrep{Cmd: osCmd, Base: *baseCmd}
}

func install() error {
	osCmd := cmd.NewCommand()
	return osCmd.InstallPackage("ripgrep")
}

func (r *Ripgrep) ForceInstall() error {
	return install()
}

func (r *Ripgrep) SoftInstall() error {
	return r.Cmd.MaybeInstallPackage("ripgrep")
}

func configure() error {
	return files.CopyDir(paths.RipgrepConfigAppDir, paths.RipgrepConfigLocalDir)
}

func (r *Ripgrep) ForceConfigure() error {
	return configure()
}

func (r *Ripgrep) SoftConfigure() error {
	configFile := filepath.Join(paths.RipgrepConfigLocalDir, ".ripgreprc")
	if files.FileAlreadyExist(configFile) {
		return nil
	}
	return configure()
}

func (r *Ripgrep) Uninstall() error {
	return r.Cmd.UninstallPackage("ripgrep")
}

func (r *Ripgrep) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		IsSudo:      false,
		Command:     constants.Ripgrep,
		Args:        args,
	}
	if _, err := r.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run ripgrep command: %w", err)
	}
	return nil
}
```

2. **Add constants** to `pkg/constants/constants.go`:
```go
Ripgrep = "ripgrep"
```

3. **Add paths** to `pkg/paths/paths.go`:
```go
RipgrepConfigAppDir = GetAppDir(constants.ConfigAppDirName, constants.Ripgrep)
RipgrepConfigLocalDir = GetConfigDir(constants.Ripgrep)
```

4. **Create config** in `configs/ripgrep/.ripgreprc`:
```
# Default ripgrep configuration
--smart-case
--hidden
--glob=!.git/*
```

5. **Create tests** with simple mocks following the template above.

This provides a complete, standardized app module ready for integration into devgita.