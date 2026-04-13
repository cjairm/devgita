# Data Model: Debian/Ubuntu Package Installation

**Feature**: 002-debian-package-fixes  
**Date**: 2026-04-09  
**Context**: This document defines the data structures, relationships, and state transitions for the Debian/Ubuntu package installation system.

---

## Core Entities

### 1. PackageMapping

**Purpose**: Maps package names across platforms (macOS Homebrew ↔ Debian/Ubuntu apt)

**Fields**:
- `MacOS` (string): Homebrew package name
- `Debian` (string): Debian/Ubuntu package name

**Validation Rules**:
- Both fields MUST NOT be empty strings
- MacOS and Debian names MAY be identical (e.g., "curl" → "curl")
- Debian names SHOULD follow Debian naming conventions (lib*-dev for libraries)

**Relationships**:
- Used by: `DebianCommand.InstallPackage()`
- Stored in: `pkg/constants/package_mappings.go`
- Lookup: By constant key (e.g., `constants.Gdbm`)

**Example**:
```go
PackageMapping{
    MacOS:  "gdbm",
    Debian: "libgdbm-dev",
}
```

---

### 2. RetryConfig

**Purpose**: Configures retry behavior for network operations (downloads, API calls)

**Fields**:
- `MaxRetries` (int): Maximum number of retry attempts (default: 3)
- `InitialWait` (time.Duration): Initial backoff delay (default: 1s)
- `MaxWait` (time.Duration): Maximum backoff delay cap (default: 10s)
- `Multiplier` (float64): Exponential backoff multiplier (default: 2.0)
- `Jitter` (float64): Randomization factor 0.0-1.0 (default: 0.2 for ±20%)

**Validation Rules**:
- `MaxRetries` MUST be >= 0
- `InitialWait` MUST be > 0
- `Multiplier` MUST be >= 1.0
- `Jitter` MUST be between 0.0 and 1.0

**Calculated Properties**:
- `CalculateBackoff(attempt int) time.Duration`: Returns wait duration for given attempt

**Relationships**:
- Used by: `downloader.DownloadFileWithRetry()`
- Instantiated by: `downloader.DefaultRetryConfig()`

**State Behavior**:
- Immutable after creation
- Each retry uses exponential backoff: `initialWait * (multiplier ^ attempt) ± jitter`

**Example**:
```go
RetryConfig{
    MaxRetries:  3,
    InitialWait: 1 * time.Second,
    MaxWait:     10 * time.Second,
    Multiplier:  2.0,
    Jitter:      0.2,
}
// Attempt 1: wait ~1s, Attempt 2: wait ~2s, Attempt 3: wait ~4s
```

---

### 3. PPAConfig

**Purpose**: Defines configuration for a Debian/Ubuntu PPA (Personal Package Archive)

**Fields**:
- `Name` (string): PPA identifier (e.g., "mise")
- `KeyURL` (string): GPG public key URL
- `RepoURL` (string): Repository base URL
- `Distribution` (string): Distribution codename or version (e.g., "stable")
- `Component` (string): Repository component (e.g., "main")
- `Architecture` (string, optional): Target architecture (auto-detected if empty)

**Validation Rules**:
- `Name` MUST be a valid filename (alphanumeric, hyphens, underscores)
- `KeyURL` MUST be a valid HTTP/HTTPS URL
- `RepoURL` MUST be a valid HTTP/HTTPS URL
- `Distribution` and `Component` MUST NOT contain spaces

**Derived Paths**:
- Keyring path: `/etc/apt/keyrings/{Name}-archive-keyring.gpg`
- Sources file: `/etc/apt/sources.list.d/{Name}.list`

**Relationships**:
- Used by: `apt.PPAManager.AddPPA()`
- Referenced by: `PPAStrategy.ppaConfig`

**Example**:
```go
PPAConfig{
    Name:         "mise",
    KeyURL:       "https://mise.jdx.dev/gpg-key.pub",
    RepoURL:      "https://mise.jdx.dev/deb",
    Distribution: "stable",
    Component:    "main",
    Architecture: "", // Auto-detected via dpkg --print-architecture
}
```

---

### 4. InstallationStrategy (Interface)

**Purpose**: Defines contract for different package installation methods

