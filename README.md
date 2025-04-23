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

---

## Questions

- Should we have `devgita` renamed to `dg`?
- Do we need to create shortcuts?
- What's that best way to handle mise? it's useful, but documentation is difficult to follow
- (?) Git related - Move all git to `devgita git clean --flags`, `devgita git revert --flags`
- (?) npm related - fully clean? fresh-installs? `devgita npm clean`
- Command to update (we may need to solve issues). If we want to update apps can be difficult, we need to handle breaking changes
- Maybe an TUI?

---

- Command to update devgita - cli (download latest version)
- Add alias of `devgita` to `cli` (for easy access)

```bash
# For M chips
GOOS=darwin GOARCH=amd64 go build -o cli_mac_amd64

# For intel chips
GOOS=darwin GOARCH=arm64 go build -o cli_mac_arm64
```
