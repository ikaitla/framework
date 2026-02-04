package profile

import (
	"github.com/spf13/cobra"
)

// SharedCommandsRegistry holds shared commands that can be added to any profile
type SharedCommandsRegistry struct {
	commands []*cobra.Command
}

// NewSharedCommandsRegistry creates a new shared commands registry
func NewSharedCommandsRegistry() *SharedCommandsRegistry {
	return &SharedCommandsRegistry{
		commands: make([]*cobra.Command, 0),
	}
}

// Register adds a shared command to the registry
func (r *SharedCommandsRegistry) Register(cmd *cobra.Command) {
	r.commands = append(r.commands, cmd)
}

// AddToRoot adds all registered shared commands to a root command
func (r *SharedCommandsRegistry) AddToRoot(root *cobra.Command) {
	for _, cmd := range r.commands {
		root.AddCommand(cmd)
	}
}

// Global shared commands registry
var sharedRegistry = NewSharedCommandsRegistry()

// RegisterSharedCommand registers a command to be available in all profiles
func RegisterSharedCommand(cmd *cobra.Command) {
	sharedRegistry.Register(cmd)
}

// AddSharedCommands adds all registered shared commands to the root command
func AddSharedCommands(root *cobra.Command) {
	sharedRegistry.AddToRoot(root)
}
