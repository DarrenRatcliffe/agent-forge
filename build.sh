#!/bin/sh
# Cross-compile agent-install for every supported OS/architecture.
# Used both locally and by .github/workflows/build.yml, so CI and local
# builds produce identical file names.
set -eu

OUT_DIR="${OUT_DIR:-./bin}"

# Static binaries, no local paths embedded, symbols stripped.
export CGO_ENABLED=0
LDFLAGS="-s -w"

# Asset names (macos, not darwin) match the existing GitHub releases.
build() {
	goos=$1
	goarch=$2
	name=$3
	ext=""
	[ "$goos" = "windows" ] && ext=".exe"
	echo "building $name-$goarch$ext"
	GOOS=$goos GOARCH=$goarch go build -trimpath -ldflags "$LDFLAGS" \
		-o "$OUT_DIR/agent-install-$name-$goarch$ext" .
}

mkdir -p "$OUT_DIR"

# macOS
build darwin arm64 macos
build darwin amd64 macos

# Linux
build linux amd64 linux
build linux arm64 linux

# Windows
build windows amd64 windows
build windows arm64 windows

# Checksums for release verification (macOS ships shasum, not sha256sum)
if command -v sha256sum >/dev/null 2>&1; then
	sha256="sha256sum"
else
	sha256="shasum -a 256"
fi
(cd "$OUT_DIR" && $sha256 agent-install-* > SHA256SUMS)
