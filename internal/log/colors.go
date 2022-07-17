package log

var (
	none = "\033[0m"
	red  = "\033[31m"
	// green  = "\033[32m".
	yellow = "\033[33m"
	blue   = "\033[34m"
	// purple = "\033[35m".
	cyan = "\033[36m"
	gray = "\033[37m"
	// white  = "\033[97m".
)

var (
	// ColorNone is the color of normal stdout.
	ColorNone = none

	// ColorDebug is the color sequence uses for
	// debug output.
	ColorDebug = blue

	// ColorInfo is the color sequence uses for
	// info output.
	ColorInfo = cyan

	// ColorWarn is the color sequence uses for
	// warn output.
	ColorWarn = yellow

	// ColorError is the color sequence uses for
	// error output.
	ColorError = red

	// ColorNotice is the color sequence uses for
	// notice output.
	ColorNotice = gray
)
