//go:build windows

package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	qm := parseFlags()
	
	ui := NewUI(os.Stdout, os.Stdin)
	color, utf8 := setupConsole()
	ui.color = color
	ui.utf8 = utf8
	lang := LangPT
	rep := newReport()
	rep.setLang(lang)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		fmt.Fprintln(os.Stdout)
		_, _ = rep.Write()
		os.Exit(130)
	}()

	if !elevated() {
		ui.paragraph(T(lang, "need_admin"))
		if !qm.Enabled {
			ui.Ask("\n  " + T(lang, "press_enter_exit") + " ")
		}
		os.Exit(1)
	}

	if qm.Enabled {
		runQuietMode(qm, lang, ui, rep)
		return
	}

	for {
		ui.Banner(lang)
		ui.Menu(lang)
		line := ui.Ask("\n> ")
		if isQuit(line) {
			break
		}
		if isLangSwitch(line) {
			lang = otherLang(lang)
			rep.setLang(lang)
			continue
		}
		acts, err := ParseSelection(line)
		if err != nil {
			ui.warn(T(lang, "bad_selection"))
			continue
		}
		if Mutates(acts) {
			fmt.Fprintln(ui.out)
			ui.paragraph(confirmBody(lang, acts))
			if !isYes(ui.Ask("  " + T(lang, "confirm_ask"))) {
				ui.warn(T(lang, "cancelled"))
				continue
			}
		}
		execute(lang, acts, ui, rep)
		rep.setLang(lang)
		path, err := rep.Write()
		if err != nil {
			ui.err(T(lang, "report_fail", err.Error()))
		} else {
			ui.ok(T(lang, "report_at", path))
		}
		ui.Ask("\n  " + T(lang, "press_enter") + " ")
	}
	if rep.touched() {
		rep.setLang(lang)
		if path, err := rep.Write(); err == nil {
			ui.ok(T(lang, "report_at", path))
		}
	}
	ui.info(T(lang, "bye"))
	ui.Ask("  " + T(lang, "press_enter_exit") + " ")
}

func runQuietMode(qm quietMode, lang Lang, ui *UI, rep *Report) {
	selection := packageToSelection(qm.Package)
	if selection == "" {
		fmt.Fprintf(os.Stderr, "Invalid package: %s. Use quick/deep/look/full or 1-4.\n", qm.Package)
		os.Exit(1)
	}

	acts, err := ParseSelection(selection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse selection: %v\n", err)
		os.Exit(1)
	}

	collectInfo(rep)
	execute(lang, acts, ui, rep)
	
	rep.setLang(lang)
	path, err := rep.Write()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write report: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Fprintf(os.Stdout, "Report: %s\n", path)
	
	if len(rep.Failed) > 0 {
		os.Exit(1)
	}
}
