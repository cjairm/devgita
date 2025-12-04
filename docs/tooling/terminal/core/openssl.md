# OpenSSL Module Documentation

## Overview

The OpenSSL module provides installation and command execution management for OpenSSL with devgita integration. It follows the standardized devgita app interface while providing openssl-specific operations for SSL/TLS protocols, cryptographic functions, certificate management, and secure communications.

## App Purpose

OpenSSL is a robust, commercial-grade, and full-featured toolkit for the Transport Layer Security (TLS) and Secure Sockets Layer (SSL) protocols. It is also a general-purpose cryptography library that provides implementations of various cryptographic algorithms and protocols. OpenSSL is widely used in Internet servers, including most HTTPS websites, and provides essential tools for generating keys, creating certificate requests, signing certificates, and performing encryption/decryption operations. This module ensures openssl is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for cryptographic operations and certificate management.

## Lifecycle Summary

1. **Installation**: Install openssl package via platform package managers (Homebrew/apt)
2. **Configuration**: openssl typically doesn't require separate configuration files - operations are handled via command-line arguments or system-level openssl.cnf files
3. **Execution**: Provide high-level openssl operations for cryptographic tasks, certificate operations, and SSL/TLS functions

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new OpenSSL instance with platform-specific commands         |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install openssl                           |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute openssl           | Runs openssl with provided arguments                                 |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
openssl := openssl.New()
err := openssl.Install()
```

- **Purpose**: Standard openssl installation
- **Behavior**: Uses `InstallPackage()` to install openssl package
- **Use case**: Initial openssl installation or explicit reinstall

### ForceInstall()

```go
openssl := openssl.New()
err := openssl.ForceInstall()
```

- **Purpose**: Force openssl installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh openssl installation or fix corrupted installation

### SoftInstall()

```go
openssl := openssl.New()
err := openssl.SoftInstall()
```

- **Purpose**: Install openssl only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing openssl installations

### Uninstall()

```go
err := openssl.Uninstall()
```

- **Purpose**: Remove openssl installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: openssl is a critical system dependency managed at the OS level

### Update()

```go
err := openssl.Update()
```

- **Purpose**: Update openssl installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: openssl updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := openssl.ForceConfigure()
err := openssl.SoftConfigure()
```

- **Purpose**: Apply openssl configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: openssl doesn't require separate configuration files in devgita; configuration is handled via command-line arguments or system-level openssl.cnf files managed by the OS

## Execution Methods

### ExecuteCommand()

```go
err := openssl.ExecuteCommand("version")
err := openssl.ExecuteCommand("genrsa", "-out", "private.key", "2048")
err := openssl.ExecuteCommand("req", "-new", "-key", "private.key")
```

- **Purpose**: Execute openssl commands with provided arguments
- **Parameters**: Variable arguments passed directly to openssl binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### OpenSSL-Specific Operations

The openssl CLI provides extensive cryptographic and certificate management capabilities:

#### Version and Information

```bash
# Show openssl version
openssl version

# Show all version information
openssl version -a

# Show configuration directory
openssl version -d

# List available commands
openssl list -commands

# List available ciphers
openssl list -cipher-commands

# List available digest algorithms
openssl list -digest-commands
```

#### Key Generation

```bash
# Generate RSA private key (2048 bits)
openssl genrsa -out private.key 2048

# Generate RSA private key (4096 bits)
openssl genrsa -out private.key 4096

# Generate encrypted RSA private key
openssl genrsa -aes256 -out private.key 2048

# Generate DSA parameters and key
openssl dsaparam -genkey -out dsaparam.pem 2048
openssl gendsa -out dsa-private.key dsaparam.pem

# Generate EC private key
openssl ecparam -genkey -name secp384r1 -out ec-private.key

# Extract public key from private key
openssl rsa -in private.key -pubout -out public.key
```

#### Certificate Requests (CSR)

```bash
# Create certificate signing request
openssl req -new -key private.key -out request.csr

# Create CSR with subject information
openssl req -new -key private.key -out request.csr \
  -subj "/C=US/ST=State/L=City/O=Organization/CN=example.com"

# Create self-signed certificate
openssl req -new -x509 -key private.key -out certificate.crt -days 365

# Generate private key and CSR in one command
openssl req -newkey rsa:2048 -nodes -keyout private.key -out request.csr
```

