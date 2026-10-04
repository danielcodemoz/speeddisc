//go:build windows

package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
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
		ui.Ask("\n  " + T(lang, "press_enter_exit") + " ")
		os.Exit(1)
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
