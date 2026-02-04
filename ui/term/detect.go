package term

import (
	"os"
	"runtime"
)

type ColorMode int

const (
	ColorAuto ColorMode = iota
	ColorAlways
	ColorNever
)

func SupportsColor() bool {
	// NO_COLOR disables
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// FORCE_COLOR enables
	if os.Getenv("FORCE_COLOR") != "" {
		return true
	}

	// Conservative on Windows by default
	if runtime.GOOS == "windows" {
		return false
	}

	// Must be a TTY
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return false
	}

	return true
}
