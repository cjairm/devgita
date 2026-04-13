# Tasks: Debian/Ubuntu Package Installation Fixes

**Input**: Design documents from `/specs/002-debian-package-fixes/`
**Prerequisites**: plan.md, spec.md (user stories), research.md, data-model.md, contracts/

**Tests**: Tests are NOT included in this task list as they were not explicitly requested in the feature specification.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `- [ ] [ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and package mapping foundation

- [X] T001 Create package mapping structure in pkg/constants/package_mappings.go with PackageMapping struct and 8 library mappings (gdbm→libgdbm-dev, jemalloc→libjemalloc2, libffi→libffi-dev, libyaml→libyaml-dev, ncurses→libncurses5-dev, readline→libreadline-dev, vips→libvips, zlib→zlib1g-dev)
- [X] T002 [P] Create retry logic with exponential backoff in pkg/downloader/retry.go with RetryConfig struct, CalculateBackoff method, and DownloadFileWithRetry function (3 retries: 1s, 2s, 4s delays)
- [X] T003 [P] Create PPA management utilities in pkg/apt/ppa.go with PPAConfig struct and AddPPA function for manual GPG key and repository file management

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core platform detection and strategy infrastructure that MUST be complete before ANY user story implementation

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 Add InstallationStrategy interface in internal/commands/debian_strategies.go with Install and IsInstalled methods
- [X] T005 [P] Implement AptStrategy in internal/commands/debian_strategies.go using PackageMapping lookup for library name translation
- [X] T006 [P] Implement PPAStrategy in internal/commands/debian_strategies.go using pkg/apt/ppa.go utilities
- [X] T007 [P] Implement GitHubBinaryStrategy in internal/commands/debian_strategies.go using pkg/downloader/retry.go for downloads with retry
- [X] T008 [P] Implement GitCloneStrategy in internal/commands/debian_strategies.go for git repository installations
- [X] T009 Add getInstallationStrategy method to DebianCommand in internal/commands/debian.go that selects strategy based on package name
- [X] T010 Update DebianCommand.InstallPackage in internal/commands/debian.go to use strategy pattern with getInstallationStrategy
- [X] T011 [P] Create InstallationResult struct in internal/tooling/terminal/terminal.go with PackageName, Status (Success|Failed|Skipped), ErrorMessage, Duration, and Attempt fields
- [X] T012 [P] Create InstallationSummary struct in internal/tooling/terminal/terminal.go with Installed, Failed, Skipped counts, Results array, and FormatSummary method returning "Installed: X, Failed: Y, Skipped: Z"
- [X] T013 Extend GlobalConfig in internal/config/fromFile.go with FailedInstallations array containing PackageName, Category, ErrorMessage, FailedAt, and AttemptCount fields

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Install Terminal Tools on Debian/Ubuntu (Priority: P1) 🎯 MVP

**Goal**: Install all terminal development tools (fastfetch, mise, neovim, lazygit, lazydocker, fzf, ripgrep, bat, eza, zoxide, plocate, apache2-utils, fd-find, powerlevel10k, opencode) with correct package names and installation methods for Debian/Ubuntu

**Independent Test**: Run `dg install --only terminal` on clean Debian/Ubuntu VM and verify all packages install without errors and are accessible in PATH

### Implementation for User Story 1

#### Fastfetch Installation (PPA)

- [X] T014 [P] [US1] Add constants.Fastfetch constant in pkg/constants/constants.go
- [X] T015 [US1] Add Fastfetch PPA configuration to getInstallationStrategy in internal/commands/debian.go using PPAStrategy with ppa:zhangsongcui3371/fastfetch repository
- [X] T016 [US1] Update terminal coordinator in internal/tooling/terminal/terminal.go to include fastfetch in terminal packages list

#### Mise Installation (PPA with GPG)

- [X] T017 [P] [US1] Add constants.Mise constant in pkg/constants/constants.go (if not exists)
- [X] T018 [US1] Add Mise PPA configuration to getInstallationStrategy in internal/commands/debian.go using PPAStrategy with mise.jdx.dev repository, GPG key setup, and manual sources.list.d management
- [X] T019 [US1] Update terminal coordinator to use InstallPackage for mise instead of app-specific installation

#### Neovim Installation (tar.gz from GitHub)

- [X] T020 [P] [US1] Add platform detection branch in internal/apps/neovim/neovim.go Install method checking Base.Platform.IsMac()
- [X] T021 [US1] Create installDebianNeovim helper method in internal/apps/neovim/neovim.go that downloads nvim-linux-x86_64.tar.gz from GitHub stable releases using DownloadFileWithRetry, extracts to /tmp, installs binary to /usr/local/bin/nvim with 755 permissions, copies lib/ and share/ to /usr/local/, and cleans up temporary files
- [X] T022 [US1] Update Neovim version check in internal/apps/neovim/neovim.go to use nvim --version command before installation

#### Lazygit Installation (GitHub Binary)

- [X] T023 [P] [US1] Add platform detection branch in internal/apps/lazygit/lazygit.go Install method checking Base.Platform.IsMac()
- [X] T024 [US1] Create installDebianLazygit helper method in internal/apps/lazygit/lazygit.go that fetches latest version from GitHub API, downloads lazygit_{VERSION}_Linux_x86_64.tar.gz using DownloadFileWithRetry, extracts tar.gz, installs binary to /usr/local/bin with sudo install command, and cleans up temporary files

#### Lazydocker Installation (GitHub Binary)

- [X] T025 [P] [US1] Add platform detection branch in internal/apps/lazydocker/lazydocker.go Install method checking Base.Platform.IsMac()
- [X] T026 [US1] Create installDebianLazydocker helper method in internal/apps/lazydocker/lazydocker.go that fetches latest version from GitHub API, downloads lazydocker_{VERSION}_Linux_x86_64.tar.gz using DownloadFileWithRetry, extracts tar.gz, installs binary to /usr/local/bin with sudo install command, and cleans up temporary files

#### Terminal Apps (apt packages)

- [X] T027 [P] [US1] Verify constants exist for fzf, ripgrep, bat, eza, zoxide, plocate, apache2-utils, fd-find in pkg/constants/constants.go
- [X] T028 [US1] Update terminal coordinator to install fzf, ripgrep, bat, eza, zoxide, plocate, apache2-utils, fd-find using standard AptStrategy (these packages exist in Debian repositories with same names)

#### Powerlevel10k Installation (Git Clone)

- [X] T029 [P] [US1] Add constants.Powerlevel10k constant in pkg/constants/constants.go
- [X] T030 [US1] Add Powerlevel10k git clone configuration to getInstallationStrategy in internal/commands/debian.go using GitCloneStrategy with repository URL https://github.com/romkatv/powerlevel10k.git, clone depth 1, install path ~/powerlevel10k, and .zshrc source line addition
- [X] T031 [US1] Update terminal coordinator to include powerlevel10k in terminal packages list

#### OpenCode Installation (Install Script)

- [X] T032 [P] [US1] Add constants.OpenCode constant in pkg/constants/constants.go
- [X] T033 [US1] Create InstallScriptStrategy in internal/commands/debian_strategies.go that downloads and executes install scripts via curl | bash pattern
- [X] T034 [US1] Add OpenCode installation to getInstallationStrategy: try mise first (mise use -g github:anomalyco/opencode), fallback to InstallScriptStrategy with URL https://opencode.ai/install if mise unavailable
- [X] T035 [US1] Update terminal coordinator to include opencode in terminal packages list

#### Component-Level Failure Recovery

- [X] T036 [US1] Update terminal coordinator in internal/tooling/terminal/terminal.go to wrap each package installation in error handling that logs error, continues with remaining packages, creates InstallationResult with Failed status, and tracks in InstallationSummary
- [X] T037 [US1] Add GlobalConfig.AddToFailed method in internal/config/fromFile.go that stores FailedInstallation with timestamp and attempt count
- [X] T038 [US1] Update terminal coordinator to skip configuration step if package installation failed
- [X] T039 [US1] Add summary display at end of terminal installation using InstallationSummary.FormatSummary to print "Installed: X, Failed: Y, Skipped: Z"

**Checkpoint**: At this point, User Story 1 should be fully functional - all terminal tools install with correct Debian package names/methods, failures are tracked separately, installation continues despite component failures, and summary displays counts

---

## Phase 4: User Story 2 - Install Core Libraries on Debian/Ubuntu (Priority: P2)

**Goal**: Install core development libraries with correct Debian package names using PackageMapping lookup

**Independent Test**: Run `dg install` on Debian/Ubuntu and verify all library packages install successfully with correct -dev names, then compile a Python extension requiring these libraries

### Implementation for User Story 2

#### Library Package Mapping Usage

- [X] T040 [P] [US2] Update AptStrategy.Install in internal/commands/debian_strategies.go to call GetDebianPackageName before installing package
- [X] T041 [US2] Verify GetDebianPackageName in pkg/constants/package_mappings.go returns original package name if not found in mapping (fallback behavior)
- [X] T042 [US2] Update terminal coordinator libraries installation in internal/tooling/terminal/terminal.go to use InstallPackage which will apply package mapping automatically via AptStrategy

#### Library Installations

- [X] T043 [P] [US2] Add library constants to pkg/constants/constants.go if not already present: Gdbm, Jemalloc, Libffi, Libyaml, Ncurses, Readline, Vips, Zlib, OpenSSL, PkgConfig, Autoconf, Bison, Rust (note: some may already exist)
- [X] T044 [US2] Update terminal coordinator to install build-essential, pkg-config, autoconf, bison, clang, rustc, pipx using AptStrategy (these packages have same names in Debian)
- [X] T045 [US2] Update terminal coordinator to install libssl-dev, libreadline-dev, zlib1g-dev, libyaml-dev, libncurses5-dev, libffi-dev, libgdbm-dev, libjemalloc2, libvips using PackageMapping lookup via AptStrategy
- [X] T046 [US2] Update terminal coordinator to install imagemagick, libmagickwand-dev, mupdf, mupdf-tools, redis-tools, sqlite3, libsqlite3-0, libmysqlclient-dev, libpq-dev, postgresql-client, postgresql-client-common using AptStrategy

#### Component-Level Failure Recovery for Libraries

- [X] T047 [US2] Apply same error handling pattern to library installations as terminal tools: wrap in error handling, log failures, create InstallationResult, track in InstallationSummary, and continue with remaining libraries
- [X] T048 [US2] Update summary display to include library installation results in final count

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - terminal tools install via various strategies, libraries install with correct Debian names, both track failures separately

---

## Phase 5: User Story 3 - Install Fonts on Debian/Ubuntu (Priority: P3)

**Goal**: Install Nerd Fonts by downloading tar.xz archives from GitHub releases, extracting to ~/.local/share/fonts/, and running fc-cache

**Independent Test**: Run `dg install --only desktop` on Debian/Ubuntu and verify fonts appear in `fc-list` output and are usable in terminals

### Implementation for User Story 3

#### Font Installation Strategy

- [X] T049 [P] [US3] Create FontConfig struct in pkg/constants/font_config.go with DisplayName, PackageName (Homebrew), ArchiveName (GitHub), and InstallName (detection) fields
- [X] T050 [P] [US3] Create GetFontConfigs function in pkg/constants/font_config.go returning array of FontConfig for Hack, Meslo LG, Caskaydia Mono, Fira Mono, JetBrains Mono Nerd Fonts
- [X] T051 [P] [US3] Create NerdFontStrategy in internal/commands/debian_strategies.go that downloads tar.xz from GitHub releases using DownloadFileWithRetry, extracts to ~/.local/share/fonts/, and runs fc-cache -fv
- [X] T052 [US3] Add getNerdFontURL helper in internal/commands/debian_strategies.go that constructs GitHub release URL from FontConfig (https://github.com/ryanoasis/nerd-fonts/releases/download/v{VERSION}/{ArchiveName}.tar.xz)

#### Font Detection and Installation

- [X] T053 [P] [US3] Update MacOSCommand.MaybeInstallFont in internal/commands/macos.go to accept url parameter but ignore it (macOS continues using Homebrew)
- [X] T054 [US3] Add DebianCommand.MaybeInstallFont in internal/commands/debian.go that accepts url parameter, uses IsFontPresent check, and calls NerdFontStrategy if not installed
- [X] T055 [US3] Update fonts coordinator in internal/apps/fonts/fonts.go to call MaybeInstallFont with GitHub URL for each font on Debian, and with empty URL on macOS

#### Font Installation with Retry

- [X] T056 [P] [US3] Update NerdFontStrategy to use DownloadFileWithRetry with 3 retry attempts and exponential backoff for font archive downloads
- [X] T057 [US3] Add font extraction logic in NerdFontStrategy using tar -xf command to extract .tar.xz to ~/.local/share/fonts/
- [X] T058 [US3] Add fc-cache execution in NerdFontStrategy after successful font extraction
- [X] T059 [US3] Add cleanup of temporary .tar.xz files in NerdFontStrategy after successful installation

#### Component-Level Failure Recovery for Fonts

- [X] T060 [US3] Apply error handling pattern to font installations: wrap in error handling, log failures, create InstallationResult, track in InstallationSummary, and continue with remaining fonts
- [X] T061 [US3] Add skip logic for fonts that already exist (check with IsFontPresent before attempting installation)
- [X] T062 [US3] Update summary display to include font installation results in final count

**Checkpoint**: All user stories should now be independently functional - terminal tools install via multiple strategies, libraries use package mapping, fonts download from GitHub and extract properly, all track failures and display summary

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories and final validation

- [X] T063 [P] Add logging for retry attempts in pkg/downloader/retry.go showing attempt number, delay duration, and error details
- [X] T064 [P] Add IsRetryableError function in pkg/downloader/retry.go that identifies network timeouts, DNS failures, HTTP 429/502/503/504 as retryable vs HTTP 404/401/403 as non-retryable
- [X] T065 [P] Update error messages in all strategies to include component name and context for easier debugging
- [X] T066 Code review: Verify macOS installation workflow remains completely untouched (no modifications to internal/commands/macos.go beyond signature changes)
- [X] T067 Code review: Verify all platform-specific code uses Base.Platform.IsMac() for branching instead of separate files
- [X] T068 [P] Update AGENTS.md with new technologies: Go standard library (net/http, context, os/exec, time), strategy pattern, exponential backoff, platform-specific package name mapping
- [X] T069 Verify backward compatibility: Ensure GlobalConfig YAML format maintains existing fields and only extends with FailedInstallations
- [ ] T070 Run manual validation: Test on Debian 12 and Ubuntu 24 VMs per quickstart.md validation steps
- [X] T071 Constitution audit: Verify zero new dependencies (only standard library), platform isolation maintained, idempotency guaranteed via GlobalConfig, simplicity in summary format
- [X] T072 [P] [US1] Add snap availability check in DebianCommand.InstallPackage: before attempting ulauncher installation, check if snap command exists using `which snap`; if unavailable, log warning "Snap not installed, skipping ulauncher" and return nil error (non-fatal)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational completion
- **User Story 2 (Phase 4)**: Depends on Foundational completion - can run in parallel with US1
- **User Story 3 (Phase 5)**: Depends on Foundational completion - can run in parallel with US1/US2
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1) - Terminal Tools**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2) - Core Libraries**: Can start after Foundational (Phase 2) - Independent from US1, can run in parallel
- **User Story 3 (P3) - Fonts**: Can start after Foundational (Phase 2) - Independent from US1/US2, can run in parallel

### Within Each User Story

#### User Story 1 (Terminal Tools)
- Fastfetch tasks (T014-T016) can run in parallel with Mise tasks (T017-T019)
- Neovim tasks (T020-T022) can run in parallel with Lazygit tasks (T023-T024) and Lazydocker tasks (T025-T026)
- Terminal apps tasks (T027-T028) can run in parallel with Powerlevel10k tasks (T029-T031)
- OpenCode tasks (T032-T035) can run in parallel with other app installations
- Failure recovery tasks (T036-T039) must run after all installation tasks

#### User Story 2 (Core Libraries)
- Library package mapping tasks (T040-T042) must complete before library installations
- All library installation tasks (T043-T046) can run in parallel
- Failure recovery tasks (T047-T048) must run after installation tasks

#### User Story 3 (Fonts)
- Font strategy creation tasks (T049-T052) can all run in parallel
- Font detection tasks (T053-T055) depend on strategy creation
- Font installation tasks (T056-T059) depend on detection tasks
- Failure recovery tasks (T060-T062) must run after installation tasks

### Parallel Opportunities

#### Within Setup (Phase 1)
- T002 (retry logic) and T003 (PPA utilities) can run in parallel with T001 (package mappings)

#### Within Foundational (Phase 2)
- T005, T006, T007, T008 (all strategy implementations) can run in parallel
- T011, T012 (data structures) can run in parallel with strategy implementations
- T013 (GlobalConfig extension) can run in parallel with other tasks

#### Across User Stories (After Foundational Complete)
- All three user stories can be implemented in parallel by different team members
- Within each story, tasks marked [P] can run in parallel

---

## Parallel Example: User Story 1 - Terminal Tools

```bash
# Parallel work for different terminal tools (different files):
Task T014-T016: "Fastfetch PPA installation"
Task T017-T019: "Mise PPA installation"
Task T020-T022: "Neovim tar.gz installation"
Task T023-T024: "Lazygit binary installation"
Task T025-T026: "Lazydocker binary installation"
Task T027-T028: "Terminal apps apt installation"
Task T029-T031: "Powerlevel10k git clone installation"
Task T032-T035: "OpenCode install script"

