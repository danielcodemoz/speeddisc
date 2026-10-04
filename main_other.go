//go:build !windows

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "SpeedDisc is a Windows console program. Build it with GOOS=windows GOARCH=amd64.")
	os.Exit(1)
}
