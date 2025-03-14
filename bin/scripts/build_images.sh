#!/bin/bash

set -e

# Default values
DEFAULT_IMAGE_NAME="gambling-api"
DEFAULT_TAG="latest"

# Parse command line arguments
IMAGE_NAME=${1:-$DEFAULT_IMAGE_NAME}
TAG=${2:-$DEFAULT_TAG}

echo "🔨 Building Docker images for $IMAGE_NAME"
echo "--------------------------------------"

# Function to build image for specific platform
build_image() {
    local os=$1
    local arch=$2
    local tag_suffix=$3
    
    echo "⚙️  Building image for $os/$arch..."
    
    docker build \
        --build-arg TARGETOS=$os \
        --build-arg TARGETARCH=$arch \
        -t "$IMAGE_NAME:$TAG-$tag_suffix" \
        -f Dockerfile .
    
    echo "✅ Successfully built $IMAGE_NAME:$TAG-$tag_suffix"
}

# Function to detect and build for current platform
build_current_platform() {
    echo "⚙️  Detecting and building for current platform..."
    
    if [[ $(uname -s) == "Darwin" ]]; then
        local os="darwin"
        if [[ $(uname -m) == "arm64" ]]; then
            local arch="arm64"
            local suffix="mac-arm64"
        else
            local arch="amd64"
            local suffix="mac-amd64"
        fi
    else
        local os="linux"
        if [[ $(uname -m) == "aarch64" || $(uname -m) == "arm64" ]]; then
            local arch="arm64"
            local suffix="linux-arm64"
        else
            local arch="amd64"
            local suffix="linux-amd64"
        fi
    fi
    
    docker build \
        --build-arg TARGETOS=$os \
        --build-arg TARGETARCH=$arch \
        -t "$IMAGE_NAME:$TAG" \
        -t "$IMAGE_NAME:$TAG-$suffix" \
        -f Dockerfile .
    
    echo "✅ Successfully built $IMAGE_NAME:$TAG (for $os/$arch)"
}

# Display build menu
show_menu() {
    echo ""
    echo "Select platform to build for:"
    echo "1) Current platform (auto-detect)"
    echo "2) Linux AMD64"
    echo "3) Linux ARM64"
    echo "4) macOS AMD64"
    echo "5) macOS ARM64"
    echo "6) All platforms"
    echo "0) Exit"
    echo ""
    read -p "Enter your choice: " choice
    
    case $choice in
        1)
            build_current_platform
            ;;
        2)
            build_image "linux" "amd64" "linux-amd64"
            ;;
        3)
            build_image "linux" "arm64" "linux-arm64"
            ;;
        4)
            build_image "darwin" "amd64" "mac-amd64"
            ;;
        5)
            build_image "darwin" "arm64" "mac-arm64"
            ;;
        6)
            build_image "linux" "amd64" "linux-amd64"
            build_image "linux" "arm64" "linux-arm64"
            build_image "darwin" "amd64" "mac-amd64"
            build_image "darwin" "arm64" "mac-arm64"
            ;;
        0)
            echo "Exiting..."
            exit 0
            ;;
        *)
            echo "❌ Invalid option. Please try again."
            show_menu
            ;;
    esac
}

# Check if running in interactive mode
if [ -t 0 ]; then
    show_menu
else
    # Non-interactive mode: build for current platform
    build_current_platform
fi

echo ""
echo "🚀 All builds completed!"