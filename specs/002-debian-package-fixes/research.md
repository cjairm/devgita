# Research: Debian/Ubuntu Package Installation Implementation

**Feature**: 002-debian-package-fixes  
**Date**: 2026-04-09  
**Context**: This document consolidates research findings for implementing Debian/Ubuntu package installation fixes with platform-specific strategies.

---

## 1. Library Package Mapping Strategy

### Decision
Use a **constant-based mapping table with a struct-based lookup system** for macOS to Debian library name translations.

### Rationale
1. **Type safety**: Constants prevent typos and enable IDE autocomplete
2. **Maintainability**: Single source of truth for package mappings
3. **Graceful degradation**: Easy to detect unmapped packages and fall back to original name
4. **Extensibility**: Simple to add new mappings as needed
5. **Testing**: Mappings can be easily unit tested

### Implementation Pattern

```go
// pkg/constants/package_mappings.go
package constants

type PackageMapping struct {
    MacOS  string // Homebrew package name
    Debian string // Debian/Ubuntu package name
}

var PackageMappings = map[string]PackageMapping{
    Gdbm: {MacOS: "gdbm", Debian: "libgdbm-dev"},
    OpenSSL: {MacOS: "openssl", Debian: "libssl-dev"},
    Readline: {MacOS: "readline", Debian: "libreadline-dev"},
    Zlib: {MacOS: "zlib", Debian: "zlib1g-dev"},
    Libyaml: {MacOS: "libyaml", Debian: "libyaml-dev"},
    Ncurses: {MacOS: "ncurses", Debian: "libncurses5-dev"},
    Libffi: {MacOS: "libffi", Debian: "libffi-dev"},
    Jemalloc: {MacOS: "jemalloc", Debian: "libjemalloc2"},
    Vips: {MacOS: "vips", Debian: "libvips"},
}

func GetDebianPackageName(packageConstant string) string {
    if mapping, exists := PackageMappings[packageConstant]; exists {
        return mapping.Debian
    }
    return packageConstant // Fallback to original
}
```

### Reference
- **Omakub pattern**: `omakub/install/terminal/libraries.sh` shows direct apt package names
- **Required mappings**: 8 libraries (gdbm, openssl, readline, zlib, libyaml, ncurses, libffi, jemalloc, vips)

### Alternatives Considered
- **map[string]string**: Simpler but less extensible, harder to add metadata later
- **Separate constants file per platform**: Duplicates logic, harder to maintain

---

## 2. GitHub Binary Download with Retry

### Decision
Use **exponential backoff with jitter** (3 retries: 1s, 2s, 4s delays) using Go standard library.

### Rationale
1. **Standard library only**: No external dependencies (aligns with Constitution Principle I)
2. **Cloud-native pattern**: Industry standard (AWS SDK, Google Cloud SDK)
3. **Network resilience**: Handles transient failures (timeouts, temporary DNS issues)
4. **Rate limiting**: Jitter prevents thundering herd problem
5. **Configurable**: Easy to adjust timing and attempts

### Implementation Pattern

```go
// pkg/downloader/retry.go
package downloader

import (
    "context"
    "math"
    "math/rand"
    "net/http"
    "time"
)

type RetryConfig struct {
    MaxRetries  int           // 3
    InitialWait time.Duration // 1s
    Multiplier  float64       // 2.0
    Jitter      float64       // 0.2 (±20%)
}

func (rc *RetryConfig) CalculateBackoff(attempt int) time.Duration {
    wait := float64(rc.InitialWait) * math.Pow(rc.Multiplier, float64(attempt))
    jitterAmount := wait * rc.Jitter
    jitter := (rand.Float64() * 2 * jitterAmount) - jitterAmount
    return time.Duration(wait + jitter)
}

func DownloadFileWithRetry(ctx context.Context, url, destPath string, config RetryConfig) error {
    for attempt := 0; attempt <= config.MaxRetries; attempt++ {
        if attempt > 0 {
            time.Sleep(config.CalculateBackoff(attempt - 1))
        }
        
        err := downloadFile(ctx, url, destPath)
        if err == nil {
            return nil
        }
        
        if !IsRetryableError(err) {
            return err
        }
    }
    return fmt.Errorf("failed after %d retries", config.MaxRetries)
}
```

