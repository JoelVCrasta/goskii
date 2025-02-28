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

# Download dependencies
echo "Downloading dependencies..."

# Check if FFmpeg is installed
if ! command -v ffmpeg &> /dev/null; then
    echo "FFmpeg not found, downloading latest FFmpeg version..."
    
    if [[ "$OS" == "linux" ]]; then
        if command -v apt &> /dev/null; then
            sudo apt update
            sudo apt install -y ffmpeg
        elif command -v pacman &> /dev/null; then
            sudo pacman -Syu --noconfirm ffmpeg
        elif command -v dnf &> /dev/null; then
            sudo dnf install -y ffmpeg
        elif command -v zypper &> /dev/null; then
            sudo zypper install -y ffmpeg
        else 
            echo "Unsupported package manager."
            exit 1
        fi
    elif [[ "$OS" == "darwin" ]]; then
        if command -v brew &> /dev/null; then
            brew install ffmpeg
        else 
            echo "Homebrew not found. Please install Homebrew and try again."
            exit 1
        fi
    fi

    if command -v ffmpeg &> /dev/null; then
        echo "FFmpeg installed successfully."
    else
        echo "Error: FFmpeg installation failed."
        exit 1
    fi
else 
    echo "FFmpeg is already installed."
fi

# Check if yt-dlp is installed
if ! command -v yt-dlp &> /dev/null; then
    echo "yt-dlp not found, downloading latest yt-dlp version..."

    if [[ "$OS" == "linux" ]]; then
        YT_DLP_URL=$(curl -s https://api.github.com/repos/yt-dlp/yt-dlp/releases/latest | jq -r '.assets[] | select(.name == "yt-dlp_linux") | .browser_download_url')
    elif [[ "$OS" == "darwin" ]]; then
        YT_DLP_URL=$(curl -s https://api.github.com/repos/yt-dlp/yt-dlp/releases/latest | jq -r '.assets[] | select(.name == "yt-dlp_macos") | .browser_download_url')
    fi

    if [[ -z "$YT_DLP_URL" ]]; then
        echo "Failed to fetch the latest version of yt-dlp."
        exit 1
    fi

    if ! sudo curl -sL $YT_DLP_URL -o $BIN_DIR/yt-dlp; then
        echo "Failed to download yt-dlp."
        exit 1
    fi

    sudo chmod +x $BIN_DIR/yt-dlp

    if command -v yt-dlp &> /dev/null; then
        echo "yt-dlp installed successfully."
    else
        echo "Error: yt-dlp installation failed."
        exit 1
    fi
else
    echo "yt-dlp is already installed."
fi

# Build Download URL
BINARY="${APP_NAME}-${OS}-amd64"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY}"

# Download and Install
echo "Downloading $APP_NAME from $DOWNLOAD_URL..."
if ! sudo curl -sL $DOWNLOAD_URL -o $BIN_DIR/$APP_NAME; then
    echo "Failed to download $APP_NAME."
    exit 1
fi

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
