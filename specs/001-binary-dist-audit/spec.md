# Feature Specification: Binary Distribution with Embedded Configs

**Feature Branch**: `001-binary-dist-audit`
**Created**: 2026-03-29
**Status**: Draft
**Input**: Migrate devgita from repository-cloning to pre-built binary distribution with embedded configuration files. Audit all app modules, fix platform gaps, ensure zero-dependency end-to-end installation.

## Clarifications

### Session 2026-03-29

- Q: Desktop coordinator has `InstallDesktopAppsWithoutConfiguration()` (Docker, GIMP, Brave, Flameshot, Raycast) which is never called by `InstallAndConfigure()`. Are these dead code or a bug? → A: Bug — wire them into the desktop install flow with proper platform gating.
- Q: What are the Debian equivalents for macOS-only desktop apps (Raycast, Aerospace)? → A: Ulauncher replaces Raycast on Debian (install via `apt install ulauncher`). i3 replaces Aerospace on Debian (install via `apt install i3`). Both require new app modules.
- Q: How should maintainers test locally built binaries without uploading to GitHub? → A: `install.sh --local <path>` flag installs from a local file, skipping download. Document in README.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Zero-Dependency Install on Fresh Machine (Priority: P1)

A developer sets up a brand new macOS or Debian/Ubuntu machine. They have nothing installed — no Git, no Go, no Homebrew. They run a single `curl` command from the README, which downloads the devgita binary and places it on their PATH. They then run `dg install` and the tool bootstraps everything: package manager, terminal tools, languages, databases, desktop apps. The entire environment is ready without any manual steps.

**Why this priority**: This is the core value proposition. Without this, the tool requires pre-installed dependencies, defeating its purpose.

**Independent Test**: Run install.sh on a clean macOS VM and a clean Debian VM. Verify the binary is placed correctly, PATH is configured, and `dg install` runs successfully end to end.

**Acceptance Scenarios**:

1. **Given** a fresh macOS (arm64) machine with no developer tools, **When** the user runs `curl -fsSL https://raw.githubusercontent.com/cjairm/devgita/main/install.sh | bash`, **Then** the devgita binary is downloaded to `~/.local/bin/`, PATH is updated in the shell config, and `dg` is callable from a new terminal session.
2. **Given** a fresh Debian 12 (amd64) machine with no developer tools, **When** the user runs the same curl command, **Then** the devgita binary is downloaded, PATH is configured, and `dg` is callable.
3. **Given** the binary is installed, **When** the user runs `dg install`, **Then** Homebrew (macOS) or apt (Debian) is bootstrapped first, then all selected tools are installed without errors.
4. **Given** the user already has devgita installed, **When** they run the install.sh script again, **Then** the binary is replaced with the latest version and PATH is not duplicated in shell config.

---

### User Story 2 - Embedded Config Extraction (Priority: P2)

A developer runs `dg install` and the tool extracts configuration templates (Alacritty, Neovim, tmux, Aerospace/i3, fastfetch, OpenCode, shell templates) from the binary to `~/.config/devgita/`. The configs are then applied to each app. No repository cloning happens.

**Why this priority**: Without embedded configs, the binary cannot configure any apps. This replaces the current git-clone approach and is required for all app configuration to work.

**Independent Test**: Build the binary with embedded configs, run `dg install` on a clean system, verify `~/.config/devgita/configs/` contains all expected template files and that apps receive their configurations.

**Acceptance Scenarios**:

1. **Given** the binary is built with `configs/` embedded, **When** `dg install` runs, **Then** all files from the embedded `configs/` directory are extracted to `~/.config/devgita/configs/` preserving directory structure.
2. **Given** configs are already extracted from a previous run, **When** `dg install` runs again, **Then** existing configs are preserved (idempotent behavior).
3. **Given** an Alacritty config template contains a path placeholder like `{{.ConfigPath}}`, **When** the config is applied, **Then** the placeholder is resolved to the user's actual home directory path.

---

### User Story 3 - Platform-Correct App Installation (Priority: P3)

