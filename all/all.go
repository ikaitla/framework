// Package all pulls in the full Ikaitla framework surface area.
// Import it for side-effects when you want to vendor all built-in packages.
//
// Example:
//
//	import _ "github.com/ikaitla/framework/all"
package all

import (
	_ "github.com/ikaitla/framework/cli/ikaitla"
	_ "github.com/ikaitla/framework/cli/shared"
	_ "github.com/ikaitla/framework/profile"
	_ "github.com/ikaitla/framework/ui"
	_ "github.com/ikaitla/framework/ui/components"
	_ "github.com/ikaitla/framework/ui/output"
	_ "github.com/ikaitla/framework/ui/term"
	_ "github.com/ikaitla/framework/ui/theme"
)

// Touch is a no-op to make the package be intentional.
func Touch() {}
