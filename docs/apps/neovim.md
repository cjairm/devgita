# Neovim Configuration Guide

## Installation

- Platform-specific installation steps
- Dependencies (ripgrep, fd, etc.)
- Version requirements

## Git Integration

- Setting as default editor
- Commit message configuration
- Git-specific Neovim settings

To set Neovim as the default editor for Git, you can configure it in your Git settings. Here’s how to do it:

1. Open your terminal.

Set Neovim as the default editor for Git: You can do this by running the following command:

```bash
git config --global core.editor "nvim"
```

This command sets Neovim (nvim) as the default editor for all your Git repositories.

2. Verify the configuration: To check if the configuration was set correctly, you can run:

```bash
git config --global --get core.editor
```

This should output `nvim`.

Using Neovim with Git: Now, whenever you run a Git command that requires an editor (like git commit), Neovim will open by default.

Additional Configuration (Optional)
If you want to set specific options for Neovim when it opens, you can modify the command slightly. For example, to open Neovim in a specific mode or with certain settings, you can do:

```bash
git config --global core.editor "nvim -c 'set ft=gitcommit'"
```

This command sets the file type to gitcommit when opening a commit message.

### Summary

By following these steps, you can easily set Neovim as your default editor for Git, enhancing your workflow with a powerful text editor.

## Configuration

- LSP server setup
- Plugin management with lazy.nvim
- Custom keybindings
- Theme integration

## Troubleshooting

- Common installation issues
- Path problems
- Plugin conflicts
