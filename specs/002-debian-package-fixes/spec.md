# Feature Specification: Debian/Ubuntu Package Installation Fixes

**Feature Branch**: `002-debian-package-fixes`  
**Created**: 2026-04-05  
**Status**: Draft  
**Input**: User description: "Fix Debian/Ubuntu package installation errors with platform-specific package names and installation methods"

## Clarifications

### Session 2026-04-09

- Q: Installation Failure Behavior → A: Stop installation of that component only; continue with remaining components, log error, track as failed in GlobalConfig
- Q: Library Package Existence Check Strategy → A: Check if macOS library name exists in constants; if yes, use Debian mapping, if no, install as-is
- Q: GitHub Binary Download Retry Strategy → A: Retry 3 times with exponential backoff; if all fail, log error and continue
- Q: Neovim Installation Method for Debian → A: Download official tar.gz from GitHub releases, extract to /usr/local, install binary to /usr/local/bin/nvim
- Q: Installation Summary Report Format → A: Simple count: "Installed: 12, Failed: 2, Skipped: 1"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Install Terminal Tools on Debian/Ubuntu (Priority: P1)

A developer runs `dg install --only terminal` on a fresh Debian 12 or Ubuntu 24 system. The system installs all terminal development tools (fastfetch, mise, neovim, lazygit, lazydocker, eza, powerlevel10k) with correct package names and installation methods for the Debian/Ubuntu platform.

**Why this priority**: Terminal tools are the foundation of the development environment. Without these working, developers cannot use the core devgita functionality. This affects the primary value proposition of the tool.

**Independent Test**: Can be fully tested by running `dg install --only terminal` on a clean Debian/Ubuntu VM and verifying all packages install without errors and are accessible in PATH.

**Acceptance Scenarios**:

1. **Given** a Debian 12 system without devgita installed, **When** user runs `dg install --only terminal`, **Then** fastfetch installs via PPA without "Unable to locate package" error
2. **Given** an Ubuntu 24 system, **When** user runs `dg install --only terminal`, **Then** mise installs via mise.jdx.dev repository with GPG key setup
3. **Given** a Debian system, **When** terminal installation completes, **Then** neovim 0.11.1+ is installed via tar.gz download from GitHub releases instead of outdated apt version
4. **Given** a Debian system, **When** lazygit installation runs, **Then** binary is downloaded from GitHub releases and installed to /usr/local/bin
5. **Given** a Debian system, **When** lazydocker installation runs, **Then** binary is downloaded from GitHub releases and installed to /usr/local/bin
6. **Given** an Ubuntu system, **When** eza installation runs, **Then** package installs successfully from standard apt repositories
7. **Given** a Debian system, **When** powerlevel10k installation runs, **Then** theme is cloned from GitHub to ~/powerlevel10k and sourced in .zshrc
8. **Given** a Debian system, **When** opencode installation runs, **Then** system tries mise first, falls back to install script if mise unavailable
9. **Given** a Debian system, **When** fastfetch installation fails, **Then** system logs error, continues with remaining terminal tools, and tracks fastfetch as failed in GlobalConfig
10. **Given** a Debian system, **When** lazygit download fails temporarily, **Then** system retries up to 3 times with exponential backoff (1s, 2s, 4s) before marking as failed
11. **Given** a Debian system, **When** neovim tar.gz installation runs, **Then** archive extracts to /usr/local with bin/, lib/, and share/ directories preserved
12. **Given** a Debian system, **When** neovim installation completes, **Then** nvim binary is executable at /usr/local/bin/nvim and tar.gz temporary file is cleaned up
13. **Given** a Debian system, **When** terminal installation completes with 2 failures, **Then** summary displays "Installed: 6, Failed: 2, Skipped: 0"

---

### User Story 2 - Install Core Libraries on Debian/Ubuntu (Priority: P2)

A developer runs `dg install` and the system installs core development libraries with correct Debian package names (libgdbm-dev instead of gdbm, libffi-dev instead of libffi, etc.).

**Why this priority**: Core libraries are required for language runtime compilation (Python, Ruby, etc.). Without correct library names, language installations fail. This is second priority because it blocks language installation but not terminal tools.

**Independent Test**: Can be tested by running `dg install` on Debian/Ubuntu and verifying all library packages install successfully, then compiling a Python extension that requires these libraries.

**Acceptance Scenarios**:

1. **Given** a Debian system, **When** core libraries install, **Then** libgdbm-dev installs instead of "gdbm" which doesn't exist in apt
2. **Given** a Debian system, **When** core libraries install, **Then** libjemalloc-dev installs instead of "jemalloc"
3. **Given** a Debian system, **When** core libraries install, **Then** libffi-dev installs instead of "libffi"
4. **Given** a Debian system, **When** core libraries install, **Then** libyaml-dev installs instead of "libyaml"
5. **Given** a Debian system, **When** core libraries install, **Then** libncurses-dev installs instead of "ncurses"
6. **Given** a Debian system, **When** core libraries install, **Then** libreadline-dev installs instead of "readline"
7. **Given** a Debian system, **When** core libraries install, **Then** libvips-dev installs instead of "vips"
8. **Given** a Debian system, **When** core libraries install, **Then** zlib1g-dev installs instead of "zlib"
9. **Given** a Debian system, **When** a library installation fails, **Then** system logs error, continues with remaining libraries, and tracks failed library in GlobalConfig
10. **Given** a Debian system, **When** a library not in mapping is requested (e.g., "curl"), **Then** system installs "curl" as-is without attempting mapping

