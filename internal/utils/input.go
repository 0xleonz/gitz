package utils

import (
	"bufio"
	"os"
	"strings"
)

func Ask() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(text))
}