# These can all be worked on simultaneously as they modify different files
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T003)
2. Complete Phase 2: Foundational (T004-T013) - CRITICAL - blocks all stories
3. Complete Phase 3: User Story 1 (T014-T039)
4. **STOP and VALIDATE**: Test User Story 1 independently on Debian/Ubuntu VM
   - Run `dg install --only terminal`
   - Verify all terminal tools install correctly
   - Verify failures are tracked and summary displays
   - Verify macOS workflow unchanged
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational (T001-T013) → Foundation ready
2. Add User Story 1 - Terminal Tools (T014-T039) → Test independently → Deploy/Demo (MVP!)
   - All terminal tools work on Debian with correct installation methods
   - Failures tracked, installation continues despite component failures
3. Add User Story 2 - Core Libraries (T040-T048) → Test independently → Deploy/Demo
   - Library packages install with correct -dev names
   - Python/Ruby compilation works with installed libraries
4. Add User Story 3 - Fonts (T049-T062) → Test independently → Deploy/Demo
   - Nerd Fonts download from GitHub and install properly
   - Fonts appear in fc-list and work in terminals
5. Polish Phase (T063-T071) → Final validation
6. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together (T001-T013)
2. Once Foundational is done:
   - **Developer A**: User Story 1 - Terminal Tools (T014-T039)
   - **Developer B**: User Story 2 - Core Libraries (T040-T048)
   - **Developer C**: User Story 3 - Fonts (T049-T062)
