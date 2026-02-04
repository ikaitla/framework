package theme

// ResolveANSI maps Tailwind-like tokens to ANSI foreground colors.
// This is intentionally approximate (terminal palettes differ).
// Extend as needed.
func ResolveANSI(t Token) string {
	switch t {

	// Reds / Danger
	case Red50, Red100, Red200:
		return "\033[91m" // bright red
	case Red300, Red400, Red500:
		return "\033[31m" // red
	case Red600, Red700, Red800, Red900, Red950:
		return "\033[31m" // red (no "darker" ANSI)

	// Yellows / Warning
	case Yellow50, Yellow100, Yellow200:
		return "\033[93m" // bright yellow
	case Yellow300, Yellow400, Yellow500:
		return "\033[33m" // yellow
	case Yellow600, Yellow700, Yellow800, Yellow900, Yellow950:
		return "\033[33m"

	// Greens / Success
	case Green50, Green100, Green200:
		return "\033[92m" // bright green
	case Green300, Green400, Green500:
		return "\033[32m" // green
	case Green600, Green700, Green800, Green900, Green950:
		return "\033[32m"

	// Cyans / Info
	case Cyan50, Cyan100, Cyan200:
		return "\033[96m" // bright cyan
	case Cyan300, Cyan400, Cyan500:
		return "\033[36m" // cyan
	case Cyan600, Cyan700, Cyan800, Cyan900, Cyan950:
		return "\033[36m"

	// Magenta
	case Magenta50, Magenta100, Magenta200:
		return "\033[95m" // bright magenta
	case Magenta300, Magenta400, Magenta500, Magenta600, Magenta700, Magenta800, Magenta900, Magenta950:
		return "\033[35m" // magenta

	// Slate -> neutral
	case Slate50, Slate100, Slate200, Slate300:
		return "\033[97m" // bright white
	case Slate400, Slate500, Slate600:
		return "\033[37m" // white
	case Slate700, Slate800, Slate900, Slate950:
		return "\033[90m" // bright black (grey)

	default:
		return "" // transparent / unknown
	}
}
