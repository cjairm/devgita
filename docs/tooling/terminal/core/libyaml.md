# Libyaml Module Documentation

## Overview

The Libyaml module provides installation and management for the libyaml YAML parser library with devgita integration. It follows the standardized devgita app interface while providing libyaml-specific operations for system library installation and management.

## App Purpose

Libyaml is a YAML 1.1 parser and emitter library written in C. It provides a low-level C API for parsing and emitting YAML documents, which is used as the foundation for YAML support in many programming languages including Ruby, Python, and others. This module ensures libyaml development libraries are properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems for building software that depends on YAML parsing functionality.

## Lifecycle Summary

1. **Installation**: Install libyaml package via platform package managers (Homebrew/apt)
2. **Configuration**: Libyaml typically doesn't require separate configuration files - it's a system library used by other applications
3. **Execution**: Provides placeholder operations for interface compliance (libyaml is a library, not a CLI tool)

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Libyaml instance with platform-specific commands         |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install libyaml                           |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute libyaml commands  | **Library only** - no CLI commands available                         |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
libyaml := libyaml.New()
err := libyaml.Install()
```

- **Purpose**: Standard libyaml installation
- **Behavior**: Uses `InstallPackage()` to install libyaml package
- **Platform differences**:
  - macOS: Installs `libyaml` via Homebrew
  - Debian/Ubuntu: Installs `libyaml-dev` via apt
- **Use case**: Initial libyaml installation or explicit reinstall

### ForceInstall()

```go
libyaml := libyaml.New()
err := libyaml.ForceInstall()
```

- **Purpose**: Force libyaml installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh libyaml installation or fix corrupted installation
- **Note**: Will fail since `Uninstall()` is not supported

### SoftInstall()

```go
libyaml := libyaml.New()
err := libyaml.SoftInstall()
```

- **Purpose**: Install libyaml only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing libyaml installations

### Uninstall()

```go
err := libyaml.Uninstall()
```

- **Purpose**: Remove libyaml installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: System libraries are typically managed at the system level and should not be uninstalled via devgita

### Update()

```go
err := libyaml.Update()
```

- **Purpose**: Update libyaml installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Libyaml updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := libyaml.ForceConfigure()
err := libyaml.SoftConfigure()
```

- **Purpose**: Apply libyaml configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Libyaml is a system library without separate configuration files. It's configured at compile-time by applications that use it.

## Execution Methods

### ExecuteCommand()

```go
err := libyaml.ExecuteCommand(args...)
```

- **Purpose**: Execute libyaml commands with provided arguments
- **Behavior**: Attempts to execute commands but libyaml has no direct CLI interface
- **Note**: Libyaml is a library, not a standalone CLI tool. This method exists for interface compliance but may not have practical use.
- **Parameters**: Variable arguments that would be passed to a hypothetical libyaml command
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Terminal Category**: Installed as part of core system libraries in the terminal tooling setup

## Constants and Paths

### Relevant Constants

- **Package name**: Must be defined in `pkg/constants/constants.go`
  ```go
  Libyaml = "libyaml"
  ```
- Used by all installation methods for consistent package reference

### Platform-Specific Package Names

- **macOS**: `libyaml` (via Homebrew)
- **Debian/Ubuntu**: `libyaml-dev` (via apt)
- **Note**: Platform abstraction handled by the command factory pattern

### Configuration Approach

- **No traditional config files**: Libyaml doesn't use runtime configuration files
- **Compile-time configuration**: Applications that use libyaml configure YAML parsing settings at compile-time
- **System library**: Installed to standard system library locations by package managers

## Implementation Notes

- **System Library Nature**: Unlike typical CLI applications, libyaml is a system library without direct user interaction
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since libyaml uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since libyaml doesn't use config files
- **ExecuteCommand**: Included for interface compliance but has limited practical use for a library
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as libyaml updates should be handled by system package managers

## Usage in Development

### Purpose as System Library

