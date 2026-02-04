package ikaitla

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ikaitla/framework/ui"
	"github.com/spf13/cobra"
)

// NewInitCmd creates the init command
func NewInitCmd() *cobra.Command {
	var modulePath string
	var profileName string

	cmd := &cobra.Command{
		Use:   "init [directory]",
		Short: "Initialize a new Ikaitla project",
		Long: `Bootstrap a new Ikaitla-based CLI project with proper structure.
Creates go.mod, main.go, Makefile, and a sample profile.`,
		Example: `  ikaitla init my-cli
  ikaitla init . --module=github.com/me/my-cli
  ikaitla init my-project --profile=deploy`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}
			return runInit(dir, modulePath, profileName)
		},
	}

	cmd.Flags().StringVarP(&modulePath, "module", "m", "", "Go module path (e.g., github.com/user/project)")
	cmd.Flags().StringVarP(&profileName, "profile", "p", "myprofile", "Name of the initial profile")

	return cmd
}

func runInit(dir, modulePath, profileName string) error {
	// Create directory if needed
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	ui.Info("Initializing Ikaitla project in: %s", absDir)
	ui.Print("")

	// Infer module path if not provided
	if modulePath == "" {
		baseName := filepath.Base(absDir)
		modulePath = "github.com/user/" + baseName
		ui.Warning("Module path not specified, using: %s", modulePath)
		ui.Info("You can change this in go.mod")
		ui.Print("")
	}

	// Create structure
	spinner := ui.NewSpinner("Creating project structure")
	spinner.Start()

	steps := []struct {
		name string
		fn   func() error
	}{
		{"go.mod", func() error { return createGoMod(absDir, modulePath) }},
		{"main.go", func() error { return createMainGo(absDir) }},
		{"Makefile", func() error { return createMakefile(absDir) }},
		{"README.md", func() error { return createReadme(absDir, modulePath) }},
		{fmt.Sprintf("cmd/%s/", profileName), func() error { return createProfile(absDir, profileName) }},
	}

	for _, step := range steps {
		if err := step.fn(); err != nil {
			spinner.Stop(false)
			ui.Error("Failed to create %s: %v", step.name, err)
			return err
		}
	}

	spinner.Stop(true)

	ui.Print("")
	ui.Success("Project initialized successfully!")
	ui.Print("")
	ui.Info("Next steps:")
	ui.Print("  1. cd %s", dir)
	ui.Print("  2. go mod download")
	ui.Print("  3. make build")
	ui.Print("  4. ./bin/%s --help", profileName)
	ui.Print("")

	return nil
}

func createGoMod(dir, modulePath string) error {
	content := fmt.Sprintf(`module %s

go 1.24

require github.com/emyassine/ikaitla latest

require (
	github.com/spf13/cobra v1.10.2
)
`, modulePath)

	return os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0644)
}

func createMainGo(dir string) error {
	content := `package main

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
		fmt.Fprintf(os.Stderr, "Error: Unknown profile '%s'\n", profileName)
		fmt.Fprintf(os.Stderr, "Available profiles: %s\n", strings.Join(listProfiles(), ", "))
		os.Exit(1)
	}

	executor()
}

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

var profileRegistry = map[string]func(){}
`

	return os.WriteFile(filepath.Join(dir, "main.go"), []byte(content), 0644)
}

func createMakefile(dir string) error {
	content := `.PHONY: all generate build install clean test

all: build

generate:
	@echo "Generating profile wiring..."
	@go run github.com/ikaitla/framework/autodiscovery

build: generate
	@echo "Building CLI..."
	@mkdir -p bin
	@go build -o bin/$(shell basename $(CURDIR)) .
	@echo "Build complete: ./bin/"

install: build
	@echo "Installing to ~/.local/bin/..."
	@mkdir -p ~/.local/bin
	@cp bin/* ~/.local/bin/
	@echo "Installed successfully"

test:
	@go test -v ./...

clean:
	@rm -rf bin/
	@rm -f generated_profiles.go
	@find cmd -name generated.go -delete 2>/dev/null || true
`

	return os.WriteFile(filepath.Join(dir, "Makefile"), []byte(content), 0644)
}

func createReadme(dir, modulePath string) error {
	content := fmt.Sprintf(`# %s

A CLI tool built with [Ikaitla Framework](https://github.com/emyassine/ikaitla)

## Installation

`+"```bash"+`
go get %s
`+"```"+`

## Usage

`+"```bash"+`
make build
./bin/myprofile --help
`+"```"+`

## Development

`+"```bash"+`
# Add a new command
cat > cmd/myprofile/newcmd.go

# Regenerate wiring
make generate

# Build
make build
`+"```"+`

## License

[Your License]
`, filepath.Base(dir), modulePath)

	return os.WriteFile(filepath.Join(dir, "README.md"), []byte(content), 0644)
}

func createProfile(dir, profileName string) error {
	profileDir := filepath.Join(dir, "cmd", profileName)
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return err
	}

	// metadata.go
	metadata := fmt.Sprintf(`package %s

import "github.com/ikaitla/framework/profile"

var Metadata = profile.ProfileMetadata{
	Name:        "%s",
	Version:     "0.1.0",
	Description: "My awesome CLI profile",
	LongDesc:    "Detailed description of what this profile does.",
	Aliases:     []string{},
	Brand: profile.Brand{
		Name:    "%s",
		Tagline: "Your tagline here",
		Color:   "\033[36m", // Cyan
	},
	License: "MIT",
}
`, profileName, profileName, profileName)

	if err := os.WriteFile(filepath.Join(profileDir, "metadata.go"), []byte(metadata), 0644); err != nil {
		return err
	}

	// hello.go
	hello := fmt.Sprintf(`package %s

import (
	"github.com/ikaitla/framework/ui"
	"github.com/spf13/cobra"
)

func NewHelloCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hello",
		Short: "Say hello",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui.Success("Hello from %s!")
			return nil
		},
	}
}
`, profileName, profileName)

	return os.WriteFile(filepath.Join(profileDir, "hello.go"), []byte(hello), 0644)
}