#### Certificate Operations

```bash
# View certificate details
openssl x509 -in certificate.crt -text -noout

# View certificate dates
openssl x509 -in certificate.crt -dates -noout

# View certificate subject
openssl x509 -in certificate.crt -subject -noout

# Verify certificate
openssl verify certificate.crt

# Convert certificate formats
openssl x509 -in certificate.crt -out certificate.pem -outform PEM
openssl x509 -in certificate.pem -out certificate.der -outform DER

# Sign a certificate request
openssl x509 -req -in request.csr -CA ca-cert.pem -CAkey ca-key.pem \
  -CAcreateserial -out certificate.crt -days 365
```

#### Encryption and Decryption

```bash
# Encrypt a file with AES-256
openssl enc -aes-256-cbc -salt -in plaintext.txt -out encrypted.txt

# Decrypt a file
openssl enc -aes-256-cbc -d -in encrypted.txt -out decrypted.txt

# Encrypt with base64 encoding
openssl enc -aes-256-cbc -a -salt -in plaintext.txt -out encrypted.txt

# List available ciphers
openssl enc -list
```

#### Message Digests (Hashing)

```bash
# Calculate SHA-256 hash
openssl dgst -sha256 file.txt

# Calculate MD5 hash
openssl dgst -md5 file.txt

# Calculate SHA-512 hash
openssl dgst -sha512 file.txt

# Sign a file with private key
openssl dgst -sha256 -sign private.key -out signature.sig file.txt

# Verify signature
openssl dgst -sha256 -verify public.key -signature signature.sig file.txt
```

#### SSL/TLS Testing

```bash
# Test SSL/TLS connection
openssl s_client -connect example.com:443

# Show certificate chain
openssl s_client -connect example.com:443 -showcerts

# Test specific TLS version
openssl s_client -connect example.com:443 -tls1_2

# Test with SNI (Server Name Indication)
openssl s_client -connect example.com:443 -servername example.com

# Check certificate expiration
echo | openssl s_client -connect example.com:443 2>/dev/null | \
  openssl x509 -noout -dates
```

#### Random Data Generation

```bash
# Generate random bytes
openssl rand 32

# Generate random hex string
openssl rand -hex 32

# Generate random base64 string
openssl rand -base64 32

# Generate random bytes to file
openssl rand -out random.bin 1024
```

#### PKCS#12 Operations

```bash
# Create PKCS#12 file (.pfx/.p12)
openssl pkcs12 -export -out certificate.pfx \
  -inkey private.key -in certificate.crt -certfile ca-cert.pem

# Extract private key from PKCS#12
openssl pkcs12 -in certificate.pfx -nocerts -out private.key

# Extract certificate from PKCS#12
openssl pkcs12 -in certificate.pfx -clcerts -nokeys -out certificate.crt

# View PKCS#12 contents
openssl pkcs12 -in certificate.pfx -info
```

#### Password Operations

```bash
# Generate password hash (bcrypt)
openssl passwd -6 -salt saltstring mypassword

# Generate MD5 password hash
openssl passwd -1 mypassword

# Generate salted password
openssl passwd -salt ab mypassword
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **OpenSSL Operations**: `New()` → `SoftInstall()` → `ExecuteCommand()` with openssl arguments
4. **Version Check**: `New()` → `ExecuteCommand("version")`

## Constants and Paths

### Relevant Constants

- **Package name**: `"openssl"` used directly for installation
- Referenced via `constants.OpenSSL` in the codebase
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No traditional config files**: OpenSSL operations are configured via command-line arguments
- **System configuration**: OpenSSL uses system-level configuration files (e.g., `/etc/ssl/openssl.cnf` on Linux, `/usr/local/etc/openssl/openssl.cnf` on macOS)
- **Runtime configuration**: Parameters passed directly to `ExecuteCommand()`
- **Environment variables**: OpenSSL respects standard environment variables like `OPENSSL_CONF`

## Implementation Notes

- **Cryptographic Library Nature**: Unlike typical applications, openssl is a comprehensive cryptographic toolkit without persistent configuration
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since openssl uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since openssl doesn't require config files in devgita
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as openssl updates should be handled by system package managers
- **Security Critical**: OpenSSL is a critical security component; updates should be managed carefully through system package managers

## Usage Examples

### Basic OpenSSL Operations

```go
openssl := openssl.New()

