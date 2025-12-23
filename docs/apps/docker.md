# Docker Module Documentation

## Overview

The Docker module provides installation and configuration management for Docker Desktop with devgita integration. It follows the standardized devgita app interface while providing Docker-specific operations for container platform setup and command execution.

## App Purpose

Docker Desktop is a containerization platform that enables developers to build, ship, and run distributed applications using containers. This module ensures Docker Desktop is properly installed across macOS (Homebrew cask) and Debian/Ubuntu (apt) systems and provides high-level operations for Docker command execution and container management.

## Lifecycle Summary

1. **Installation**: Install Docker Desktop desktop application via platform package managers (Homebrew cask/apt)
2. **Configuration**: Docker Desktop uses GUI-based configuration and daemon.json (no default configuration applied by devgita)
3. **Execution**: Provide high-level Docker operations for container management and image operations

## Exported Functions

| Function           | Purpose                   | Behavior                                                       |
| ------------------ | ------------------------- | -------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Docker instance with platform-specific commands    |
| `Install()`        | Standard installation     | Uses `InstallDesktopApp()` to install Docker Desktop           |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error), then `Install()`    |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallDesktopApp()` to check before installing     |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                               |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                               |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error with manual cleanup steps    |
| `ExecuteCommand()` | Execute docker commands   | Runs docker with provided arguments                            |
| `Update()`         | Update installation       | **Not implemented** - returns error                            |

## Installation Methods

### Install()

```go
docker := docker.New()
err := docker.Install()
```

- **Purpose**: Standard Docker Desktop installation
- **Behavior**: Uses `InstallDesktopApp()` to install Docker Desktop application
- **Use case**: Initial Docker Desktop installation or explicit reinstall

### ForceInstall()

```go
docker := docker.New()
err := docker.ForceInstall()
```

- **Purpose**: Force Docker Desktop installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error since uninstall is not supported), then `Install()`
- **Use case**: Ensure fresh Docker Desktop installation

### SoftInstall()

```go
docker := docker.New()
err := docker.SoftInstall()
```

- **Purpose**: Install Docker Desktop only if not already present
- **Behavior**: Uses `MaybeInstallDesktopApp()` to check before installing
- **Use case**: Standard installation that respects existing Docker Desktop installations

### Uninstall()

```go
err := docker.Uninstall()
```

- **Purpose**: Remove Docker Desktop installation
- **Behavior**: **Not supported** - returns error with manual cleanup instructions
- **Rationale**: Docker Desktop uninstallation requires elevated privileges, interactive confirmation, and comprehensive cleanup across multiple system locations
- **Manual cleanup steps** (macOS):
  1. Quit Docker Desktop: Right-click Docker icon in menu bar → "Quit Docker Desktop"
  2. Uninstall cask: `brew uninstall --cask docker`
  3. Remove binaries:
     ```bash
     sudo rm -f /usr/local/bin/docker
     sudo rm -f /usr/local/bin/docker-compose
     sudo rm -f /usr/local/bin/docker-credential-desktop
     sudo rm -f /usr/local/bin/docker-credential-ecr-login
     sudo rm -f /usr/local/bin/docker-credential-osxkeychain
     sudo rm -f /usr/local/bin/hub-tool
     sudo rm -f /usr/local/bin/kubectl.docker
     ```
  4. Remove containers and data:
     ```bash
     sudo rm -rf ~/Library/Containers/com.docker.docker
     sudo rm -rf ~/Library/Application\ Support/Docker\ Desktop
     sudo rm -rf ~/.docker
     ```
  5. Remove shell completions:
     ```bash
     sudo rm -f /usr/local/etc/bash_completion.d/docker
     sudo rm -f /usr/local/share/zsh/site-functions/_docker
     sudo rm -f /usr/local/share/fish/vendor_completions.d/docker.fish
     ```

### Update()

```go
err := docker.Update()
```

- **Purpose**: Update Docker Desktop installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Docker Desktop updates are typically handled by the system package manager or Docker Desktop's built-in update mechanism

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := docker.ForceConfigure()
err := docker.SoftConfigure()
```

- **Purpose**: Apply Docker Desktop configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Docker Desktop configuration is managed through the Docker Desktop GUI or `~/.docker/daemon.json` file. Devgita does not apply default Docker configuration to allow users to configure based on their specific needs.

### Configuration Options

While devgita doesn't apply default configuration, users can customize Docker Desktop via:

- **GUI Settings**: Docker Desktop → Preferences
- **daemon.json**: `~/.docker/daemon.json` (user-managed)

Example `daemon.json`:

```json
{
  "builder": {
    "gc": {
      "defaultKeepStorage": "20GB",
      "enabled": true
    }
  },
  "experimental": false,
  "features": {
    "buildkit": true
  },
  "insecure-registries": [],
  "registry-mirrors": []
}
```

## Execution Methods

### ExecuteCommand()

```go
err := docker.ExecuteCommand("--version")
err := docker.ExecuteCommand("ps", "-a")
err := docker.ExecuteCommand("images")
err := docker.ExecuteCommand("run", "-it", "ubuntu", "bash")
```

- **Purpose**: Execute docker commands with provided arguments
- **Parameters**: Variable arguments passed directly to docker binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Docker-Specific Operations

The docker CLI provides extensive container and image management capabilities:

#### Container Operations

```bash
# List containers
docker ps
docker ps -a

# Run containers
docker run -it ubuntu bash
docker run -d -p 80:80 nginx

# Container lifecycle
docker start <container>
docker stop <container>
docker restart <container>
docker rm <container>

# Container inspection
docker logs <container>
docker inspect <container>
docker exec -it <container> bash
```

