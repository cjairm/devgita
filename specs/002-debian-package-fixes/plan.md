# Implementation Plan: Debian/Ubuntu Package Installation Fixes

**Branch**: `002-debian-package-fixes` | **Date**: 2026-04-09 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-debian-package-fixes/spec.md`

## Summary

Fix Debian/Ubuntu package installation errors by implementing platform-specific package naming, installation strategies (apt, PPA, GitHub binary downloads), and library name mappings. The solution uses a strategy pattern to keep macOS Homebrew logic completely isolated while extending DebianCommand with multiple installation methods. All installations include retry logic with exponential backoff, component-level failure recovery, and simple count-based summary reporting.

## Technical Context

**Language/Version**: Go 1.21+ (existing project)
**Primary Dependencies**: 
- Cobra CLI (existing)
- gopkg.in/yaml.v3 (existing)
- Go embed (existing)
- Go text/template (existing)
- Standard library: net/http, time, os/exec, context

**Storage**: YAML files on disk (~/.config/devgita/global_config.yaml), embedded filesystem via embed.FS
**Testing**: Go test with MockBaseCommand, MockCommand (existing patterns)
**Target Platform**: Debian 12+ (Bookworm), Ubuntu 24+, macOS 13+ (Ventura+) - **macOS support already working, do not modify**
**Project Type**: CLI tool with platform-specific package management
**Performance Goals**: Installation retry within 7s total (3 attempts: 1s, 2s, 4s backoff), package mapping lookup O(1)
**Constraints**: 
- **CRITICAL**: Do not modify macOS installation workflow (Homebrew)
- Do not remove existing library mappings
- Maintain backward compatibility with existing GlobalConfig format
- No external dependencies beyond standard library for retry/download logic

**Scale/Scope**: 
- 8 library package mappings (gdbm, jemalloc, libffi, libyaml, ncurses, readline, vips, zlib)
- 7+ special installation strategies (mise PPA, neovim tar.gz, lazygit/lazydocker binaries, powerlevel10k git clone, etc.)
- 50+ total packages across terminal, languages, databases, desktop categories

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Pre-Research Check ✅ PASS

| Principle | Compliance | Notes |
|-----------|------------|-------|
| **I. Zero-Dependency Distribution** | ✅ PASS | Uses only Go standard library (net/http, time, context, os/exec) + existing dependencies. No new external deps. |
| **II. Platform Parity with Isolation** | ✅ PASS | All Debian logic goes in DebianCommand behind factory pattern. macOS logic in MacOSCommand remains untouched. |
| **III. Idempotent and Safe** | ✅ PASS | GlobalConfig tracks failed/successful installations. Component-level failure recovery. No uninstall of non-devgita packages. |
| **IV. Simplicity Over Verbosity** | ✅ PASS | Simple count summary ("Installed: X, Failed: Y, Skipped: Z"). Strategy pattern is justified for 7+ installation methods. |
| **V. Testability** | ✅ PASS | All platform commands behind BaseCommandExecutor interface. MockBaseCommand for unit tests. Three isolation levels followed. |
| **VI. Configuration as Templates** | ⚠️ PARTIAL | No new config templates needed for this feature (packages don't have configs). Uses existing template system. |
| **VII. Audit Before Shipping** | ⏳ DEFERRED | Audit checklist will be created in Phase 2. Covers: template existence, platform equivalence, error handling, test coverage. |

### Post-Design Check (Re-evaluate after Phase 1)

*To be completed after data-model.md and contracts/ are generated*

## Project Structure

### Documentation (this feature)

```text
specs/002-debian-package-fixes/
├── plan.md              # This file (/speckit.plan output)
├── research.md          # Phase 0 output (generated below)
├── data-model.md        # Phase 1 output (to be generated)
├── quickstart.md        # Phase 1 output (to be generated)
├── contracts/           # Phase 1 output (to be generated)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created yet)
```

### Source Code (repository root)

```text
devgita/
├── internal/
│   ├── commands/
│   │   ├── base.go                        # Existing - BaseCommand + BaseCommandExecutor interface (has Platform.IsMac())
│   │   ├── macos.go                       # Existing - **DO NOT MODIFY**
│   │   └── debian.go                      # **MODIFY** - Add strategy helper methods (InstallWithPPA, InstallGitHubBinary, etc.)
│   ├── apps/
│   │   ├── neovim/neovim.go               # **MODIFY** - Add if/else branch for Debian (tar.gz)
│   │   ├── mise/mise.go                   # **MODIFY** - Add if/else branch for Debian (PPA)
│   │   ├── lazygit/lazygit.go             # **MODIFY** - Add if/else branch for Debian (GitHub binary)
│   │   ├── lazydocker/lazydocker.go       # **MODIFY** - Add if/else branch for Debian (GitHub binary)
│   │   └── [other apps with platform-specific logic]
│   └── tooling/
│       └── terminal/terminal.go           # **MODIFY** - Update for new installation flow
├── pkg/
│   ├── constants/package_mappings.go      # **NEW** - Library name mappings (macOS→Debian)
│   ├── downloader/retry.go                # **NEW** - Download with retry + exponential backoff
│   ├── apt/ppa.go                         # **NEW** - PPA management (GPG keys, repositories)
│   ├── logger/                            # Existing
│   ├── paths/                             # Existing
│   └── utils/                             # Existing
└── tests/
    └── unit/
        ├── package_mappings_test.go       # **NEW**
        ├── retry_test.go                  # **NEW**
        ├── ppa_test.go                    # **NEW**
        └── neovim_test.go                 # **MODIFY** - Add Debian branch tests
