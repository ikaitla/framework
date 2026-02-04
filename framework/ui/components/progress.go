package components

import (
	"fmt"
	"strings"

	"github.com/emyassine/ikaitla/framework/ui/output"
	"github.com/emyassine/ikaitla/framework/ui/theme"
)

type ProgressBar struct {
	out     *output.Output
	total   int
	current int
	width   int
	prefix  string
}

func NewProgressBar(out *output.Output, total int, prefix string) *ProgressBar {
	return &ProgressBar{
		out:    out,
		total:  total,
		width:  40,
		prefix: prefix,
	}
}

func (p *ProgressBar) Update(current int) {
	p.current = current
	p.render()
}

func (p *ProgressBar) Increment() {
	p.current++
	p.render()
}

func (p *ProgressBar) Finish() {
	p.current = p.total
	p.render()
	p.out.Printf("") // newline
}

func (p *ProgressBar) render() {
	if !p.out.ColorsEnabled() {
		return
	}

	if p.total <= 0 {
		p.total = 1
	}
	if p.current < 0 {
		p.current = 0
	}
	if p.current > p.total {
		p.current = p.total
	}

	percent := float64(p.current) / float64(p.total)
	filled := int(percent * float64(p.width))
	empty := p.width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	percentStr := fmt.Sprintf("%.0f%%", percent*100)

	fmt.Fprintf(p.out.Out, "\r%s [%s] %s %d/%d",
		p.prefix,
		p.out.Stylize(bar, theme.Info600),
		percentStr,
		p.current,
		p.total,
	)
}