---

### User Story 3 - Install Fonts on Debian/Ubuntu (Priority: P3)

A developer runs `dg install` and the system installs Nerd Fonts by downloading tar.xz archives from GitHub releases, extracting to ~/.local/share/fonts/, and running fc-cache to update the font cache.

**Why this priority**: Fonts enhance the developer experience but are not critical for core functionality. Terminal emulators work without custom fonts, just with reduced visual appeal.

**Independent Test**: Can be tested by running `dg install --only desktop` on Debian/Ubuntu and verifying fonts appear in `fc-list` output and are usable in Alacritty/terminals.

**Acceptance Scenarios**:

1. **Given** a Debian system, **When** font installation runs, **Then** Hack Nerd Font downloads from GitHub releases as Hack.tar.xz
2. **Given** a Debian system, **When** font archive downloads, **Then** fonts extract to ~/.local/share/fonts/
3. **Given** fonts are extracted, **When** installation completes, **Then** fc-cache runs to refresh system font cache
4. **Given** fonts are installed, **When** user runs `fc-list | grep "Hack"`, **Then** Hack Nerd Font appears in output
5. **Given** a Debian system, **When** fonts are already installed, **Then** system skips re-download and extraction
6. **Given** a Debian system, **When** a font download fails, **Then** system logs error, continues with remaining fonts, and tracks failed font in GlobalConfig
7. **Given** a Debian system, **When** font download fails with network timeout, **Then** system retries up to 3 times with exponential backoff before marking as failed

---

### Edge Cases

