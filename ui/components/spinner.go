package components

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ikaitla/framework/ui/output"
	"github.com/ikaitla/framework/ui/theme"
)

type Spinner struct {
	out     *output.Output
	message string

	mu     sync.Mutex
	active bool
	done   chan struct{}
	frames []string
	i      int
}

func NewSpinner(out *output.Output, message string) *Spinner {
	return &Spinner{
		out:     out,
		message: message,
		done:    make(chan struct{}),
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
}

func (s *Spinner) Start() {
	if !s.out.ColorsEnabled() {
		s.out.Printf("%s...", s.message)
		return
	}

	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()

	go s.animate()
}

func (s *Spinner) Stop(success bool) {
	if !s.out.ColorsEnabled() {
		return
	}

	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	s.mu.Unlock()

	<-s.done

	// clear line
	clear := "\r" + strings.Repeat(" ", len(s.message)+12) + "\r"
	fmt.Fprint(s.out.Out, clear)

	if success {
		fmt.Fprintf(s.out.Out, "%s %s\n", s.out.Stylize("[✓]", theme.Success600), s.message)
	} else {
		fmt.Fprintf(s.out.Out, "%s %s\n", s.out.Stylize("[✗]", theme.Danger600), s.message)
	}
}

func (s *Spinner) UpdateMessage(message string) {
	s.mu.Lock()
	s.message = message
	s.mu.Unlock()
}

func (s *Spinner) animate() {
	t := time.NewTicker(80 * time.Millisecond)
	defer t.Stop()

	for {
		s.mu.Lock()
		active := s.active
		msg := s.message
		frame := s.frames[s.i]
		s.i = (s.i + 1) % len(s.frames)
		s.mu.Unlock()

		if !active {
			break
		}

		fmt.Fprintf(s.out.Out, "\r%s %s", s.out.Stylize(frame, theme.Info600), msg)
		<-t.C
	}

	close(s.done)
}
