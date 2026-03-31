# Tasks: Binary Distribution with Embedded Configs

**Input**: Design documents from `/specs/001-binary-dist-audit/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/

**Tests**: Not explicitly requested in spec. Test tasks omitted. Unit tests should be written alongside implementation per existing testing-patterns guide.

**Organization**: Tasks grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1-US6)
- Include exact file paths in descriptions

## Path Conventions

Existing Go CLI project structure. Paths relative to repository root.

---

## Phase 1: Setup

**Purpose**: Foundation for embed infrastructure and new constants

- [X] T001 Add `Ulauncher` and `I3` constants to `pkg/constants/constants.go`
- [X] T002 [P] Add i3 config paths (`Paths.App.Configs.I3`, `Paths.Config.I3`) to `pkg/paths/paths.go`
- [X] T003 [P] Create `embedded.go` at repository root with `//go:embed all:configs` directive and `ExtractEmbeddedConfigs(destDir string) error` function using `fs.WalkDir` pattern from research.md R1

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Fix existing bugs and create config templates that ALL user stories depend on

**CRITICAL**: No user story work can begin until this phase is complete.

- [X] T004 Add `IsMac bool` field with `yaml:"is_mac"` tag to `ShellFeatures` struct in `internal/config/fromFile.go`
- [X] T005 [P] Create `configs/git/.gitconfig` template with sensible defaults (user name/email placeholders, default branch = main, pull rebase, color auto, common aliases)
- [X] T006 [P] Create `configs/i3/config` with development defaults (Mod=Super, terminal=alacritty, gaps=10px, vim navigation $Mod+h/j/k/l, workspaces $Mod+1-9, reload $Mod+Shift+r)
- [X] T007 Update `configs/templates/devgita.zsh.tmpl` to use `{{if .IsMac}}` conditionals for zsh-autosuggestions, zsh-syntax-highlighting, and powerlevel10k paths (macOS: `$(brew --prefix)/share/...`, Debian: `/usr/share/...`)
- [X] T008 Remove broken `source` exec call (lines 67-74) from `internal/tooling/terminal/terminal.go` in `InstallAndConfigure()`
- [X] T009 [P] Add `SoftConfigure()` call for OpenCode after `SoftInstall()` in `InstallTerminalApps()` at `internal/tooling/terminal/terminal.go` (move from install-only to install+configure loop)
- [X] T010 [P] Add `SoftConfigure()` call for Mise after `SoftInstall()` in `InstallTerminalApps()` at `internal/tooling/terminal/terminal.go`

**Checkpoint**: Foundation ready — all config templates exist, bugs fixed, user story implementation can begin.

---

## Phase 3: User Story 2 - Embedded Config Extraction (Priority: P2)

**Goal**: Binary embeds configs and extracts them to `~/.config/devgita/configs/` replacing git clone.

**Independent Test**: Build binary, run `dg install`, verify configs extracted to correct paths.

**Note**: US2 before US1 because US1 (install.sh) depends on having a working binary with embedded configs.

### Implementation for User Story 2

- [X] T011 [US2] Refactor `internal/apps/devgita/devgita.go`: remove `Git git.Git` field and `git` import from struct, replace `git.Clone()` in `Install()` with `embedded.ExtractEmbeddedConfigs()`, update `ForceInstall()` to work with embed model (Uninstall extracts dir, then re-extract).
- [X] T012 [US2] Refactor `SoftInstall()` in `internal/apps/devgita/devgita.go` to check whether `configs/` subdirectory exists and is non-empty in `~/.config/devgita/` (not whether git repo exists)
- [X] T013 [US2] Update `ForceConfigure()` in `internal/apps/devgita/devgita.go` to set `gc.Shell.IsMac = runtime.GOOS == "darwin"` before saving GlobalConfig
- [X] T014 [US2] Update `Uninstall()` in `internal/apps/devgita/devgita.go` to clean up extracted configs directory
- [X] T015 [US2] Update `cmd/install.go`: remove `git.New()` and `g.SoftInstall()` from prerequisites (lines 85-86), update Long description to remove "Clones the devgita repository"
- [X] T016 [US2] Update `internal/apps/devgita/devgita_test.go` to test embed extraction instead of git clone, using `testutil.SetupCompleteTest()`

**Checkpoint**: `dg install` extracts embedded configs instead of cloning. Binary is self-contained.

---

## Phase 4: User Story 3 - Platform-Correct App Installation (Priority: P3)

**Goal**: macOS-only apps gated, Debian equivalents (Ulauncher, i3) created and wired in.

