package ikaitla

import (
	"fmt"
	"runtime"

	"github.com/emyassine/ikaitla/framework"
	"github.com/emyassine/ikaitla/framework/ui"
	"github.com/spf13/cobra"
)

// NewVersionCmd creates the version command
func NewVersionCmd() *cobra.Command {
	var detailed bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show Ikaitla framework version",
		Long:  "Display version information for the Ikaitla framework and runtime.",
		Example: `  ikaitla version
  ikaitla version --detailed`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if detailed {
				showDetailedVersion()
			} else {
				showSimpleVersion()
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&detailed, "detailed", "d", false, "Show detailed version info")

	return cmd
}

func showSimpleVersion() {
	ui.Info("%s %s", framework.EngineName, framework.EngineVersion)
	if framework.EngineTagline != "" {
		ui.Print("%s", framework.EngineTagline)
	}
}

func showDetailedVersion() {
	ui.Print("")
	ui.Info("Ikaitla Framework Information")
	ui.Print("")

	ui.RenderKeyValue(map[string]string{
		"Engine":      framework.EngineName,
		"Version":     framework.EngineVersion,
		"Tagline":     framework.EngineTagline,
		"Go Version":  runtime.Version(),
		"OS/Arch":     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		"Compiler":    runtime.Compiler,
	})
	
	ui.Print("")
}
