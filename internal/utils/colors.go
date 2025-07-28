// Catpuccin colors
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