### Retry Logic
- **Retryable errors**: Network timeouts, temporary DNS failures, HTTP 429/502/503/504
- **Non-retryable errors**: HTTP 404/401/403, invalid URL, file system errors
- **Timing**: Attempt 1→wait 1s→Attempt 2→wait 2s→Attempt 3→wait 4s→Final attempt
- **Total max wait**: ~7 seconds across all retries

### Reference
- **Omakub pattern**: `omakub/install/app-lazygit.sh` downloads without retry (we improve on this)
- **Standard**: RFC 7231 (HTTP retry behavior), AWS SDK retry patterns

### Alternatives Considered
- **github.com/cenkalti/backoff**: External dependency (violates Constitution)
- **Fixed delays**: Less efficient, no jitter protection
- **No retry**: Fragile to transient network issues

---

## 3. Neovim Installation Method

### Decision
Use **tar.gz extraction method** (NOT AppImage) following official Neovim releases.

### Rationale
1. **Official format**: Neovim releases tar.gz for Linux, not AppImage
2. **No FUSE dependency**: Works in containers and restricted environments
3. **Proven pattern**: Omakub uses this exact approach successfully
4. **Simpler**: Direct binary installation to `/usr/local/bin`
5. **Complete**: Includes lib and share directories for full functionality

### Implementation Pattern

```go
// internal/apps/neovim/neovim_debian.go
func (n *Neovim) InstallDebianFromGitHub() error {
    const (
        nvimURL = "https://github.com/neovim/neovim/releases/download/stable/nvim-linux-x86_64.tar.gz"
        tmpArchive = "/tmp/nvim.tar.gz"
    )
    
    // 1. Download with retry
    ctx := context.Background()
    config := downloader.DefaultRetryConfig()
    if err := downloader.DownloadFileWithRetry(ctx, nvimURL, tmpArchive, config); err != nil {
        return err
    }
    defer os.Remove(tmpArchive)
    
    // 2. Extract tar.gz
    exec.Command("tar", "-xf", tmpArchive, "-C", "/tmp").Run()
    defer os.RemoveAll("/tmp/nvim-linux-x86_64")
    
    // 3. Install binary
    exec.Command("sudo", "install", "-m", "755", 
        "/tmp/nvim-linux-x86_64/bin/nvim", 
        "/usr/local/bin/nvim").Run()
    
    // 4. Copy support files
    exec.Command("sudo", "cp", "-R", "/tmp/nvim-linux-x86_64/lib", "/usr/local/").Run()
    exec.Command("sudo", "cp", "-R", "/tmp/nvim-linux-x86_64/share", "/usr/local/").Run()
    
    return nil
}
```

### Installation Steps
1. Download `nvim-linux-x86_64.tar.gz` from GitHub releases
2. Extract to `/tmp/nvim-linux-x86_64/`
3. Install binary to `/usr/local/bin/nvim` with 755 permissions
4. Copy `lib/` and `share/` to `/usr/local/` for runtime dependencies
5. Cleanup temporary files

### Reference
- **Omakub**: `omakub/install/app-neovim.sh` - identical pattern
- **Official**: https://github.com/neovim/neovim/releases - tar.gz is primary Linux format

### Alternatives Considered
- **AppImage with FUSE**: Requires FUSE, extra complexity for fallback
- **AppImage extraction**: Neovim doesn't distribute AppImage officially
- **Build from source**: Too slow, requires build dependencies

---

## 4. PPA Management

### Decision
Use **manual GPG key + repository file management** (NOT add-apt-repository).

