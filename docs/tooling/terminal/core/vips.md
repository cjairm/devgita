# Vips Module Documentation

## Overview

The Vips module provides installation and command execution management for libvips image processing library with devgita integration. It follows the standardized devgita app interface while providing vips-specific operations for fast image processing with low memory requirements.

## App Purpose

libvips (vips) is a demand-driven, horizontally threaded image processing library with operations for arithmetic, histograms, convolution, morphological operations, filtering, DFT, profile operations, and file format support. This module ensures vips is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for image processing tasks.

## Lifecycle Summary

1. **Installation**: Install vips package via platform package managers (Homebrew/apt)
2. **Configuration**: vips doesn't require separate configuration files (no default configuration applied by devgita)
3. **Execution**: Provide high-level vips operations for image processing and format conversion

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Vips instance with platform-specific commands            |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install vips                              |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute vips commands     | Runs vips with provided arguments                                    |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
vips := vips.New()
err := vips.Install()
```

- **Purpose**: Standard vips installation
- **Behavior**: Uses `InstallPackage()` to install vips package
- **Use case**: Initial vips installation or explicit reinstall

### ForceInstall()

```go
vips := vips.New()
err := vips.ForceInstall()
```

- **Purpose**: Force vips installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh vips installation or fix corrupted installation

### SoftInstall()

```go
vips := vips.New()
err := vips.SoftInstall()
```

- **Purpose**: Install vips only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing vips installations

### Uninstall()

```go
err := vips.Uninstall()
```

- **Purpose**: Remove vips installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: System image processing libraries are typically managed at the system level

### Update()

```go
err := vips.Update()
```

- **Purpose**: Update vips installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Vips updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := vips.ForceConfigure()
err := vips.SoftConfigure()
```

- **Purpose**: Apply vips configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Vips is a system library without separate configuration files. Configuration is managed through runtime options and environment variables.

## Execution Methods

### ExecuteCommand()

```go
err := vips.ExecuteCommand("--version")
err := vips.ExecuteCommand("resize", "input.jpg", "output.jpg")
err := vips.ExecuteCommand("thumbnail", "image.jpg", "thumb.jpg", "200")
```

- **Purpose**: Execute vips commands with provided arguments
- **Parameters**: Variable arguments passed directly to vips binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Vips-Specific Operations

The vips CLI provides extensive image processing capabilities:

#### Image Information

```bash
# Show image information
vips --version
vips --list classes
vips --list operations

# Image properties
vips image-get input.jpg width
vips image-get input.jpg height
vips image-get input.jpg bands
```

#### Image Conversion

```bash
# Format conversion
vips copy input.jpg output.png
vips copy input.png output.webp
vips copy input.tiff output.jpg

# Quality settings
vips jpegsave input.jpg output.jpg --Q 90
vips webpsave input.jpg output.webp --Q 80
```

#### Image Resizing

```bash
# Resize operations
vips resize input.jpg output.jpg 0.5
vips thumbnail input.jpg output.jpg 200
vips thumbnail input.jpg output.jpg 200 --height 150

# Smart cropping
vips smartcrop input.jpg output.jpg 300 200
vips gravity input.jpg output.jpg north 300 200
```

#### Image Processing

```bash
# Basic operations
vips invert input.jpg output.jpg
vips flip input.jpg output.jpg vertical
vips rot input.jpg output.jpg d90

# Color operations
vips colourspace input.jpg output.jpg srgb
vips extract_band input.jpg output.jpg 0

# Filtering
vips gaussblur input.jpg output.jpg 5
vips sharpen input.jpg output.jpg
vips conv input.jpg output.jpg kernel.mat
```

#### Batch Processing

```bash
# Process multiple images
for img in *.jpg; do
    vips thumbnail "$img" "thumb_$img" 200
done

# Format conversion batch
for img in *.png; do
    vips copy "$img" "${img%.png}.webp"
done
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Image Operations**: `New()` → `SoftInstall()` → `ExecuteCommand()` with vips arguments
4. **Version Check**: `New()` → `ExecuteCommand("--version")`

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.Vips` (typically "vips")
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No traditional config files**: Vips operations are configured via command-line arguments and environment variables
- **Runtime configuration**: Parameters passed directly to `ExecuteCommand()`
- **Environment variables**:
  - `VIPS_CONCURRENCY`: Control thread pool size
  - `VIPS_DISC_THRESHOLD`: Control memory vs disc usage
  - `VIPS_WARNING`: Enable/disable warnings

