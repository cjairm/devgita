# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Cross-platform development environment manager
- Support for macOS (Homebrew) and Debian/Ubuntu (apt)
- Automated installation for terminal tools, languages, databases
- Configuration management with global state tracking

### Categories Supported

- **Terminal**: Alacritty, Tmux, Neovim, shell enhancements
- **Languages**: Node.js, Python, Go, Rust with interactive selection
- **Databases**: PostgreSQL, Redis, MongoDB
- **Desktop**: Aerospace window manager

## [0.1.0] - Initial Release

### Added

- Basic CLI structure with Cobra framework
- Platform detection and abstraction layer
- Configuration templates system
- Smart installation with conflict detection
