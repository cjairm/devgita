# Adding New Basic Apps to Devgita

This guide provides a step-by-step process for adding new applications to the devgita development environment manager, following the established patterns and conventions.

## Overview

Devgita follows a consistent module pattern where each app is a self-contained package with standardized methods for installation, configuration, and execution. This guide covers adding basic apps (simple package installations with optional configurations).

## Prerequisites

- Understanding of Go modules and interfaces
- Familiarity with the devgita codebase structure
- Knowledge of the target application's installation requirements

## Step-by-Step Process

### 1. Create the App Module Directory

Create a new directory under `internal/apps/` for your application:

```
internal/apps/your-app-name/
└── your-app-name.go
```

**Naming Convention**: Use lowercase, hyphenated names that match the actual package name when possible.

### 2. Create the App Module File

Use the following template based on the fastfetch example:

```go
// -------------------------
// NOTE: Brief description of the app and its purpose
// - **Documentation**: https://link-to-official-docs
// - **Devgita documentation**: <basic-documentation>
// -------------------------

package yourappname

import (
	"fmt"
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type YourAppName struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *YourAppName {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &YourAppName{Cmd: osCmd, Base: *baseCmd}
}

func (y *YourAppName) Install() error {
	return y.Cmd.InstallPackage("package-name")
}

func (y *YourAppName) MaybeInstall() error {
	return y.Cmd.MaybeInstallPackage("package-name")
}

func configure() error {
	return files.CopyDir(paths.YourAppConfigAppDir, paths.YourAppConfigLocalDir)
}

func (y *YourAppName) ForceConfigure() error {
	return configure()
}

func (y *YourAppName) SoftConfigure() error {
	configFile := filepath.Join(paths.YourAppConfigLocalDir, "config-filename")
	if files.FileAlreadyExist(configFile) {
		return nil
	}
	return configure()
}

func (y *YourAppName) Uninstall() error {
	return y.Cmd.UninstallPackage("package-name")
}

func (y *YourAppName) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		IsSudo:      false,
		Command:     constants.YourAppName,
		Args:        args,
	}
	if _, err := y.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run %s command: %w", constants.YourAppName, err)
	}
	return nil
}
```

### 3. Update Constants

Add your app constants to `pkg/constants/constants.go`:

```go
// App names
YourAppName = "your-app-name"
```

### 4. Update Paths

Add path constants to `pkg/paths/paths.go`:

```go
// Configs from Devgita app
YourAppConfigAppDir = GetAppDir(constants.ConfigAppDirName, constants.YourAppName)

// Config in local (usually `.config` folder)
YourAppConfigLocalDir = GetConfigDir(constants.YourAppName)

// Cache local (usually `.cache` folder)
YourAppCache = GetCacheDir(constants.YourAppName)

// see pkg/paths/paths.go for helper functions
```

### 5. Create Configuration Files (Optional)

If your app requires configuration:

1. Create a directory under `configs/your-app-name/`
2. Add default configuration files
3. Ensure the configuration structure matches what the app expects

Example structure:

```
configs/your-app-name/
├── config.json
├── themes/
│   └── default.conf
└── settings/
    └── defaults.yml
```

### 6. Choose Installation Method

Select the appropriate installation method in your `Install()` function based on app type:

- **Regular Package**: `y.Cmd.InstallPackage("package-name")`
- **Desktop Application**: `y.Cmd.InstallDesktopApp("package-name")`
- **Custom Installation**: `y.Cmd.InstallCustom("custom-install-command")`

### 7. Required Methods

Every app module must implement these core methods:

#### Installation Methods

- `Install() error` - Force install the package
- `MaybeInstall() error` - Install only if not already present

#### Uninstallation Methods

- `Uninstall() error` - Remove the package

#### Configuration Methods

- `configure() error` - Internal helper function for configuration logic
- `ForceConfigure() error` - Force setup configuration files
- `SoftConfigure() error` - Setup only if config doesn't exist

#### Execution Methods (Optional)

- `ExecuteCommand(args ...string) error` - Execute the application with arguments

### 8. Error Handling Best Practices

```go
// Use wrapped errors with context
return fmt.Errorf("failed to install %s: %w", appName, err)

// Check for existing installations
if alreadyInstalled {
	return nil // Silent success
}

// Validate prerequisites
if !hasPrerequisite {
	return fmt.Errorf("prerequisite not found: %s", prerequisite)
}
```

