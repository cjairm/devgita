# CLI Command Contracts

**Feature**: `001-binary-dist-audit`

## Existing Commands (updated behavior)

### `dg install`

```
Usage: dg install [flags]

Flags:
  --only <categories>   Only install specific categories (comma or repeatable)
  --skip <categories>   Skip specific categories (comma or repeatable)

Categories: terminal, languages, databases, desktop

Flow:
  1. Validate OS version (macOS 13+ / Debian 12+ / Ubuntu 24+)
  2. Bootstrap package manager (Homebrew on macOS, apt on Debian)
  3. Extract embedded configs to ~/.config/devgita/configs/    ← CHANGED (was: git clone)
  4. Configure devgita (GlobalConfig + shell config)
  5. Install terminal tools (if not skipped)
  6. Choose + install languages (if not skipped)
  7. Choose + install databases (if not skipped)
  8. Install desktop apps (if not skipped, platform-gated)

Exit codes:
  0 - Success
  1 - OS validation failed
  1 - Package manager bootstrap failed
  Non-zero on fatal errors; individual app failures are non-fatal (logged, continue)
```

## Planned Commands (roadmap — not implemented)

| Command | Purpose |
|---------|---------|
| `dg update` | Download latest binary from GitHub Releases, replace current |
| `dg configure` | Apply/update configurations for installed apps |
| `dg uninstall` | Remove devgita-managed packages and configs |
| `dg list` | Show installed packages and their status |
| `dg change --theme/--font` | Switch active themes or fonts |

## install.sh Contract

```
Usage: install.sh [--local <path>]

Flags:
  --local <path>   Install from local binary file (skip GitHub download)

Flow (default):
  1. Detect OS (Darwin/Linux) and arch (arm64/x86_64)
  2. Fetch latest release tag from GitHub API
  3. Download devgita-{os}-{arch} from GitHub Releases
  4. Place binary in ~/.local/bin/devgita
  5. Make executable (chmod +x)
  6. Add ~/.local/bin to PATH in shell config (if not present)
  7. Verify installation

Flow (--local):
  1. Validate provided file exists and is executable
  2. Copy to ~/.local/bin/devgita
  3. Steps 5-7 same as above

Supported platforms:
  - macOS arm64 (Apple Silicon)
  - macOS amd64 (Intel)
  - Linux amd64 (Debian/Ubuntu)

Exit codes:
  0 - Success
  1 - Unsupported OS or architecture
  1 - Download failed / file not found
  1 - Binary verification failed
```
