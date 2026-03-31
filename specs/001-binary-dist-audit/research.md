# Research: Binary Distribution with Embedded Configs

**Feature**: `001-binary-dist-audit`
**Date**: 2026-03-29

## R1: Go embed.FS — Extracting Nested Directory Trees

**Decision**: Use `//go:embed all:configs` directive (not `configs/*`).

**Rationale**: The `*` wildcard only matches one directory level deep. The `configs/` tree has files up to 5 levels deep (e.g., `configs/neovim/lua/custom/plugins/treesitter.lua`). Using `configs/*` would silently miss all nested content. The bare name `configs` recurses but skips dotfiles. The `all:` prefix includes hidden files (`.gitconfig`, etc.).

**Extraction pattern**: `fs.WalkDir` over the `embed.FS`, stripping the `configs` prefix, writing to `destDir` preserving structure. New package: `internal/embedded/configs.go`.

**Alternatives considered**:
- `//go:embed configs/*` — Rejected: does not recurse into subdirectories
- `//go:embed configs` — Rejected: skips dotfiles like `.gitconfig`
- External archive (tar/zip) — Rejected: adds complexity per Principle IV

## R2: Go Cross-Compilation Commands

**Decision**: Three build commands from macOS, no CGO needed.

```bash
GOOS=darwin GOARCH=arm64 go build -o devgita-darwin-arm64 main.go
GOOS=darwin GOARCH=amd64 go build -o devgita-darwin-amd64 main.go
GOOS=linux  GOARCH=amd64 go build -o devgita-linux-amd64 main.go
```

**Rationale**: Pure Go project, no CGO dependencies. Go compiler natively cross-compiles all three targets from a single machine.

**Bug found**: Existing README (lines 98-101) has architecture labels swapped — "M chips" labeled as `amd64` and "Intel chips" as `arm64`. Must fix.

**Alternatives considered**:
- Docker-based build — Rejected: unnecessary complexity
- CI/CD pipeline — Rejected: per spec, local builds only

## R3: Shell Template Platform-Conditional Paths

**Decision**: Add `IsMac bool` field to `ShellFeatures` struct. Use `{{if .IsMac}}` in template for conditional plugin paths.

**Rationale**: `ShellFeatures` is already the template data context passed to `GenerateFromTemplate()`. Adding a boolean field is the simplest change — no pipeline modifications needed. The field is persisted in GlobalConfig YAML, making shell config regeneration deterministic from saved state.

**Linux plugin paths** (Debian/Ubuntu standard):

| Plugin | macOS (Homebrew) | Debian/Ubuntu |
|--------|------------------|---------------|
| zsh-autosuggestions | `$(brew --prefix)/share/zsh-autosuggestions/zsh-autosuggestions.zsh` | `/usr/share/zsh-autosuggestions/zsh-autosuggestions.zsh` |
| zsh-syntax-highlighting | `$(brew --prefix)/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh` | `/usr/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh` |
| powerlevel10k | `$(brew --prefix)/share/powerlevel10k/powerlevel10k.zsh-theme` | `/usr/share/powerlevel10k/powerlevel10k.zsh-theme` |

**Implementation**: Set `gc.Shell.IsMac = runtime.GOOS == "darwin"` when GlobalConfig is created. Template uses `{{if .IsMac}}...{{else}}...{{end}}` blocks.

**Alternatives considered**:
- Pass `Platform` interface to template — Rejected: templates take data structs, not interfaces
- Separate template files per platform — Rejected: violates DRY, doubles maintenance
- Runtime `$(uname)` check in shell script — Rejected: template should be static after generation

## R4: Ulauncher Installation on Debian

**Decision**: Install via `apt install ulauncher` on Ubuntu (PPA available). On Debian, use direct `.deb` download from GitHub releases.

**Rationale**: Ulauncher is not in Debian's official repositories. Ubuntu has a PPA (`ppa:agornostal/ulauncher`). For Debian, the most reliable approach is downloading the `.deb` from `https://github.com/Ulauncher/Ulauncher/releases`.

**Implementation**: The `DebianCommand.InstallPackage()` handles `apt install`. For Debian without PPA, may need a custom install method that downloads and installs the `.deb` file. The app module should try `apt install ulauncher` first and fall back to `.deb` download if not found.

## R5: i3 Installation and Configuration

**Decision**: Install via `apt install i3` on both Debian and Ubuntu. Create `configs/i3/config` with development-focused defaults.

**Rationale**: i3 is in official repositories for both Debian 12+ and Ubuntu 24+. The `i3` metapackage includes i3-wm, i3lock, i3status, dunst, and suckless-tools. i3-gaps has been merged with i3 since version 4.22 (Debian 12).

**Config defaults**: Mod key = Windows/Super, terminal = alacritty (if installed), gaps = 10px, vim-style navigation ($Mod+h/j/k/l), workspace keybindings $Mod+1-9.

## R6: install.sh Design

**Decision**: Single bash script. Detect OS/arch, download from GitHub Releases, install to `~/.local/bin/`, configure PATH. Support `--local <path>` flag.

**Flow**:
1. Detect `uname -s` (Darwin/Linux) and `uname -m` (arm64/x86_64)
2. Map to binary name: `devgita-{os}-{arch}`
3. If `--local <path>` flag: copy local binary to install dir
4. Else: fetch latest release tag from GitHub API, download binary
5. Make executable, move to `~/.local/bin/`
6. Detect shell config file (`.zshrc`, `.bashrc`, `.bash_profile`)
7. Add `~/.local/bin` to PATH if not present
8. Verify installation

**Rationale**: Follows rocketctl install.sh pattern. Minimal dependencies (curl + POSIX shell). `--local` flag enables testing without GitHub upload.
