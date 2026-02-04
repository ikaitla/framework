package components

import (
	"strings"

	"github.com/emyassine/ikaitla/framework/ui/output"
	"github.com/emyassine/ikaitla/framework/ui/theme"
)

type Table struct {
	headers []string
	rows    [][]string
	widths  []int
	out     *output.Output
}

func NewTable(out *output.Output, headers ...string) *Table {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	return &Table{
		out:     out,
		headers: headers,
		widths:  widths,
	}
}

func (t *Table) AddRow(cells ...string) {
	row := make([]string, len(t.headers))
	for i := 0; i < len(t.headers); i++ {
		if i < len(cells) {
			row[i] = cells[i]
			if len(cells[i]) > t.widths[i] {
				t.widths[i] = len(cells[i])
			}
		}
	}
	t.rows = append(t.rows, row)
}

func (t *Table) Render() {
	if len(t.headers) == 0 {
		return
	}

	var headerLine strings.Builder
	var separator strings.Builder

	for i, h := range t.headers {
		headerLine.WriteString(pad(h, t.widths[i]))
		headerLine.WriteString("  ")

		separator.WriteString(strings.Repeat("─", t.widths[i]))
		separator.WriteString("  ")
	}

	hline := strings.TrimRight(headerLine.String(), " ")
	sline := strings.TrimRight(separator.String(), " ")

	if t.out.ColorsEnabled() {
		t.out.Printf("%s", t.out.Stylize(hline, theme.Slate900, theme.Bold))
	} else {
		t.out.Printf("%s", hline)
	}
	t.out.Printf("%s", sline)

	for _, row := range t.rows {
		var rowLine strings.Builder
		for i, cell := range row {
			rowLine.WriteString(pad(cell, t.widths[i]))
			rowLine.WriteString("  ")
		}
		t.out.Printf("%s", strings.TrimRight(rowLine.String(), " "))
	}
}

func pad(s string, w int) string {
	if len(s) >= w {
		return s
	}
	return s + strings.Repeat(" ", w-len(s))
}
