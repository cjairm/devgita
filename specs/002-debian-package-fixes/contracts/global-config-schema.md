# GlobalConfig Schema Contract

**Feature**: 002-debian-package-fixes  
**Date**: 2026-04-09  
**Context**: This document defines the YAML schema contract for the GlobalConfig file.

---

## File Location

```
~/.config/devgita/global_config.yaml
```

---

## Schema (YAML)

### Current Schema (Existing)

```yaml
app_path: /Users/username/.config/devgita
config_path: /Users/username/.config/devgita

installed:
  packages:
    - curl
    - git
    - neovim
  dev_languages:
    - node@lts
    - python@latest
  databases:
    - postgresql
    - redis

already_installed:
  packages:
    - vim
    - wget

shell:
  mise: true
```

### Extended Schema (New)

```yaml
app_path: /Users/username/.config/devgita
config_path: /Users/username/.config/devgita

installed:
  packages:
    - curl
    - git
    - neovim
  dev_languages:
    - node@lts
    - python@latest
  databases:
    - postgresql
    - redis

already_installed:
  packages:
    - vim
    - wget

failed_installations:  # NEW - Optional field
  - package_name: lazygit
    category: package
    error_message: "download failed after 3 retries: connection timeout"
    failed_at: "2026-04-09T14:30:00Z"
    attempt_count: 4
  - package_name: fastfetch
    category: package
    error_message: "unable to add PPA: GPG key download failed"
    failed_at: "2026-04-09T14:31:15Z"
    attempt_count: 1

shell:
  mise: true
```

---

## Field Definitions

### Root Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `app_path` | string | Yes | Absolute path to devgita repository |
| `config_path` | string | Yes | Absolute path to devgita configuration directory |
| `installed` | object | Yes | Packages installed by devgita |
| `already_installed` | object | Yes | Pre-existing packages (not uninstalled by devgita) |
| `failed_installations` | array | No | Packages that failed to install (NEW) |
| `shell` | object | No | Shell feature flags |

### installed Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `packages` | array[string] | Yes | Standard packages (can be empty) |
| `dev_languages` | array[string] | Yes | Language runtimes (can be empty) |
| `databases` | array[string] | Yes | Database systems (can be empty) |

**Format**:
- `packages`: Package names (e.g., "curl", "git", "neovim")
- `dev_languages`: Language specifications (e.g., "node@lts", "python@latest")
- `databases`: Database names (e.g., "postgresql", "redis")

### already_installed Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `packages` | array[string] | Yes | Pre-existing packages (can be empty) |

**Purpose**: Tracks packages that existed before devgita installation to prevent accidental uninstallation.

### failed_installations Array (NEW)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `package_name` | string | Yes | Name of package that failed |
| `category` | enum | Yes | "package" \| "dev_language" \| "database" |
| `error_message` | string | Yes | Last error message |
| `failed_at` | string (ISO 8601) | Yes | Timestamp of failure |
| `attempt_count` | int | Yes | Number of retry attempts made |

**Validation Rules**:
- `package_name` MUST NOT exist in any `installed` arrays
- `failed_at` MUST be valid RFC3339/ISO8601 timestamp
- `attempt_count` MUST be >= 1
- `category` MUST be one of: "package", "dev_language", "database"

### shell Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `mise` | bool | No | Whether mise shell integration is enabled |

---

## State Transitions

### Package Installation Success

```yaml
# Before
installed:
  packages: []

# After
installed:
  packages:
    - neovim
```

### Package Installation Failure

```yaml
# Before
installed:
  packages: []
failed_installations: []

# After
installed:
  packages: []  # Still empty
failed_installations:
  - package_name: neovim
    category: package
    error_message: "download failed after 3 retries"
    failed_at: "2026-04-09T14:30:00Z"
    attempt_count: 4
```

### Retry After Failure (Success)

```yaml
# Before
installed:
  packages: []
failed_installations:
  - package_name: neovim
    category: package
    error_message: "download failed after 3 retries"
    failed_at: "2026-04-09T14:30:00Z"
    attempt_count: 4

# After (successful retry)
installed:
  packages:
    - neovim
failed_installations: []  # Cleared
```

### Retry After Failure (Still Fails)

```yaml
# Before
failed_installations:
  - package_name: neovim
    category: package
    error_message: "download failed after 3 retries"
    failed_at: "2026-04-09T14:30:00Z"
    attempt_count: 4

# After (retry failed again)
failed_installations:
  - package_name: neovim
    category: package
    error_message: "download failed after 3 retries: connection timeout"  # Updated
    failed_at: "2026-04-09T15:45:12Z"  # Updated
    attempt_count: 4  # Same (max attempts per run)
```

---

## Backward Compatibility

### Reading (Load)

```go
// OLD VERSION (without failed_installations field)
type GlobalConfig struct {
    AppPath           string
    ConfigPath        string
    Installed         InstalledPackages
    AlreadyInstalled  AlreadyInstalledPackages
    Shell             ShellConfig
    // FailedInstallations field doesn't exist
}

// NEW VERSION (with failed_installations field)
type GlobalConfig struct {
    AppPath           string
    ConfigPath        string
    Installed         InstalledPackages
    AlreadyInstalled  AlreadyInstalledPackages
    FailedInstallations []FailedInstallation  // NEW - Optional
    Shell             ShellConfig
}
```

