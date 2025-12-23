# Databases Module Documentation

## Overview

The Databases module coordinates database installation and management with devgita integration. It follows the standardized tooling coordinator pattern while providing database-specific operations for database management via native package managers.

## Module Purpose

This coordinator module manages the installation of database systems across development environments. It provides:

- **Interactive database selection** via TUI with installed database detection
- **Native package manager installation** for all databases (Homebrew/apt)
- **GlobalConfig tracking** of installed databases for proper state management
- **Conflict prevention** by filtering already-installed databases from selection

## Lifecycle Summary

1. **Selection**: Present TUI for database selection, filtering out already-installed databases
2. **Installation**: Install selected databases via native package manager
3. **Tracking**: Record installations in GlobalConfig for state management and future reference

## Architecture

### Coordinator Pattern

The Databases module acts as a **coordinator** (not an app), orchestrating database installations:

```
Databases (coordinator)
├── DatabaseConfig (configuration)
└── Native Package Manager (via Command interface)
```

Unlike individual apps in `internal/apps/`, coordinators:
- Manage multiple installation strategies
- Provide interactive selection UIs
- Coordinate between different tools
- Track state in GlobalConfig

### Comparison with Languages Coordinator

Similar to `internal/tooling/languages/languages.go`:

| Aspect | Languages | Databases |
|--------|-----------|-----------|
| Pattern | Coordinator | Coordinator |
| Selection | Interactive TUI | Interactive TUI |
| Installation | Mise + Native | Native only |
| Strategy | LanguageConfig array | DatabaseConfig array |
| Tracking | Dev languages | Databases |
| Pre-detection | System check on New() | System check on New() |

## Exported Functions

| Function | Purpose | Behavior |
|----------|---------|----------|
| `New()` | Factory method | Creates Databases instance with command executors |
| `GetSelectionOptions()` | Get TUI options | Returns TUI menu options including control items and database names |
| `ChooseDatabases(ctx)` | Interactive selection | Presents TUI, filters installed databases, returns context |
| `InstallChosen(ctx)` | Install selected | Installs databases from context via native package manager |

## Database Configuration

### DatabaseConfig Structure

```go
type DatabaseConfig struct {
    DisplayName string  // Human-readable name (e.g., "PostgreSQL")
    Name        string  // Package name (e.g., "postgresql")
}
```

### Available Databases

| Database | Name | Installation Method |
|----------|------|---------------------|
| MongoDB | mongodb | Native package manager |
| MySQL | mysql | Native package manager |
| PostgreSQL | postgresql | Native package manager |
| Redis | redis | Native package manager |
| SQLite | sqlite | Native package manager |

**Note**: Database names are defined as constants in `/pkg/constants/constants.go`:

```go
const (
    MongoDB    = "mongodb"
    MySQL      = "mysql"
    PostgreSQL = "postgresql"
    Redis      = "redis"
    SQLite     = "sqlite"
)
```

These constants are used throughout the databases module to ensure consistency and type safety, following the same pattern as the languages coordinator.

### Database Specifications

Databases are tracked in GlobalConfig with these specifications:

- **All databases**: `name` only (e.g., `redis`, `postgresql`, `mysql`)

## Installation Flow

### 1. Interactive Selection

```go
d := databases.New()
ctx, err := d.ChooseDatabases(ctx)
```

**Process**:
1. Load GlobalConfig to identify installed databases
2. Filter installed databases from available options
3. Display warning if databases are already installed
4. Present TUI multi-select with filtered options
5. Store selections in context

**Example Output**:
```
Already installed databases (skipped from selection): Redis, PostgreSQL
? Select databases to install
  ▸ All
    None
    Done
    MongoDB
    MySQL
    SQLite
```

### 2. Install Selected Databases

```go
d.InstallChosen(ctx)
```

**Process** for each selected database:

#### Native Installation (All Databases)

```
1. Install via package manager (MaybeInstallPackage)
2. Track in GlobalConfig as "name"
```

### 3. State Tracking

```go
gc.AddToInstalled(databaseSpec, "database")
gc.Save()
```

Tracked in `~/.config/devgita/global_config.yaml`:

```yaml
installed:
  databases:
    - redis
    - postgresql
    - mysql
```

## Implementation Details

### Installed Database Detection