### 9. Testing Your App Module

Create a test file `your-app-name_test.go`:

```go
package yourappname

import (
	"testing"
	"path/filepath"
	"os"
)

func TestNew(t *testing.T) {
	app := New()
	if app == nil {
		t.Fatal("New() returned nil")
	}
}

func TestSoftConfigure(t *testing.T) {
	tempDir := t.TempDir()
	// Test configure logic with temporary directory
}
```

### 10. Documentation

Update relevant documentation:

- Add app description to main README.md
- Create app-specific docs in `docs/apps/your-app-name.md` if complex
- Update CHANGELOG.md with the addition

## Common Patterns

### Apps with Multiple Packages

```go
func (y *YourAppName) Install() error {
	packages := []string{"package1", "package2", "package3"}
	for _, pkg := range packages {
		if err := y.Cmd.MaybeInstallPackage(pkg); err != nil {
			return fmt.Errorf("failed to install %s: %w", pkg, err)
		}
	}
	return nil
}
```

### Platform-Specific Installation

```go
func (y *YourAppName) Install() error {
	if y.Cmd.IsMacOS() {
		return y.Cmd.InstallPackage("macos-package-name")
	}
	return y.Cmd.InstallPackage("linux-package-name")
}
```

### Configuration with Templates

```go
func (y *YourAppName) ForceConfigure() error {
	// Copy base configuration
	if err := files.CopyDir(paths.YourAppConfigAppDir, paths.YourAppConfigLocalDir); err != nil {
		return err
	}

	// Apply user-specific customizations
	return y.applyUserCustomizations()
}
```

## Troubleshooting

### Common Issues

1. **Package name mismatch**: Ensure the package name in code matches the actual system package
2. **Path conflicts**: Verify path constants are unique and don't conflict with existing apps
3. **Missing dependencies**: Check if the app requires other packages to be installed first
4. **Permission issues**: Some apps may require sudo installation - use `IsSudo: true` in CommandParams

### Debugging

(Only when asked by user)

```bash
# Test individual app installation
go run main.go install --only terminal --verbose

# Run specific tests
go test ./internal/apps/your-app-name -v

# Check path resolution
go run -c "import paths; print(paths.YourAppConfigLocalDir)"
```

## Example: Adding "btop" (System Monitor)

Here's a complete example of adding the `btop` system monitor:

1. **Create the module** (`internal/apps/btop/btop.go`):

```go
// -------------------------
// NOTE: Modern system monitor with better visuals than htop
// - Documentation: https://github.com/aristocratos/btop
// - Devgita documentation: Btop is a resource monitor that shows usage and stats for processor, memory, disks, network and processes. Here some useful commands to get you started:
//   - Launch btop: `btop`
// -------------------------

package btop

import (
	"fmt"
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type Btop struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *Btop {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Btop{Cmd: osCmd, Base: *baseCmd}
}

func (b *Btop) Install() error {
	return b.Cmd.InstallPackage("btop")
}

func (b *Btop) MaybeInstall() error {
	return b.Cmd.MaybeInstallPackage("btop")
}

func configure() error {
	return files.CopyDir(paths.BtopConfigAppDir, paths.BtopConfigLocalDir)
}

func (b *Btop) ForceConfigure() error {
	return configure()
}

func (b *Btop) SoftConfigure() error {
	configFile := filepath.Join(paths.BtopConfigLocalDir, "btop.conf")
	if files.FileAlreadyExist(configFile) {
		return nil
	}
	return configure()
}

func (b *Btop) Uninstall() error {
	return b.Cmd.UninstallPackage("btop")
}

func (b *Btop) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		IsSudo:      false,
		Command:     constants.Btop,
		Args:        args,
	}
	if _, err := b.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run btop command: %w", err)
	}
	return nil
}
```

2. **Update constants**:

```go
// Add to pkg/constants/constants.go
Btop = "btop"
```

3. **Update paths**:

```go
// Add to pkg/paths/paths.go
BtopConfigAppDir = GetAppDir(constants.ConfigAppDirName, constants.Btop)
BtopConfigLocalDir = GetConfigDir(constants.Btop)
```

4. **Create config** (`configs/btop/btop.conf`):

```conf
# btop configuration
color_theme="Default"
theme_background=True
truecolor=True
```

This example demonstrates the complete flow for adding a new basic app to devgita.
