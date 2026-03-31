#!/usr/bin/env bash
set -e

# install.sh - Zero-dependency devgita installer
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/cjairm/devgita/main/install.sh | bash
#   bash install.sh --local /path/to/devgita-binary

REPO="cjairm/devgita"
INSTALL_DIR="$HOME/.local/bin"
BINARY_NAME="devgita"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_error() {
    echo -e "${RED}Error: $1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}$1${NC}"
}

print_info() {
    echo -e "${YELLOW}$1${NC}"
}

# Parse command-line arguments
LOCAL_BINARY=""
while [[ $# -gt 0 ]]; do
    case $1 in
        --local)
            LOCAL_BINARY="$2"
            shift 2
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Usage: $0 [--local /path/to/binary]"
            exit 1
            ;;
    esac
done

# Detect OS
OS=$(uname -s)
case "$OS" in
    Darwin)
        OS_NAME="darwin"
        ;;
    Linux)
        OS_NAME="linux"
        ;;
    *)
        print_error "Unsupported operating system: $OS"
        echo "Devgita only supports macOS (Darwin) and Linux (Debian/Ubuntu)."
        exit 1
        ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)
        ARCH_NAME="amd64"
        ;;
    amd64)
        ARCH_NAME="amd64"
        ;;
    arm64)
        ARCH_NAME="arm64"
        ;;
    aarch64)
        ARCH_NAME="arm64"
        ;;
    *)
        print_error "Unsupported architecture: $ARCH"
        echo "Devgita only supports amd64 (x86_64) and arm64 (aarch64) architectures."
        exit 1
        ;;
esac

BINARY_FILENAME="devgita-${OS_NAME}-${ARCH_NAME}"

print_info "Installing devgita for ${OS_NAME}/${ARCH_NAME}..."

# Create install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Install binary
DEST_PATH="$INSTALL_DIR/$BINARY_NAME"

if [ -n "$LOCAL_BINARY" ]; then
    # Local installation mode
    print_info "Installing from local file: $LOCAL_BINARY"
    
    if [ ! -f "$LOCAL_BINARY" ]; then
        print_error "Local binary not found: $LOCAL_BINARY"
        exit 1
    fi
    
    cp "$LOCAL_BINARY" "$DEST_PATH"
    chmod +x "$DEST_PATH"
    
    print_success "Copied local binary to $DEST_PATH"
else
    # Download from GitHub
    print_info "Fetching latest release from GitHub..."
    
    # Get latest release tag
    LATEST_RELEASE=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$LATEST_RELEASE" ]; then
        print_error "Failed to fetch latest release tag from GitHub"
        exit 1
    fi
    
    print_info "Latest release: $LATEST_RELEASE"
    
    # Construct download URL
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/$BINARY_FILENAME"
    
    print_info "Downloading from: $DOWNLOAD_URL"
    
    # Download binary
    if ! curl -fsSL -o "$DEST_PATH" "$DOWNLOAD_URL"; then
        print_error "Failed to download binary from GitHub"
        print_error "URL: $DOWNLOAD_URL"
        exit 1
    fi
    
    chmod +x "$DEST_PATH"
    
    print_success "Downloaded and installed to $DEST_PATH"
fi

# Detect shell and shell config file
SHELL_CONFIG=""
CURRENT_SHELL=$(basename "$SHELL")

case "$CURRENT_SHELL" in
    zsh)
        SHELL_CONFIG="$HOME/.zshrc"
        ;;
    bash)
        # Check for .bash_profile first (macOS), then .bashrc (Linux)
        if [ -f "$HOME/.bash_profile" ]; then
            SHELL_CONFIG="$HOME/.bash_profile"
        else
            SHELL_CONFIG="$HOME/.bashrc"
        fi
        ;;
    *)
        print_error "Unsupported shell: $CURRENT_SHELL"
        echo "Please manually add $INSTALL_DIR to your PATH"
        exit 1
        ;;
esac

# Add to PATH if not already present
PATH_EXPORT="export PATH=\"\$HOME/.local/bin:\$PATH\""

if [ -f "$SHELL_CONFIG" ]; then
    # Check if PATH entry already exists
    if grep -qF "$INSTALL_DIR" "$SHELL_CONFIG" 2>/dev/null; then
        print_info "PATH already configured in $SHELL_CONFIG"
    else
        print_info "Adding $INSTALL_DIR to PATH in $SHELL_CONFIG"
        echo "" >> "$SHELL_CONFIG"
        echo "# Added by devgita installer" >> "$SHELL_CONFIG"
        echo "$PATH_EXPORT" >> "$SHELL_CONFIG"
        print_success "Updated $SHELL_CONFIG"
    fi
else
    # Create shell config if it doesn't exist
    print_info "Creating $SHELL_CONFIG"
    echo "# Added by devgita installer" > "$SHELL_CONFIG"
    echo "$PATH_EXPORT" >> "$SHELL_CONFIG"
    print_success "Created $SHELL_CONFIG with PATH configuration"
fi

# Verify installation
print_info "Verifying installation..."

# Add to current PATH for verification
export PATH="$INSTALL_DIR:$PATH"

if command -v devgita &> /dev/null; then
    INSTALLED_VERSION=$(devgita --version 2>/dev/null || echo "unknown")
    print_success "✓ devgita installed successfully!"
    print_success "  Version: $INSTALLED_VERSION"
    echo ""
    print_info "Next steps:"
    echo "  1. Restart your shell or run: source $SHELL_CONFIG"
    echo "  2. Run: devgita install"
    echo ""
else
    print_error "Installation verification failed"
    echo "Binary installed to: $DEST_PATH"
    echo "Please restart your shell and try again"
    exit 1
fi
