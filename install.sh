#!/bin/bash

# GitHub Tools Installation Script
# This script installs ghtools to ~/scripts

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

echo -e "${BLUE}╔════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     GitHub Tools Installation Script      ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════╝${NC}"
echo ""

# Check if running as root
if [ "$EUID" -eq 0 ]; then
    print_warning "Please do not run this script as root (with sudo)"
    echo "The script will ask for sudo password when needed."
    exit 1
fi

# Get the directory where the install script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GHTOOLS_SCRIPT="$SCRIPT_DIR/ghtools"

# Check if script exists
if [ ! -f "$GHTOOLS_SCRIPT" ]; then
    print_error "ghtools script not found at: $GHTOOLS_SCRIPT"
    exit 1
fi

print_info "Found ghtools script at: $GHTOOLS_SCRIPT"

# Check dependencies
print_info "Checking dependencies..."

missing_deps=()

if ! command -v gh &> /dev/null; then
    missing_deps+=("gh (GitHub CLI)")
fi

if ! command -v fzf &> /dev/null; then
    missing_deps+=("fzf")
fi

if [ ${#missing_deps[@]} -ne 0 ]; then
    print_error "Missing required dependencies:"
    for dep in "${missing_deps[@]}"; do
        echo "  - $dep"
    done
    echo ""
    echo "Install them using:"
    echo "  ${BLUE}sudo pacman -S github-cli fzf${NC}"
    echo "  or"
    echo "  ${BLUE}yay -S github-cli fzf${NC}"
    exit 1
fi

print_success "All dependencies are installed"

# Create ~/scripts directory if it doesn't exist
INSTALL_DIR="$HOME/scripts"
mkdir -p "$INSTALL_DIR"

# Install script to ~/scripts
print_info "Installing script to $INSTALL_DIR..."

cp "$GHTOOLS_SCRIPT" "$INSTALL_DIR/ghtools"
chmod +x "$INSTALL_DIR/ghtools"

print_success "ghtools installed successfully!"

# Check if ~/scripts is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    print_warning "$INSTALL_DIR is not in your PATH"
    echo ""
    echo "Add the following line to your ~/.zshrc:"
    echo "  ${BLUE}export PATH=\"\$HOME/scripts:\$PATH\"${NC}"
    echo ""
    echo "Then run:"
    echo "  ${BLUE}source ~/.zshrc${NC}"
    echo ""

    # Offer to add it automatically
    read -p "Would you like to add it automatically? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        if [ -f "$HOME/.zshrc" ]; then
            # Check if it's already there but commented or something
            if grep -q "export PATH=.*scripts.*PATH" "$HOME/.zshrc"; then
                print_info "PATH configuration already exists in ~/.zshrc"
            else
                echo "" >> "$HOME/.zshrc"
                echo "# Add ~/scripts to PATH for custom scripts" >> "$HOME/.zshrc"
                echo "export PATH=\"\$HOME/scripts:\$PATH\"" >> "$HOME/.zshrc"
                print_success "Added to ~/.zshrc"
                echo "Run: ${GREEN}source ~/.zshrc${NC}"
            fi
        else
            print_error "~/.zshrc not found"
        fi
    fi
else
    print_success "Script is now available in your PATH"
    echo ""
    echo "Run ${GREEN}ghtools clone${NC} or ${GREEN}ghtools delete${NC} to start using it!"
fi

echo ""
print_info "Usage: ghtools <command>"
echo "  ${GREEN}ghtools clone${NC}  - Clone repositories interactively"
echo "  ${GREEN}ghtools delete${NC} - Delete repositories interactively"
echo "  ${GREEN}ghtools help${NC}   - Show help message"
echo ""
print_info "To uninstall, run: rm ~/scripts/ghtools"
