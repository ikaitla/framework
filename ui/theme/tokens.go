package theme

// Token is a Tailwind-like color token.
// Example: Cyan600, Slate900, Danger50.
type Token string

// Core palette (subset, extend anytime).
const (
	Transparent Token = "transparent"

	Slate50  Token = "slate-50"
	Slate100 Token = "slate-100"
	Slate200 Token = "slate-200"
	Slate300 Token = "slate-300"
	Slate400 Token = "slate-400"
	Slate500 Token = "slate-500"
	Slate600 Token = "slate-600"
	Slate700 Token = "slate-700"
	Slate800 Token = "slate-800"
	Slate900 Token = "slate-900"
	Slate950 Token = "slate-950"

	Red50  Token = "red-50"
	Red100 Token = "red-100"
	Red200 Token = "red-200"
	Red300 Token = "red-300"
	Red400 Token = "red-400"
	Red500 Token = "red-500"
	Red600 Token = "red-600"
	Red700 Token = "red-700"
	Red800 Token = "red-800"
	Red900 Token = "red-900"
	Red950 Token = "red-950"

	Yellow50  Token = "yellow-50"
	Yellow100 Token = "yellow-100"
	Yellow200 Token = "yellow-200"
	Yellow300 Token = "yellow-300"
	Yellow400 Token = "yellow-400"
	Yellow500 Token = "yellow-500"
	Yellow600 Token = "yellow-600"
	Yellow700 Token = "yellow-700"
	Yellow800 Token = "yellow-800"
	Yellow900 Token = "yellow-900"
	Yellow950 Token = "yellow-950"

	Green50  Token = "green-50"
	Green100 Token = "green-100"
	Green200 Token = "green-200"
	Green300 Token = "green-300"
	Green400 Token = "green-400"
	Green500 Token = "green-500"
	Green600 Token = "green-600"
	Green700 Token = "green-700"
	Green800 Token = "green-800"
	Green900 Token = "green-900"
	Green950 Token = "green-950"

	Cyan50  Token = "cyan-50"
	Cyan100 Token = "cyan-100"
	Cyan200 Token = "cyan-200"
	Cyan300 Token = "cyan-300"
	Cyan400 Token = "cyan-400"
	Cyan500 Token = "cyan-500"
	Cyan600 Token = "cyan-600"
	Cyan700 Token = "cyan-700"
	Cyan800 Token = "cyan-800"
	Cyan900 Token = "cyan-900"
	Cyan950 Token = "cyan-950"

	Magenta50  Token = "magenta-50"
	Magenta100 Token = "magenta-100"
	Magenta200 Token = "magenta-200"
	Magenta300 Token = "magenta-300"
	Magenta400 Token = "magenta-400"
	Magenta500 Token = "magenta-500"
	Magenta600 Token = "magenta-600"
	Magenta700 Token = "magenta-700"
	Magenta800 Token = "magenta-800"
	Magenta900 Token = "magenta-900"
	Magenta950 Token = "magenta-950"
)

// Semantic aliases (Tailwind-ish semantics).
// Danger == Red, Warning == Yellow, Success == Green, Info == Cyan.
const (
	Danger50  Token = Red50
	Danger100 Token = Red100
	Danger200 Token = Red200
	Danger300 Token = Red300
	Danger400 Token = Red400
	Danger500 Token = Red500
	Danger600 Token = Red600
	Danger700 Token = Red700
	Danger800 Token = Red800
	Danger900 Token = Red900
	Danger950 Token = Red950

	Warning50  Token = Yellow50
	Warning100 Token = Yellow100
	Warning200 Token = Yellow200
	Warning300 Token = Yellow300
	Warning400 Token = Yellow400
	Warning500 Token = Yellow500
	Warning600 Token = Yellow600
	Warning700 Token = Yellow700
	Warning800 Token = Yellow800
	Warning900 Token = Yellow900
	Warning950 Token = Yellow950

	Success50  Token = Green50
	Success100 Token = Green100
	Success200 Token = Green200
	Success300 Token = Green300
	Success400 Token = Green400
	Success500 Token = Green500
	Success600 Token = Green600
	Success700 Token = Green700
	Success800 Token = Green800
	Success900 Token = Green900
	Success950 Token = Green950

	Info50  Token = Cyan50
	Info100 Token = Cyan100
	Info200 Token = Cyan200
	Info300 Token = Cyan300
	Info400 Token = Cyan400
	Info500 Token = Cyan500
	Info600 Token = Cyan600
	Info700 Token = Cyan700
	Info800 Token = Cyan800
	Info900 Token = Cyan900
	Info950 Token = Cyan950
)