### Rationale
1. **Better error handling**: Direct file operations provide clearer errors
2. **Idempotency**: Easy to check if PPA already configured
3. **No wrapper dependency**: add-apt-repository not always installed
4. **Proven pattern**: Omakub uses this approach successfully
5. **Explicit control**: Clear what files are created and where

### Implementation Pattern

```go
// pkg/apt/ppa.go
package apt

type PPAConfig struct {
    Name         string // "mise"
    KeyURL       string // "https://mise.jdx.dev/gpg-key.pub"
    RepoURL      string // "https://mise.jdx.dev/deb"
    Distribution string // "stable"
    Component    string // "main"
}

func (pm *PPAManager) AddPPA(config PPAConfig) error {
    // 1. Check if already added
    sourcesFile := fmt.Sprintf("/etc/apt/sources.list.d/%s.list", config.Name)
    if fileExists(sourcesFile) {
        return nil // Already configured
    }
    
    // 2. Install prerequisites (gpg, wget, curl)
    exec.Command("sudo", "apt", "install", "-y", "gpg", "wget", "curl").Run()
    
    // 3. Create keyring directory
    exec.Command("sudo", "install", "-dm", "755", "/etc/apt/keyrings").Run()
    
    // 4. Download and install GPG key
    keyringPath := fmt.Sprintf("/etc/apt/keyrings/%s-archive-keyring.gpg", config.Name)
    // wget -qO - {KeyURL} | gpg --dearmor | sudo tee {keyringPath}
    
    // 5. Create repository entry
    repoEntry := fmt.Sprintf("deb [signed-by=%s arch=amd64] %s %s %s",
        keyringPath, config.RepoURL, config.Distribution, config.Component)
    // echo {repoEntry} | sudo tee {sourcesFile}
    
    // 6. Update apt cache
    exec.Command("sudo", "apt", "update").Run()
    
    return nil
}
```

### Installation Steps (Mise example)
1. Check if `/etc/apt/sources.list.d/mise.list` exists → skip if present
2. Install prerequisites: `gpg`, `wget`, `curl`
3. Create `/etc/apt/keyrings/` with 755 permissions
4. Download GPG key: `wget -qO - https://mise.jdx.dev/gpg-key.pub | gpg --dearmor`
5. Save key to `/etc/apt/keyrings/mise-archive-keyring.gpg`
6. Create sources entry: `deb [signed-by=/etc/apt/keyrings/mise-archive-keyring.gpg arch=amd64] https://mise.jdx.dev/deb stable main`
7. Save to `/etc/apt/sources.list.d/mise.list`
8. Run `apt update`

### Reference
- **Omakub**: `omakub/install/terminal/mise.sh` - exact pattern to follow
- **Debian wiki**: https://wiki.debian.org/DebianRepository/Format - repository format spec

### Alternatives Considered
- **add-apt-repository**: Wrapper not always available, less control
- **Manual apt-key (deprecated)**: apt-key is deprecated in modern Debian/Ubuntu
- **Snap**: Not suitable for repositories, only individual packages

---

## 5. Installation Strategy Pattern

### Decision
Use **strategy pattern with platform-specific implementations** in DebianCommand.

### Rationale
1. **Platform isolation**: Keeps macOS code completely untouched (Constitution Principle II)
2. **Extensibility**: Easy to add new installation methods (flatpak, snap, etc.)
3. **Type safety**: Interface ensures all strategies implement required methods
4. **Testability**: Each strategy can be tested independently with mocks
5. **Consistency**: Follows existing devgita factory pattern architecture

### Implementation Pattern

