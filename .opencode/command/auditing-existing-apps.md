# Auditing One App Module: `$ARGUMENTS`

This command audits and modernizes **exactly one** app module, `$ARGUMENTS`, to match current `devgita` patterns. It only touches `internal/apps/$ARGUMENTS/` and related `constants/paths` entries.

## What This Command Guarantees

1. **Scopes to one app**: only `internal/apps/$ARGUMENTS/**` is changed.
2. **Adds/ensures standardized methods** (without removing existing behavior):
   - Install: `ForceInstall()`, `SoftInstall()`, `install()` (private), `Uninstall()`
   - Configure: `ForceConfigure()`, `SoftConfigure()`, `configure()` (private)
   - Exec: `ExecuteCommand(args ...string)`

3. **Keeps existing logic**; only fills gaps.
4. **Creates simple, meaningful tests**:
   - Mocks install/uninstall calls (no system interaction).
   - **Copies real files in a temp dir** to validate configuration copy logic.

5. **Ensures constants/paths entries exist** for `$ARGUMENTS`.

---

## Pre-Flight (Safety & Focus)

```bash
test -n "$ARGUMENTS" || { echo "APP is required"; exit 1; }
test -d "internal/apps/$ARGUMENTS" || { echo "internal/apps/$ARGUMENTS not found"; exit 1; }

# Optional: guard against accidental wide edits
git diff --quiet || { echo "Working tree not clean. Commit or stash first."; exit 1; }
```

---

## Required Interfaces (to enable mocking & real file copies)

> If your code already uses similar abstractions, skip or adapt naming.

```go
// pkg/cmd/cmd.go
package cmd

type Commander interface {
	InstallPackage(pkg string) error
	MaybeInstallPackage(pkg string) error
	UninstallPackage(pkg string) error
}

type BaseCommander interface {
	ExecCommand(params CommandParams) (stdout string, stderr string, err error)
}

type CommandParams struct {
	PreExecMsg, PostExecMsg string
	IsSudo                  bool
	Command                 string
	Args                    []string
}

func NewCommand() Commander        { /* existing */ return nil }
func NewBaseCommand() BaseCommander { /* existing */ return nil }
```

```go
// pkg/files/files.go
package files

type FS interface {
	CopyDir(src, dst string) error
	FileAlreadyExist(path string) bool
}

func DefaultFS() FS { /* existing wrappers */ return nil }
```

This makes it trivial to **mock installs** while still **copying real files** in tests using `t.TempDir()`.

---

## Standard Structure (per app)

```go
// internal/apps/$ARGUMENTS/$ARGUMENTS.go
package $ARGUMENTS

import (
	"fmt"
	"path/filepath"

	"devgita/pkg/cmd"
	"devgita/pkg/constants"
	"devgita/pkg/files"
	"devgita/pkg/paths"
)

type App struct {
	Cmd   cmd.Commander
	Base  cmd.BaseCommander
	FS    files.FS
}

func New() *App {
	return &App{
		Cmd:  cmd.NewCommand(),
		Base: cmd.NewBaseCommand(),
		FS:   files.DefaultFS(),
	}
}

// --- install helpers ---
func install(c cmd.Commander) error {
	return c.InstallPackage(constants.$ARGUMENTS)
}

func (a *App) ForceInstall() error { return install(a.Cmd) }

func (a *App) SoftInstall() error {
	return a.Cmd.MaybeInstallPackage(constants.$ARGUMENTS)
}

func (a *App) Uninstall() error {
	return a.Cmd.UninstallPackage(constants.$ARGUMENTS)
}

// --- configure helpers ---
func configure(fs files.FS) error {
	return fs.CopyDir(paths.$ARGUMENTSConfigAppDir, paths.$ARGUMENTSConfigLocalDir)
}

func (a *App) ForceConfigure() error { return configure(a.FS) }

func (a *App) SoftConfigure() error {
	// Pick a representative file or dir that signals "already configured"
	// Adjust the marker name for the specific app if needed.
	marker := filepath.Join(paths.$ARGUMENTSConfigLocalDir, "main-config-file")
	if a.FS.FileAlreadyExist(marker) {
		return nil
	}
	return configure(a.FS)
}

// --- exec ---
func (a *App) ExecuteCommand(args ...string) error {
	params := cmd.CommandParams{
		Command: constants.$ARGUMENTS,
		Args:    args,
	}
	if _, _, err := a.Base.ExecCommand(params); err != nil {
		return fmt.Errorf("failed to run %s: %w", constants.$ARGUMENTS, err)
	}
	return nil
}
```

---

## Constants & Paths (required keys)

```go
// pkg/constants/constants.go
package constants

const (
	$ARGUMENTS = "$ARGUMENTS" // exact binary/package name
)
```

```go
// pkg/paths/paths.go
package paths

import "devgita/pkg/constants"

// "Configs from Devgita app"
var $ARGUMENTSConfigAppDir = GetAppDir(ConfigAppDirName, constants.$ARGUMENTS)

// "Config in local"
var $ARGUMENTSConfigLocalDir = GetConfigDir(constants.$ARGUMENTS)
```

> If these already exist, leave them as-is.

---

## Tests (mock installs + real file copies)