**Methods**:
- `Install(packageName string) error`: Installs the package
- `IsInstalled(packageName string) (bool, error)`: Checks if package is already installed

**Implementations**:

#### AptStrategy
- **Use case**: Standard apt packages and library mappings
- **Behavior**: Translates package name via `PackageMapping`, installs via apt
- **Examples**: curl, git, libgdbm-dev

#### PPAStrategy
- **Use case**: Packages from custom repositories
- **Fields**: `cmd *DebianCommand`, `ppaConfig PPAConfig`
- **Behavior**: Adds PPA first, then installs package via apt
- **Examples**: mise, fastfetch

#### GitHubBinaryStrategy
- **Use case**: Pre-built binaries from GitHub releases
- **Fields**: `cmd *DebianCommand`, `downloadURL string`, `installPath string`, `binaryName string`
- **Behavior**: Downloads with retry, extracts if needed, installs to /usr/local/bin
- **Examples**: neovim, lazygit, lazydocker

#### GitCloneStrategy
- **Use case**: Git repository installations
- **Fields**: `cmd *DebianCommand`, `repoURL string`, `installPath string`
- **Behavior**: Clones repository to specified path
- **Examples**: powerlevel10k

**Relationships**:
- Selected by: `DebianCommand.getInstallationStrategy(packageName)`
- Executed by: `DebianCommand.InstallPackage(packageName)`

**State Transitions**:
```
[Not Installed] --Install()--> [Downloading] ---> [Installing] ---> [Installed]
                                     |                  |
                                     v                  v
                                [Failed]           [Failed]
```

---

### 5. InstallationResult

**Purpose**: Tracks the outcome of a single package installation attempt

**Fields**:
- `PackageName` (string): Name of the package attempted
- `Status` (enum): Success | Failed | Skipped
- `ErrorMessage` (string, optional): Error details if failed
- `Duration` (time.Duration): Time taken for installation
- `Attempt` (int): Retry attempt number (1-based)

**Status Values**:
- **Success**: Package installed successfully
- **Failed**: Installation failed after all retries
- **Skipped**: Package already installed, not attempted

**Validation Rules**:
- If `Status == Failed`, `ErrorMessage` MUST NOT be empty
- If `Status == Success` or `Skipped`, `ErrorMessage` SHOULD be empty
- `Duration` MUST be >= 0
- `Attempt` MUST be >= 1

**Relationships**:
- Aggregated by: `InstallationSummary`
- Stored in: `GlobalConfig.FailedInstallations` (if failed)

**Example**:
```go
InstallationResult{
    PackageName:  "lazygit",
    Status:       Failed,
    ErrorMessage: "download failed after 3 retries: connection timeout",
    Duration:     8500 * time.Millisecond,
    Attempt:      4, // Final attempt
}
```

---

### 6. InstallationSummary

**Purpose**: Aggregates installation results for final report

**Fields**:
- `Installed` (int): Count of successfully installed packages
- `Failed` (int): Count of failed installations
- `Skipped` (int): Count of skipped installations (already present)
- `Results` ([]InstallationResult): Detailed results for each package

**Calculated Properties**:
- `Total() int`: Returns `Installed + Failed + Skipped`
- `FormatSummary() string`: Returns "Installed: X, Failed: Y, Skipped: Z"

**Validation Rules**:
- All counts MUST be >= 0
- `len(Results)` SHOULD equal `Total()`

**Relationships**:
- Built by: Installation coordinator (terminal, languages, databases)
- Displayed by: CLI summary output
- Persisted to: GlobalConfig (failed installations only)

**Example**:
```go
InstallationSummary{
    Installed: 12,
    Failed:    2,
    Skipped:   1,
    Results: []InstallationResult{
        {PackageName: "curl", Status: Success, ...},
        {PackageName: "neovim", Status: Success, ...},
        {PackageName: "lazygit", Status: Failed, ErrorMessage: "...", ...},
        // ... more results
    },
}

summary.FormatSummary() // "Installed: 12, Failed: 2, Skipped: 1"
```

---

### 7. GlobalConfig (Extended)

**Purpose**: Existing global configuration structure, extended with failure tracking