```go
func (d *Databases) getInstalledDatabases(gc *config.GlobalConfig) []string
```

**Logic**:
1. Iterate through database configurations
2. Generate database spec (name only)
3. Check in both `installed.databases` and `already_installed.databases`
4. Return display names of installed databases

**Purpose**: Filter already-installed databases from TUI selection to prevent conflicts

### Installation Strategy

```go
err = d.installNative(dbCfg)
```

**Native Installation**:
- Uses platform package manager (Homebrew/apt)
- Calls `MaybeInstallPackage(name)`
- Simple approach for all databases

### Error Handling

```go
if err != nil {
    utils.PrintError("Error: Unable to install %s: %v", name, err)
    logger.L().Errorw("Database installation failed", "database", name, "error", err)
    return // Non-fatal: continue with other databases
}
```

**Philosophy**:
- Database installation failures are **non-fatal**
- Print error to user, log details, continue with remaining databases
- GlobalConfig tracking failures are **non-fatal** but logged as warnings

## Usage Examples

### Standard Installation Flow

```go
// In cmd/install.go
func installDatabases(ctx context.Context, onlySet, skipSet map[string]bool) {
    if shouldInstall("databases", onlySet, skipSet) {
        d := databases.New()
        ctx, err := d.ChooseDatabases(ctx)
        utils.MaybeExitWithError(err)
        d.InstallChosen(ctx)
    }
}
```

### Adding a New Database

To add a new database (e.g., Cassandra):

1. **Add to configurations**:
```go
func GetDatabaseConfigs() []DatabaseConfig {
    return []DatabaseConfig{
        // ... existing databases ...
        {DisplayName: "Cassandra", Name: "cassandra"},
    }
}
```

2. **Add to version command map** (optional, for detection):
```go
func getVersionCommand(dbName string) (string, []string) {
    switch dbName {
    // ... existing cases ...
    case "cassandra":
        return "cassandra", []string{"--version"}
    default:
        return dbName, []string{"--version"}
    }
}
```

That's it! The coordinator handles the rest automatically.

## GlobalConfig Integration

### Schema

```yaml
installed:
  databases:
    - redis       # Native package manager
    - postgresql  # Native package manager
    - mysql       # Native package manager
    - sqlite      # Native package manager

already_installed:
  databases:
    - mongodb     # Pre-existing, tracked for awareness
```

### Methods Used

- `gc.AddToInstalled(databaseSpec, "database")` - Track devgita installation
- `gc.IsInstalledByDevgita(databaseSpec, "database")` - Check if installed
- `gc.IsAlreadyInstalled(databaseSpec, "database")` - Check if pre-existing

## Testing Patterns

### Test Structure

```go
func TestInstallChosen_WithSelections(t *testing.T) {
    mockApp := testutil.NewMockApp()
    mockApp.Base.SetExecCommandResult("", "", nil)
    
    d := &Databases{Cmd: mockApp.Cmd, Base: mockApp.Base}
    
    ctx := context.Background()
    config := config.ContextConfig{SelectedDbs: []string{"Redis"}}
    ctx = config.WithConfig(ctx, config)
    
    d.InstallChosen(ctx)
    
    testutil.VerifyNoRealCommands(t, mockApp.Base)
}
```

### Test Coverage

- ✅ Configuration structure validation
- ✅ Database spec generation
- ✅ Installed database detection
- ✅ Selection filtering logic
- ✅ Installation orchestration
- ✅ GlobalConfig tracking
- ✅ Pre-installed database detection
- ✅ Version command mapping

## Comparison with Individual Apps

| Aspect | Individual App (e.g., Git) | Databases Coordinator |
|--------|----------------------------|----------------------|
| Location | `internal/apps/git/` | `internal/tooling/databases/` |
| Interface | Implements standard app interface | Custom coordinator interface |
| Installation | Single package | Multiple databases |
| Configuration | Config files (e.g., .gitconfig) | GlobalConfig tracking only |
| Selection | N/A | Interactive TUI |
| Structure | Singular focus | Multi-database orchestration |

## Integration with Other Modules

### Command Interface

Uses `Command` and `BaseCommandExecutor` interfaces:

```go
type Databases struct {
    Cmd  cmd.Command              // Package management
    Base cmd.BaseCommandExecutor  // Command execution
}
```

### Context Configuration