**Independent Test**: Run `dg install --only desktop` on both platforms, verify correct apps installed per platform.

### Implementation for User Story 3

- [X] T017 [P] [US3] Create `internal/apps/ulauncher/ulauncher.go` following standard app interface: `New()`, `Install()`, `SoftInstall()`, `ForceInstall()`, `Uninstall()` (not supported), `ForceConfigure()` (nil), `SoftConfigure()` (nil), `ExecuteCommand()` (nil), `Update()` (not implemented). Use `InstallDesktopApp()` for installation.
- [X] T018 [P] [US3] Create `internal/apps/ulauncher/ulauncher_test.go` with tests for `New()`, `Install()`, `SoftInstall()`, `ForceInstall()`, `Uninstall()` using `testutil.NewMockApp()`
- [X] T019 [P] [US3] Create `internal/apps/i3/i3.go` following standard app interface. `SoftConfigure()` checks for `config` marker file in i3 config dir. `ForceConfigure()` copies `configs/i3/` to local i3 config dir. Use `InstallPackage()`.
- [X] T020 [P] [US3] Create `internal/apps/i3/i3_test.go` with tests for all methods including `ForceConfigure()` and `SoftConfigure()` using `testutil.SetupIsolatedPaths()`
- [X] T021 [US3] Refactor `InstallAndConfigure()` in `internal/tooling/desktop/desktop.go`: wrap `InstallAerospace()` in `d.Base.Platform.IsMac()` check; add `d.InstallI3()` in `else` branch for Linux
- [X] T022 [US3] Add `InstallI3()` method to `Desktop` struct in `internal/tooling/desktop/desktop.go`: calls `i3.New().SoftInstall()` then `SoftConfigure()`
- [X] T023 [US3] Refactor `InstallDesktopAppsWithoutConfiguration()` in `internal/tooling/desktop/desktop.go`: wrap Raycast in `d.Base.Platform.IsMac()` check; add Ulauncher in `else` branch for Linux
- [X] T024 [US3] Wire `InstallDesktopAppsWithoutConfiguration()` call into `InstallAndConfigure()` in `internal/tooling/desktop/desktop.go` (currently dead code — never called)
- [X] T025 [US3] Add `git` to terminal tools in `InstallTerminalApps()` at `internal/tooling/terminal/terminal.go` (moved from `cmd/install.go` prerequisite per FR-020)

**Checkpoint**: Platform-specific apps correctly gated. Debian gets Ulauncher + i3, macOS gets Raycast + Aerospace.

---

## Phase 5: User Story 4 - Audit Gaps Fixed (Priority: P4)

**Goal**: All remaining audit findings resolved. Shell template platform-aware. All configure methods called.

**Independent Test**: Run `dg install --only terminal` and verify OpenCode config, Mise shell integration, and git config all applied.

### Implementation for User Story 4

- [X] T026 [US4] Verify `configs/git/.gitconfig` created in T005 is referenced correctly by `git.ForceConfigure()` — ensure `paths.Paths.App.Configs.Git` points to the right directory and `files.CopyDir()` works with the new file
- [X] T027 [US4] Create `docs/apps/ulauncher.md` documenting the Ulauncher module following existing app doc pattern (see `docs/apps/raycast.md` as template)
- [X] T028 [P] [US4] Create `docs/apps/i3.md` documenting the i3 module following existing app doc pattern (see `docs/apps/aerospace.md` as template)
- [X] T029 [US4] Run `go test ./...` and fix any test failures from refactored code
- [X] T030 [US4] Run `go vet ./...` and `go fmt ./...` to ensure code quality

**Checkpoint**: All audit gaps resolved. Every app with configure methods has templates. All coordinators call configure.

---

## Phase 6: User Story 1 - Zero-Dependency Install (Priority: P1)

**Goal**: install.sh downloads binary, configures PATH. Full end-to-end zero-dependency install.

**Independent Test**: Run install.sh on clean system, then `dg install`.

**Note**: US1 is implemented after US2-US4 because install.sh requires a working binary with all fixes applied.

### Implementation for User Story 1

- [X] T031 [US1] Create `install.sh` at repository root: detect OS (`uname -s`) and arch (`uname -m`), map to binary name, fetch latest release tag from GitHub API, download binary to `~/.local/bin/devgita`, chmod +x, detect shell config, add PATH if not present, verify with `devgita --help`
- [X] T032 [US1] Add `--local <path>` flag support to `install.sh`: if flag present, copy local file to `~/.local/bin/devgita` instead of downloading; skip GitHub API call; proceed with PATH and verification steps
- [X] T033 [US1] Add unsupported OS/arch error handling to `install.sh`: reject non-Darwin non-Linux OS, reject non-amd64 non-arm64 arch, clear error messages with exit code 1
- [X] T034 [US1] Add idempotent PATH handling to `install.sh`: check if `.local/bin` already in shell config before appending, prevent duplicate entries

