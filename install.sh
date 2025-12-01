#!/bin/bash

# GitHub Tools Installation Script
# This script installs ghtools to ~/scripts

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
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

if ! command -v git &> /dev/null; then
    missing_deps+=("git")
fi

if ! command -v jq &> /dev/null; then
    missing_deps+=("jq")
fi

# Check for optional but recommended dependencies
optional_deps=()
if ! command -v gum &> /dev/null; then
    optional_deps+=("gum")
fi

if [ ${#missing_deps[@]} -ne 0 ]; then
    print_error "Missing required dependencies:"
    for dep in "${missing_deps[@]}"; do
        echo "  - $dep"
    done
    echo ""
    echo "Install them using:"
    echo -e "  ${BLUE}sudo pacman -S github-cli fzf git jq gum${NC}"
    echo "  or"
    echo -e "  ${BLUE}yay -S github-cli fzf git jq gum${NC}"
    exit 1
fi

print_success "All required dependencies are installed"

# Notify about optional dependencies
if [ ${#optional_deps[@]} -ne 0 ]; then
    echo ""
    print_warning "Optional dependencies for enhanced UI:"
    for dep in "${optional_deps[@]}"; do
        echo -e "  ${YELLOW}•${NC} $dep - Beautiful terminal UI components"
    done
    echo ""
    echo -e "Install for best experience: ${CYAN}sudo pacman -S gum${NC}"
    echo ""
fi

# Create ~/scripts directory if it doesn't exist
INSTALL_DIR="$HOME/scripts"
mkdir -p "$INSTALL_DIR"

# Install script to ~/scripts
print_info "Installing script to $INSTALL_DIR..."

cp "$GHTOOLS_SCRIPT" "$INSTALL_DIR/ghtools"
chmod +x "$INSTALL_DIR/ghtools"

print_success "ghtools installed successfully!"

# ============================================================================
# PATH Configuration Functions
# ============================================================================

# Get all possible zsh config files
get_zsh_config_files() {
    local files=()

    # Main zsh config files
    [ -f "$HOME/.zshrc" ] && files+=("$HOME/.zshrc")
    [ -f "$HOME/.zshrc_custom" ] && files+=("$HOME/.zshrc_custom")
    [ -f "$HOME/.zshenv" ] && files+=("$HOME/.zshenv")
    [ -f "$HOME/.zprofile" ] && files+=("$HOME/.zprofile")
    [ -f "$HOME/.zlogin" ] && files+=("$HOME/.zlogin")

    # Modular config files
    if [ -d "$HOME/.config/zshrc" ]; then
        for config in "$HOME/.config/zshrc/"*; do
            [ -f "$config" ] && files+=("$config")
        done
    fi

    printf '%s\n' "${files[@]}"
}

# Check if scripts directory is already in PATH configuration files
# Returns: array of files containing the PATH
check_scripts_in_config_files() {
    local files_with_path=()
    local pattern="(export PATH=.*scripts|PATH=.*scripts)"

    while IFS= read -r file; do
        if grep -qE "$pattern" "$file" 2>/dev/null; then
            files_with_path+=("$file")
        fi
    done < <(get_zsh_config_files)

    printf '%s\n' "${files_with_path[@]}"
}

# Remove PATH configuration from a file
remove_path_from_file() {
    local file=$1
    local backup="${file}.backup_$(date +%s)"

    # Create backup
    cp "$file" "$backup"

    # Remove lines with scripts PATH
    sed -i '/export PATH=.*scripts/d' "$file"
    sed -i '/PATH=.*scripts/d' "$file"

    # Also remove the comment line if it exists
    sed -i '/# Add ~\/scripts to PATH/d' "$file"

    print_success "Removed PATH configuration from: $file"
    print_info "Backup created at: $backup"
}

# Add PATH to selected file
add_path_to_file() {
    local file=$1

    echo "" >> "$file"
    echo "# Add ~/scripts to PATH for custom scripts" >> "$file"
    echo "export PATH=\"\$HOME/scripts:\$PATH\"" >> "$file"

    print_success "Added PATH to: $file"
}

# ============================================================================
# PATH Configuration Management
# ============================================================================

echo ""
echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo -e "${BLUE}PATH Configuration Check${NC}"
echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo ""

# Check if scripts directory is in any config files
files_with_path=()
mapfile -t files_with_path < <(check_scripts_in_config_files)

# Check if currently in PATH
if [[ ":$PATH:" == *":$INSTALL_DIR:"* ]]; then
    print_success "$INSTALL_DIR is already in your PATH"

    # Check for duplicates in config files
    if [ ${#files_with_path[@]} -gt 1 ]; then
        echo ""
        print_warning "DUPLICATE DETECTED! PATH is configured in multiple files:"
        echo ""
        for file in "${files_with_path[@]}"; do
            line_info=$(grep -n "PATH=.*scripts" "$file" 2>/dev/null | head -1)
            echo -e "  ${YELLOW}→${NC} $file"
            echo -e "    Line: ${CYAN}$line_info${NC}"
        done
        echo ""
        echo "Having duplicates can cause PATH pollution and unexpected behavior."
        echo ""

        read -p "Would you like to remove duplicates? (y/N): " -r
        echo ""

        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo "Select which file to KEEP the PATH configuration:"
            echo ""
            i=1
            for file in "${files_with_path[@]}"; do
                echo -e "  ${GREEN}$i${NC}) $file"
                ((i++))
            done
            echo -e "  ${RED}0${NC}) Cancel"
            echo ""

            read -p "Enter number: " -r file_choice

            if [[ "$file_choice" =~ ^[0-9]+$ ]] && [ "$file_choice" -ge 1 ] && [ "$file_choice" -le ${#files_with_path[@]} ]; then
                keep_file="${files_with_path[$((file_choice-1))]}"
                echo ""
                print_info "Keeping PATH configuration in: $keep_file"
                echo ""

                # Remove from all other files
                for file in "${files_with_path[@]}"; do
                    if [ "$file" != "$keep_file" ]; then
                        remove_path_from_file "$file"
                    fi
                done

                echo ""
                print_success "Duplicates removed successfully!"
                print_info "Reload your shell: ${GREEN}source ~/.zshrc${NC}"
            else
                print_info "Cancelled. No changes made."
            fi
        fi
    elif [ ${#files_with_path[@]} -eq 1 ]; then
        print_info "PATH is configured in: ${files_with_path[0]}"
    fi
else
    # Not in current PATH
    if [ ${#files_with_path[@]} -gt 0 ]; then
        print_warning "$INSTALL_DIR is configured but not in current PATH"
        echo ""
        echo "PATH is configured in the following file(s):"
        for file in "${files_with_path[@]}"; do
            echo "  → $file"
        done
        echo ""
        print_info "Reload your shell: ${GREEN}source ~/.zshrc${NC}"

        # Check for duplicates
        if [ ${#files_with_path[@]} -gt 1 ]; then
            echo ""
            print_warning "Multiple PATH configurations detected (duplicates)"
            read -p "Would you like to remove duplicates? (y/N): " -r
            echo ""

            if [[ $REPLY =~ ^[Yy]$ ]]; then
                echo "Select which file to KEEP:"
                echo ""
                i=1
                for file in "${files_with_path[@]}"; do
                    echo -e "  ${GREEN}$i${NC}) $file"
                    ((i++))
                done
                echo ""

                read -p "Enter number: " -r file_choice

                if [[ "$file_choice" =~ ^[0-9]+$ ]] && [ "$file_choice" -ge 1 ] && [ "$file_choice" -le ${#files_with_path[@]} ]; then
                    keep_file="${files_with_path[$((file_choice-1))]}"
                    echo ""

                    for file in "${files_with_path[@]}"; do
                        if [ "$file" != "$keep_file" ]; then
                            remove_path_from_file "$file"
                        fi
                    done

                    print_success "Duplicates removed!"
                fi
            fi
        fi
    else
        # Not configured anywhere
        print_warning "$INSTALL_DIR is not in your PATH"
        echo ""
        echo "To use ghtools from anywhere, add it to your PATH."
        echo ""

        read -p "Would you like to add it now? (y/N): " -r
        echo ""

        if [[ $REPLY =~ ^[Yy]$ ]]; then
            # Offer file selection
            echo "Where would you like to add the PATH?"
            echo ""
            echo -e "  ${GREEN}1${NC}) ~/.zshrc ${CYAN}(recommended for general use)${NC}"
            echo -e "  ${GREEN}2${NC}) ~/.zshrc_custom ${CYAN}(for user customizations)${NC}"

            if [ -d "$HOME/.config/zshrc" ]; then
                echo -e "  ${GREEN}3${NC}) ~/.config/zshrc/00-init ${CYAN}(modular config)${NC}"
            fi

            echo -e "  ${RED}0${NC}) Skip"
            echo ""

            read -p "Enter number: " -r choice
            echo ""

            case "$choice" in
                1)
                    if [ -f "$HOME/.zshrc" ]; then
                        add_path_to_file "$HOME/.zshrc"
                        print_info "Reload your shell: ${GREEN}source ~/.zshrc${NC}"
                    else
                        print_error "~/.zshrc not found"
                    fi
                    ;;
                2)
                    if [ -f "$HOME/.zshrc_custom" ]; then
                        add_path_to_file "$HOME/.zshrc_custom"
                        print_info "Reload your shell: ${GREEN}source ~/.zshrc${NC}"
                    else
                        print_error "~/.zshrc_custom not found"
                    fi
                    ;;
                3)
                    if [ -d "$HOME/.config/zshrc" ]; then
                        add_path_to_file "$HOME/.config/zshrc/00-init"
                        print_info "Reload your shell: ${GREEN}source ~/.zshrc${NC}"
                    else
                        print_error "~/.config/zshrc not found"
                    fi
                    ;;
                *)
                    print_info "Skipped. You can add it manually later."
                    echo ""
                    echo "Add this line to your preferred zsh config file:"
                    echo -e "  ${BLUE}export PATH=\"\$HOME/scripts:\$PATH\"${NC}"
                    ;;
            esac
        else
            print_info "Skipped. Add this line to your zsh config:"
            echo -e "  ${BLUE}export PATH=\"\$HOME/scripts:\$PATH\"${NC}"
        fi
    fi
fi

echo ""
print_success "Installation complete!"
echo ""
echo -e "Run ${GREEN}ghtools${NC} to start using the interactive menu!"

echo ""
echo -e "${CYAN}═══════════════════════════════════════════${NC}"
print_info "Usage:"
echo ""
echo -e "  ${GREEN}ghtools${NC}        - Interactive menu (recommended)"
echo -e "  ${GREEN}ghtools clone${NC}  - Clone repositories"
echo -e "  ${GREEN}ghtools delete${NC} - Delete repositories"
echo -e "  ${GREEN}ghtools help${NC}   - Show help message"
echo ""
print_info "To uninstall, run: rm ~/scripts/ghtools"
