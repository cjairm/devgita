# Cycle: `ai-tools` category with rtk as first app

**Date:** 2026-07-23
**Status:** Done
**ADR:** [ADR-0004](../../decisions/ADR-0004-ai-tools-install-category.md)
**Origin:** ROADMAP.md "AI & Development Tools" + [task-design.md](../../guides/task-design.md) §Future: rtk

## Goal

Add the planned `ai-tools` install category and ship
[rtk](https://github.com/rtk-ai/rtk) as its first app: binary installed via
Homebrew (macOS) / GitHub release binary (Debian), command-rewriting hook
(`rtk init -g`) left opt-in and documented.

## Scope

- New app module `internal/apps/rtk/` (App interface, lazygit pattern:
  brew on mac, `InstallGitHubBinary` on Debian, injectable version/download
  fns for tests).
- New coordinator `internal/tooling/aitools/` (desktop-style minimal loop,
  override-injectable app list for tests).
- `cmd/install.go`: `ai-tools` in `knownCategories`, `rtk` in
  `appToCoordinator`, `runAITools` + ai-tools filters in `installConfig`,
  `installAITools()` in `run()`.
- Registry: `Meta`, `factories`, uninstall `knownCategories` gain
  rtk / `ai-tools`.
- Docs: `docs/apps/rtk.md`, ROADMAP.md, CLAUDE.md §8 categories,
  task-design.md rtk stance updated from "candidate" to shipped.

Out of scope: automating `rtk init -g` (ADR-0004 decision 2), Ollama and the
rest of the planned ai-tools roster, an embedded rtk config (rtk's defaults
are sane and its config path differs per OS; nothing to ship yet), a new
`dg list` category (rtk tracked as `"package"`, shows under terminal_tools).

## Steps

- [x] Security review of rtk v0.43.0 (install script, release pipeline,
      telemetry) — recorded in ADR-0004
- [x] ADR-0004 written and indexed
- [x] `constants.Rtk`
- [x] `internal/apps/rtk/rtk.go` + `rtk_test.go`
- [x] `internal/tooling/aitools/aitools.go` + `aitools_test.go`
- [x] `cmd/install.go` plumbing + `install_test.go` updates
- [x] Registry `Meta`/`factories`/categories + `registry_test.go` updates
- [x] `docs/apps/rtk.md`
- [x] ROADMAP.md, CLAUDE.md §8, task-design.md currency updates
- [x] `go build ./...`, `go test ./...`, `make lint` green
- [x] Manual verification: `dg install --only rtk` golden path on macOS
- [x] Review fix: mandatory SHA-256 verification in `InstallGitHubBinary`
      (rtk + lazygit + lazydocker callers), with mismatch/refusal tests
- [x] Follow-up: rtk opt-in moved onto the AI coders
      (`dg configure claude|opencode --force --only=rtk`); claude's
      settings.json rendered from a template honoring the persisted
      `integrations.rtk_claude_hook` flag so the hook survives re-renders
