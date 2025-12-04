# Bison Module Documentation

## Overview

The Bison module provides installation and command execution management for GNU Bison with devgita integration. It follows the standardized devgita app interface while providing bison-specific operations for parser generation, grammar processing, and language implementation support.

## App Purpose

GNU Bison is a general-purpose parser generator that converts an annotated context-free grammar into a deterministic LR or generalized LR (GLR) parser employing LALR(1) parser tables. Once you are proficient with Bison, you can use it to develop a wide range of language parsers, from those used in simple desk calculators to complex programming languages. This module ensures bison is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for generating parsers from grammar specifications.

## Lifecycle Summary

1. **Installation**: Install bison package via platform package managers (Homebrew/apt)
2. **Configuration**: bison typically doesn't require separate configuration files - operations are handled via grammar files (.y) in project directories
3. **Execution**: Provide high-level bison operations for generating parsers and processing grammars

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Bison instance with platform-specific commands           |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install bison                             |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute bison             | Runs bison with provided arguments                                   |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
bison := bison.New()
err := bison.Install()
```

- **Purpose**: Standard bison installation
- **Behavior**: Uses `InstallPackage()` to install bison package
- **Use case**: Initial bison installation or explicit reinstall

### ForceInstall()

```go
bison := bison.New()
err := bison.ForceInstall()
```

- **Purpose**: Force bison installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh bison installation or fix corrupted installation

### SoftInstall()

```go
bison := bison.New()
err := bison.SoftInstall()
```

- **Purpose**: Install bison only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing bison installations

### Uninstall()

```go
err := bison.Uninstall()
```

- **Purpose**: Remove bison installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: bison is a system dependency managed at the OS level

### Update()

```go
err := bison.Update()
```

- **Purpose**: Update bison installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: bison updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := bison.ForceConfigure()
err := bison.SoftConfigure()
```

- **Purpose**: Apply bison configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: bison doesn't use traditional config files; configuration is handled via grammar files (.y) in project directories

## Execution Methods

### ExecuteCommand()

```go
err := bison.ExecuteCommand("--version")
err := bison.ExecuteCommand("parser.y")
err := bison.ExecuteCommand("-d", "parser.y") // Generate header file
```

- **Purpose**: Execute bison commands with provided arguments
- **Parameters**: Variable arguments passed directly to bison binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Bison-Specific Operations

The bison CLI provides extensive parser generation capabilities:

#### Version and Help

```bash
# Show bison version
bison --version

# Show help information
bison --help

# Show usage information
bison --usage
```

#### Generate Parsers

```bash
# Generate parser from grammar file
bison parser.y

# Generate parser with header file
bison -d parser.y
bison --defines parser.y

# Specify output file name
bison -o output.c parser.y
bison --output=output.c parser.y

# Generate verbose report
bison -v parser.y
bison --verbose parser.y
```

#### Debug and Analysis

```bash
# Enable debug mode
bison --debug parser.y

# Generate graphical representation
bison --graph parser.y

# Generate XML report
bison --xml parser.y

# Report grammar conflicts
bison -r all parser.y
bison --report=all parser.y
```

#### Advanced Options

```bash
# Generate GLR parser
bison --glr parser.y

# Generate C++ parser
bison -o parser.cpp parser.yy

# Specify language
bison --language=c++ parser.y
bison --language=java parser.y

# Define file locations
bison --defines=parser.h --output=parser.c parser.y

# Specify skeleton file
bison --skeleton=lalr1.cc parser.y
```

#### Error Handling

```bash
# Enable warnings
bison -W parser.y
bison --warnings=all parser.y

# Treat warnings as errors
bison -Werror parser.y

# Show specific warnings
bison -Wconflicts-sr parser.y
bison -Wconflicts-rr parser.y
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Parser Generation**: `New()` → `SoftInstall()` → `ExecuteCommand()` to generate parser
4. **Version Checks**: `New()` → `ExecuteCommand("--version")`

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.Bison` (typically "bison")
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No traditional config files**: bison operations are configured via grammar files (.y or .yy) in project directories
- **Runtime configuration**: Parameters passed directly to `ExecuteCommand()`
- **Project-specific**: Each project has its own grammar files defining language syntax
- **Grammar language**: Uses Bison's own grammar specification language

## Implementation Notes

