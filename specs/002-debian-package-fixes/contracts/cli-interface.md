# CLI Interface Contract

**Feature**: 002-debian-package-fixes  
**Date**: 2026-04-09  
**Context**: This document defines the command-line interface contract for Debian/Ubuntu package installation.

---

## Command Structure

### Existing Commands (Unchanged)

```bash
dg install                    # Install all categories
dg install --only terminal   # Install only terminal tools
dg install --only languages  # Install only languages (interactive)
dg install --only databases  # Install only databases (interactive)
dg install --only desktop    # Install only desktop apps
dg install --skip desktop    # Install all except desktop
```

**Behavior**:
- All existing commands continue to work exactly as before
- macOS workflow remains unchanged
- Debian/Ubuntu now uses platform-specific installation strategies

---

## Output Contract

### Installation Progress Output

```text
Installing neovim...
✓ neovim installed successfully

Installing lazygit...
Retrying download (attempt 1/3, wait 1s)...
Retrying download (attempt 2/3, wait 2s)...
✓ lazygit installed successfully

Installing mise...
✗ mise installation failed: unable to add PPA

Installing curl...
⊘ curl already installed, skipping
```

**Format**:
- `Installing {package}...` - Initial message
- `Retrying download (attempt X/Y, wait Zs)...` - Retry messages (if applicable)
- `✓ {package} installed successfully` - Success
- `✗ {package} installation failed: {error}` - Failure
- `⊘ {package} already installed, skipping` - Skipped

**Status Symbols**:
- `✓` - Success
- `✗` - Failure
- `⊘` - Skipped

---

### Installation Summary Output

```text
============================================================
Installation Summary
============================================================
Installed: 12, Failed: 2, Skipped: 1
============================================================
```

**Format Contract**:
- **Required**: Simple count format exactly as shown
- **Fields**: `Installed`, `Failed`, `Skipped` (comma-separated)
- **Order**: Always in this order (Installed, Failed, Skipped)
- **Separator**: Commas with single space after
- **Uppercase**: Field names capitalized

**Examples**:
```text
Installed: 10, Failed: 0, Skipped: 2
Installed: 5, Failed: 3, Skipped: 0
Installed: 0, Failed: 1, Skipped: 0
```

**Invalid formats** (DO NOT use):
```text
✗ Installed 10 packages, 2 failed     # Wrong format
✗ 10 installed | 2 failed | 0 skipped # Wrong separator
✗ Success: 10, Errors: 2              # Wrong field names
```

---

### Error Output

```text
Error: Unable to install lazygit: download failed after 3 retries: connection timeout
```

**Format**:
- Prefix: `Error: `
- Package context: `Unable to install {package}: `
- Error details: Specific error message
- No stack traces in user output (logged separately)

**Component-Level Failures**:
- Errors are logged and reported
- Installation continues with remaining components
- Failed packages tracked in GlobalConfig

---

## Exit Codes

| Code | Meaning | When Used |
|------|---------|-----------|
| 0 | Success | All packages installed successfully OR some packages installed (with failures tracked) |
| 1 | Error | Fatal error (e.g., invalid command, platform not supported) |

**Important**:
- Partial failures (some packages fail) exit with code 0
- Only fatal errors (unable to proceed) exit with code 1
- This allows scripts to continue after devgita installation

---

## Environment Variables

None. All configuration is internal.

---

## Configuration Files

### Input (Read)

**~/.config/devgita/global_config.yaml**
- Read to check pre-existing installations
- Used to determine skipped packages

### Output (Write)

**~/.config/devgita/global_config.yaml**
- Updated with successfully installed packages
- Updated with failed installations
- Maintains backward compatibility

---

## Backward Compatibility

### Guaranteed

✅ Existing commands work without modification
✅ macOS installation workflow unchanged
✅ GlobalConfig format preserved (new optional field only)
✅ Success/failure exit codes unchanged

### Changes

- **Output**: Summary format changed from "Success: X packages installed" to "Installed: X, Failed: Y, Skipped: Z"
- **Behavior**: Failures are non-fatal (continue with remaining packages)
- **Tracking**: New `failed_installations` field in GlobalConfig (optional)

---

## Platform-Specific Behavior

### macOS
- Uses Homebrew for all packages
- No retry logic needed (Homebrew handles this)
- Existing implementation unchanged

### Debian/Ubuntu
- Uses multiple installation strategies (apt, PPA, GitHub, git clone)
- Implements retry logic for downloads (3 attempts)
- Package name translation for libraries (gdbm → libgdbm-dev)

**Common Interface**:
- Both platforms use same `InstallPackage(packageName)` method
- Platform detection automatic via factory pattern
- User sees no difference in command usage

---

## Examples

### Successful Installation

```bash
$ dg install --only terminal
Installing curl...
✓ curl installed successfully
Installing neovim...
✓ neovim installed successfully
Installing lazygit...
✓ lazygit installed successfully
# ... more packages ...
============================================================
Installation Summary
============================================================
Installed: 15, Failed: 0, Skipped: 0
============================================================
$ echo $?
0
```

### Partial Failure

```bash
$ dg install --only terminal
Installing curl...
✓ curl installed successfully
Installing neovim...
Retrying download (attempt 1/3, wait 1s)...
Retrying download (attempt 2/3, wait 2s)...
Retrying download (attempt 3/3, wait 4s)...
✗ neovim installation failed: download failed after 3 retries: connection timeout
Installing lazygit...
✓ lazygit installed successfully
# ... more packages ...
============================================================
Installation Summary
============================================================
Installed: 14, Failed: 1, Skipped: 0
============================================================
$ echo $?
0
```

### Already Installed Packages

```bash
$ dg install --only terminal
Installing curl...
⊘ curl already installed, skipping
Installing neovim...
⊘ neovim already installed, skipping
Installing lazygit...
✓ lazygit installed successfully
# ... more packages ...
============================================================
Installation Summary
============================================================
Installed: 5, Failed: 0, Skipped: 10
============================================================
$ echo $?
0
```

---

## Testing Contract

### Manual Testing Checklist

- [ ] Run `dg install --only terminal` on fresh Debian 12
- [ ] Run `dg install --only terminal` on fresh Ubuntu 24
- [ ] Verify summary format matches exactly
- [ ] Test with network disconnected (should fail gracefully)
- [ ] Test with some packages pre-installed (should skip)
- [ ] Verify exit code 0 even with partial failures
- [ ] Verify macOS workflow unchanged

### Automated Testing

```go
func TestInstallSummaryFormat(t *testing.T) {
    summary := InstallationSummary{
        Installed: 12,
        Failed:    2,
        Skipped:   1,
    }
    
    expected := "Installed: 12, Failed: 2, Skipped: 1"
    actual := summary.FormatSummary()
    
    if actual != expected {
        t.Errorf("Expected %q, got %q", expected, actual)
    }
}
```

---

## Version Compatibility

| Version | CLI Contract | Notes |
|---------|--------------|-------|
| < 0.2.0 | Old format | "Success: X packages installed" |
| >= 0.2.0 | New format | "Installed: X, Failed: Y, Skipped: Z" |

**Migration**: No user action required. Output format change only.

---

This CLI interface contract ensures consistent, predictable behavior across platforms while maintaining backward compatibility with existing devgita installations.
