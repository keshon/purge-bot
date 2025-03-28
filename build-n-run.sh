#!/bin/bash

# BUILD

# Get Go version
GO_VERSION=$(go version | awk '{print $3}')

# Get the build date
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build command
go run -ldflags "-X github.com/keshon/purge-bot/internal/version.BuildDate=$BUILD_DATE -X github.com/keshon/purge-bot/internal/version.GoVersion=$GO_VERSION" cmd/purgebot/purgebot.go