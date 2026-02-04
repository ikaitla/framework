# Ikaitla Engine – Multi-Profile CLI System

<p align="center">
  <img src="art/ikaitla-blue-dual-mode.png" alt="Ikaitla Engine - Multi-Profile CLI System" width="400">
</p>

## Get Ikaitla

```bash
mkdir -p my-cli-project && cd my-cli-project
curl -fsSL "https://raw.githubusercontent.com/ikaitla/framework/main/install.sh?cachebust=$(date +%s)" | bash
```

## Perspective

1. One binary with infinite profiles.
2. Add a profile by creating a folder in `cmd/<profile-name>/` with commands.
3. We build once and get all the profiles automatically.

## Core Principles

1. **No boilerplate**: The filesystem structure IS the configuration
2. **Auto-discovery**: Commands and profiles discovered at generation time
3. **Single build**: One `go build` command generates everything
4. **Profile isolation**: Each profile is independent and self-contained
5. **No manual registration**: No init(), no explicit wiring, no central declarations

## Design Philosophy

1. **Convention over configuration**: Structure implies behavior
2. **Fail fast**: Invalid profiles fail at generation, not runtime
3. **Developer ergonomics**: Minimal cognitive load, maximum clarity
4. **Production ready**: Type-safe, testable, maintainable

## Architecture

```bash
.
├── bin
├── build.sh
├── cmd
│   └── <profile>                  # Profile discovered by name
│       ├── generated.go           # list of commands (generated on build)
│       ├── metadata.go            # Profile config (version, brand, settings)
│       └── <command-name>.go      # Command: New<OperationName>Cmd()
├── framework
│   ├── autodiscovery.go           # Auto-discovery generator
│   ├── framework.go               # config (EngineName, EngineVersion, EngineTagline)
│   ├── build-info.go              # config (RepoModulePath, ProfilePackageImportPath,
│   │                                      # CmdProfilesImportBase, ProfileWiringOutputFile,
│   │                                      # MainRegistryOutputFile)
│   ├── profile
│   │   └── profile-metadata.go
│   └── ui                         # UI Library
│       ├── colors.go
│       ├── output.go
│       ├── progress.go
│       ├── table.go
│       └── ui_test.go
├── generated_profiles.go          # list of profiles (generated on build)
├── go.mod
├── go.sum
├── LICENSE
├── main.go                        # Entry point: binary name → profile selector
└── Makefile
```

## Commands

```bash
# Development
go generate ./...     # Generate all wiring code
go build             # Build binary + create symlinks
go test ./...        # Run all tests

# Single-command build
go build             # Runs generation automatically via build tags

# Installation
make install         # Copy to ~/.local/bin/
```

## Adding a New Profile

1. Create `cmd/newprofile/metadata.go`
2. Add command files with `NewXxxCmd()` functions
3. Run `go build`
4. Done. Binary `newprofile` is ready.

## Future Extensions

- Plugin system for external profiles
- Dynamic command loading
- Profile inheritance and composition
- Built-in update mechanism
- Telemetry and crash reporting

---

- **License**: Dual (Numerimondes + EPL-2.0)
- **Creator**: El Moumen Yassine
