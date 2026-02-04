package profile

import (
	"fmt"
	"os"

	"github.com/emyassine/ikaitla/framework/ui/theme"
	"github.com/spf13/cobra"
)

// ProfileMetadata defines all configuration for a profile
type ProfileMetadata struct {
	Name        string
	Version     string
	Description string
	LongDesc    string
	Aliases     []string
	Brand       Brand
	Hidden      bool

	// Optional settings
	DocsURL    string
	IssuesURL  string
	License    string
	Contact    string
	ConfigPath string
}

// Brand contains visual identity information
type Brand struct {
	Name    string
	Tagline string
	Color   theme.Token
}


// NewRootCommand creates the root cobra command for a profile
func NewRootCommand(meta ProfileMetadata) *cobra.Command {
	cmd := &cobra.Command{
		Use:     meta.Name,
		Short:   meta.Description,
		Long:    buildLongDescription(meta),
		Version: meta.Version,
		Aliases: meta.Aliases,
		Hidden:  meta.Hidden,
	}

	// Add global flags
	cmd.PersistentFlags().StringP("output", "o", "text", "Output format (text|json|yaml)")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
	cmd.PersistentFlags().Bool("no-color", false, "Disable colored output")

	return cmd
}

// ExecuteProfile runs the profile's root command
func ExecuteProfile(root *cobra.Command) {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

// buildLongDescription constructs the full long description
func buildLongDescription(meta ProfileMetadata) string {
	desc := meta.LongDesc
	if desc == "" {
		desc = meta.Description
	}

	if meta.Brand.Name != "" && meta.Brand.Tagline != "" {
		desc = fmt.Sprintf("%s - %s\n\n%s", meta.Brand.Name, meta.Brand.Tagline, desc)
	}

	return desc
}
