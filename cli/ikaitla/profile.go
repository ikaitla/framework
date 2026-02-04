package ikaitla

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewProfileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage profiles",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all profiles",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Available profiles:")
			fmt.Println("  - ikaitla")
		},
	})

	return cmd
}
