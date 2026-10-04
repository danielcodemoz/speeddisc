#!/bin/sh
set -eu
cd "$(dirname "$0")"

if command -v go-winres >/dev/null 2>&1; then
  WR=go-winres
else
  WR="$(go env GOPATH)/bin/go-winres"
fi

"$WR" make --arch amd64 --in winres/winres.json --out rsrc
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o SpeedDisc.exe .
python3 verify_pe.py SpeedDisc.exe
