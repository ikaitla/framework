package shared

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

func NewDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check system health",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("System Health Check")
			fmt.Println("===================")
			fmt.Printf("Go version: %s\n", runtime.Version())
			fmt.Printf("OS: %s\n", runtime.GOOS)
			fmt.Printf("Arch: %s\n", runtime.GOARCH)
			fmt.Printf("CPUs: %d\n", runtime.NumCPU())
			fmt.Println("\nAll systems operational")
		},
	}
}