A developer runs `dg install` on Debian. macOS-only apps (Aerospace, Raycast, Xcode tools) are silently skipped and their Debian equivalents are installed instead: Ulauncher replaces Raycast, i3 replaces Aerospace. On macOS, all macOS apps install normally and Debian-only equivalents are skipped.

**Why this priority**: Without platform gating, the tool attempts to install unavailable packages and produces confusing errors. Platform parity is a constitution principle.

**Independent Test**: Run `dg install` on macOS and verify Aerospace and Raycast install. Run on Debian and verify Ulauncher and i3 install instead, with no errors related to macOS-only apps.

**Acceptance Scenarios**:

1. **Given** a Debian system, **When** `dg install` runs the desktop category, **Then** Aerospace and Raycast are not attempted; Ulauncher and i3 are installed instead.
2. **Given** a macOS system, **When** `dg install` runs, **Then** Aerospace and Raycast install normally via Homebrew cask; Ulauncher and i3 are not attempted.
3. **Given** a Debian system, **When** shell configuration is generated, **Then** paths use apt/system paths instead of `$(brew --prefix)` references.

---

### User Story 4 - Audit Gaps Fixed (Priority: P4)

All app modules pass an audit: every app that configures files has templates in `configs/`, the git config directory is populated, OpenCode and Mise configurations are applied during `dg install --only terminal`, the shell template works on both platforms, desktop apps without configuration (Docker, GIMP, Brave, Flameshot, Raycast/Ulauncher) are wired into the desktop install flow, and the `source` builtin call in terminal coordinator is removed.

**Why this priority**: Existing bugs prevent correct operation. Fixing them is prerequisite to reliable binary distribution.

**Independent Test**: Run `dg install --only terminal` and verify OpenCode config is generated, Mise shell integration is active, and git config is applied. Run `dg install --only desktop` and verify Docker, Brave, GIMP, Flameshot install. On macOS verify Raycast installs; on Debian verify Ulauncher and i3 install. Run on both platforms.

**Acceptance Scenarios**:

