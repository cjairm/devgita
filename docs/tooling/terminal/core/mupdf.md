# MuPDF Module Documentation

## Overview

The MuPDF module provides installation and command execution management for MuPDF document viewer and toolkit with devgita integration. It follows the standardized devgita app interface while providing mupdf-specific operations for PDF viewing, rendering, and manipulation.

## App Purpose

MuPDF is a lightweight PDF, XPS, and E-book viewer with a small footprint. It consists of a software library, command line tools (mutool), and viewers for various platforms. MuPDF is designed for portability and includes support for many document formats including PDF, XPS, OpenXPS, CBZ, EPUB, and FictionBook. This module ensures mupdf is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for document processing tasks.

## Lifecycle Summary

1. **Installation**: Install mupdf package via platform package managers (Homebrew/apt)
2. **Configuration**: mupdf doesn't require separate configuration files (no default configuration applied by devgita)
3. **Execution**: Provide high-level mupdf operations for document viewing and manipulation via mutool

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Mupdf instance with platform-specific commands           |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install mupdf                             |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute mupdf commands    | Runs mupdf/mutool with provided arguments                            |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
mupdf := mupdf.New()
err := mupdf.Install()
```

- **Purpose**: Standard mupdf installation
- **Behavior**: Uses `InstallPackage()` to install mupdf package
- **Use case**: Initial mupdf installation or explicit reinstall

### ForceInstall()

```go
mupdf := mupdf.New()
err := mupdf.ForceInstall()
```

- **Purpose**: Force mupdf installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh mupdf installation or fix corrupted installation

### SoftInstall()

```go
mupdf := mupdf.New()
err := mupdf.SoftInstall()
```

- **Purpose**: Install mupdf only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing mupdf installations

### Uninstall()

```go
err := mupdf.Uninstall()
```

- **Purpose**: Remove mupdf installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Document processing tools are typically managed at the system level

### Update()

```go
err := mupdf.Update()
```

- **Purpose**: Update mupdf installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: MuPDF updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := mupdf.ForceConfigure()
err := mupdf.SoftConfigure()
```

- **Purpose**: Apply mupdf configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: MuPDF is a document viewer/toolkit without separate configuration files. Configuration is managed through command-line arguments when invoking mutool.

## Execution Methods

### ExecuteCommand()

```go
err := mupdf.ExecuteCommand("--version")
err := mupdf.ExecuteCommand("draw", "-o", "output.png", "input.pdf")
err := mupdf.ExecuteCommand("info", "document.pdf")
```

- **Purpose**: Execute mupdf/mutool commands with provided arguments
- **Parameters**: Variable arguments passed directly to mupdf/mutool binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### MuPDF-Specific Operations

The mupdf CLI (primarily through `mutool`) provides extensive document processing capabilities:

#### Document Information

```bash
# Get document info
mutool info document.pdf

# Show version
mutool --version

# List available commands
mutool help
```

#### PDF Rendering

```bash
# Render PDF to PNG
mutool draw -o output.png input.pdf

# Render specific page
mutool draw -o page3.png input.pdf 3

# Render with resolution
mutool draw -r 300 -o high-res.png input.pdf

# Render all pages
mutool draw -o page%d.png input.pdf
```

#### Text Extraction

```bash
# Extract text from PDF
mutool draw -F txt -o output.txt input.pdf

# Extract structured text (preserves layout)
mutool draw -F stext -o output.xml input.pdf

# Extract text from specific page
mutool draw -F txt input.pdf 5
```

#### PDF Manipulation

```bash
# Clean and optimize PDF
mutool clean input.pdf output.pdf

# Merge PDF files
mutool merge output.pdf file1.pdf file2.pdf file3.pdf

# Extract pages
mutool merge output.pdf input.pdf 1-5,10,15-20

# Convert to grayscale
mutool clean -g input.pdf output.pdf

# Remove annotations
mutool clean -a input.pdf output.pdf
```

#### Format Conversion

```bash
# Convert to SVG
mutool draw -F svg -o output.svg input.pdf

# Convert to HTML
mutool draw -F html -o output.html input.pdf

# Convert to PostScript
mutool draw -F ps -o output.ps input.pdf

# Convert to text
mutool draw -F txt -o output.txt input.pdf
```

#### Document Signing and Encryption

