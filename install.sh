#!/bin/bash

set -e 

APP_NAME="goskii"
REPO="JoelVCrasta/goskii"
BIN_DIR="/usr/local/bin"

# Determine OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

if [[ "$OS" != "linux" && "$OS" != "darwin" ]]; then
    echo "Unsupported OS: $OS"
    exit 1
fi

# Determine Architecture
ARCH=$(uname -m)
if [[ "$ARCH" != "x86_64" ]]; then
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

# Fetch Latest Version
echo "Fetching the latest version of $APP_NAME..."
LATEST_VERSION=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep -oP '"tag_name": "\K(.*)(?=")')

if [[ -z "$LATEST_VERSION" ]]; then
    echo "Failed to fetch the latest version of $APP_NAME."
    exit 1
fi

# Build Download URL
BINARY="${APP_NAME}-${OS}-amd64"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY}"

# Download and Install
echo "Downloading $APP_NAME from $DOWNLOAD_URL..."
sudo curl -sL $DOWNLOAD_URL -o $BIN_DIR/$APP_NAME

# Make Binary Executable
echo "Granting executable permissions to $APP_NAME..."
sudo chmod +x $BIN_DIR/$APP_NAME

# Verify Installation
if command -v $APP_NAME >/dev/null; then
    echo "Successfully installed $APP_NAME $LATEST_VERSION to $BIN_DIR/$APP_NAME."
    echo "Run '$APP_NAME --help' or '$APP_NAME -h' to get started."
else
    echo "Failed to install $APP_NAME."
    exit 1
fi

