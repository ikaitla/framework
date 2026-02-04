package shared

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/emyassine/ikaitla/framework/ui"
	"github.com/spf13/cobra"
)

// NewDoctorCmd creates a diagnostic command for checking environment
func NewDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check environment and configuration",
		Long:  "Run diagnostics to verify your environment is properly configured.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctor()
		},
	}
}

func runDoctor() error {
	ui.Info("Running diagnostics...")
	ui.Print("")

	checks := []struct {
		name string
		fn   func() (string, bool)
	}{
		{"Go installation", checkGo},
		{"Project structure", checkProjectStructure},
		{"Dependencies", checkDependencies},
		{"Build tools", checkBuildTools},
	}

	table := ui.NewTable("Check", "Status", "Details")

	allPassed := true
	for _, check := range checks {
		detail, passed := check.fn()
		status := "✓"
		if !passed {
			status = "✗"
			allPassed = false
		}
		table.AddRow(check.name, status, detail)
	}

	table.Render()

	ui.Print("")
	if allPassed {
		ui.Success("All checks passed!")
	} else {
		ui.Warning("Some checks failed. See details above.")
	}

	return nil
}

func checkGo() (string, bool) {
	version := runtime.Version()
	return version, true
}

func checkProjectStructure() (string, bool) {
	// Check for key files
	required := []string{"go.mod", "main.go", "cmd"}
	for _, f := range required {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			return fmt.Sprintf("Missing: %s", f), false
		}
	}
	return "OK", true
}

func checkDependencies() (string, bool) {
	cmd := exec.Command("go", "mod", "verify")
	if err := cmd.Run(); err != nil {
		return "Failed", false
	}
	return "OK", true
}

func checkBuildTools() (string, bool) {
	// Check for make
	if _, err := exec.LookPath("make"); err != nil {
		return "make not found (optional)", true // Not critical
	}
	return "OK", true
}
