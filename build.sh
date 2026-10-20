#!/bin/bash

# build script for different OS and architectures
#
# Usage:
#   ./build.sh                    build all platform binaries into ./bin
#   ./build.sh --install [DIR]    also copy this machine's binary into DIR as "agent-install"
#                                 (DIR defaults to $AGENT_INSTALL_DIR, then ~/.local/bin)

INSTALL=false
INSTALL_DIR="${AGENT_INSTALL_DIR:-$HOME/.local/bin}"
if [ "${1:-}" = "--install" ]; then
  INSTALL=true
  if [ -n "${2:-}" ]; then
    INSTALL_DIR="$2"
  fi
elif [ -n "${1:-}" ]; then
  echo "usage: $0 [--install [DIR]]" >&2
  exit 2
fi

# macOS
GOOS=darwin GOARCH=arm64 go build -o ./bin/agent-install-macos-arm64 .
GOOS=darwin GOARCH=amd64 go build -o ./bin/agent-install-macos-amd64 .

# Linux
GOOS=linux GOARCH=amd64 go build -o ./bin/agent-install-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o ./bin/agent-install-linux-arm64 .

# Windows
GOOS=windows GOARCH=amd64 go build -o ./bin/agent-install-windows-amd64.exe .

if [ "$INSTALL" = true ]; then
  HOST_OS=$(go env GOOS)
  HOST_ARCH=$(go env GOARCH)
  EXT=""
  case "$HOST_OS" in
    darwin)  NAME=macos ;;
    windows) NAME=windows; EXT=.exe ;;
    *)       NAME=$HOST_OS ;;
  esac

  SRC="./bin/agent-install-$NAME-$HOST_ARCH$EXT"
  DEST="$INSTALL_DIR/agent-install$EXT"

  if ! { mkdir -p "$INSTALL_DIR" && cp "$SRC" "$DEST" && chmod 755 "$DEST"; }; then
    echo "error: could not install $SRC to $DEST" >&2
    exit 1
  fi
  echo "installed $SRC -> $DEST"

  case ":$PATH:" in
    *":$INSTALL_DIR:"*) ;;
    *)
      echo "note: $INSTALL_DIR is not on your PATH. Add it with:"
      echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
      ;;
  esac
fi