```bash
# Show certificate info
mutool sign -v document.pdf

# Encrypt PDF
mutool clean -e -p password input.pdf output.pdf

# Decrypt PDF
mutool clean -d -p password input.pdf output.pdf
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Document Operations**: `New()` → `SoftInstall()` → `ExecuteCommand()` with mutool arguments
4. **Version Check**: `New()` → `ExecuteCommand("--version")`

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.Mupdf` (typically "mupdf")
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No traditional config files**: MuPDF operations are configured via command-line arguments
- **Runtime configuration**: Parameters passed directly to `ExecuteCommand()`
- **Environment variables**: Standard system environment variables for locale and display

## Implementation Notes

- **Document Tool Nature**: Unlike typical applications, mupdf is a document viewer and processing toolkit without traditional config file templates
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since mupdf uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since mupdf uses runtime configuration
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as mupdf updates should be handled by system package managers
- **Main Tool**: Primary interface is through `mutool` command-line utility

## Usage Examples

### Basic MuPDF Operations

```go
mupdf := mupdf.New()

// Install mupdf
err := mupdf.SoftInstall()
if err != nil {
    return err
}

// Check version
err = mupdf.ExecuteCommand("--version")

// Get document info
err = mupdf.ExecuteCommand("info", "document.pdf")

// Render PDF to PNG
err = mupdf.ExecuteCommand("draw", "-o", "output.png", "input.pdf")
```

### Advanced Operations

```go
// Extract text from PDF
err := mupdf.ExecuteCommand("draw", "-F", "txt", "-o", "output.txt", "input.pdf")

// Clean and optimize PDF
err = mupdf.ExecuteCommand("clean", "input.pdf", "output.pdf")

// Merge multiple PDFs
err = mupdf.ExecuteCommand("merge", "combined.pdf", "file1.pdf", "file2.pdf")

// Render at high resolution
err = mupdf.ExecuteCommand("draw", "-r", "300", "-o", "high-res.png", "input.pdf")

// Convert to SVG
err = mupdf.ExecuteCommand("draw", "-F", "svg", "-o", "output.svg", "input.pdf")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Command Not Found**: Verify mupdf/mutool is installed and accessible in PATH
3. **Permission Issues**: Check file permissions for input/output operations
4. **Format Not Supported**: Some formats may require additional system libraries
5. **Rendering Errors**: Check PDF file integrity and format compatibility

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager
- **Format Support**: Different platforms may have different format support based on dependencies
- **Performance**: Large PDF files may require significant memory for rendering

### Document Format Support

MuPDF supports various document formats:

- **PDF**: Full support including annotations, forms, encryption
- **XPS**: Microsoft XML Paper Specification
- **OpenXPS**: Open XML Paper Specification
- **CBZ**: Comic book archive (ZIP)
- **EPUB**: Electronic publication
- **FictionBook**: FB2 format

### Performance Optimization

- **Caching**: MuPDF caches rendered pages for faster display
- **Resolution**: Use appropriate resolution (-r flag) for output quality vs. size
- **Memory**: Large documents benefit from streaming rendering
- **Compression**: Use clean command to optimize PDF size

### Command-Line Tools

MuPDF provides several command-line tools:

- **mutool**: Main command-line utility (draw, clean, merge, info, sign)
- **muraster**: Rasterization tool
- **mupdf**: Document viewer (GUI on some platforms)

## Integration with Devgita

MuPDF integrates with devgita's terminal category:

- **Installation**: Installed as part of core system tools
- **Configuration**: Runtime configuration via command-line arguments
- **Usage**: Available system-wide for document processing operations
- **Updates**: Managed through system package manager
- **Dependencies**: Minimal dependencies make it lightweight

## External References

- **MuPDF Official Site**: https://mupdf.com/
- **Documentation**: https://mupdf.com/docs/
- **Repository**: https://git.ghostscript.com/?p=mupdf.git
- **mutool Manual**: https://mupdf.com/docs/manual-mutool-run.html
- **API Documentation**: https://mupdf.com/docs/api/

## Key Features

### Lightweight Design

- Small binary size and minimal dependencies
- Fast startup and rendering
- Low memory footprint

### Format Support

- Wide range of document formats
- PDF with full feature support
- E-book formats (EPUB, FictionBook)

### Command-Line Tools

- Comprehensive mutool for document manipulation
- Text extraction and conversion
- PDF optimization and cleaning

### Rendering Quality

- High-quality rendering engine
- Configurable resolution
- Multiple output formats (PNG, SVG, PS, HTML)

### PDF Operations

- Merge and split documents
- Extract pages and content
- Optimize and compress
- Sign and encrypt

This module provides essential document viewing and manipulation capabilities for development workflows, documentation processing, and automated PDF operations within the devgita ecosystem.
