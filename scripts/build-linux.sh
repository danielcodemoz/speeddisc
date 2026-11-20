#!/bin/sh
# Builds (but never releases) the portable Linux archives.
set -eu
cd "$(dirname "$0")/.."
out="${1:-dist}"
mkdir -p "$out"
out="$(cd "$out" && pwd -P)"
stage="$(mktemp -d)"
trap 'rm -rf "$stage"' EXIT HUP INT TERM
for arch in amd64 arm64; do
  dir="$stage/SpeedDisc-Linux-$arch"
  mkdir -p "$dir/licenses"
  CGO_ENABLED=0 GOOS=linux GOARCH="$arch" go build -buildvcs=false -trimpath -ldflags='-s -w' -o "$dir/SpeedDisc-Linux" ./cmd/speeddisc-linux
  chmod 755 "$dir/SpeedDisc-Linux"
  cp LINUX.md "$dir/README.md"
  cp LICENSE "$dir/licenses/SpeedDisc-MIT.txt"
  cp "$(go env GOROOT)/LICENSE" "$dir/licenses/Go.txt"
  cp "$(go list -m -f '{{.Dir}}' golang.org/x/sys)/LICENSE" "$dir/licenses/golang-x-sys.txt"
  cp cmd/speeddisc-linux/web/OFL-Outfit.txt cmd/speeddisc-linux/web/OFL-Fraunces.txt "$dir/licenses/"
  tar -C "$stage" -czf "$out/SpeedDisc-Linux-$arch.tar.gz" "SpeedDisc-Linux-$arch"
done
(cd "$out" && sha256sum SpeedDisc-Linux-amd64.tar.gz SpeedDisc-Linux-arm64.tar.gz > SHA256SUMS.txt)
