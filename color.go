package main

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

var useColor bool

func init() {
	out := os.Stdout
	// Disable color when output is not a terminal or when the user opted out.
	fi, err := out.Stat()
	isTTY := err == nil && (fi.Mode()&os.ModeCharDevice) != 0
	useColor = isTTY && !isDumbTerminal()
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

func wrap(code, msg string) string {
	if !useColor {
		return msg
	}
	return code + msg + colorReset
}

func printInfo(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Print(wrap(colorBlue, "ℹ  "+msg) + "\n")
}

func printSuccess(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Print(wrap(colorGreen, "✓  "+msg) + "\n")
}

func printWarn(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Print(wrap(colorYellow, "⚠  "+msg) + "\n")
}

func printError(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprint(os.Stderr, wrap(colorRed, "✗  "+msg)+"\n")
}

func printHeader(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Print(wrap(colorBold+colorCyan, msg) + "\n")
}

func printStep(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Print(wrap(colorCyan, "→  "+msg) + "\n")
}

func printPlain(format string, a ...any) {
	fmt.Printf(format, a...)
}