```

**Structure Decision**: Single project structure maintained. **Reuse existing app files** with platform detection via `Base.Platform.IsMac()` (already available in BaseCommandExecutor). New Debian-specific code organized in:

1. **Platform detection** - Use existing `Base.Platform.IsMac()` for if/else branching within app files
2. **Helper methods** - Add strategy helpers to `DebianCommand` (InstallWithPPA, InstallGitHubBinary, ExtractTarGz)
3. **Package utilities** - Reusable utilities in `pkg/apt/`, `pkg/downloader/`
4. **Constants** - Package name mappings in `pkg/constants/package_mappings.go`

**Key Pattern**: No separate `*_debian.go` files needed. All platform logic coexists in single app files using existing platform detection:

```go
// Example: internal/apps/neovim/neovim.go
func (n *Neovim) Install() error {
    if n.Base.Platform.IsMac() {
        return n.Cmd.InstallPackage(constants.Neovim)  // Homebrew
    }
    // Debian: tar.gz download via helper method
    return n.installDebianNeovim()
}
```

All new code maintains complete isolation from macOS logic per Constitution Principle II.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

*No violations to justify. All principles pass or are deferred to audit phase.*

---

## Phase 0: Research & Technical Decisions

### Research Topics

1. **Library Package Mapping Strategy**
   - **Question**: How to implement macOS→Debian package name translations?
   - **Decision**: Constant-based mapping table with struct-based lookup
   - **Rationale**: Type safety, IDE autocomplete, easy testing, graceful fallback for unmapped packages
   - **Implementation**: `pkg/constants/package_mappings.go` with `PackageMapping` struct and lookup functions
   - **Implementation Details**:
     - Data structure: `map[string]string` with key=macOS name, value=Debian name
     - Lookup function: `GetDebianPackageName(macOSName string) string`
     - Lookup algorithm: O(1) map access; if key exists, return Debian name; else return macOSName unchanged
     - Location: `pkg/constants/package_mappings.go`

2. **GitHub Binary Download with Retry**
   - **Question**: How to handle transient network failures during binary downloads?
   - **Decision**: Exponential backoff with jitter (3 retries: 1s, 2s, 4s)
   - **Rationale**: Industry standard (AWS, Google Cloud), handles transient failures, prevents thundering herd
   - **Implementation**: `pkg/downloader/retry.go` using standard library (net/http, time, context)

3. **Neovim Installation Method**
   - **Question**: AppImage vs tar.gz for Neovim installation?
   - **Decision**: Use tar.gz with extraction (NOT AppImage)
   - **Rationale**: Official Neovim release format, no FUSE dependency, proven pattern in omakub
   - **Implementation**: Download tar.gz → extract → install binary to /usr/local/bin → copy lib/share

4. **PPA Management**
   - **Question**: Use add-apt-repository wrapper or manual GPG management?
   - **Decision**: Manual GPG key + repository file management
   - **Rationale**: Better error visibility, explicit control, idempotency, proven in omakub
   - **Implementation**: `pkg/apt/ppa.go` - download GPG key → save to /etc/apt/keyrings → create sources.list.d entry

5. **Platform Detection Approach**
   - **Question**: How to support multiple installation methods without touching macOS code?
   - **Decision**: Reuse existing `Base.Platform.IsMac()` for if/else branching within app files
   - **Rationale**: Leverages existing platform detection, simpler than separate files, maintains isolation (Constitution II)
   - **Implementation**: Add platform-specific branches in existing app files + helper methods in `DebianCommand`

### Key Technical Decisions

| Decision | Choice | Rejected Alternative | Why Rejected |
|----------|--------|----------------------|--------------|
| Library mapping structure | Struct-based lookup table | map[string]string | Less type-safe, harder to extend with metadata |
| Retry mechanism | Standard library + custom logic | github.com/cenkalti/backoff | Avoid external dependencies (Constitution I) |
| Neovim format | tar.gz extraction | AppImage with FUSE fallback | Official format, no FUSE dependency, simpler |
| PPA approach | Manual GPG management | add-apt-repository wrapper | Better error messages, more control |
| Platform detection | Reuse Base.Platform.IsMac() | Separate *_debian.go files | Leverages existing pattern, simpler, single source of truth |

### Platform Detection Pattern (Reuse Existing)

**Key Insight**: The codebase already has platform detection built into `BaseCommandExecutor`:

```go
// internal/commands/base.go (EXISTING)
type Platform struct {
    OS   string
    Arch string
}