1. **Given** the `configs/git/` directory, **When** inspected, **Then** it contains a `.gitconfig` template with sensible defaults.
2. **Given** a user runs `dg install --only terminal`, **When** OpenCode is installed, **Then** its configuration (template, theme, agents) is also applied.
3. **Given** a user runs `dg install --only terminal`, **When** Mise is installed, **Then** shell integration is also configured.
4. **Given** a Debian system, **When** shell config is generated with Zsh features enabled, **Then** plugin paths reference `/usr/share/` locations instead of `$(brew --prefix)`.
5. **Given** a user runs `dg install --only desktop` on Debian, **When** the desktop category executes, **Then** Docker, GIMP, Brave, Flameshot, Ulauncher, and i3 are installed. Raycast and Aerospace are skipped.
6. **Given** a user runs `dg install --only desktop` on macOS, **When** the desktop category executes, **Then** Docker, GIMP, Brave, Flameshot, Raycast, and Aerospace are installed. Ulauncher and i3 are skipped.
7. **Given** the terminal coordinator, **When** it completes installation, **Then** it does not attempt to call `source` as a subprocess (shell builtin cannot be exec'd).

---

### User Story 5 - Local Binary Build and Release (Priority: P5)

A maintainer builds binaries locally for all supported platforms (darwin-arm64, darwin-amd64, linux-amd64) and can test them locally before uploading to GitHub Releases. The install.sh script fetches the correct binary based on detected OS and architecture. Local testing MUST be possible without uploading to GitHub.

**Why this priority**: This enables the distribution mechanism. Lower priority because it's a maintainer workflow, not end-user facing.

**Independent Test**: Build all three binaries locally. Test the local binary by running it directly (`./devgita-darwin-arm64 install`). Once validated, create a GitHub Release and verify install.sh downloads correctly.

**Acceptance Scenarios**:

1. **Given** a maintainer on macOS, **When** they run the documented build commands, **Then** three binaries are produced: `devgita-darwin-arm64`, `devgita-darwin-amd64`, `devgita-linux-amd64`.
2. **Given** a locally built binary, **When** the maintainer runs it directly (e.g., `./devgita-darwin-arm64 install`), **Then** it works identically to an installed binary — extracting configs, bootstrapping packages, and completing the full install flow.
3. **Given** binaries are uploaded to a GitHub Release with a tag, **When** install.sh runs on a macOS arm64 machine, **Then** it downloads `devgita-darwin-arm64`.
4. **Given** no releases exist on GitHub, **When** install.sh runs, **Then** it prints a clear error message and exits.
5. **Given** install.sh is passed a `--local <path>` flag, **When** the maintainer runs it, **Then** it installs from the local file instead of downloading from GitHub Releases.

---

### User Story 6 - README Roadmap (Priority: P6)

The README documents the current installation method (install.sh + `dg install`) and lists planned commands (`dg update`, `dg configure`, `dg uninstall`, `dg list`, `dg change`) as a roadmap. Users understand what exists today and what's coming.

**Why this priority**: Documentation. Important but does not block functionality.

**Independent Test**: Read the README and verify install instructions work and roadmap commands are listed.

**Acceptance Scenarios**:

1. **Given** the README, **When** a user reads the installation section, **Then** they find a single curl command to install devgita.
2. **Given** the README, **When** a user reads the commands section, **Then** `dg install` is documented as available and `dg update`, `dg configure`, `dg uninstall`, `dg list`, `dg change` are listed as planned/roadmap.

---

### Edge Cases

- What happens when install.sh is run on an unsupported OS (e.g., Fedora, Windows)? It prints a clear error and exits with non-zero code.
- What happens when the user's `~/.config/devgita/` already has manually edited configs? `SoftInstall`/`SoftConfigure` preserve existing files. `ForceInstall`/`ForceConfigure` overwrite.
- What happens when GitHub Releases are unreachable (no internet)? install.sh fails with a clear network error message.
- What happens when `configs/` is embedded but a specific app's config directory is missing? The app's `SoftConfigure` detects missing source and returns an error that is logged but non-fatal.
- What happens on macOS Intel (amd64) vs Apple Silicon (arm64)? install.sh detects `uname -m` and downloads the correct binary variant.
- What happens when `dg install --only terminal` is run without `--only languages`? Mise is installed AND shell-configured so it's ready for manual use.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The binary MUST embed all files under `configs/` at compile time.
- **FR-002**: `dg install` MUST extract embedded configs to `~/.config/devgita/configs/` on first run.
- **FR-003**: `SoftInstall()` in devgita module MUST detect whether configs are already extracted (not whether a git repo exists).
- **FR-004**: `install.sh` MUST detect OS (`Darwin`/`Linux`) and architecture (`arm64`/`x86_64`), download the matching binary from GitHub Releases, and place it in `~/.local/bin/`.
- **FR-005**: `install.sh` MUST configure PATH in the user's shell config file (`.zshrc`, `.bashrc`, or `.bash_profile`).
- **FR-006**: `install.sh` MUST NOT require any pre-installed tools beyond `curl` and a POSIX shell.
- **FR-007**: Desktop coordinator MUST install Aerospace and Raycast only on macOS, and i3 and Ulauncher only on Debian/Ubuntu.
- **FR-008**: Shell template (`devgita.zsh.tmpl`) MUST use platform-conditional paths for Zsh plugins (Homebrew paths on macOS, system paths on Debian).
- **FR-009**: `configs/git/` MUST contain a `.gitconfig` template with sensible defaults.
- **FR-010**: Terminal coordinator MUST call `SoftConfigure()` for OpenCode after `SoftInstall()`.
- **FR-011**: Terminal coordinator MUST call `SoftConfigure()` for Mise after `SoftInstall()`.
- **FR-012**: Binary builds MUST be producible locally with documented cross-compilation commands for darwin-arm64, darwin-amd64, and linux-amd64.
- **FR-013**: GlobalConfig schema MUST remain backward-compatible. The additive `shell.is_mac` field MUST NOT break deserialization of existing configs. Config paths (`~/.config/devgita/`) MUST remain unchanged.
- **FR-014**: README MUST document the install.sh installation method and list planned commands as roadmap.
- **FR-015**: `install.sh` on unsupported platforms MUST exit with a clear error message.
- **FR-016**: Config templates that reference paths MUST use template variables and resolve correctly after extraction from embedded files.
- **FR-017**: Desktop coordinator MUST call `InstallDesktopAppsWithoutConfiguration()` within `InstallAndConfigure()` with platform-specific apps gated appropriately.
- **FR-018**: Terminal coordinator MUST NOT attempt to exec `source` as a subprocess (it is a shell builtin). Remove the broken `source` call.
- **FR-019**: `cmd/install.go` Long description MUST be updated to remove reference to "Clones the devgita repository" and reflect the new embedded config extraction approach.
- **FR-020**: Git installation MUST move from essential prerequisite in `cmd/install.go` to the terminal tools category, since repository cloning is no longer required for devgita setup.
- **FR-021**: New app modules MUST be created for Ulauncher (`internal/apps/ulauncher/`) and i3 (`internal/apps/i3/`) following the standard app interface pattern.
- **FR-022**: i3 configuration template MUST be created in `configs/i3/` with sensible defaults for development workflows.
- **FR-023**: Locally built binaries MUST be directly executable for testing without requiring upload to GitHub Releases.
- **FR-024**: `install.sh` MUST support a `--local <path>` flag to install from a local binary file instead of downloading from GitHub.

### Key Entities

- **Embedded Configs**: The `configs/` directory tree compiled into the binary. Contains templates for alacritty, aerospace, fastfetch, git, i3, neovim, opencode, templates (shell), and tmux.
- **GlobalConfig**: YAML state file at `~/.config/devgita/global_config.yaml` tracking installed packages, shell features, and paths. Schema unchanged.
- **Install Script**: Shell script (`install.sh`) at repository root that downloads and installs the binary. Supports `--local` flag for testing.
- **Platform Gate**: Runtime check (`Platform.IsMac()` / `Platform.IsLinux()`) that determines whether a platform-specific app should be installed or skipped.
- **Platform Equivalents**: Apps that serve the same function on different platforms: Raycast (macOS) / Ulauncher (Debian), Aerospace (macOS) / i3 (Debian).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user on a clean macOS or Debian machine can go from zero to fully configured development environment using only `curl` + `dg install` with no other prerequisites.
- **SC-002**: All config directories under `configs/` contain valid templates (no empty directories).
- **SC-003**: `dg install` on Debian installs Ulauncher and i3 instead of Raycast and Aerospace, with zero macOS-only errors.
- **SC-004**: Shell configuration generated on Debian correctly sources Zsh plugins from system paths.
- **SC-005**: All app modules that have meaningful `ForceConfigure` methods are called during the standard `dg install` flow.
- **SC-006**: Configs extracted from the binary are byte-identical to the source files in `configs/` (no corruption during embed/extract).
- **SC-007**: Binary size remains under 50MB for each platform variant.
- **SC-008**: `dg install --only desktop` installs Docker, GIMP, Brave, Flameshot on both platforms, Raycast and Aerospace on macOS only, Ulauncher and i3 on Debian only.
- **SC-009**: A locally built binary can be tested end-to-end without uploading to GitHub Releases.

## Assumptions

- Users have `curl` and a POSIX-compatible shell (`/bin/bash` or `/bin/sh`) available on fresh systems. Both macOS and Debian include these by default.
- `~/.local/bin/` is the standard user-local binary directory and is appropriate for both macOS and Linux.
- The GitHub repository is public and GitHub Releases are accessible without authentication.
- Debian users have `sudo` access for `apt install` operations.
- macOS users are on Ventura (13+), Debian users on Bookworm (12+), Ubuntu users on Noble (24+), matching existing version validation.
- TUI features (`internal/tui/`) are excluded from this scope — existing functionality preserved but not enhanced.
- No CI/CD pipeline is needed. Binaries are built locally and uploaded manually.
- Ulauncher is available via `apt install ulauncher` on Ubuntu (PPA) and via .deb download on Debian.
- i3 is available via `apt install i3` on both Debian and Ubuntu from official repositories.
