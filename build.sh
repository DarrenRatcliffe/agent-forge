#!/bin/bash

# build script for different OS and architectures   

# macOS
GOOS=darwin GOARCH=arm64 go build -o ./bin/agent-install-macos-arm64 .
GOOS=darwin GOARCH=amd64 go build -o ./bin/agent-install-macos-amd64 .

# Linux
GOOS=linux GOARCH=amd64 go build -o ./bin/agent-install-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o ./bin/agent-install-linux-arm64 .

# Windows
GOOS=windows GOARCH=amd64 go build -o ./bin/agent-install-windows-amd64.exe .