## Implementation Notes

- **System Library Nature**: Unlike typical applications, vips is an image processing library without traditional config file templates
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since vips uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since vips uses runtime configuration
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as vips updates should be handled by system package managers

## Usage Examples

### Basic Vips Operations

```go
vips := vips.New()

// Install vips
err := vips.SoftInstall()
if err != nil {
    return err
}

// Check version
err = vips.ExecuteCommand("--version")

// Resize image
err = vips.ExecuteCommand("resize", "input.jpg", "output.jpg", "0.5")

// Create thumbnail
err = vips.ExecuteCommand("thumbnail", "photo.jpg", "thumb.jpg", "200")
```

### Advanced Operations

```go
// Format conversion
err := vips.ExecuteCommand("copy", "input.png", "output.webp")

// Smart crop
err = vips.ExecuteCommand("smartcrop", "input.jpg", "output.jpg", "300", "200")

// Blur effect
err = vips.ExecuteCommand("gaussblur", "input.jpg", "blurred.jpg", "10")

// List available operations
err = vips.ExecuteCommand("--list", "operations")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Memory Issues**: Large images may require adjusting `VIPS_DISC_THRESHOLD`
3. **Format Support**: Some formats require additional system libraries
4. **Permission Issues**: Check file permissions for input/output operations
5. **Commands Don't Work**: Verify vips is installed and accessible in PATH

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager
- **Format Support**: Different platforms may have different format support based on dependencies
- **Performance**: Multi-threaded operations benefit from multiple CPU cores

### Image Format Support

Vips supports numerous image formats:

- **Input/Output**: JPEG, PNG, TIFF, WebP, GIF, SVG, PDF, HEIF
- **Input Only**: RAW formats (via libraw), OpenSlide formats
- **Specialized**: Scientific formats (FITS), Medical formats (DICOM)

### Performance Optimization

- **Streaming**: Vips processes images in small chunks (demand-driven)
- **Threading**: Automatically uses multiple threads for operations
- **Memory**: Efficiently manages memory through tiling
- **Caching**: Smart caching of intermediate results

### Environment Variables

```bash
# Control concurrency
export VIPS_CONCURRENCY=8

# Control memory usage
export VIPS_DISC_THRESHOLD=500m

# Enable warnings
export VIPS_WARNING=1

# Progress reporting
export VIPS_PROGRESS=1
```

## Integration with Devgita

Vips integrates with devgita's terminal category:

- **Installation**: Installed as part of core system libraries
- **Configuration**: Runtime configuration via environment variables
- **Usage**: Available system-wide for image processing operations
- **Updates**: Managed through system package manager
- **Dependencies**: Often used by other tools requiring image processing

## External References

- **Vips Documentation**: https://www.libvips.org/
- **Vips Repository**: https://github.com/libvips/libvips
- **API Documentation**: https://www.libvips.org/API/current/
- **Operations Reference**: https://www.libvips.org/API/current/using-cli.html
- **Benchmarks**: https://github.com/libvips/libvips/wiki/Speed-and-memory-use

## Key Features

### Speed and Efficiency

- Demand-driven architecture processes only necessary pixels
- Multi-threaded operations for parallel processing
- Smart memory management with streaming

### Format Support

- Wide range of input/output formats
- Scientific and medical image formats
- RAW camera formats support

### Operations

- Resize, crop, rotate, flip operations
- Color space conversions
- Filtering and convolution
- Histogram operations
- Morphological operations

### Integration

- Command-line interface for scripting
- C API for integration
- Bindings for many languages (Python, Ruby, PHP, etc.)

This module provides essential image processing capabilities for development workflows, build systems, and automated image optimization within the devgita ecosystem.
