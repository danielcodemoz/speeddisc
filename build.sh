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

# Code signing (optional, deferred until OV certificate is available)
# Set CODESIGN_PFX=/path/to/cert.pfx and CODESIGN_PASSWORD=yourpass to sign
# Requires osslsigncode (Linux/macOS) or signtool (Windows)
#
# Code signing is DEFERRED by choice:
# - OV (Organization Validation) certificates cost hundreds of euros per year
# - Trust comes from open source (MIT license, GitHub) instead of a paid cert
# - SmartScreen warning is expected and documented for users
# - When a certificate is acquired in the future, set the env vars and rebuild
#
# DO NOT fake-sign or use self-signed certificates. Leave unsigned until
# a legitimate paid Authenticode certificate exists.

if [ -n "${CODESIGN_PFX:-}" ] && [ -n "${CODESIGN_PASSWORD:-}" ]; then
  echo "Code signing enabled, attempting to sign SpeedDisc.exe..."
  if command -v osslsigncode >/dev/null 2>&1; then
    osslsigncode sign \
      -pkcs12 "$CODESIGN_PFX" \
      -pass "$CODESIGN_PASSWORD" \
      -n "SpeedDisc" \
      -i "https://speeddisc.danielpro.dev" \
      -t http://timestamp.digicert.com \
      -in SpeedDisc.exe \
      -out SpeedDisc-signed.exe
    mv SpeedDisc-signed.exe SpeedDisc.exe
    echo "Signed with osslsigncode."
  elif command -v signtool >/dev/null 2>&1; then
    signtool sign \
      /f "$CODESIGN_PFX" \
      /p "$CODESIGN_PASSWORD" \
      /n "SpeedDisc" \
      /d "SpeedDisc maintenance tool" \
      /du "https://speeddisc.danielpro.dev" \
      /t http://timestamp.digicert.com \
      SpeedDisc.exe
    echo "Signed with signtool."
  else
    echo "WARNING: CODESIGN_PFX set but neither osslsigncode nor signtool found."
    echo "Build continues unsigned."
  fi
else
  echo "No code signing configured (CODESIGN_PFX not set). Build is unsigned by choice."
  echo "See build.sh comments: signing is deferred until an OV certificate exists."
fi

python3 verify_pe.py SpeedDisc.exe
