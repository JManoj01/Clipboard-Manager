package ui

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
)

func colorize(color, text string) string {
	return color + text + ColorReset
}

func bold(text string) string {
	return ColorBold + text + ColorReset
}

func success(text string) string {
	return ColorGreen + "✓ " + text + ColorReset
}

func errText(text string) string { 
	return ColorRed + "✗ " + text + ColorReset
}

func info(text string) string {
	return ColorBlue + "ℹ " + text + ColorReset
}

func warning(text string) string {
	return ColorYellow + "⚠ " + text + ColorReset
}
