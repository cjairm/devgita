# Development Guides

Detailed guides for implementing features, tests, and commands in devgita. These supplement [CLAUDE.md](../../CLAUDE.md) with practical examples and deep dives.

---

## Guide Index

### [CLI Patterns](cli-patterns.md) — Building commands with Cobra

Everything about designing and building devgita commands:

- Command hierarchy and structure
- Flag handling (string slices, booleans, validation)
- Subcommands and aliases
- Error handling and output patterns
- Context usage
- Testing commands

**When to read:** Before building a new command or modifying existing command structure.

**Referenced by:** CLAUDE.md section 7 (adding new commands)

---

### [Testing Patterns](testing-patterns.md) — Isolation, mocks, and reliability

Complete testing architecture for devgita:

- Three levels of test isolation
- Mock injection patterns
- `testutil` package reference
- Test examples for different app types
- Running tests with coverage
- Common issues and solutions

**When to read:** When writing tests or adding test infrastructure.

**Referenced by:** CLAUDE.md section 6 (testing requirements)

---

### [Error Handling](error-handling.md) — Consistent error messages and exit behavior

Principles and patterns for error handling:

- Return errors instead of panicking
- Add context when propagating errors
- Centralize CLI exit logic
- Rollback on failure
- Verbose/debug logging
- User-friendly messages

**When to read:** When designing error flows or implementing error handling in commands.

**Referenced by:** CLAUDE.md section 6 (error handling)

---

### [Cross-Platform Installation](cross-platform-installation.md) — Strategy pattern and package mappings

Technical deep dive into cross-platform installation:

- Package name translations (Homebrew ↔ apt)
- Installation strategy pattern
- Available strategies (Apt, PPA, script download, git clone, etc.)
- When to add new strategies
- Platform detection and fallback logic

**When to read:** When adding support for a new package, tool, or platform.

**Referenced by:** CLAUDE.md section 11 (architecture patterns)

---

### [Releasing](releasing.md) — Version management and GitHub Actions

Complete release workflow:

- Semantic versioning scheme
- Creating release tags
- GitHub Actions workflow automation
- Verifying releases
- Publishing binaries

**When to read:** When preparing a new release.

---

## Quick Start by Task

### I'm adding a new command

1. Read [CLI Patterns](cli-patterns.md) — overall structure
2. Check [CLAUDE.md](../../CLAUDE.md) section 12 — where to add code
3. Implement using patterns from [CLI Patterns](cli-patterns.md)
4. Test using [Testing Patterns](testing-patterns.md)

### I'm adding a new app/tool installer

1. Read [Cross-Platform Installation](cross-platform-installation.md) — strategy pattern
2. Check [CLAUDE.md](../../CLAUDE.md) section 6 — testing requirement (always use mocks)
3. Implement using patterns from [Testing Patterns](testing-patterns.md)

### I'm fixing error handling

1. Read [Error Handling](error-handling.md) — principles
2. Use `utils.MaybeExitWithError()` in commands
3. Reference [CLAUDE.md](../../CLAUDE.md) section 6 for patterns

### I'm releasing a new version

1. Read [Releasing](releasing.md) — complete workflow
2. Follow semantic versioning scheme
3. Create tag and push to trigger GitHub Actions

---

## How Guides Fit Into Devgita's Documentation

```
CLAUDE.md (Source of Truth)
├── Section 6: Implementation behavior (brief overview)
│   └─> Links to [Error Handling](error-handling.md) for details
├── Section 6: Testing requirements (brief overview)
│   └─> Links to [Testing Patterns](testing-patterns.md) for details
├── Section 7: Command patterns (brief overview)
│   └─> Links to [CLI Patterns](cli-patterns.md) for details
├── Section 11: Architecture patterns
│   └─> Links to [Cross-Platform Installation](cross-platform-installation.md) for details
└── Section 12: Code landmarks (brief overview)

docs/spec.md (What features exist)
└─> Links to ROADMAP.md for future features

ROADMAP.md (What's planned)
└─> Links to docs/plans/cycles/ for active work

docs/decisions/ (How & why we decided)
└─> Individual ADRs for significant technical choices
```

---

## Keeping Guides Current

Guides are the detailed implementation reference. When you:

1. **Change a pattern** — Update the relevant guide immediately
2. **Discover a better practice** — Document it here, then update CLAUDE.md summary
3. **Add a new pattern** — Create a new guide subsection with examples
4. **Find duplication** — Link to the authoritative guide instead of duplicating

Stale guides are worse than no guides — they mislead developers. If a guide describes something that changed, update or remove it in the same PR.

---

## Modern Practices (2026)

These guides follow current industry standards for CLI tool development:

✓ **Type safety** — Go's type system enforces patterns  
✓ **Composition over inheritance** — Interface-based design  
✓ **Dependency injection** — Testable via mocks, not global state  
✓ **Progressive disclosure** — Help shows what matters, examples show workflows  
✓ **Fast feedback** — Validation early, clear error messages  
✓ **Observability** — Verbose/debug modes, structured logging  
✓ **Consistency** — Shared utilities for output, flags, context  
✓ **Zero-downtime** — No breaking changes without deprecation plan

---

## Contributing a Guide

If you're adding new guides:

1. **Use markdown with clear hierarchy** — H2/H3 only, descriptive anchors
2. **Start with context** — Why does this matter?
3. **Show patterns with examples** — Not theory, concrete code
4. **Include a checklist** — Quick reference for "did I do this?"
5. **Link to related guides** — Help developers navigate
6. **Add to this README** — Update the index and quick start

See [TEMPLATE.md](./TEMPLATE.md) for structure.
