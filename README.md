## Pending work

- Debian/Ubuntu compatibility
- Global document to see what we have installed already
  - Check if languages are already installed - to avoid duplicates/collisions
- Command `re-configure` app
  - Reinstall configuration files / should we uninstall and install app?
- Command `uninstall` app, font or theme.
  - make sure we only uninstall stuff related in the config file.
- Command to select a new `theme`
  - Update background images, too?
- Command to select a new `font`
- Help to display available commands
- MANUAL - Super important
- Revert installed programs if any error
  - We need to confirm if program was installed by us. If update leave the program there
- REMOVE GIT STUFF or any unrelated to installing new apps
- Add `verbose` option for all commands

### Commands

- `dg install ...` (`--soft` that does `maybeInstall`)
- `dg reinstall ...` (does configure and install)
- `dg configure ...`
- `dg re-configure ..` (or `configure --force`)
- `dg uninstall ...`
- `gd update ... [--neovim=[...options] --aerospace=[...options] ...flags] ...apps`
- `dg list` or `dg installed`
- `dg check-updates`
- `dg backup ...` - This would allow users to create backups of their current configurations before making changes.
- `dg restore ...` - This would allow users to revert to a previous configuration if needed.
- `dg validate ...` - This could check if the current configuration is valid and if all dependencies are met.
- `dg change --theme=[...options] --font=[...options]`

Note. We should optionally be able to pass `--app` or `--package` to only do it for one app/package

#### Considerations

- **Error Handling**: Ensure that all commands have robust error handling to provide meaningful feedback to users in case something goes wrong.
- **Logging**: Implement logging for actions taken by the commands, which can help in troubleshooting and understanding user actions.
- **User Prompts**: For commands that make significant changes (like uninstalling or reinstalling), consider adding user prompts to confirm actions.
- **Dependency Management**: Ensure that when installing or uninstalling, dependencies are also managed appropriately to avoid leaving orphaned packages.
- **Documentation**: Make sure that the MANUAL is comprehensive and includes examples for each command to help users understand how to use them effectively.

---

## Optional apps/packages

- Postman
- fc-list ([maybe required](https://github.com/cjairm/devgita/commit/c01797defb5e95a5ccce4206d46f435f9c513215)?)
- Music
- Emails (?)

---

## Questions

- Do we need to create shortcuts?
- What's that best way to handle mise? it's useful, but the documentation is difficult to follow
- (?) Git related - Move all git to `dg git clean --flags`, `dg git revert --flags`, or `dg clean-branch` ?
- (?) npm related - fully clean? fresh-installs? `dg npm clean`
- Command to update (we may need to solve issues). If we want to update apps can be difficult, we need to handle breaking changes
- Maybe a TUI?

---

- Command to update devgita - cli (download latest version)
- Add alias of `devgita` to `cli` (for easy access) (-or `dg`-)

```bash
# For M chips
GOOS=darwin GOARCH=amd64 go build -o cli_mac_amd64

# For intel chips
GOOS=darwin GOARCH=arm64 go build -o cli_mac_arm64
```
