package shared

import (
	"fmt"
	"runtime"

	"github.com/emyassine/ikaitla/framework"
	"github.com/emyassine/ikaitla/framework/ui"
	"github.com/spf13/cobra"
)

// NewVersionCmd creates a reusable version command for any profile
// Usage in profile: root.AddCommand(shared.NewVersionCmd(metadata))
func NewVersionCmd(profileVersion string) *cobra.Command {
	var detailed bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if detailed {
				showDetailedVersion(profileVersion)
			} else {
				showSimpleVersion(profileVersion)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&detailed, "detailed", "d", false, "Show detailed version info")

	return cmd
}

func showSimpleVersion(profileVersion string) {
	ui.Info("Version: %s", profileVersion)
	ui.Print("Framework: %s %s", framework.EngineName, framework.EngineVersion)
}

func showDetailedVersion(profileVersion string) {
	ui.Print("")
	ui.Info("Version Information")
	ui.Print("")

	ui.RenderKeyValue(map[string]string{
		"Profile Version": profileVersion,
		"Framework":       framework.EngineName + " " + framework.EngineVersion,
		"Go Version":      runtime.Version(),
		"OS/Arch":         fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	})

	ui.Print("")
}
