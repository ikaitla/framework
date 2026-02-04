package ikaitla

import "github.com/emyassine/ikaitla/framework/profile"

// Metadata defines the ikaitla internal CLI profile
var Metadata = profile.ProfileMetadata{
	Name:        "ikaitla",
	Version:     "2026.02.04",
	Description: "Ikaitla framework management CLI",
	LongDesc: `Internal CLI for managing the Ikaitla framework itself.
Use this to update, configure, and manage your Ikaitla installation.`,
	Aliases: []string{"ik"},
	Brand: profile.Brand{
		Name:    "Ikaitla",
		Tagline: "Framework management",
		Color:   "\033[34m", // Blue
	},
	DocsURL:   "https://ikaitla.dev/docs",
	IssuesURL: "https://github.com/emyassine/ikaitla/issues",
	License:   "DUAL: Numerimondes + EPL-2.0",
	Contact:   "support@numerimondes.com",
}
