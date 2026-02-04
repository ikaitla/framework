package ui_test

import (
	"testing"

	"github.com/emyassine/ikaitla/framework/ui"
	"github.com/emyassine/ikaitla/framework/ui/theme"
)

func TestColorize_NoColor(t *testing.T) {
	ui.SetColorMode(2) // ColorNever
	got := ui.Colorize("test", theme.Danger600)
	if got != "test" {
		t.Fatalf("expected plain text, got %q", got)
	}
}

func TestTableCreation(t *testing.T) {
	ui.SetColorMode(2) // ColorNever
	table := ui.NewTable("Col1", "Col2", "Col3")
	table.AddRow("A", "B", "C")
	table.AddRow("X", "Y", "Z")
	// no panic => ok
}

func TestProgressBar(t *testing.T) {
	ui.SetColorMode(2) // ColorNever
	pb := ui.NewProgressBar(100, "Testing")
	pb.Update(50)
	pb.Finish()
}
