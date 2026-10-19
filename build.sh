#!/bin/bash

# build script for different OS and architectures   

# macOS
GOOS=darwin GOARCH=arm64 go build -o ./bin/claude-install-macos-arm64 .
GOOS=darwin GOARCH=amd64 go build -o ./bin/claude-install-macos-amd64 .

# Linux
GOOS=linux GOARCH=amd64 go build -o ./bin/claude-install-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o ./bin/claude-install-linux-arm64 .

# Windows
GOOS=windows GOARCH=amd64 go build -o ./bin/claude-install-windows-amd64.exe .