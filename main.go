package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ikaitla/framework"
)

func main() {
	binaryName := filepath.Base(os.Args[0])
	profileName := strings.ToLower(strings.TrimSuffix(binaryName, filepath.Ext(binaryName)))

	executor, found := getProfileExecutor(profileName)
	if !found {
		fmt.Fprintf(os.Stderr, "%s %s\n", framework.EngineName, framework.EngineVersion)
		if framework.EngineTagline != "" {
			fmt.Fprintf(os.Stderr, "%s\n\n", framework.EngineTagline)
		}
		fmt.Fprintf(os.Stderr, "Error: Unknown profile '%s'\n", profileName)
		fmt.Fprintf(os.Stderr, "Available profiles: %s\n", strings.Join(listProfiles(), ", "))
		os.Exit(1)
	}

	executor()
}

// getProfileExecutor returns the executor function for a given profile name
// This function is populated by generated_profiles.go
func getProfileExecutor(name string) (func(), bool) {
	executor, found := profileRegistry[name]
	return executor, found
}

func listProfiles() []string {
	profiles := make([]string, 0, len(profileRegistry))
	for name := range profileRegistry {
		profiles = append(profiles, name)
	}
	return profiles
}

var profileRegistry = map[string]func(){} // populated by generated code
