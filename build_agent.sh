#!/usr/bin/env sh

CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o cocoon-agent cmd/agent/main.go