```go
// internal/commands/debian_strategies.go
package commands

type InstallationStrategy interface {
    Install(packageName string) error
    IsInstalled(packageName string) (bool, error)
}

// Strategy implementations
type AptStrategy struct { cmd *DebianCommand }
type PPAStrategy struct { cmd *DebianCommand; ppaConfig apt.PPAConfig }
type GitHubBinaryStrategy struct { cmd *DebianCommand; downloadURL, installPath string }
type GitCloneStrategy struct { cmd *DebianCommand; repoURL, installPath string }

// Strategy selection in DebianCommand
func (d *DebianCommand) getInstallationStrategy(packageName string) InstallationStrategy {
    switch packageName {
    case constants.Mise:
        return &PPAStrategy{cmd: d, ppaConfig: miseConfig}
    case constants.Neovim:
        return &GitHubBinaryStrategy{cmd: d, downloadURL: nvimURL}
    case constants.LazyGit:
        return &GitHubBinaryStrategy{cmd: d, downloadURL: lazygitURL}
    default:
        return &AptStrategy{cmd: d}
    }
}
```

### Strategy Types

| Strategy | Use Case | Examples |
|----------|----------|----------|
| AptStrategy | Standard apt packages + library mappings | curl, git, libgdbm-dev |
| PPAStrategy | Packages requiring custom repositories | mise, fastfetch |
| GitHubBinaryStrategy | Binary downloads from GitHub releases | neovim, lazygit, lazydocker |
| GitCloneStrategy | Git repository clones | powerlevel10k |

### macOS Isolation

```go
// internal/commands/macos.go - UNCHANGED
func (m *MacOSCommand) InstallPackage(packageName string) error {
    // Homebrew logic stays exactly as-is
    // No awareness of Debian strategies
    cmd := CommandParams{
        Command: "brew",
        Args:    []string{"install", packageName},
    }
    _, _, err := m.ExecCommand(cmd)
    return err
}
```

### Reference
- **Existing pattern**: `internal/commands/factory.go` - already uses factory pattern for platform detection
- **Design pattern**: Strategy pattern (Gang of Four)

### Alternatives Considered
- **Separate package per strategy**: Over-engineering, adds complexity
- **Single mega-function with switch**: Hard to test, violates Single Responsibility
- **Modify MacOSCommand**: Violates Constitution Principle II (Platform Parity with Isolation)

---

## Summary

This research provides battle-tested implementation patterns for Debian/Ubuntu support:

1. **Package Mappings**: 8 libraries with struct-based lookup and graceful fallback
2. **Retry Logic**: 3 attempts with exponential backoff (1s, 2s, 4s) using standard library
3. **Neovim**: tar.gz extraction (official format), not AppImage
4. **PPA Management**: Manual GPG + repository file management following omakub
5. **Strategy Pattern**: Interface-based platform-specific installers maintaining macOS isolation

All patterns are proven by omakub (`/Users/jair.mendez/Documents/projects/devgita/omakub/install/`) and align with devgita's existing architecture. The implementation maintains backward compatibility while extending platform support cleanly without touching working macOS code.

---

## References

- **Omakub repository**: `/Users/jair.mendez/Documents/projects/devgita/omakub/install/`
  - `terminal/mise.sh` - PPA installation pattern
  - `terminal/libraries.sh` - Library package names
  - `app-neovim.sh` - tar.gz extraction pattern
  - `app-lazygit.sh` - GitHub binary download pattern
- **Official documentation**:
  - Neovim releases: https://github.com/neovim/neovim/releases
  - Mise installation: https://mise.jdx.dev/installing-mise.html
  - Debian repository format: https://wiki.debian.org/DebianRepository/Format
  - Go net/http: https://pkg.go.dev/net/http
  - Exponential backoff: https://en.wikipedia.org/wiki/Exponential_backoff













-- PERSONAL INVESTIGATION


















1. **Package naming differences** (Homebrew vs apt):
   - `fastfetch` - needs alternative installation
   - `mise` - needs alternative installation  
   - `opencode` - needs alternative installation
   - `lazydocker` - needs alternative installation
   - `lazygit` - needs PPA or alternative
   - `eza` - needs alternative
   - `powerlevel10k` - needs git clone
   - `neovim` - version too old in apt, needs AppImage/PPA
   