Leverages context for passing selections:

```go
// Store selections
initialConfig := config.ContextConfig{SelectedDbs: selections}
ctx = config.WithConfig(ctx, initialConfig)

// Retrieve selections
selections, ok := config.GetConfig(ctx)
```

## Platform Considerations

### macOS (Homebrew)

- **Redis**: Installed via `brew install redis`
- **PostgreSQL**: Installed via `brew install postgresql`
- **MySQL**: Installed via `brew install mysql`
- **SQLite**: Installed via `brew install sqlite`
- **MongoDB**: Installed via `brew install mongodb-community`

### Linux (Debian/Ubuntu)

- **Redis**: Installed via `apt install redis`
- **PostgreSQL**: Installed via `apt install postgresql`
- **MySQL**: Installed via `apt install mysql-server`
- **SQLite**: Installed via `apt install sqlite3`
- **MongoDB**: Requires MongoDB repository setup

## Future Enhancements

### Potential Additions

1. **Service management**: Start/stop database services after installation
2. **Configuration templates**: Apply devgita-optimized database configs
3. **Uninstall support**: Remove databases and clean up GlobalConfig
4. **Update support**: Update database versions
5. **Alternative databases**: Support for CouchDB, DynamoDB local, etc.
6. **Container-based**: Option to install databases via Docker
7. **Version selection**: Allow users to specify database versions

### Planned Databases

- CouchDB (via native package manager)
- DynamoDB Local (via download)
- Elasticsearch (via native package manager)
- InfluxDB (via native package manager)

## Troubleshooting

### Common Issues

1. **Installation fails**: Ensure package manager is available and updated
2. **Service not starting**: Check platform-specific service management
3. **Permission errors**: Some databases require sudo for installation
4. **GlobalConfig errors**: Non-fatal but logged; check `~/.config/devgita/global_config.yaml`

### Debugging

```bash
# Check installed databases
brew list | grep -E "redis|postgresql|mysql|sqlite|mongodb"  # macOS
dpkg -l | grep -E "redis|postgresql|mysql|sqlite|mongodb"   # Linux

# Check database services
brew services list  # macOS
systemctl list-units --type=service | grep -E "redis|postgres|mysql"  # Linux

# Check devgita global config
cat ~/.config/devgita/global_config.yaml

# Run with verbose logging
dg install --only databases --verbose
```

## Key Design Decisions

### Why Coordinator Pattern?

1. **Complexity**: Managing multiple databases with similar strategies
2. **Selection**: Interactive TUI requires coordination logic
3. **State management**: GlobalConfig tracking spans multiple installations
4. **Consistency**: Similar to languages coordinator for maintainability

### Why Native Package Manager for All Databases?

1. **Service integration**: Native packages include service management
2. **System compatibility**: Platform-optimized builds and configurations
3. **Dependencies**: Automatically handles system library dependencies
4. **Updates**: Integrated with system package manager updates
5. **Simplicity**: Single installation strategy for all databases

### Why No Version Management?

1. **Service stability**: Databases often require specific system-level integration
2. **Migration complexity**: Database version upgrades require data migration
3. **Platform packages**: Native packages provide tested, stable versions
4. **Development focus**: Development environments typically use single versions
5. **Container alternative**: Docker provides isolated version management

## External References

- **Languages Coordinator**: `internal/tooling/languages/languages.go`
- **GlobalConfig**: `docs/project-overview.md#configuration-management`
- **Testing Patterns**: `docs/guides/testing-patterns.md`
- **Terminal Coordinator**: `internal/tooling/terminal/terminal.go`

## Key Features

### Pre-Installation Detection
- Automatically detects databases installed before devgita
- Tracks pre-existing installations in GlobalConfig
- Prevents reinstalling already-present databases
- Runs detection on module initialization (New())

### Interactive Selection
- Multi-select TUI for choosing databases
- Filters already-installed databases from options
- Clear warnings about skipped databases
- Control options (All, None, Done)

### State Management
- Tracks devgita-installed vs pre-existing databases
- Maintains GlobalConfig for safe uninstall (future)
- Non-fatal error handling for resilient installation

### Simple Installation
- Native package manager for all databases
- Consistent installation strategy
- Platform-appropriate package management

This coordinator provides a robust, extensible foundation for managing database installations within the devgita development environment ecosystem.
