// framework/autodiscovery.go
//
// Auto-discovery generator for Ikaitla multi-profile CLI
// Supports:
// - User profiles in cmd/*
// - Master profile (ikaitla) in framework/cli/ikaitla
// - Automatic shared commands injection
// - .disabled marker support

//go:build ignore

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/ikaitla/framework"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	root, err := findRepoRoot()
	if err != nil {
		return err
	}

	paths, err := resolveFrameworkPaths(root)
	if err != nil {
		return err
	}

	// Generate profile wiring for user profiles (cmd/*)
	userProfiles, err := generateProfileCommands(root, paths, "cmd")
	if err != nil {
		return err
	}

	// Generate profile wiring for ikaitla master profile (framework/cli/ikaitla)
	masterProfile, err := generateMasterProfile(root, paths)
	if err != nil {
		return err
	}

	// Combine all profiles
	allProfiles := append(userProfiles, masterProfile)

	// Generate main registry
	if err := generateMainRegistry(root, allProfiles, paths); err != nil {
		return err
	}

	fmt.Printf("✓ Generated code for %d profiles (%d user + 1 master)\n", len(allProfiles), len(userProfiles))
	fmt.Printf("✓ Shared commands automatically injected in all profiles\n")
	return nil
}

// ----------------------------------------------------------------------------
// Master Profile (ikaitla)
// ----------------------------------------------------------------------------

func generateMasterProfile(root string, paths FrameworkPaths) (ProfileInfo, error) {
	ikaitlaDir := filepath.Join(root, "framework", "cli", "ikaitla")

	// Check if exists
	if _, err := os.Stat(ikaitlaDir); os.IsNotExist(err) {
		fmt.Println("⚠ Ikaitla master profile not found, skipping")
		return ProfileInfo{}, nil
	}

	// Parse the ikaitla profile
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ikaitlaDir, func(fi os.FileInfo) bool {
		name := fi.Name()
		return !strings.HasSuffix(name, "_test.go") && name != framework.ProfileWiringOutputFile
	}, 0)
	if err != nil {
		return ProfileInfo{}, fmt.Errorf("parsing ikaitla profile: %w", err)
	}

	if len(pkgs) == 0 {
		return ProfileInfo{}, fmt.Errorf("no package found in framework/cli/ikaitla")
	}

	// Get package name
	var pkgName string
	for name := range pkgs {
		pkgName = name
		break
	}

	// Discover commands
	commands := discoverCommands(pkgs)
	if len(commands) == 0 {
		return ProfileInfo{}, fmt.Errorf("no commands found in ikaitla profile")
	}

	profile := ProfileInfo{
		Name:        "ikaitla",
		Package:     pkgName,
		Commands:    commands,
		ImportPath:  paths.ModulePath + "/framework/cli/ikaitla",
		ImportAlias: "p_ikaitla",
		Aliases:     []string{"ik"},
	}

	// Generate wiring code for ikaitla
	if err := generateProfileWiring(root, profile, paths, "framework/cli/ikaitla"); err != nil {
		return ProfileInfo{}, err
	}

	fmt.Printf("✓ Master profile 'ikaitla' registered with %d commands\n", len(commands))
	return profile, nil
}

// ----------------------------------------------------------------------------
// framework-derived values
// ----------------------------------------------------------------------------

type FrameworkPaths struct {
	ModulePath        string
	ProfileImportPath string
	SharedImportPath  string
	CmdImportBase     string
}

func resolveFrameworkPaths(repoRoot string) (FrameworkPaths, error) {
	// Module path
	modulePath := strings.TrimSpace(framework.RepoModulePath)
	if modulePath == "" {
		var err error
		modulePath, err = readModulePath(filepath.Join(repoRoot, "go.mod"))
		if err != nil {
			return FrameworkPaths{}, err
		}
	}
	modulePath = strings.TrimRight(modulePath, "/")

	// Profile import path
	profileImport := strings.TrimSpace(framework.ProfilePackageImportPath)
	if profileImport == "" {
		profileImport = modulePath + "/framework/profile"
	}

	// Shared commands import path
	sharedImport := modulePath + "/framework/cli/shared"

	// Cmd import base
	cmdBase := strings.TrimSpace(framework.CmdProfilesImportBase)
	if cmdBase == "" {
		cmdBase = modulePath + "/cmd"
	}

	return FrameworkPaths{
		ModulePath:        modulePath,
		ProfileImportPath: strings.TrimRight(profileImport, "/"),
		SharedImportPath:  strings.TrimRight(sharedImport, "/"),
		CmdImportBase:     strings.TrimRight(cmdBase, "/"),
	}, nil
}

