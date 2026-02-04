package framework

/**
 * Defaults:
 * -> can be overridden at build time
 * -> with -ldflags -X)
 *
 *  How to Override:
 *
 *  go build -ldflags "\
 *  -X 'github.com/emyassine/ikaitla/framework.EngineName=Ikaitla Engine' \
 *  -X 'github.com/emyassine/ikaitla/framework.EngineVersion=2026.02.03' \
 *  -X 'github.com/emyassine/ikaitla/framework.EngineTagline=Composable multi-profile CLI runtime' \
 *  " -o bin/ikaitla .
 */

var (
	EngineName    = "Ikaitla Engine"
	EngineVersion = "2026.02.03"
	EngineTagline = "Composable multi-profile CLI runtime"
)
