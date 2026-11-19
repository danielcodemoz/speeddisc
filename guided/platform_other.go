//go:build !windows

package main

import (
	"context"
	"fmt"
	"os"
)

func unsafeEntry(p string) bool { i, e := os.Lstat(p); return e != nil || i.Mode()&os.ModeSymlink != 0 }
func localRoot(p string) error {
	if unsafeEntry(p) {
		return fmt.Errorf("link or inaccessible root excluded")
	}
	return nil
}
func startup(context.Context) (any, error) {
	return nil, fmt.Errorf("Windows-only collector; not tested here")
}