// ----------------------------------------------------------------------------
// repo root + go.mod
// ----------------------------------------------------------------------------

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("cannot find go.mod (not in a Go module)")
}

func readModulePath(goModPath string) (string, error) {
	f, err := os.Open(goModPath)
	if err != nil {
		return "", fmt.Errorf("opening go.mod: %w", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if strings.HasPrefix(line, "module ") {
			mod := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			if mod == "" {
				return "", fmt.Errorf("invalid module directive in go.mod")
			}
			return mod, nil
		}
	}
	if err := sc.Err(); err != nil {
		return "", fmt.Errorf("reading go.mod: %w", err)
	}
	return "", fmt.Errorf("module directive not found in go.mod")
}

// ----------------------------------------------------------------------------
// discovery types
// ----------------------------------------------------------------------------

type ProfileInfo struct {
	Name        string
	Package     string
	Commands    []string
	ImportPath  string
	ImportAlias string
	Aliases     []string
}

// ----------------------------------------------------------------------------
// generate: profile commands + per-profile wiring
// ----------------------------------------------------------------------------

func generateProfileCommands(root string, paths FrameworkPaths, baseDir string) ([]ProfileInfo, error) {
	cmdDir := filepath.Join(root, baseDir)
	entries, err := os.ReadDir(cmdDir)
	if err != nil {
		return nil, fmt.Errorf("reading %s directory: %w", baseDir, err)
	}

	var profiles []ProfileInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		profileName := entry.Name()
		profileDir := filepath.Join(cmdDir, profileName)

		// Check for .disabled marker
		disabledMarker := filepath.Join(profileDir, ".disabled")
		if _, err := os.Stat(disabledMarker); err == nil {
			fmt.Printf("⊘ Skipping disabled profile: %s\n", profileName)
			continue
		}

		// Parse profile package
		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, profileDir, func(fi os.FileInfo) bool {
			name := fi.Name()
			return !strings.HasSuffix(name, "_test.go") && name != framework.ProfileWiringOutputFile
		}, 0)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", profileName, err)
		}

		if len(pkgs) == 0 {
			continue
		}

		// Get package name (first one)
		var pkgName string
		for name := range pkgs {
			pkgName = name
			break
		}

		// Discover command factories
		commands := discoverCommands(pkgs)
		if len(commands) == 0 {
			fmt.Printf("⚠ No commands found in profile: %s\n", profileName)
			continue
		}

		// Extract aliases from metadata if exists
		aliases := extractProfileAliases(pkgs)

		profile := ProfileInfo{
			Name:        profileName,
			Package:     pkgName,
			Commands:    commands,
			ImportPath:  paths.CmdImportBase + "/" + profileName,
			ImportAlias: makeImportAlias(profileName, pkgName),
			Aliases:     aliases,
		}

		// Generate wiring code for this profile
		outputDir := filepath.Join(baseDir, profileName)
		if err := generateProfileWiring(root, profile, paths, outputDir); err != nil {
			return nil, err
		}

		profiles = append(profiles, profile)
		fmt.Printf("✓ Profile '%s' registered with %d commands\n", profileName, len(commands))
	}

	// stable order
	sort.Slice(profiles, func(i, j int) bool { return profiles[i].Name < profiles[j].Name })
	return profiles, nil
}

func discoverCommands(pkgs map[string]*ast.Package) []string {
	commands := make(map[string]bool)
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				fn, ok := n.(*ast.FuncDecl)
				if !ok || fn.Recv != nil || fn.Name == nil {
					return true
				}
				name := fn.Name.Name
				if strings.HasPrefix(name, "New") && strings.HasSuffix(name, "Cmd") {
					commands[name] = true
				}
				return true
			})
		}
	}

	result := make([]string, 0, len(commands))
	for cmd := range commands {
		result = append(result, cmd)
	}
	sort.Strings(result)
	return result
}

