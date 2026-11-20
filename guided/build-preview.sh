#!/bin/sh
set -eu
cd "$(dirname "$0")"
go test ./...
go vet ./...
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -o SpeedDisc-Guided-preview.exe .
printf '%s\n' 'PRIVATE PREVIEW ONLY: no Windows runtime acceptance, no publication.'