#### Image Operations

```bash
# List images
docker images

# Pull images
docker pull ubuntu:latest
docker pull nginx:alpine

# Build images
docker build -t myapp:latest .
docker build -f Dockerfile.dev -t myapp:dev .

# Image management
docker rmi <image>
docker tag <image> <tag>
docker push <image>
```

#### Docker Compose

```bash
# Start services
docker-compose up
docker-compose up -d

# Stop services
docker-compose down
docker-compose down -v

# View logs
docker-compose logs
docker-compose logs -f service-name

# Execute commands
docker-compose exec service-name bash
```

#### System Operations

```bash
# System info
docker info
docker version

# Cleanup
docker system prune
docker system prune -a

# Resource usage
docker stats
docker system df
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` (fails if docker installed) → `ForceConfigure()` (no-op)
3. **Docker Operations**: `New()` → `SoftInstall()` → `ExecuteCommand()` with docker arguments

## Constants and Paths

### Relevant Constants

- `constants.Docker`: Package name ("docker") for installation and command execution
- Used by all installation methods for consistent desktop app reference

### Configuration Approach

- **GUI-based**: Primary configuration through Docker Desktop GUI
- **daemon.json**: Optional advanced configuration (`~/.docker/daemon.json`)
- **No default config**: Devgita does not apply default configuration for Docker Desktop
- **User customization**: Users configure based on their specific development needs

## Implementation Notes

- **Desktop App Installation**: Uses `InstallDesktopApp()` instead of `InstallPackage()` for GUI applications
- **ForceInstall Logic**: Calls `Uninstall()` first but will fail since Docker uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since Docker uses GUI/daemon.json configuration
- **Uninstall Complexity**: Docker Desktop uninstall requires manual steps due to elevated privileges and comprehensive cleanup requirements
- **Error Handling**: All methods return errors that should be checked by callers
- **Cross-Platform**: Works on both macOS and Linux through desktop app installation methods
- **Update Method**: Not implemented as Docker Desktop updates should be handled by system package managers or Docker Desktop's built-in updater

## Usage Examples

### Basic Docker Operations

```go
docker := docker.New()

// Install Docker Desktop
err := docker.SoftInstall()
if err != nil {
    return err
}

// Check version
err = docker.ExecuteCommand("--version")

// List containers
err = docker.ExecuteCommand("ps", "-a")

// Pull image
err = docker.ExecuteCommand("pull", "nginx:alpine")

// Run container
err = docker.ExecuteCommand("run", "-d", "-p", "8080:80", "nginx:alpine")
```

### Advanced Operations

```go
// Build image
err := docker.ExecuteCommand("build", "-t", "myapp:latest", ".")

// Run with environment variables
err = docker.ExecuteCommand("run", "-e", "ENV=production", "myapp:latest")

// System cleanup
err = docker.ExecuteCommand("system", "prune", "-a", "-f")

// View container logs
err = docker.ExecuteCommand("logs", "-f", "container-name")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Docker daemon not running**: Start Docker Desktop application from Applications folder
3. **Permission errors**: Add user to docker group (Linux): `sudo usermod -aG docker $USER`
4. **Port conflicts**: Check for conflicting services using same ports
5. **Commands Don't Work**: Verify Docker Desktop is installed and daemon is running

### Platform Considerations

- **macOS**: Installed via Homebrew cask as Docker.app
- **Linux**: Installed via apt or other package managers
- **Docker daemon**: Must be running for docker commands to work
- **Resource limits**: Configured in Docker Desktop preferences

### Prerequisites

Before using Docker Desktop:
- **macOS**: macOS 11 (Big Sur) or newer
- **Linux**: 64-bit Linux distribution with kernel 3.10+
- **Virtualization**: Hardware virtualization support enabled in BIOS

### Docker Desktop Features

- **Kubernetes**: Optional local Kubernetes cluster
- **Docker Compose**: Bundled for multi-container applications
- **Volume mounting**: Easy file sharing between host and containers
- **Networking**: Built-in networking for container communication
- **Resource controls**: CPU, memory, and disk space limits

### Security Considerations

- **Privileged containers**: Use with caution
- **Image scanning**: Docker Desktop includes vulnerability scanning
- **Registry authentication**: Use `docker login` for private registries
- **Network isolation**: Use Docker networks for service isolation

## External References

- **Docker Documentation**: https://docs.docker.com/
- **Docker Desktop**: https://www.docker.com/products/docker-desktop/
- **Docker Hub**: https://hub.docker.com/
- **Docker Compose**: https://docs.docker.com/compose/
- **Dockerfile Reference**: https://docs.docker.com/engine/reference/builder/

## Integration with Devgita

Docker Desktop integrates with devgita's desktop category:

- **Installation**: Installed as part of desktop applications setup
- **Configuration**: User-managed through Docker Desktop GUI
- **Usage**: Available system-wide after installation
- **Updates**: Managed through system package manager or Docker Desktop updater
- **Dependencies**: Lazydocker (terminal UI) requires Docker Desktop to be installed

## Key Features

### Container Management
- Build, run, and manage containers
- Interactive container shell access
- Container lifecycle management
- Resource monitoring and statistics

### Image Management
- Pull images from Docker Hub and private registries
- Build custom images with Dockerfiles
- Tag and push images to registries
- Multi-stage builds for optimized images

### Docker Compose
- Define multi-container applications
- Orchestrate service dependencies
- Network and volume management
- Environment-specific configurations

### Development Workflow
- Local development with containers
- Consistent development environments
- Easy microservices testing
- CI/CD pipeline integration

This module provides essential containerization platform capabilities for modern development workflows within the devgita development environment ecosystem.