func extractProfileAliases(pkgs map[string]*ast.Package) []string {
	// Try to extract aliases from Metadata.Aliases
	// This is a simplified version - could be more sophisticated
	var aliases []string

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				// Look for: Aliases: []string{"alias1", "alias2"}
				comp, ok := n.(*ast.CompositeLit)
				if !ok {
					return true
				}

				for _, elt := range comp.Elts {
					kv, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}

					key, ok := kv.Key.(*ast.Ident)
					if !ok || key.Name != "Aliases" {
						continue
					}

					// Extract string literals
					if compLit, ok := kv.Value.(*ast.CompositeLit); ok {
						for _, item := range compLit.Elts {
							if lit, ok := item.(*ast.BasicLit); ok && lit.Kind == token.STRING {
								alias := strings.Trim(lit.Value, `"`)
								aliases = append(aliases, alias)
							}
						}
					}
				}
				return true
			})
		}
	}

	return aliases
}

func makeImportAlias(profileName, pkgName string) string {
	base := profileName
	if base == "" {
		base = pkgName
	}
	var b strings.Builder
	b.WriteString("p_")
	for _, r := range base {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune('_')
		}
	}
	alias := strings.Trim(b.String(), "_")
	if alias == "p_" {
		return "p_profile"
	}
	return alias
}

func generateProfileWiring(root string, profile ProfileInfo, paths FrameworkPaths, outputDir string) error {
	// Template with shared commands injection
	tmplText := `// Code generated by ikaitla. DO NOT EDIT.
//
// Engine: {{.EngineName}} {{.EngineVersion}}
// Profile: {{.ProfileName}}
package {{.Package}}

import (
	profile "{{.ProfileImportPath}}"
	shared "{{.SharedImportPath}}"
)

// Execute runs the {{.ProfileName}} profile
func Execute() {
	root := profile.NewRootCommand(Metadata)

	// Auto-discovered commands
{{- range .Commands}}
	root.AddCommand({{.}}())
{{- end}}

	// Shared commands (available in all profiles)
	root.AddCommand(shared.NewVersionCmd(Metadata.Version))
	root.AddCommand(shared.NewDoctorCmd())

	profile.ExecuteProfile(root)
}
`

	tmpl := template.Must(template.New("profile").Parse(tmplText))

	data := struct {
		EngineName        string
		EngineVersion     string
		EngineTagline     string
		ProfileName       string
		Package           string
		Commands          []string
		ProfileImportPath string
		SharedImportPath  string
	}{
		EngineName:        framework.EngineName,
		EngineVersion:     framework.EngineVersion,
		EngineTagline:     framework.EngineTagline,
		ProfileName:       profile.Name,
		Package:           profile.Package,
		Commands:          profile.Commands,
		ProfileImportPath: paths.ProfileImportPath,
		SharedImportPath:  paths.SharedImportPath,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("template execution for %s: %w", profile.Name, err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("formatting %s: %w\n%s", profile.Name, err, buf.String())
	}

	outputPath := filepath.Join(root, outputDir, framework.ProfileWiringOutputFile)
	if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
		return fmt.Errorf("writing %s: %w", outputPath, err)
	}

	return nil
}

// ----------------------------------------------------------------------------
// generate: main registry
// ----------------------------------------------------------------------------

func generateMainRegistry(root string, profiles []ProfileInfo, paths FrameworkPaths) error {
	tmplText := `// Code generated by ikaitla. DO NOT EDIT.
//
// Engine: {{.EngineName}} {{.EngineVersion}}
package main

import (
{{- range .Profiles}}
	{{.ImportAlias}} "{{.ImportPath}}"
{{- end}}
)

func init() {
{{- range $profile := .Profiles}}
	profileRegistry["{{$profile.Name}}"] = {{$profile.ImportAlias}}.Execute
{{- if $profile.Aliases}}
{{- range $alias := $profile.Aliases}}
	profileRegistry["{{$alias}}"] = {{$profile.ImportAlias}}.Execute
{{- end}}
{{- end}}
{{- end}}
}
`

	tmpl := template.Must(template.New("registry").Parse(tmplText))

	data := struct {
		EngineName    string
		EngineVersion string
		EngineTagline string
		Profiles      []ProfileInfo
	}{
		EngineName:    framework.EngineName,
		EngineVersion: framework.EngineVersion,
		EngineTagline: framework.EngineTagline,
		Profiles:      profiles,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("template execution for main registry: %w", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("formatting registry: %w\n%s", err, buf.String())
	}

	outputPath := filepath.Join(root, framework.MainRegistryOutputFile)
	if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
		return fmt.Errorf("writing %s: %w", outputPath, err)
	}

	return nil
}