**Checkpoint**: Users can install devgita with a single curl command on macOS and Debian.

---

## Phase 7: User Story 5 - Local Binary Build and Release (Priority: P5)

**Goal**: Documented build commands, local testing workflow, release process.

**Independent Test**: Build all three binaries, test locally with `--local` flag.

### Implementation for User Story 5

- [X] T035 [US5] Add `Makefile` or `build.sh` at repository root with three cross-compilation commands from research.md R2 (`GOOS=darwin GOARCH=arm64`, `GOOS=darwin GOARCH=amd64`, `GOOS=linux GOARCH=amd64`) producing `devgita-{os}-{arch}` binaries
- [X] T036 [US5] Fix architecture label swap in existing README.md (lines 98-101): "M chips" should be arm64, "Intel chips" should be amd64

**Checkpoint**: Maintainers can build, test locally, and release.

---

## Phase 8: User Story 6 - README Roadmap (Priority: P6)

**Goal**: README documents install.sh, current commands, and planned commands.

**Independent Test**: Read README, follow instructions, verify they work.

### Implementation for User Story 6

- [X] T037 [US6] Update `README.md` installation section: replace current build-from-source instructions with `curl -fsSL ... | bash` one-liner, document `--local` flag for maintainers
- [X] T038 [US6] Add commands section to `README.md`: document `dg install` with `--only`/`--skip` flags as available; list `dg update`, `dg configure`, `dg uninstall`, `dg list`, `dg change` as planned/roadmap

**Checkpoint**: README reflects current state and future roadmap.

---

## Phase 9: Polish & Cross-Cutting Concerns

**Purpose**: Final validation across all user stories

- [X] T039 Run `go build -o devgita main.go` and verify binary includes embedded configs (check binary size > configs directory size)
- [X] T040 Run `go test ./...` full test suite — all tests must pass
- [X] T041 Run `go vet ./...` and `go fmt ./...` — zero warnings
- [X] T042 Build all three platform binaries and verify each produces a working binary
- [X] T043 Test `install.sh --local ./devgita-darwin-arm64` on macOS — verify full install flow end to end
- [X] T044 Verify `~/.config/devgita/configs/` contains all expected directories after `dg install` (aerospace, alacritty, fastfetch, git, i3, neovim, opencode, templates, tmux)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup (T001-T003) — BLOCKS all user stories
- **US2 (Phase 3)**: Depends on Foundational — embed extraction core
- **US3 (Phase 4)**: Depends on Foundational — can run parallel with US2
- **US4 (Phase 5)**: Depends on US2 + US3 — validates all fixes
- **US1 (Phase 6)**: Depends on US2 + US3 + US4 — needs working binary
- **US5 (Phase 7)**: Depends on US1 — documents build process
- **US6 (Phase 8)**: Depends on US1 + US5 — documents everything
- **Polish (Phase 9)**: Depends on all phases

### User Story Dependencies

- **US2 (Embed Extraction)**: Foundational only — can start first
- **US3 (Platform Gating)**: Foundational only — can start parallel with US2
- **US4 (Audit Fixes)**: Depends on US2 + US3 being complete
- **US1 (Install Script)**: Depends on US2 + US3 + US4 (needs working binary)
- **US5 (Build/Release)**: Depends on US1
- **US6 (README)**: Depends on US1 + US5

### Parallel Opportunities

- T001, T002, T003 can all run in parallel (Setup phase)
- T004, T005, T006, T007 can run in parallel (different files)
- T008, T009, T010 modify same file (`terminal.go`) — run sequentially
- T017, T018, T019, T020 can all run in parallel (new files)
- T027, T028 can run in parallel (new doc files)

---

## Implementation Strategy

### MVP First (US2 + US3)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL)
3. Complete Phase 3: US2 - Embed extraction works
4. Complete Phase 4: US3 - Platform gating works
5. **STOP and VALIDATE**: Build binary, run `dg install`, verify configs extracted and platform apps correct
6. Continue with US4 → US1 → US5 → US6

### Incremental Delivery

1. Setup + Foundational → Foundation ready
2. US2 → Binary self-contained with embedded configs
3. US3 → Platform parity achieved
4. US4 → All bugs fixed, full audit pass
5. US1 → install.sh works, zero-dependency install
6. US5 → Build/release documented
7. US6 → README complete