```go
// internal/apps/$ARGUMENTS/$ARGUMENTS_test.go
package $ARGUMENTS

import (
	"os"
	"path/filepath"
	"testing"

	"devgita/pkg/cmd"
	"devgita/pkg/constants"
	"devgita/pkg/files"
	"devgita/pkg/paths"
)

// --- simple mocks for install/uninstall ---
type mockCmd struct {
	installedPkg   string
	maybeInstalled string
	uninstalledPkg string
}

func (m *mockCmd) InstallPackage(pkg string) error     { m.installedPkg = pkg; return nil }
func (m *mockCmd) MaybeInstallPackage(pkg string) error { m.maybeInstalled = pkg; return nil }
func (m *mockCmd) UninstallPackage(pkg string) error   { m.uninstalledPkg = pkg; return nil }

// --- stub BaseCommander (we don’t assert exec behaviour here) ---
type noopBase struct{}
func (n *noopBase) ExecCommand(_ cmd.CommandParams) (string, string, error) { return "", "", nil }

// --- in-memory FS that delegates to real OS for copy/exists, but lets us point paths to temp dirs ---
type realFS struct{}
func (realFS) CopyDir(src, dst string) error                  { return copyDir(src, dst) }
func (realFS) FileAlreadyExist(p string) bool                 { _, err := os.Stat(p); return err == nil }

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(p string, info os.FileInfo, err error) error {
		if err != nil { return err }
		rel, _ := filepath.Rel(src, p)
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, rerr := os.ReadFile(p)
		if rerr != nil { return rerr }
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil { return err }
		return os.WriteFile(target, data, 0o644)
	})
}

func TestNew(t *testing.T) {
	app := New()
	if app == nil { t.Fatalf("New() returned nil") }
}

func TestForceInstallAndUninstall(t *testing.T) {
	mc := &mockCmd{}
	app := &App{Cmd: mc, Base: &noopBase{}, FS: realFS{}}

	if err := app.ForceInstall(); err != nil {
		t.Fatalf("ForceInstall error: %v", err)
	}
	if mc.installedPkg != constants.$ARGUMENTS {
		t.Fatalf("expected InstallPackage(%s), got %q", constants.$ARGUMENTS, mc.installedPkg)
	}

	if err := app.Uninstall(); err != nil {
		t.Fatalf("Uninstall error: %v", err)
	}
	if mc.uninstalledPkg != constants.$ARGUMENTS {
		t.Fatalf("expected UninstallPackage(%s), got %q", constants.$ARGUMENTS, mc.uninstalledPkg)
	}
}

func TestSoftInstall(t *testing.T) {
	mc := &mockCmd{}
	app := &App{Cmd: mc, Base: &noopBase{}, FS: realFS{}}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mc.maybeInstalled != constants.$ARGUMENTS {
		t.Fatalf("expected MaybeInstallPackage(%s), got %q", constants.$ARGUMENTS, mc.maybeInstalled)
	}
}

func TestConfigureCopiesRealFilesInTemp(t *testing.T) {
	// Create temp "app config" dir with a fake file as source
	src := t.TempDir()
	dst := t.TempDir()

	// Override global paths for the duration of the test
	oldAppDir, oldLocalDir := paths.$ARGUMENTSConfigAppDir, paths.$ARGUMENTSConfigLocalDir
	paths.$ARGUMENTSConfigAppDir, paths.$ARGUMENTSConfigLocalDir = src, dst
	t.Cleanup(func() {
		paths.$ARGUMENTSConfigAppDir, paths.$ARGUMENTSConfigLocalDir = oldAppDir, oldLocalDir
	})

	// Seed source
	if err := os.MkdirAll(filepath.Join(src, "nested"), 0o755); err != nil { t.Fatal(err) }
	if err := os.WriteFile(filepath.Join(src, "nested", "main-config-file"), []byte("ok"), 0o644); err != nil { t.Fatal(err) }

	app := &App{Cmd: &mockCmd{}, Base: &noopBase{}, FS: realFS{}}

	// ForceConfigure should copy the files
	if err := app.ForceConfigure(); err != nil {
		t.Fatalf("ForceConfigure error: %v", err)
	}
	check := filepath.Join(dst, "nested", "main-config-file")
	if _, err := os.Stat(check); err != nil {
		t.Fatalf("expected copied file at %s: %v", check, err)
	}

	// SoftConfigure should detect marker and do nothing (still succeeds)
	if err := app.SoftConfigure(); err != nil {
		t.Fatalf("SoftConfigure error: %v", err)
	}
}
```

**Why this works**

- **Install/Uninstall** are fully mocked (no system calls).
- **Configuration** does **real file copies** but safely inside `t.TempDir()` paths.
- The “marker” (`main-config-file`) exercises the **soft** branch without over-engineering.

---

## Developer Checklist (Step-by-Step)

1. **Inspect module**

   ```bash
   ls -la internal/apps/$ARGUMENTS
   sed -n '1,200p' internal/apps/$ARGUMENTS/$ARGUMENTS.go
   ```

2. **Add the standard struct + methods** above (preserve existing functions).
3. **Ensure constants/paths** entries exist (add if missing).
4. **Add tests** exactly as shown (adapt names only).
5. **Run**

   ```bash
   go test ./internal/apps/$ARGUMENTS -v
   go test ./...      # sanity
   go build ./...     # compile check
   ```

---

## Do / Don’t

**Do**

- Limit edits strictly to `internal/apps/$ARGUMENTS/**` plus the one constant/path entry.
- Keep existing behavior; only add missing pieces.
- Use tiny mocks and `t.TempDir()` for file operations.

**Don’t**

- Refactor global patterns across the repo in this command.
- Add heavy test harnesses or integration tests.
- Overwrite any user configs outside temp dirs.