**Behavior**:
- Old config files (without `failed_installations`) load successfully
- `FailedInstallations` field is empty array if not present
- No migration required

### Writing (Save)

**When `FailedInstallations` is empty**:
```yaml
# Field is omitted (cleaner YAML)
installed:
  packages: []
# failed_installations not present
```

**When `FailedInstallations` has entries**:
```yaml
installed:
  packages: []
failed_installations:
  - package_name: neovim
    # ... fields
```

---

## Example Configurations

### Fresh Installation (No Failures)

```yaml
app_path: /Users/username/.config/devgita
config_path: /Users/username/.config/devgita

installed:
  packages:
    - curl
    - git
    - neovim
    - lazygit
    - lazydocker
  dev_languages:
    - node@lts
    - python@latest
    - go@latest
  databases:
    - postgresql
    - redis

already_installed:
  packages: []

shell:
  mise: true
```

### Partial Installation (With Failures)

```yaml
app_path: /Users/username/.config/devgita
config_path: /Users/username/.config/devgita

installed:
  packages:
    - curl
    - git
  dev_languages:
    - node@lts
  databases: []

already_installed:
  packages:
    - vim

failed_installations:
  - package_name: neovim
    category: package
    error_message: "download failed after 3 retries: connection timeout"
    failed_at: "2026-04-09T14:30:00Z"
    attempt_count: 4
  - package_name: lazygit
    category: package
    error_message: "failed to extract archive: invalid tar format"
    failed_at: "2026-04-09T14:31:45Z"
    attempt_count: 1
  - package_name: python@latest
    category: dev_language
    error_message: "mise installation failed: unable to add PPA"
    failed_at: "2026-04-09T14:32:10Z"
    attempt_count: 1

shell:
  mise: false  # Mise failed, shell integration not enabled
```

### System with Pre-existing Packages

```yaml
app_path: /Users/username/.config/devgita
config_path: /Users/username/.config/devgita

installed:
  packages:
    - neovim
    - lazygit
  dev_languages:
    - node@lts
  databases:
    - postgresql

already_installed:
  packages:
    - git      # Was already installed
    - curl     # Was already installed
    - vim      # Was already installed
    - docker   # Was already installed

shell:
  mise: true
```

---

## Validation Rules

### Schema Validation

```yaml
# VALID
installed:
  packages:
    - curl
    - git

# INVALID - packages must be array
installed:
  packages: curl  # ERROR: must be array
```

### Package Name Uniqueness

```yaml
# VALID
installed:
  packages:
    - curl
    - git

# INVALID - duplicates not allowed
installed:
  packages:
    - curl
    - curl  # ERROR: duplicate
```

### Failed Installation Validation

```yaml
# VALID
failed_installations:
  - package_name: neovim
    category: package
    error_message: "download failed"
    failed_at: "2026-04-09T14:30:00Z"
    attempt_count: 4

# INVALID - missing required field
failed_installations:
  - package_name: neovim
    category: package
    # ERROR: missing error_message, failed_at, attempt_count

# INVALID - invalid category
failed_installations:
  - package_name: neovim
    category: tool  # ERROR: must be "package", "dev_language", or "database"
    error_message: "download failed"
    failed_at: "2026-04-09T14:30:00Z"
    attempt_count: 4

# INVALID - package in both installed and failed
installed:
  packages:
    - neovim
failed_installations:
  - package_name: neovim  # ERROR: cannot be in both
    category: package
    error_message: "download failed"
    failed_at: "2026-04-09T14:30:00Z"
    attempt_count: 4
```

---

## Go Struct Mapping

```go
// internal/config/fromFile.go (extended)

type GlobalConfig struct {
    AppPath             string                      `yaml:"app_path"`
    ConfigPath          string                      `yaml:"config_path"`
    Installed           InstalledPackages           `yaml:"installed"`
    AlreadyInstalled    AlreadyInstalledPackages    `yaml:"already_installed"`
    FailedInstallations []FailedInstallation        `yaml:"failed_installations,omitempty"`  // NEW
    Shell               ShellConfig                 `yaml:"shell,omitempty"`
}

type FailedInstallation struct {
    PackageName  string    `yaml:"package_name"`
    Category     string    `yaml:"category"`  // "package" | "dev_language" | "database"
    ErrorMessage string    `yaml:"error_message"`
    FailedAt     time.Time `yaml:"failed_at"`
    AttemptCount int       `yaml:"attempt_count"`
}

type InstalledPackages struct {
    Packages     []string `yaml:"packages"`
    DevLanguages []string `yaml:"dev_languages"`
    Databases    []string `yaml:"databases"`
}

type AlreadyInstalledPackages struct {
    Packages []string `yaml:"packages"`
}

type ShellConfig struct {
    Mise bool `yaml:"mise,omitempty"`
}
```

---

## File Permissions

```bash
-rw-r--r--  1 user  user  1234 Apr  9 14:30 global_config.yaml
```

**Permissions**: 0644 (readable by all, writable by owner)

---

This schema contract ensures consistent state tracking across installations while maintaining backward compatibility with existing devgita configurations.