// Install openssl
err := openssl.SoftInstall()
if err != nil {
    return err
}

// Check version
err = openssl.ExecuteCommand("version")

// Generate RSA private key
err = openssl.ExecuteCommand("genrsa", "-out", "private.key", "2048")

// Create certificate request
err = openssl.ExecuteCommand("req", "-new", "-key", "private.key", "-out", "request.csr")
```

### Advanced Operations

```go
// Generate self-signed certificate
err := openssl.ExecuteCommand("req", "-new", "-x509", "-key", "private.key",
    "-out", "certificate.crt", "-days", "365")

// Encrypt a file
err = openssl.ExecuteCommand("enc", "-aes-256-cbc", "-salt",
    "-in", "plaintext.txt", "-out", "encrypted.txt")

// Calculate SHA-256 hash
err = openssl.ExecuteCommand("dgst", "-sha256", "file.txt")

// Test SSL connection
err = openssl.ExecuteCommand("s_client", "-connect", "example.com:443")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Certificate Errors**: Verify certificate format and paths are correct
3. **Permission Issues**: Some operations may require proper file permissions
4. **Version Conflicts**: Multiple OpenSSL versions may be installed on the system
5. **Path Issues**: Ensure openssl binary is in PATH and accessible

### Platform Considerations

- **macOS**: Installed via Homebrew package manager, system openssl also available
- **Linux**: Installed via apt package manager, usually pre-installed
- **SSL/TLS Support**: Full support for modern TLS versions (1.2, 1.3)
- **Algorithm Support**: Comprehensive suite of cryptographic algorithms
- **FIPS Mode**: Some platforms support FIPS 140-2 validated cryptography

### Security Notes

- **Keep Updated**: OpenSSL security updates are critical; use system package manager for updates
- **Strong Keys**: Use at least 2048-bit RSA keys for production environments
- **Password Protection**: Protect private keys with strong passwords
- **Certificate Validation**: Always verify certificates in production environments
- **Secure Storage**: Store private keys securely with appropriate file permissions (e.g., 600)

### Common Use Cases

#### Development HTTPS Server

```bash
# Generate self-signed certificate for local development
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

#### Certificate Chain Verification

```bash
# Verify certificate against CA
openssl verify -CAfile ca-bundle.crt server-cert.crt
```

#### Convert Certificate Formats

```bash
# PEM to DER
openssl x509 -in cert.pem -outform DER -out cert.der

# DER to PEM
openssl x509 -in cert.der -inform DER -out cert.pem -outform PEM

# PFX to PEM
openssl pkcs12 -in certificate.pfx -out certificate.pem -nodes
```

## Integration with Devgita

OpenSSL integrates with devgita's terminal category:

- **Installation**: Installed as part of terminal tools setup (system libraries)
- **Configuration**: No default configuration applied by devgita
- **Usage**: Available system-wide after installation
- **Updates**: Managed through system package manager
- **Dependencies**: Required by many development tools and programming languages

## External References

- **OpenSSL Website**: https://www.openssl.org/
- **OpenSSL Documentation**: https://www.openssl.org/docs/
- **OpenSSL GitHub**: https://github.com/openssl/openssl
- **OpenSSL Wiki**: https://wiki.openssl.org/
- **OpenSSL Cookbook**: https://www.feistyduck.com/books/openssl-cookbook/
- **SSL/TLS Best Practices**: https://ssl-config.mozilla.org/

## Versioning Notes

- **OpenSSL 1.1.1**: LTS version, widely deployed
- **OpenSSL 3.0+**: Current stable version with improved API
- **LibreSSL**: Alternative fork used on some BSD systems
- **BoringSSL**: Google's fork, used in Chrome and Android

This module provides essential SSL/TLS and cryptographic capabilities for securing communications and handling certificates within the devgita development environment.