**Existing Fields** (from `internal/config/fromFile.go`):
- `AppPath` (string)
- `ConfigPath` (string)
- `Installed.Packages` ([]string)
- `Installed.DevLanguages` ([]string)
- `Installed.Databases` ([]string)
- `AlreadyInstalled.Packages` ([]string)
- `Shell.Mise` (bool)

**New Fields**:
- `FailedInstallations` ([]FailedInstallation): Packages that failed to install

**FailedInstallation Structure**:
- `PackageName` (string): Name of failed package
- `Category` (string): "package" | "dev_language" | "database"
- `ErrorMessage` (string): Last error message
- `FailedAt` (time.Time): Timestamp of failure
- `AttemptCount` (int): Number of retry attempts made

**Validation Rules**:
- Failed installations MUST NOT be in `Installed` arrays
- Failed installations MAY be retried on subsequent runs
- `FailedInstallations` SHOULD be cleared on successful installation

**State Transitions**:
```
[Package Not Installed] 
    --Install Succeeds--> [Package in Installed.Packages]
    --Install Fails--> [Package in FailedInstallations]
    
[Package in FailedInstallations]
    --Retry Succeeds--> [Package in Installed.Packages, removed from FailedInstallations]
    --Retry Fails--> [Package remains in FailedInstallations, error updated]
```

**Example YAML**:
```yaml
installed:
  packages:
    - curl
    - git
    - neovim
  dev_languages:
    - node@lts
  databases:
    - postgresql

already_installed:
  packages:
    - vim

failed_installations:
  - package_name: lazygit
    category: package
    error_message: "download failed after 3 retries: connection timeout"
    failed_at: "2026-04-09T14:30:00Z"
    attempt_count: 4
```

---

## Relationships Diagram

```
PackageMapping
    └─> Used by: DebianCommand.InstallPackage()
            └─> Selects: InstallationStrategy
                    ├─> AptStrategy
                    ├─> PPAStrategy (uses PPAConfig)
                    ├─> GitHubBinaryStrategy (uses RetryConfig)
                    └─> GitCloneStrategy

InstallationStrategy.Install()
    └─> Returns: error
            └─> Aggregated into: InstallationResult
                    └─> Collected in: InstallationSummary
                            ├─> Displayed: CLI summary output
                            └─> Persisted: GlobalConfig.FailedInstallations
```

---

## State Transitions

### Package Installation Lifecycle

```
┌──────────────────┐
│  Not Installed   │
└────────┬─────────┘
         │
         v
    ┌────────────────────┐
    │ Check IsInstalled  │
    └─────┬──────────┬───┘
          │          │
    Yes   │          │ No
          │          │
          v          v
    ┌─────────┐  ┌─────────────────┐
    │ Skipped │  │ Start Install   │
    └─────────┘  └────────┬────────┘
                          │
                          v
                    ┌──────────────┐
                    │ Downloading  │<───┐
                    └──────┬───────┘    │
                           │            │
                     Success│ Fail      │
                           │            │
                           v            │
                    ┌──────────────┐    │
                    │ Installing   │    │
                    └──────┬───────┘    │
                           │            │
                     Success│ Fail      │ Retry
                           │            │ (up to 3x)
                           v            │
                    ┌──────────────┐    │
                    │  Installed   │    │
                    └──────────────┘    │
                           │            │
                           v            │
                    ┌──────────────┐    │
                    │ Track in GC  │    │
                    └──────────────┘    │
                                        │
                    ┌──────────────┐    │
                    │   Failed     │────┘
                    └──────┬───────┘
                           │
                           v
                    ┌──────────────────────┐
                    │ Track in GC.Failed   │
                    └──────────────────────┘
```

### PPA Installation Sequence

```
┌─────────────────────┐
│ Check PPA Exists    │
└──────┬──────────────┘
       │
  Yes  │  No
       │
       v
┌──────────────────────┐
│ Skip (Already Added) │
└──────────────────────┘

       │ No
       v
┌────────────────────────┐
│ Install Prerequisites  │
│ (gpg, wget, curl)      │
└────────┬───────────────┘
         v
┌────────────────────────┐
│ Create Keyring Dir     │
│ (/etc/apt/keyrings)    │
└────────┬───────────────┘
         v
┌────────────────────────┐
│ Download GPG Key       │
└────────┬───────────────┘
         v
┌────────────────────────┐
│ Save Key (dearmored)   │
└────────┬───────────────┘
         v
┌────────────────────────┐
│ Create Sources Entry   │
│ (.list.d/{name}.list)  │
└────────┬───────────────┘
         v
┌────────────────────────┐
│ Update Apt Cache       │
└────────┬───────────────┘
         v
┌────────────────────────┐
│ PPA Ready for Use      │
└────────────────────────┘
```

