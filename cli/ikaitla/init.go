package ikaitla

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init [profile-name]",
		Short: "Initialize a new profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			profileName := args[0]
			fmt.Printf("Initializing profile: %s\n", profileName)
			fmt.Println("Profile created successfully")
		},
	}
}
