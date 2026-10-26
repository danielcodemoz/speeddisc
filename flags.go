package main

import (
	"flag"
	"os"
	"strings"
)

type quietMode struct {
	Enabled bool
	Package string
}

func parseFlags() quietMode {
	var qm quietMode
	var quietFlag bool
	var packageName string

	fs := flag.NewFlagSet("SpeedDisc", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.BoolVar(&quietFlag, "q", false, "quiet mode (non-interactive)")
	fs.BoolVar(&quietFlag, "quiet", false, "quiet mode (non-interactive)")
	fs.StringVar(&packageName, "p", "", "package to run (quick/deep/look/full or 1-4)")
	fs.StringVar(&packageName, "package", "", "package to run (quick/deep/look/full or 1-4)")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return qm
	}

	if !quietFlag && packageName == "" {
		return qm
	}

	qm.Enabled = true
	if packageName != "" {
		qm.Package = strings.ToLower(strings.TrimSpace(packageName))
	} else if len(fs.Args()) > 0 {
		qm.Package = strings.ToLower(strings.TrimSpace(fs.Args()[0]))
	}

	return qm
}

func packageToSelection(pkg string) string {
	pkg = strings.ToLower(strings.TrimSpace(pkg))
	switch pkg {
	case "quick", "1":
		return "1"
	case "deep", "2":
		return "2"
	case "look", "3":
		return "3"
	case "full", "4":
		return "4"
	default:
		return ""
	}
}
