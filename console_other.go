//go:build !windows

package main

func consoleColumns() int { return 0 }

func consoleRows() int { return 0 }
