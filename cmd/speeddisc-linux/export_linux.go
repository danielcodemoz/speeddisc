//go:build linux

package main

import (
	"bytes"

	"github.com/danielcodemoz/speeddisc/internal/report"
)

func renderReport(format string, result report.Result) ([]byte, error) {
	var buffer bytes.Buffer
	if err := report.Write(&buffer, format, result); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