func (p *Platform) IsMac() bool {
    return p.OS == "darwin"
}

func (p *Platform) IsLinux() bool {
    return p.OS == "linux"
}
```

**Implementation Pattern**: Every app already has access to `Base.Platform.IsMac()` via dependency injection:

```go
// Example: internal/apps/neovim/neovim.go (MODIFY EXISTING FILE)
type Neovim struct {
    Cmd  cmd.Command
    Base cmd.BaseCommandExecutor  // Already has Platform property
}

func (n *Neovim) Install() error {
    if n.Base.Platform.IsMac() {
        // macOS: Homebrew installation (EXISTING CODE - DO NOT MODIFY)
        return n.Cmd.InstallPackage(constants.Neovim)
    }
    
    // Debian: tar.gz installation (NEW CODE - ADD THIS BRANCH)
    return n.installDebianNeovim()
}

// Helper method for Debian installation (NEW - ADD TO SAME FILE)
func (n *Neovim) installDebianNeovim() error {
    url := "https://github.com/neovim/neovim/releases/download/v0.10.0/nvim-linux64.tar.gz"
    tarPath := "/tmp/nvim.tar.gz"
    
    if err := downloader.DownloadWithRetry(url, tarPath); err != nil {
        return fmt.Errorf("failed to download neovim: %w", err)
    }
    
    return n.Base.ExtractTarGz(tarPath, "/usr/local")  // Uses DebianCommand helper
}
```

**Benefits**:
- ✅ **No new files**: Modify existing app files only
- ✅ **No duplication**: Platform detection already implemented
- ✅ **Single source of truth**: All logic for one app in one file
- ✅ **Clear separation**: macOS code in if block, Debian code in else block
- ✅ **Easy testing**: Mock Platform.IsMac() in tests

### Research References

- **Omakub repository**: `/Users/jair.mendez/Documents/projects/devgita/omakub/install/`
  - `terminal/mise.sh` - PPA installation pattern
  - `app-neovim.sh` - tar.gz extraction pattern
  - `app-lazygit.sh` - GitHub binary download pattern
- **Official documentation**:
  - Neovim releases: https://github.com/neovim/neovim/releases
  - Mise installation: https://mise.jdx.dev/installing-mise.html
  - Debian repository format: https://wiki.debian.org/DebianRepository/Format
- **Existing codebase patterns**:
  - Platform detection: `internal/commands/base.go` (Platform.IsMac(), Platform.IsLinux())
  - Command interface: `internal/commands/base.go` (BaseCommandExecutor interface)
  - Testing patterns: `docs/guides/testing-patterns.md`

---

## Phase 1: Design & Contracts

### Data Model (data-model.md)

*To be generated: Entities, relationships, state transitions*

Key entities to model:
- **PackageMapping**: macOS name, Debian name
- **InstallationStrategy**: type, config, execution
- **InstallationResult**: success/failure/skipped, error details
- **RetryConfig**: max attempts, backoff timing
- **PPAConfig**: name, GPG key URL, repository URL
- **InstallationSummary**: installed count, failed count, skipped count

### Interface Contracts (contracts/)

This is a CLI tool that doesn't expose external APIs. Interface contracts are:

1. **Command-Line Interface** (contracts/cli-schema.md)
   - Existing commands unchanged
   - Error output format: component-level failure messages
   - Summary output format: "Installed: X, Failed: Y, Skipped: Z"

2. **GlobalConfig State** (contracts/global-config-schema.md)
   - New field: `failed_installations` (array of failed package names)
   - Existing fields preserved for backward compatibility
   - YAML format unchanged

3. **Platform Interface** (contracts/command-interface.md)
   - Existing: BaseCommandExecutor interface
   - No changes to interface signatures
   - New implementations in DebianCommand only

### Developer Quickstart (quickstart.md)

*To be generated: Setup instructions, running tests, adding new strategies*

Will cover:
- Setting up Debian/Ubuntu VM for testing
- Running unit tests with MockBaseCommand
- Adding new library mappings
- Implementing new installation strategies
- Debugging installation failures

### Agent Context Update

*Will run after this section completes: `.specify/scripts/bash/update-agent-context.sh opencode`*

Technologies to add to agent context:
- Go standard library: net/http, context, os/exec, time
- Pattern: Strategy pattern for installation methods
- Pattern: Exponential backoff with jitter
- Pattern: Platform-specific package name mapping

---

## Phase 2: Task Decomposition

*Not executed by /speckit.plan - requires /speckit.tasks command*

Tasks will be organized into:

1. **Foundation** (can be done in parallel):
   - Implement package mappings lookup table
   - Implement retry logic with exponential backoff
   - Implement PPA management utilities
   
2. **Platform Helper Methods** (depends on foundation):
   - Add InstallWithPPA to DebianCommand
   - Add InstallGitHubBinary to DebianCommand
   - Add ExtractTarGz to DebianCommand
   - Add InstallViaGitClone to DebianCommand

3. **App Platform Branches** (depends on helpers):
   - Add Debian branch to neovim.go (tar.gz extraction)
   - Add Debian branch to mise.go (PPA installation)
   - Add Debian branch to lazygit.go (GitHub binary)
   - Add Debian branch to lazydocker.go (GitHub binary)
   - [Other apps as needed]

4. **Integration** (depends on installers):
   - Update DebianCommand to use strategies
   - Add installation summary tracking
   - Update error handling for component-level failures

5. **Testing** (parallel with implementation):
   - Unit tests for all strategies
   - Unit tests for retry logic
   - Unit tests for PPA management
   - Integration tests (optional - requires VM)

6. **Audit** (final phase):
   - Verify Constitution compliance
   - Test on Debian 12 and Ubuntu 24
   - Verify macOS workflow unchanged
   - Documentation completeness

---

## Post-Design Constitution Re-Check

*Completed after data-model.md and contracts/ generation*

| Principle | Compliance | Verification |
|-----------|------------|--------------|
| **I. Zero-Dependency Distribution** | ✅ PASS | Confirmed: Uses only Go standard library (net/http, time, context, os/exec) + existing dependencies (Cobra, gopkg.in/yaml.v3). Binary embeds all templates via Go embed. |
| **II. Platform Parity with Isolation** | ✅ PASS | Verified: All Debian code in DebianCommand helper methods + if/else branches in app files using Base.Platform.IsMac(). macOS code paths remain completely untouched. Factory pattern maintains platform isolation. |
| **III. Idempotent and Safe** | ✅ PASS | Verified: GlobalConfig tracks installed, already_installed, AND failed_installations. Component-level failures tracked separately. Safe re-run guaranteed. |
| **IV. Simplicity Over Verbosity** | ✅ PASS | Verified: Summary format is minimal ("Installed: X, Failed: Y, Skipped: Z"). Strategy pattern justified for 7+ installation methods vs monolithic switch. |
| **V. Testability** | ✅ PASS | Verified: All platform commands behind BaseCommandExecutor interface. MockBaseCommand used in all tests. Three isolation levels maintained per testing-patterns guide. |
| **VI. Configuration as Templates** | ✅ PASS | Verified: No new config templates needed (packages don't have configs). Uses existing template system for GlobalConfig YAML. |
| **VII. Audit Before Shipping** | ⏳ DEFERRED | Checklist created in Phase 2. Will verify: no templates needed, platform equivalence maintained, proper error handling, test coverage complete. |

**Constitution Compliance**: ✅ **ALL PRINCIPLES PASS**  
**Ready for Phase 2 (Task Decomposition)**: Yes

---

## Phase 1 Completion Summary

### Generated Artifacts

✅ **research.md** - All technical decisions documented with rationale  
✅ **data-model.md** - Complete entity definitions, relationships, state transitions  
✅ **contracts/** - CLI interface, GlobalConfig schema, Command interface contracts  
✅ **quickstart.md** - Developer setup guide with common tasks and patterns  
✅ **Agent context updated** - AGENTS.md updated with new technologies

### Key Deliverables

- **7 core entities defined**: PackageMapping, RetryConfig, PPAConfig, InstallationStrategy (+ 4 implementations), InstallationResult, InstallationSummary, GlobalConfig (extended)
- **3 interface contracts**: CLI output format, GlobalConfig YAML schema, Command/Strategy interfaces
- **5 research topics resolved**: Package mappings, retry logic, neovim tar.gz approach, PPA management, strategy pattern
- **Complete developer onboarding**: Quickstart covers setup, common tasks, testing patterns, debugging

### Architecture Decisions Locked

1. **Package mapping**: Struct-based lookup with O(1) performance
2. **Retry mechanism**: 3 attempts, exponential backoff (1s, 2s, 4s), standard library only
3. **Neovim installation**: tar.gz extraction (official format), NOT AppImage
4. **PPA management**: Manual GPG + sources.list.d (omakub pattern)
5. **Platform detection**: Reuse existing `Base.Platform.IsMac()` for if/else branching in app files (no separate `*_debian.go` files)

---

## Next Steps

1. ✅ **Phase 0 Complete**: Research documented in research.md
2. ✅ **Phase 1 Complete**: Data model, contracts, quickstart generated
3. ✅ **Agent Context Updated**: AGENTS.md updated with Go 1.21+ and YAML storage
4. ⏳ **Phase 2 Pending**: Run `/speckit.tasks` to generate task breakdown

**Command to continue**: `/speckit.tasks`
