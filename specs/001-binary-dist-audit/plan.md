# Implementation Plan: Binary Distribution with Embedded Configs

**Branch**: `001-binary-dist-audit` | **Date**: 2026-03-29 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-binary-dist-audit/spec.md`

## Summary

Migrate devgita from a repository-cloning installation model to pre-built binary distribution with Go-embedded configuration templates. The binary is self-contained: users download it via a simple install.sh script and run `dg install` which bootstraps the package manager, extracts embedded configs, and installs the full development environment. This plan also covers auditing all existing app modules for platform parity (macOS/Debian equivalents), fixing workflow bugs (dead desktop app code, broken `source` call, missing configure calls), and creating new app modules for Debian equivalents (Ulauncher, i3).

## Technical Context

**Language/Version**: Go 1.21+ (existing project, uses `embed` package from Go 1.16+)
**Primary Dependencies**: Cobra CLI, gopkg.in/yaml.v3, Go `embed`, Go `text/template`
**Storage**: YAML files on disk (`~/.config/devgita/global_config.yaml`), embedded filesystem via `embed.FS`
**Testing**: `go test` with `MockBaseCommand` and `MockCommand` interfaces, 3 isolation levels per `docs/guides/testing-patterns.md`
**Target Platform**: macOS 13+ (Ventura) via Homebrew, Debian 12+ (Bookworm) / Ubuntu 24+ via apt
**Project Type**: CLI tool
**Performance Goals**: N/A (one-shot installer, not latency-sensitive)
**Constraints**: Binary size < 50MB per variant, zero pre-installed dependencies for end user
**Scale/Scope**: ~17 app modules, ~8 config directories, 3 binary targets (darwin-arm64, darwin-amd64, linux-amd64)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Evidence |
|-----------|--------|----------|
| I. Zero-Dependency Distribution | PASS | Binary embeds configs via `embed.FS`. install.sh downloads binary. No git/Go needed. |
| II. Platform Parity with Isolation | PASS (after fixes) | Audit identified gaps: Aerospace/Raycast not gated, no Debian equivalents. Plan adds i3, Ulauncher, platform gates. |
| III. Idempotent and Safe | PASS | `SoftInstall()`/`SoftConfigure()` preserve existing state. GlobalConfig tracks installed vs pre-existing. |
| IV. Simplicity Over Verbosity | PASS | install.sh is minimal. No CI/CD. Local builds. `--local` flag for testing. |
| V. Testability | PASS | All new modules use `BaseCommandExecutor` interface. New apps (i3, Ulauncher) follow mock patterns. |
| VI. Configuration as Templates | PASS (after fixes) | Audit found `configs/git/` empty. Plan populates it. All apps with configure methods will have templates. |
| VII. Audit Before Shipping | PASS | Full audit completed during spec phase. Bugs documented. Fixes planned. |

**Gate result**: PASS — no violations. Proceed to Phase 0.

## Project Structure

### Documentation (this feature)

```text
specs/001-binary-dist-audit/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (CLI command interface)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
devgita/
├── main.go                          # Entry point
├── configs/                         # Embedded config templates
│   ├── aerospace/                   # macOS-only tiling WM config
│   │   └── aerospace.toml
│   ├── alacritty/                   # Cross-platform terminal config
│   │   ├── alacritty.toml.tmpl
│   │   └── starter.sh
│   ├── fastfetch/                   # Cross-platform system info config
│   │   └── config.jsonc
│   ├── git/                         # Cross-platform git config (currently EMPTY — to fix)
│   │   └── .gitconfig               # NEW: sensible defaults template
│   ├── i3/                          # NEW: Debian-only tiling WM config
│   │   └── config                   # NEW: i3 config with dev defaults
│   ├── neovim/                      # Cross-platform editor config
│   │   ├── init.lua
│   │   └── lua/...
│   ├── opencode/                    # Cross-platform editor config
│   │   ├── opencode.json.tmpl
│   │   ├── themes/
│   │   ├── agents/
│   │   └── commands/
│   ├── templates/                   # Shell config templates
│   │   ├── devgita.zsh.tmpl         # FIX: platform-conditional paths
│   │   └── global_config.yaml
│   └── tmux/                        # Cross-platform multiplexer config
│       └── tmux.conf
├── cmd/
│   ├── root.go
│   └── install.go                   # FIX: remove git prerequisite, update Long desc
├── internal/
│   ├── apps/
│   │   ├── aerospace/               # macOS-only (platform-gated)
│   │   ├── alacritty/               # Cross-platform
│   │   ├── brave/                   # Cross-platform
│   │   ├── devgita/                 # REFACTOR: embed extraction instead of git clone
│   │   ├── docker/                  # Cross-platform
│   │   ├── fastfetch/               # Cross-platform
│   │   ├── flameshot/               # Cross-platform
│   │   ├── fonts/                   # Cross-platform
│   │   ├── gimp/                    # Cross-platform
│   │   ├── git/                     # Cross-platform
│   │   ├── i3/                      # NEW: Debian-only tiling WM
│   │   ├── lazydocker/              # Cross-platform
│   │   ├── lazygit/                 # Cross-platform
│   │   ├── mise/                    # Cross-platform
│   │   ├── neovim/                  # Cross-platform
│   │   ├── opencode/                # Cross-platform
│   │   ├── raycast/                 # macOS-only (platform-gated)
│   │   ├── tmux/                    # Cross-platform
│   │   └── ulauncher/               # NEW: Debian-only launcher
│   ├── commands/
│   │   ├── base.go
│   │   ├── debian.go
│   │   ├── factory.go
│   │   ├── macos.go
│   │   ├── mock.go
│   │   └── platform.go
│   ├── config/
│   │   ├── fromContext.go
│   │   └── fromFile.go
│   ├── embedded/                    # NEW: embed.FS declaration + extraction logic
│   │   └── configs.go
│   ├── tooling/
│   │   ├── databases/
│   │   ├── desktop/
│   │   │   └── desktop.go           # FIX: wire desktop apps, platform gates
│   │   ├── languages/
│   │   └── terminal/
│   │       └── terminal.go          # FIX: add SoftConfigure for Mise/OpenCode, remove source call
│   └── testutil/
│       └── testutil.go
├── pkg/
│   ├── constants/
│   │   └── constants.go             # ADD: Ulauncher, I3 constants
│   ├── files/
│   ├── logger/
│   ├── paths/
│   │   └── paths.go                 # ADD: i3 config paths
│   ├── promptui/
│   └── utils/
├── install.sh                       # NEW: download binary + configure PATH
└── README.md                        # UPDATE: install instructions + roadmap
```

**Structure Decision**: Existing Go CLI project structure. No new top-level directories except `internal/embedded/` for the embed declaration. New app modules follow existing pattern under `internal/apps/`. Config templates under `configs/`.

## Complexity Tracking

No constitution violations to justify. All changes follow existing patterns.