---

## Indexing & Lookup

### PackageMapping Lookup
- **Key**: Package constant (e.g., `constants.Gdbm`)
- **Structure**: `map[string]PackageMapping`
- **Complexity**: O(1) lookup
- **Fallback**: Returns original constant if not found

### InstallationStrategy Selection
- **Key**: Package name (string)
- **Structure**: Switch statement in `getInstallationStrategy()`
- **Complexity**: O(1) lookup
- **Default**: `AptStrategy` if not matched

### GlobalConfig Package Tracking
- **Key**: Package name (string)
- **Structure**: Linear search in arrays (small N <100)
- **Complexity**: O(n) but acceptable for small datasets
- **Operations**: Add, Remove, Contains check

---

## Validation Rules Summary

| Entity | Required Fields | Constraints |
|--------|-----------------|-------------|
| PackageMapping | MacOS, Debian | Both non-empty |
| RetryConfig | MaxRetries, InitialWait, Multiplier | MaxRetries≥0, InitialWait>0, Multiplier≥1.0 |
| PPAConfig | Name, KeyURL, RepoURL, Distribution, Component | Valid URLs, no spaces in dist/component |
| InstallationResult | PackageName, Status | If Failed, ErrorMessage required |
| InstallationSummary | Installed, Failed, Skipped | All ≥0, len(Results)==Total() |
| GlobalConfig | AppPath, ConfigPath | Valid absolute paths |

---

## Concurrency Considerations

- **Package installation**: Sequential (one at a time per category)
- **Download retry**: Single-threaded per file
- **GlobalConfig updates**: Not concurrent-safe (single writer during install)
- **No locking needed**: CLI runs single-threaded

---

## Performance Characteristics

| Operation | Time Complexity | Notes |
|-----------|-----------------|-------|
| PackageMapping lookup | O(1) | Hash map lookup |
| Strategy selection | O(1) | Switch statement |
| Download with retry | O(retries) | Max 4 attempts (3 retries) |
| GlobalConfig save | O(n) | YAML serialization, n=package count |
| InstallationSummary format | O(1) | Simple string formatting |

---

## Data Persistence

| Entity | Persisted | Location | Format |
|--------|-----------|----------|--------|
| PackageMapping | No | In-memory constant | Go code |
| RetryConfig | No | In-memory config | Go code |
| PPAConfig | Yes (implicit) | /etc/apt/sources.list.d/ | Debian sources format |
| InstallationResult | Partial | GlobalConfig (failures only) | YAML |
| InstallationSummary | No | CLI output only | String |
| GlobalConfig | Yes | ~/.config/devgita/global_config.yaml | YAML |

---

## Example Data Flow

```
1. User runs: `dg install --only terminal`

2. Terminal coordinator:
   - Initializes InstallationSummary
   - Iterates through terminal packages

3. For each package (e.g., "neovim"):
   - DebianCommand.InstallPackage("neovim")
       ↓
   - getInstallationStrategy("neovim") 
       → Returns GitHubBinaryStrategy
       ↓
   - GitHubBinaryStrategy.Install("neovim")
       → DownloadFileWithRetry(url, config)
           → Attempt 1: Success → extract + install
       ↓
   - Returns nil (success)
       ↓
   - Create InstallationResult{Status: Success}
       ↓
   - Add to InstallationSummary.Results
       ↓
   - Increment InstallationSummary.Installed
       ↓
   - GlobalConfig.AddToInstalled("neovim", "package")

4. Display summary: "Installed: 7, Failed: 0, Skipped: 1"
```

---

This data model provides the complete structure for implementing Debian/Ubuntu package installation with retry logic, failure tracking, and cross-platform package name mapping while maintaining separation from macOS code.
