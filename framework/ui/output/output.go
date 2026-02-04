package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/emyassine/ikaitla/framework/ui/term"
	"github.com/emyassine/ikaitla/framework/ui/theme"
)

type Output struct {
	mu sync.Mutex

	Out io.Writer
	Err io.Writer

	Format    Format
	ColorMode term.ColorMode
}

func New() *Output {
	return &Output{
		Out:       os.Stdout,
		Err:       os.Stderr,
		Format:    Text,
		ColorMode: term.ColorAuto,
	}
}

func (o *Output) ColorsEnabled() bool {
	switch o.ColorMode {
	case term.ColorAlways:
		return true
	case term.ColorNever:
		return false
	default:
		return term.SupportsColor()
	}
}

func (o *Output) Printf(format string, args ...any) {
	o.mu.Lock()
	defer o.mu.Unlock()
	fmt.Fprintf(o.Out, format+"\n", args...)
}

func (o *Output) Errorf(format string, args ...any) {
	o.mu.Lock()
	defer o.mu.Unlock()
	fmt.Fprintf(o.Err, format+"\n", args...)
}

func (o *Output) PrintJSON(v any) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	enc := json.NewEncoder(o.Out)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func (o *Output) Stylize(s string, token theme.Token, attrs ...string) string {
	if !o.ColorsEnabled() {
		return s
	}
	ansi := theme.ResolveANSI(token)
	prefix := ""
	for _, a := range attrs {
		prefix += a
	}
	if ansi == "" {
		return prefix + s + theme.Reset
	}
	return prefix + ansi + s + theme.Reset
}
