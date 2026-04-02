# Cleanup Guide - Remove Pre-Release Installation

**⚠️ DEVELOPER USE ONLY** - This guide is for cleaning up the old git-based devgita installation before the v1.0.0 binary release.

## Why Clean Up?

The old devgita installation model used:
- Git clone to `~/.config/devgita/`
- Manual PATH configuration
- Repository-based config files

The new binary distribution uses:
- Pre-built binaries in `~/.local/bin/`
- Embedded configuration templates
- Automatic PATH and alias setup

## Cleanup Steps

### 1. Remove Old Repository Installation

```bash
# Remove the git repository
rm -rf ~/.config/devgita/

# Remove old global config (will be recreated by new version)
rm -f ~/.config/devgita/global_config.yaml
```

### 2. Clean Up Shell Configuration

Remove old devgita references from your shell config file (`~/.zshrc` or `~/.bashrc`):

```bash
# Open your shell config
vim ~/.zshrc  # or ~/.bashrc

# Look for and REMOVE these lines (if present):
# export PATH="$HOME/.config/devgita:$PATH"
# source ~/.config/devgita/devgita.zsh
# alias dg='...'  # Any old dg alias

# Keep only:
# export PATH="$HOME/.local/bin:$PATH"
# alias dg='devgita'
```

**Quick automated cleanup (Zsh):**
```bash
# Backup your config
cp ~/.zshrc ~/.zshrc.backup

# Remove old devgita lines
sed -i.bak '/\.config\/devgita/d' ~/.zshrc
```

**Quick automated cleanup (Bash):**
```bash
# Backup your config
cp ~/.bashrc ~/.bashrc.backup

# Remove old devgita lines
sed -i.bak '/\.config\/devgita/d' ~/.bashrc
```

### 3. Remove Old Binaries (if manually installed)

```bash
# Check if old binary exists
which devgita

# If it points to something OTHER than ~/.local/bin/devgita, remove it
# Example: if it was in /usr/local/bin
sudo rm -f /usr/local/bin/devgita
```

### 4. Verify Cleanup

```bash
# Check PATH doesn't reference old location
echo $PATH | grep -o '\.config/devgita'
# Should return nothing

# Check no old config sourcing
grep -n 'devgita.zsh' ~/.zshrc ~/.bashrc 2>/dev/null
# Should return nothing

# Check no old repository
ls -la ~/.config/devgita
# Should show "No such file or directory" OR only contain the new global_config.yaml
```

### 5. Fresh Install

After cleanup, install the new binary version:

```bash
# Install using the new installer
curl -fsSL https://raw.githubusercontent.com/cjairm/devgita/main/install.sh | bash

# Or for local testing:
bash install.sh --local ./devgita-darwin-arm64  # or your platform binary

# Restart shell
exec $SHELL

# Verify installation
dg --version
which dg  # Should show: dg: aliased to devgita
which devgita  # Should show: ~/.local/bin/devgita
```

## What the New Installation Creates

### Files Created
- `~/.local/bin/devgita` - The binary executable
- `~/.config/devgita/global_config.yaml` - Global state tracking (created on first run)
- Shell configs extracted to `~/.config/{alacritty,neovim,tmux,etc.}/` on first install

### Shell Configuration Added
```bash
# In ~/.zshrc or ~/.bashrc
export PATH="$HOME/.local/bin:$PATH"
alias dg='devgita'
```

### On First `dg install` Run
The binary will:
1. Extract embedded configs to `~/.config/devgita/` (config templates)
2. Create `global_config.yaml` for state tracking
3. Install tools and copy configurations as selected

## Troubleshooting

### "command not found: dg"

```bash
# Restart shell
exec $SHELL

# Or manually source
source ~/.zshrc  # or ~/.bashrc
```

### "dg points to old installation"

```bash
# Check current alias
type dg

# Check PATH order (new should be first)
echo $PATH | tr ':' '\n' | grep -n local

# Re-run installer
curl -fsSL https://raw.githubusercontent.com/cjairm/devgita/main/install.sh | bash
```

### "Config files still reference old paths"

This is expected! The new installation extracts configs on first run:

```bash
# Run install to trigger config extraction
dg install

# Configs will be extracted to proper locations
```

## Developer Testing Checklist

After cleanup and fresh install:

- [ ] `which devgita` shows `~/.local/bin/devgita`
- [ ] `which dg` shows `dg: aliased to devgita`
- [ ] `dg --version` works
- [ ] `~/.config/devgita/` does NOT contain git repository
- [ ] `~/.zshrc` or `~/.bashrc` has correct PATH and alias
- [ ] `dg install` runs successfully
- [ ] Configs are extracted to `~/.config/{app}/` directories
- [ ] `global_config.yaml` is created and tracks installations

## Complete Uninstall (if needed)

To completely remove devgita:

```bash
# Remove binary
rm -f ~/.local/bin/devgita

# Remove configs (WARNING: removes ALL devgita-managed configs)
rm -rf ~/.config/devgita/
rm -rf ~/.config/alacritty/  # if installed by devgita
rm -rf ~/.config/neovim/     # if installed by devgita
# ... etc for other apps

# Remove shell config lines
# Edit ~/.zshrc or ~/.bashrc and remove:
# export PATH="$HOME/.local/bin:$PATH"
# alias dg='devgita'

# Restart shell
exec $SHELL
```

## Questions?

If you encounter issues during cleanup:
1. Check the cleanup verification steps above
2. Review shell config for duplicate entries
3. Ensure old repository is completely removed
4. Try fresh install with `--local` flag for testing
