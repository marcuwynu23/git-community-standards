package app

import (
	"fmt"
	"os"
)

const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

// detectTTY reports whether the process stdout is an interactive terminal and
// the user has not opted out of color (NO_COLOR / dumb TERM).
func detectTTY() bool {
	if isDumbTerminal() {
		return false
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func isDumbTerminal() bool {
	if term := os.Getenv("TERM"); term == "" || term == "dumb" {
		return true
	}
	if os.Getenv("NO_COLOR") != "" {
		return true
	}
	return false
}

func (a *App) wrap(code, msg string) string {
	if !a.useColor {
		return msg
	}
	return code + msg + colorReset
}

func (a *App) printInfo(format string, val ...any) {
	msg := fmt.Sprintf(format, val...)
	fmt.Fprint(a.Stdout, a.wrap(colorBlue, "ℹ  "+msg)+"\n")
}

func (a *App) printSuccess(format string, val ...any) {
	msg := fmt.Sprintf(format, val...)
	fmt.Fprint(a.Stdout, a.wrap(colorGreen, "✓  "+msg)+"\n")
}

func (a *App) printWarn(format string, val ...any) {
	msg := fmt.Sprintf(format, val...)
	fmt.Fprint(a.Stdout, a.wrap(colorYellow, "⚠  "+msg)+"\n")
}

func (a *App) printError(format string, val ...any) {
	msg := fmt.Sprintf(format, val...)
	fmt.Fprint(a.Stderr, a.wrap(colorRed, "✗  "+msg)+"\n")
}

func (a *App) printHeader(format string, val ...any) {
	msg := fmt.Sprintf(format, val...)
	fmt.Fprint(a.Stdout, a.wrap(colorBold+colorCyan, msg)+"\n")
}

func (a *App) printStep(format string, val ...any) {
	msg := fmt.Sprintf(format, val...)
	fmt.Fprint(a.Stdout, a.wrap(colorCyan, "→  "+msg)+"\n")
}

func (a *App) printPlain(format string, val ...any) {
	fmt.Fprintf(a.Stdout, format, val...)
}