Libyaml is typically installed as a dependency for:

- Building Ruby interpreters and gems that use YAML
- Python YAML libraries (PyYAML uses libyaml for C extension)
- Language runtimes that require YAML configuration support
- Development tools that parse YAML configuration files
- Package managers that use YAML for configuration

### Developer Impact

While libyaml itself has no direct CLI commands, it's essential infrastructure for:

- Installing programming language runtimes (Ruby, Python)
- Building software that uses YAML for configuration
- Ensuring configuration management tools function correctly
- Supporting YAML parsing in custom software development

### YAML Use Cases

YAML (YAML Ain't Markup Language) is widely used for:

- Configuration files (Docker Compose, Kubernetes, CI/CD)
- Data serialization
- Infrastructure as Code (Ansible, CloudFormation)
- Application settings and metadata

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Missing Development Headers**: On Linux, ensure `libyaml-dev` is installed (not just `libyaml`)
3. **Build Failures**: Many Ruby gems and Python packages fail to build if libyaml is not present
4. **Version Conflicts**: Use system package manager to handle version dependencies

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager, requires dev package for headers
- **System Integration**: Installed to standard system library paths
- **Ruby Integration**: Required for native YAML support in Ruby (Psych gem)
- **Python Integration**: Required for PyYAML C extension for performance

### Common Build Errors

```bash
# Ruby gem installation failing
# Error: "libyaml is missing. Please install libyaml"
# Solution: Install libyaml-dev (Linux) or libyaml (macOS)

# Python PyYAML failing to build
# Error: "yaml.h: No such file or directory"
# Solution: Install libyaml development headers
```

## References

- **Libyaml Homepage**: https://pyyaml.org/wiki/LibYAML
- **Source Repository**: https://github.com/yaml/libyaml
- **YAML Specification**: https://yaml.org/spec/1.1/
- **YAML Official Site**: https://yaml.org/
- **PyYAML Documentation**: https://pyyaml.org/
- **Ruby Psych Documentation**: https://ruby-doc.org/stdlib/libdoc/psych/rdoc/Psych.html

## Integration with Devgita

Libyaml integrates with devgita's terminal core category:

- **Installation**: Installed as part of core system libraries setup
- **Configuration**: No configuration needed (system library)
- **Usage**: Provides YAML parsing functionality to language runtimes and tools
- **Updates**: Managed through system package manager
- **Dependencies**: Required by Ruby, Python, and other language runtime installations

### Installation Order

Libyaml should be installed before:

- Ruby (via Mise or system package manager)
- Python packages that require PyYAML
- Any tools or applications that parse YAML files

### Related Libraries

Often installed alongside:

- **readline**: For interactive shell support
- **openssl**: For SSL/TLS functionality
- **zlib**: For compression support
- **pkg-config**: For build system integration

## Version Information

### Current Stable Version

- **Latest**: 0.2.5 (as of documentation creation)
- **Minimum Required**: Varies by dependent software
- **Compatibility**: YAML 1.1 specification

### Version Checking

```bash
# Check if libyaml is installed (macOS)
brew list libyaml

# Check if libyaml is installed (Linux)
dpkg -l | grep libyaml

# Find libyaml version (pkg-config)
pkg-config --modversion yaml-0.1
```

## Language-Specific Integration

### Ruby

```ruby
# Ruby uses Psych gem (built on libyaml)
require 'yaml'
data = YAML.load_file('config.yml')
```

### Python

```python
# Python uses PyYAML (can use libyaml for C extension)
import yaml
with open('config.yml', 'r') as file:
    data = yaml.safe_load(file)
```

### C/C++

```c
// Direct usage of libyaml C API
#include <yaml.h>
yaml_parser_t parser;
yaml_parser_initialize(&parser);
```

This module provides essential YAML parsing library support for building and running development tools within the devgita ecosystem, ensuring that programming language runtimes and configuration management tools requiring YAML functionality can be installed and executed successfully.