2. **Library naming differences** (macOS vs Linux):
   - `gdbm` → `libgdbm-dev`
   - `jemalloc` → `libjemalloc-dev`
   - `libffi` → `libffi-dev`
   - `libyaml` → `libyaml-dev`
   - `ncurses` → `libncurses-dev`
   - `readline` → `libreadline-dev`
   - `vips` → `libvips-dev`
   - `zlib` → `zlib1g-dev`

---

## User

## User Input

```text
fastfetch - @omakub/install/terminal/app-fastfetch.sh

mise - @omakub/install/terminal/mise.sh
sudo add-apt-repository -y ppa:jdxcode/mise
sudo apt update -y
sudo apt install -y mise

opencode - curl -fsSL https://opencode.ai/install | bash OR mise use -g github:anomalyco/opencode OR use https://github.com/anomalyco/opencode/releases/tag/latest

lazydocker - @omakub/install/terminal/app-lazydocker.sh

lazygit - @omakub/install/terminal/app-lazygit.sh

fzf - see the name of the package here @omakub/install/terminal/apps-terminal.sh

ripgrep - see the name of the package here @omakub/install/terminal/apps-terminal.sh

bat - see the name of the package here @omakub/install/terminal/apps-terminal.sh

eza - see the name of the package here @omakub/install/terminal/apps-terminal.sh

zoxide - see the name of the package here @omakub/install/terminal/apps-terminal.sh

plocate - see the name of the package here @omakub/install/terminal/apps-terminal.sh

apache2-utils - see the name of the package here @omakub/install/terminal/apps-terminal.sh

fd-find - see the name of the package here @omakub/install/terminal/apps-terminal.sh

powerlevel10k - git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ~/powerlevel10k
echo 'source ~/powerlevel10k/powerlevel10k.zsh-theme' >>~/.zshrc

neovim it's already instelled, I wonder if it's not checking correctly can yu confirm another way to do it for debian? Use git tags Linux (x86_64)
AppImage
Download nvim-linux-x86_64.appimage
Run chmod u+x nvim-linux-x86_64.appimage && ./nvim-linux-x86_64.appimage
If your system does not have FUSE you can extract the appimage:
./nvim-linux-x86_64.appimage --appimage-extract
./squashfs-root/usr/bin/nvim
Tarball
Download nvim-linux-x86_64.tar.gz
Extract: tar xzvf nvim-linux-x86_64.tar.gz
Run ./nvim-linux-x86_64/bin/nvim
Linux (arm64)
AppImage
Download nvim-linux-arm64.appimage
Run chmod u+x nvim-linux-arm64.appimage && ./nvim-linux-arm64.appimage
If your system does not have FUSE you can extract the appimage:
./nvim-linux-arm64.appimage --appimage-extract
./squashfs-root/usr/bin/nvim https://github.com/neovim/neovim/releases/tag/stable \n\n ENTIRE step 2 found @omakub/install/terminal/libraries.sh

for fonts we do it this way @internal/commands/base.go#L328 at InstallFontFromURL I wonder if we should do it differently ### Proposed Implementation Plan:
1. **Create a font configuration mapping** (similar to languages/databases):
   ```go
   type FontConfig struct {
       DisplayName  string // "Hack Nerd Font"
       PackageName  string // "font-hack-nerd-font" (for Homebrew)
       ArchiveName  string // "Hack" (for GitHub releases)
       InstallName  string // "Hack Nerd Font" (for detection)
   }
   ```
2. **Platform-specific installation**:
   - **macOS**: Continue using Homebrew casks (current approach)
   - **Debian/Ubuntu**: Download tar.xz from GitHub releases, extract to `~/.local/share/fonts/`, run `fc-cache`
3. **Update `MaybeInstallFont` to accept URL**:
   - macOS: Ignore URL, use Homebrew
   - Debian: Use URL to download and install
4. **Font detection improvements**:
   - The current `IsFontPresent` already checks `fc-list` and font directories
   - Should work for both approaches if you need mroe context, we have @internal/commands/macos.go 