- **Parser Generator Nature**: Unlike typical applications, bison is a parser generator tool without persistent configuration
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since bison uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since bison doesn't use config files
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as bison updates should be handled by system package managers

## Usage Examples

### Basic Parser Generation

```go
bison := bison.New()

// Install bison
err := bison.SoftInstall()
if err != nil {
    return err
}

// Check version
err = bison.ExecuteCommand("--version")

// Generate parser from grammar file
err = bison.ExecuteCommand("parser.y")

// Generate with header file
err = bison.ExecuteCommand("-d", "parser.y")
```

### Advanced Operations

```go
// Generate C++ parser
err := bison.ExecuteCommand("-o", "parser.cpp", "grammar.yy")

// Verbose output with report
err = bison.ExecuteCommand("-v", "-r", "all", "parser.y")

// GLR parser generation
err = bison.ExecuteCommand("--glr", "parser.y")

// Debug mode
err = bison.ExecuteCommand("--debug", "parser.y")

// Generate graph visualization
err = bison.ExecuteCommand("--graph", "parser.y")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Grammar File Not Found**: Run bison in directory containing .y grammar file
3. **Shift/Reduce Conflicts**: Use `-v` flag to generate report and analyze conflicts
4. **Reduce/Reduce Conflicts**: Review grammar rules for ambiguities
5. **Header File Issues**: Use `-d` or `--defines` flag to generate header

### Platform Considerations

- **macOS**: Installed via Homebrew package manager, may include flex for lexer generation
- **Linux**: Installed via apt package manager, often pre-installed on development systems
- **Dependencies**: Often used with flex (lexical analyzer generator)
- **Related Tools**: Part of compiler toolchain with flex, gcc/clang

### Best Practices

- **Version Control**: Commit grammar files (.y) but not generated parsers
- **Regeneration**: Run bison when grammar changes
- **Conflict Resolution**: Address all shift/reduce and reduce/reduce conflicts
- **Documentation**: Document grammar structure and parser interface
- **Testing**: Test generated parser with representative input
- **Language Choice**: Use appropriate skeleton for target language (C, C++, Java)

### Parser Development Workflow

1. **Define Grammar**: Create .y file with grammar rules
2. **Run Bison**: Generate parser from grammar file
3. **Integrate Lexer**: Combine with flex-generated lexer
4. **Compile**: Build parser with compiler
5. **Test**: Validate parser behavior with test inputs

```bash
# Typical bison/flex workflow
flex lexer.l              # Generate lexer (lex.yy.c)
bison -d parser.y         # Generate parser (parser.tab.c, parser.tab.h)
gcc -o myparser lex.yy.c parser.tab.c -lfl
./myparser < input.txt    # Run parser
```

## Integration with Devgita

bison integrates with devgita's terminal category:

- **Installation**: Installed as part of core terminal tools setup
- **Configuration**: No configuration files - uses project-specific grammar files
- **Usage**: Available system-wide after installation for parser generation
- **Updates**: Managed through system package manager
- **Dependencies**: Often paired with flex for complete lexer/parser generation

## Grammar File Structure

Bison grammar files (.y) typically contain:

### Declarations Section

```yacc
%{
#include <stdio.h>
int yylex(void);
void yyerror(char *s);
%}

%token NUMBER
%token PLUS MINUS MULTIPLY DIVIDE
%token LPAREN RPAREN
```

### Grammar Rules Section

```yacc
%%
expression: term
    | expression PLUS term
    | expression MINUS term
    ;

term: factor
    | term MULTIPLY factor
    | term DIVIDE factor
    ;

factor: NUMBER
    | LPAREN expression RPAREN
    ;
%%
```

### Additional C Code Section

```c
void yyerror(char *s) {
    fprintf(stderr, "Error: %s\n", s);
}

int main(void) {
    return yyparse();
}
```

## External References

- **GNU Bison Manual**: https://www.gnu.org/software/bison/manual/
- **Bison Documentation**: https://www.gnu.org/software/bison/
- **Flex and Bison**: https://web.iitd.ac.in/~sumeet/flex__bison.pdf
- **Parsing Techniques**: https://dickgrune.com/Books/PTAPG_2nd_Edition/
- **Compiler Construction**: https://www.gnu.org/software/bison/manual/html_node/index.html

This module provides essential parser generation support for building compilers, interpreters, and language processors within the devgita ecosystem.