- What happens when GitHub API rate limits are hit during binary downloads? → System retries with exponential backoff, logs rate limit error after 3 attempts, continues with remaining installations
- How does system handle network failures during PPA addition or binary downloads? → System retries network operations up to 3 times with exponential backoff (1s, 2s, 4s)
- What happens when neovim AppImage fails to execute due to missing FUSE? → System extracts AppImage manually using --appimage-extract and installs nvim binary to /usr/local/bin
- How does system handle when snap is not available for ulauncher installation?
- What happens when font URLs change or GitHub releases structure changes?
- How does system handle when user's ~/.local/share/fonts directory doesn't exist?
- When a package installation fails, does the system skip its configuration step? → Yes, configuration is skipped for failed installations
- How are failed installations reported to the user at the end of the installation process? → Summary displays simple count: "Installed: X, Failed: Y, Skipped: Z"
- What happens when a library package name exists in both the mapping and as a direct apt package?
- How does the system handle future library additions without breaking existing mappings?
- What happens when AppImage extraction fails or produces unexpected directory structure?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST map macOS package names to correct Debian/Ubuntu equivalents (gdbm → libgdbm-dev, libffi → libffi-dev, etc.)
- **FR-002**: System MUST install fastfetch via PPA (ppa:zhangsongcui3371/fastfetch) instead of attempting apt install
- **FR-003**: System MUST install mise via mise.jdx.dev repository with GPG key setup instead of attempting direct apt install
- **FR-004**: System MUST download and install neovim 0.11.1+ tar.gz from GitHub releases instead of using outdated apt version (0.7.2)
- **FR-005**: System MUST download lazygit/lazydocker binaries from GitHub releases instead of attempting apt install
- **FR-006**: System MUST clone powerlevel10k from GitHub and add source line to .zshrc instead of attempting apt install
- **FR-007**: System MUST install opencode via mise first (mise use -g github:anomalyco/opencode); if mise is unavailable or fails, fallback to official install script (https://opencode.ai/install)
- **FR-008**: System MUST install fonts by downloading Nerd Font archives from GitHub and extracting to ~/.local/share/fonts/
- **FR-009**: System MUST run fc-cache after font installation to refresh system font cache
- **FR-010**: System MUST check neovim version using `nvim --version` before attempting installation
- **FR-011**: System MUST create platform-specific installation methods in DebianCommand that differ from macOS Homebrew approach
- **FR-012**: System MUST handle ulauncher installation failure gracefully when snap is unavailable
- **FR-013**: System MUST preserve existing error handling and GlobalConfig tracking for all installation methods
- **FR-014**: System MUST NOT configure a package if its installation failed
- **FR-015**: System MUST continue installing remaining components when a single component installation fails
- **FR-016**: System MUST log installation failures with component name and error details
- **FR-017**: System MUST track failed installations separately in GlobalConfig to differentiate from successful installations
- **FR-018**: System MUST check if a library package name exists in the Debian mapping using O(1) map lookup before applying transformation; if not found, use original package name unchanged
- **FR-019**: System MUST install unmapped library packages as-is without transformation (e.g., "curl" remains "curl")
- **FR-020**: System MUST NOT remove or modify existing library mappings when adding new ones
- **FR-021**: System MUST retry GitHub binary downloads exactly 3 times with exponential backoff delays (1s, 2s, 4s)
- **FR-022**: System MUST retry font archive downloads exactly 3 times with exponential backoff delays (1s, 2s, 4s)
- **FR-023**: System MUST log each retry attempt with attempt number and delay duration
- **FR-024**: System MUST mark download as failed only after all 3 retry attempts are exhausted
- **FR-025**: System MUST download neovim tar.gz archive from official GitHub releases (nvim-linux64.tar.gz)
- **FR-026**: System MUST extract tar.gz to /usr/local preserving directory structure (bin/, lib/, share/)
- **FR-027**: System MUST verify extracted nvim binary is executable at /usr/local/bin/nvim
- **FR-028**: System MUST clean up temporary tar.gz file after successful installation
- **FR-029**: System MUST display installation summary at completion in format: "Installed: X, Failed: Y, Skipped: Z"
- **FR-030**: System MUST track installation counts (installed, failed, skipped) throughout installation process
- **FR-031**: System MUST increment "Installed" count only for successfully completed installations
- **FR-032**: System MUST increment "Failed" count for installations that failed after all retry attempts
- **FR-033**: System MUST increment "Skipped" count for components not attempted due to already being installed

### Key Entities

- **PackageMapping**: Represents the mapping between generic package name (used in constants), macOS package name (Homebrew), and Debian package name (apt)
- **InstallationStrategy**: Represents different installation approaches (apt, PPA, GitHub binary download, git clone, install script)
- **FontConfig**: Represents font metadata including display name, Homebrew cask name, GitHub archive name, and installation name for detection
- **InstallationResult**: Represents the outcome of a component installation (success, failure, skipped) with error details and tracking info
- **LibraryMapping**: Maps macOS library names (gdbm, libffi, etc.) to Debian package names (libgdbm-dev, libffi-dev, etc.) with lookup capability
- **RetryConfig**: Defines retry behavior including max attempts (3), backoff strategy (exponential), and delay sequence (1s, 2s, 4s)
- **TarGzExtractor**: Handles tar.gz archive extraction to /usr/local with directory structure preservation
- **InstallationSummary**: Tracks counts of installed, failed, and skipped components for final summary report

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All terminal tools install without "Unable to locate package" errors on Debian 12 and Ubuntu 24
- **SC-002**: Core library installations succeed with correct -dev package names, allowing Python/Ruby compilation
- **SC-003**: Neovim version 0.11.1 or higher is installed and accessible, replacing the outdated 0.7.2 apt version
- **SC-004**: Lazygit and lazydocker binaries are executable from PATH after installation
- **SC-005**: All installed fonts appear in fc-list output and are selectable in terminal emulators
- **SC-006**: Installation process completes without fatal errors, only warnings for non-critical packages like ulauncher
- **SC-007**: All successfully installed packages are tracked in GlobalConfig for future uninstall capability
- **SC-008**: When component installation fails, subsequent components continue installing and successful installations are tracked correctly
- **SC-009**: Failed component installations do NOT proceed to configuration step
- **SC-010**: Installation summary displays count of successful, failed, and skipped components at completion in format "Installed: X, Failed: Y, Skipped: Z"
- **SC-011**: Library packages without Debian mappings install successfully using their original names
- **SC-012**: All 8 mapped libraries (gdbm, jemalloc, libffi, libyaml, ncurses, readline, vips, zlib) install with correct Debian equivalents
- **SC-013**: Transient network failures (1-2 failures followed by success) do NOT cause installation failures due to retry mechanism
- **SC-014**: GitHub binary downloads succeed on retry after initial network timeout, completing within 3 attempts
- **SC-015**: Neovim tar.gz extracts successfully to /usr/local with all required directories (bin/, lib/, share/)
- **SC-016**: Installed neovim binary is executable from /usr/local/bin/nvim and reports version 0.11.1 or higher
- **SC-017**: Installation summary counts accurately reflect actual installation outcomes (verified against GlobalConfig entries)

## Assumptions

- Users have stable internet connectivity to download binaries from GitHub and add PPAs
- Users have sudo privileges to install packages and add apt repositories
- GitHub releases API structure remains stable for lazygit/lazydocker/neovim downloads
- Font structure in Nerd Fonts GitHub releases remains consistent (tar.xz archives)
- macOS installations continue using Homebrew without changes (no regression)
- The omakub scripts referenced provide working installation patterns that can be adapted
- Neovim GitHub releases maintain tar.gz format for Linux installations
- The mise.jdx.dev repository remains available and stable for mise installation
- The library mapping is maintained as a static lookup table in code, not a dynamic configuration file
- Network failures are transient and resolve within the 3-retry window for most cases
- Exponential backoff delays (1s, 2s, 4s totaling 7s) are acceptable wait times for users during installation
- Systems have sufficient disk space for tar.gz extraction to /usr/local (~50MB)
- /usr/local/bin is in user's PATH for nvim binary access
- Simple count summary format provides sufficient visibility into installation outcomes for users
