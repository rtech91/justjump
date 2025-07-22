#!/usr/bin/env bash
set -e

# JustJump installation script
#
# The script requires:
# - curl to download the .justjumprc file
# - Go compiler to install JustJump from GitHub
# - Appropriate permissions to write to /usr/local/bin (root or sudo privileges may be required)

# Check if Go is installed
if ! command -v go &>/dev/null; then
    echo "Go compiler not found. Please install Go first."
    exit 1
fi

echo "Installing JustJump from git..."
GOBIN=/usr/local/bin go install github.com/rtech91/justjump@latest

# Detect shell and download the corresponding .justjumprc file
detect_shell() {
    # Prefer current process if possible
    if [ -n "$BASH_VERSION" ]; then
        echo "bash"
    elif [ -n "$ZSH_VERSION" ]; then
        echo "zsh"
    elif [ -n "$SHELL" ]; then
        basename "$SHELL"
    else
        ps -p $$ -o comm= | awk -F'-' '{print $NF}'
    fi
}

SHELL_NAME=$(detect_shell)
if [ "$SHELL_NAME" = "bash" ]; then
    echo "Detected bash shell."
    curl -sSf -o "$HOME/.justjumprc" https://raw.githubusercontent.com/rtech91/justjump/refs/heads/main/misc/rc/bash/.justjumprc
    RC_FILE="$HOME/.bashrc"
elif [ "$SHELL_NAME" = "zsh" ]; then
    echo "Detected zsh shell."
    curl -sSf -o "$HOME/.justjumprc" https://raw.githubusercontent.com/rtech91/justjump/refs/heads/main/misc/rc/zsh/.justjumprc
    RC_FILE="$HOME/.zshrc"
else
    echo "Could not reliably detect your shell (detected: '$SHELL_NAME')."
    echo "Please manually configure your shell integration."
    echo "Supported shells: bash, zsh."
    exit 1
fi

# Append the integration line if not already present
if ! grep -qF '[ -f "$HOME/.justjumprc" ] && source "$HOME/.justjumprc"' "$RC_FILE"; then
    echo "Adding JustJump integration to $RC_FILE..."
    echo '[ -f "$HOME/.justjumprc" ] && source "$HOME/.justjumprc"' >> "$RC_FILE"
else
    echo "Shell integration already exists in $RC_FILE."
fi

echo "Installation complete. Please reload your shell (or run: source $RC_FILE) to apply changes."