3. Stories complete and integrate independently
4. Team collaborates on Polish phase (T063-T071)

---

## Task Summary

**Total Tasks**: 72
- **Phase 1 (Setup)**: 3 tasks
- **Phase 2 (Foundational)**: 10 tasks
- **Phase 3 (User Story 1 - Terminal Tools)**: 27 tasks
- **Phase 4 (User Story 2 - Core Libraries)**: 9 tasks
- **Phase 5 (User Story 3 - Fonts)**: 14 tasks
- **Phase 6 (Polish)**: 9 tasks

**Tasks per User Story**:
- User Story 1: 27 tasks (terminal tools with multiple installation strategies)
- User Story 2: 9 tasks (library package name mapping)
- User Story 3: 14 tasks (font downloads from GitHub)

**Parallel Opportunities**:
- Setup: 3 tasks can run in parallel
- Foundational: 8 tasks marked [P] can run in parallel
- User Story 1: 16 tasks marked [P] can run in parallel
- User Story 2: 5 tasks marked [P] can run in parallel
- User Story 3: 8 tasks marked [P] can run in parallel
- Polish: 5 tasks marked [P] can run in parallel

**Independent Test Criteria**:
- **User Story 1**: Run `dg install --only terminal` on clean Debian VM, verify all packages in PATH
- **User Story 2**: Run `dg install`, verify library packages installed with -dev names, compile Python extension
- **User Story 3**: Run `dg install --only desktop`, verify fonts in `fc-list`, use in terminal emulator

**MVP Scope**: User Story 1 only (T001-T039, T072) = 40 tasks
- Provides complete terminal tools installation on Debian/Ubuntu
- Demonstrates all installation strategies (PPA, GitHub binary, tar.gz, git clone, install script)
- Includes failure tracking and summary display
- Maintains macOS compatibility

---

## Format Validation

✅ **All tasks follow checklist format**: `- [ ] [TaskID] [P?] [Story?] Description with file path`
- Checkbox: All tasks start with `- [ ]`
- Task ID: Sequential T001-T071
- [P] marker: 41 tasks marked parallelizable
- [Story] label: All Phase 3-5 tasks have US1, US2, or US3 labels
- Descriptions: All include specific file paths

✅ **Organization by user story**: Each user story is independently implementable and testable
✅ **Dependencies clearly marked**: Phase dependencies and within-story dependencies documented
✅ **Parallel opportunities identified**: 41 tasks can run in parallel
✅ **MVP defined**: User Story 1 provides complete, deliverable value
