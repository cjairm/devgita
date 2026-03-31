# Quickstart: Binary Distribution with Embedded Configs

**Feature**: `001-binary-dist-audit`

## For Users (end-to-end install)

### macOS

```bash
curl -fsSL https://raw.githubusercontent.com/cjairm/devgita/main/install.sh | bash
source ~/.zshrc
dg install
```

### Debian/Ubuntu

```bash
curl -fsSL https://raw.githubusercontent.com/cjairm/devgita/main/install.sh | bash
source ~/.bashrc  # or ~/.zshrc
dg install
```

### Selective Install

```bash
dg install --only terminal          # Terminal tools only
dg install --only languages         # Languages only (interactive)
dg install --skip desktop           # Everything except desktop apps
dg install --only terminal,desktop  # Terminal + desktop
```

## For Maintainers (build and release)

### Build All Binaries

```bash
GOOS=darwin GOARCH=arm64 go build -o devgita-darwin-arm64 main.go
GOOS=darwin GOARCH=amd64 go build -o devgita-darwin-amd64 main.go
GOOS=linux  GOARCH=amd64 go build -o devgita-linux-amd64 main.go
```

### Test Locally (no upload needed)

```bash
# Test the binary directly
./devgita-darwin-arm64 install

# Or test via install.sh with --local flag
bash install.sh --local ./devgita-darwin-arm64
```

### Release to GitHub

```bash
# Tag and push
git tag v0.1.0
git push origin v0.1.0

# Upload binaries to the release via GitHub web UI or gh CLI
gh release create v0.1.0 \
  devgita-darwin-arm64 \
  devgita-darwin-amd64 \
  devgita-linux-amd64
```

## Verify Installation

```bash
# Check binary is installed
which dg

# Check configs were extracted
ls ~/.config/devgita/configs/

# Check global config exists
cat ~/.config/devgita/global_config.yaml

# Run with verbose logging
dg install --verbose
```
