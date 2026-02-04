package ui

import (
	"fmt"

	"github.com/ikaitla/framework/ui/components"
	"github.com/ikaitla/framework/ui/output"
	"github.com/ikaitla/framework/ui/term"
	"github.com/ikaitla/framework/ui/theme"
)

var defaultOut = output.New()

// SetFormat lets your root command wire `--output`.
func SetFormat(f output.Format) { defaultOut.Format = f }

// SetColorMode wires `--no-color` and/or future flags.
func SetColorMode(m term.ColorMode) { defaultOut.ColorMode = m }

// ColorsEnabled exposes current state
func ColorsEnabled() bool { return defaultOut.ColorsEnabled() }

// Colorize applies a token, not raw ANSI.
func Colorize(text string, t theme.Token, attrs ...string) string {
	return defaultOut.Stylize(text, t, attrs...)
}

// Print matches your old API
func Print(format string, args ...any) { defaultOut.Printf(format, args...) }

// PrintJSON matches old API
func PrintJSON(v any) error { return defaultOut.PrintJSON(v) }

func Success(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if defaultOut.ColorsEnabled() {
		defaultOut.Printf("%s %s", defaultOut.Stylize("[✓]", theme.Success600), msg)
		return
	}
	defaultOut.Printf("[✓] %s", msg)
}

func Error(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if defaultOut.ColorsEnabled() {
		defaultOut.Errorf("%s %s", defaultOut.Stylize("[✗]", theme.Danger600), msg)
		return
	}
	defaultOut.Errorf("[✗] %s", msg)
}

func Warning(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if defaultOut.ColorsEnabled() {
		defaultOut.Printf("%s %s", defaultOut.Stylize("[!]", theme.Warning600), msg)
		return
	}
	defaultOut.Printf("[!] %s", msg)
}

func Info(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if defaultOut.ColorsEnabled() {
		defaultOut.Printf("%s %s", defaultOut.Stylize("[i]", theme.Info600), msg)
		return
	}
	defaultOut.Printf("[i] %s", msg)
}

// Components
type Spinner = components.Spinner
type ProgressBar = components.ProgressBar
type Table = components.Table

func NewSpinner(message string) *Spinner {
	return components.NewSpinner(defaultOut, message)
}

func NewProgressBar(total int, prefix string) *ProgressBar {
	return components.NewProgressBar(defaultOut, total, prefix)
}

func NewTable(headers ...string) *Table {
	return components.NewTable(defaultOut, headers...)
}

func RenderKeyValue(pairs map[string]string) {
	components.RenderKeyValue(defaultOut, pairs)
}
