#!/bin/bash
set -e

REPO="gopherguides/hype"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

main() {
    local version="${1:-latest}"
    local os arch archive_name download_url

    os=$(uname -s)
    arch=$(uname -m)

    case "$os" in
        Linux)  os="Linux" ;;
        Darwin) os="Darwin" ;;
        *)
            echo "Error: Unsupported operating system: $os"
            exit 1
            ;;
    esac

    case "$arch" in
        x86_64|amd64) arch="x86_64" ;;
        arm64|aarch64) arch="arm64" ;;
        i386|i686) arch="i386" ;;
        *)
            echo "Error: Unsupported architecture: $arch"
            exit 1
            ;;
    esac

    if [ "$version" = "latest" ]; then
        version=$(curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
        if [ -z "$version" ]; then
            echo "Error: Failed to determine latest version"
            exit 1
        fi
    fi

    echo "Installing hype ${version} for ${os}/${arch}..."

    archive_name="hype_${os}_${arch}.tar.gz"
    download_url="https://github.com/${REPO}/releases/download/${version}/${archive_name}"

    tmpdir=$(mktemp -d)
    trap 'rm -rf "$tmpdir"' EXIT

    echo "Downloading ${download_url}..."
    if ! curl -sL "$download_url" -o "${tmpdir}/${archive_name}"; then
        echo "Error: Failed to download ${archive_name}"
        exit 1
    fi

    echo "Extracting..."
    tar -xzf "${tmpdir}/${archive_name}" -C "$tmpdir"

    if [ ! -w "$INSTALL_DIR" ]; then
        echo "Installing to ${INSTALL_DIR} (requires sudo)..."
        sudo mv "${tmpdir}/hype" "$INSTALL_DIR/hype"
    else
        echo "Installing to ${INSTALL_DIR}..."
        mv "${tmpdir}/hype" "$INSTALL_DIR/hype"
    fi

    chmod +x "$INSTALL_DIR/hype"

    echo ""
    echo "hype installed successfully!"
    echo ""
    if command -v hype >/dev/null 2>&1; then
        hype version
    else
        echo "Note: ${INSTALL_DIR} may not be in your PATH"
        echo "Add it with: export PATH=\"\$PATH:${INSTALL_DIR}\""
    fi
}

main "$@"
