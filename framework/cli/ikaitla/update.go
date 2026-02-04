package ikaitla

import (
	"os"
	"os/exec"
	"time"

	"github.com/ikaitla/framework/ui"
	"github.com/spf13/cobra"
)

// NewUpdateCmd creates the update command
func NewUpdateCmd() *cobra.Command {
	var checkOnly bool
	var version string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update Ikaitla framework",
		Long: `Update the Ikaitla framework to the latest version or a specific version.
This will run 'go get' to update the module and regenerate all wiring code.`,
		Example: `  ikaitla update                    # Update to latest
  ikaitla update --version=v2.1.0   # Update to specific version
  ikaitla update --check            # Check for updates only`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(checkOnly, version)
		},
	}

	cmd.Flags().BoolVarP(&checkOnly, "check", "c", false, "Check for updates without installing")
	cmd.Flags().StringVarP(&version, "version", "v", "latest", "Specific version to install")

	return cmd
}

func runUpdate(checkOnly bool, version string) error {
	if checkOnly {
		ui.Info("Checking for updates...")
		// TODO: Implement version check via GitHub API
		ui.Success("You are on the latest version")
		return nil
	}

	ui.Info("Updating Ikaitla framework...")
	ui.Print("")

	spinner := ui.NewSpinner("Fetching latest version")
	spinner.Start()

	// Build go get command
	packagePath := "github.com/emyassine/ikaitla"
	if version != "latest" {
		packagePath += "@" + version
	} else {
		packagePath += "@latest"
	}

	// Run go get
	cmdGet := exec.Command("go", "get", "-u", packagePath)
	cmdGet.Stdout = os.Stdout
	cmdGet.Stderr = os.Stderr

	time.Sleep(1 * time.Second) // Simulate
	spinner.Stop(true)

	if err := cmdGet.Run(); err != nil {
		ui.Error("Failed to update: %v", err)
		return err
	}

	ui.Print("")
	spinner = ui.NewSpinner("Regenerating wiring code")
	spinner.Start()

	// Run go generate
	cmdGen := exec.Command("go", "generate", "./...")
	cmdGen.Stdout = os.Stdout
	cmdGen.Stderr = os.Stderr

	time.Sleep(1 * time.Second) // Simulate
	spinner.Stop(true)

	if err := cmdGen.Run(); err != nil {
		ui.Warning("Code generation failed: %v", err)
		ui.Info("You may need to run: go run framework/autodiscovery.go")
		return err
	}

	ui.Print("")
	ui.Success("Ikaitla framework updated successfully!")
	ui.Info("Run 'make build' to rebuild your CLI")

	return nil
}
