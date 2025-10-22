# Why: Cobra framework usage in cmd/ Content:

- Command structure and flag handling
- Help system customization

## `dg install` Command Structure

### Command Definition

```go
var installCmd = &cobra.Command{
    Use:   "install",
    Short: "Install devgita and all required tools",
    Long:  `[Comprehensive help text with steps and platforms]`,
    Run:   run,
}
```

### Flag Handling

- **StringSlice flags**: `--only` and `--skip` accept multiple values
- **Usage patterns**:
  - `dg install --only terminal,languages`
  - `dg install --skip desktop`
  - `dg install --only terminal --skip databases` (flags can be repeated)

### Installation Flow

1. **Validation Phase**
   - OS version validation
   - Package manager installation (Homebrew/apt)
   - Essential tools (git, fc-\*)

2. **Devgita Setup**
   - Clone devgita repository to `~/.config/devgita/`
   - Initialize global configuration
   - Setup tracking for installed vs pre-existing packages

3. **Category-Based Installation**
   - **Terminal**: Shell tools, tmux, neovim, fonts
   - **Languages**: Node.js, Python, Go, Rust (interactive selection)
   - **Databases**: PostgreSQL, Redis, MongoDB (interactive selection)
   - **Desktop**: GUI applications like Aerospace

### Help System Features

- **Multi-line Long description** with bullet points
- **Platform support** clearly documented
- **Flag examples** in help text
- **Step-by-step process** explanation
- **Error handling** with `utils.MaybeExitWithError()`

### Interactive Elements

- Uses **context** to pass user selections between steps
- **TUI selection** for languages and databases via `ChooseLanguages(ctx)`
- **Progress indicators** with `utils.PrintInfo()`
- **Bold headers** with `utils.PrintBold()`

This command exemplifies Cobra's power for complex, multi-step CLI operations with category filtering and interactive selection.
