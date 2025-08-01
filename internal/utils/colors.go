// Catpuccin mocha colors
package utils

import "fmt"

var (
	Purple = "\033[38;5;140m"
	Cyan   = "\033[38;5;117m"
	Yellow = "\033[38;5;179m"
	Green  = "\033[38;5;114m"
	Pink   = "\033[38;5;175m"
	Reset  = "\033[0m"
	Red    = "\033[38;5;204m"
	Blue   = "\033[38;5;110m"
)

// Apply color to text
func Colorize(text, color string) string {
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}

// Print a list of strings in the same line with color, wrapping at maxColumn characters
func OnlineWithColor(text []string, color string, separator string, maxColumn int) string {
	var result string
	lineLen := 0

	for i, word := range text {
		colored := Colorize(word, color)
		sep := ""
		if i != len(text)-1 {
			sep = separator
		}
		segment := colored + sep
		segmentLen := len(word) + len(separator)

		if lineLen+segmentLen > maxColumn && lineLen > 0 {
			result += "\n"
			lineLen = 0
		}
		result += segment
		lineLen += segmentLen
	}

	return result
}
