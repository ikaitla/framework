package ikaitla

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/emyassine/ikaitla/framework/ui"
	"github.com/spf13/cobra"
)

// NewProfileCmd creates the profile management command
func NewProfileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage CLI profiles",
		Long:  "Create, list, enable, or disable CLI profiles in your project.",
	}

	cmd.AddCommand(newProfileListCmd())
	cmd.AddCommand(newProfileNewCmd())
	cmd.AddCommand(newProfileEnableCmd())
	cmd.AddCommand(newProfileDisableCmd())

	return cmd
}

func newProfileListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listProfiles()
		},
	}
}

func newProfileNewCmd() *cobra.Command {
	var description string

	cmd := &cobra.Command{
		Use:   "new [name]",
		Short: "Create a new profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createNewProfile(args[0], description)
		},
	}

	cmd.Flags().StringVarP(&description, "description", "d", "", "Profile description")

	return cmd
}

func newProfileEnableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enable [name]",
		Short: "Enable a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return toggleProfile(args[0], false)
		},
	}
}

func newProfileDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable [name]",
		Short: "Disable a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return toggleProfile(args[0], true)
		},
	}
}

func listProfiles() error {
	cmdDir := "cmd"
	entries, err := os.ReadDir(cmdDir)
	if err != nil {
		return fmt.Errorf("failed to read cmd directory: %w", err)
	}

	ui.Info("Available Profiles:")
	ui.Print("")

	table := ui.NewTable("Profile", "Status", "Path")

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		metadataPath := filepath.Join(cmdDir, name, "metadata.go")
		
		status := "enabled"
		if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
			status = "invalid"
		}

		// Check if disabled (via .disabled file or similar)
		disabledMarker := filepath.Join(cmdDir, name, ".disabled")
		if _, err := os.Stat(disabledMarker); err == nil {
			status = "disabled"
		}

		table.AddRow(name, status, filepath.Join("cmd", name))
	}

	table.Render()
	return nil
}

func createNewProfile(name, description string) error {
	profileDir := filepath.Join("cmd", name)

	// Check if exists
	if _, err := os.Stat(profileDir); err == nil {
		return fmt.Errorf("profile '%s' already exists", name)
	}

	ui.Info("Creating new profile: %s", name)

	spinner := ui.NewSpinner("Setting up profile structure")
	spinner.Start()

	// Create directory
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		spinner.Stop(false)
		return err
	}

	// Create metadata.go
	if description == "" {
		description = fmt.Sprintf("%s CLI profile", name)
	}

	metadata := fmt.Sprintf(`package %s

import "github.com/emyassine/ikaitla/framework/profile"

var Metadata = profile.ProfileMetadata{
	Name:        "%s",
	Version:     "0.1.0",
	Description: "%s",
	LongDesc:    "Detailed description goes here.",
	Aliases:     []string{},
	Brand: profile.Brand{
		Name:    "%s",
		Tagline: "Your tagline",
		Color:   "\033[36m",
	},
	License: "MIT",
}
`, name, name, description, name)

	if err := os.WriteFile(filepath.Join(profileDir, "metadata.go"), []byte(metadata), 0644); err != nil {
		spinner.Stop(false)
		return err
	}

	// Create sample command
	sampleCmd := fmt.Sprintf(`package %s

import (
	"github.com/emyassine/ikaitla/framework/ui"
	"github.com/spf13/cobra"
)

func NewExampleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "example",
		Short: "Example command",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui.Success("This is an example command in %s profile!")
			return nil
		},
	}
}
`, name, name)

	if err := os.WriteFile(filepath.Join(profileDir, "example.go"), []byte(sampleCmd), 0644); err != nil {
		spinner.Stop(false)
		return err
	}

	spinner.Stop(true)

	ui.Print("")
	ui.Success("Profile '%s' created successfully!", name)
	ui.Print("")
	ui.Info("Next steps:")
	ui.Print("  1. Edit cmd/%s/metadata.go to customize", name)
	ui.Print("  2. Add commands in cmd/%s/*.go", name)
	ui.Print("  3. Run: make generate && make build")
	ui.Print("")

	return nil
}

func toggleProfile(name string, disable bool) error {
	profileDir := filepath.Join("cmd", name)
	disabledMarker := filepath.Join(profileDir, ".disabled")

	// Check if profile exists
	if _, err := os.Stat(profileDir); os.IsNotExist(err) {
		return fmt.Errorf("profile '%s' not found", name)
	}

	if disable {
		// Create .disabled marker
		if err := os.WriteFile(disabledMarker, []byte(""), 0644); err != nil {
			return err
		}
		ui.Success("Profile '%s' disabled", name)
		ui.Info("Run 'make generate' to apply changes")
	} else {
		// Remove .disabled marker
		if err := os.Remove(disabledMarker); err != nil && !os.IsNotExist(err) {
			return err
		}
		ui.Success("Profile '%s' enabled", name)
		ui.Info("Run 'make generate' to apply changes")
	}

	return nil
}
