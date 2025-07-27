package git

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// AppendToGitignore agrega una entrada a .gitignore si no existe aún.
func AppendToGitignore(entry string) error {
	entry = strings.TrimSpace(entry)
	if entry == "" {
		return nil
	}

	// Leer entradas actuales
	lines := make(map[string]bool)
	if f, err := os.Open(".gitignore"); err == nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lines[strings.TrimSpace(scanner.Text())] = true
		}
		f.Close()
	}

	if lines[entry] {
		fmt.Println("⚠️ La entrada ya está en .gitignore:", entry)
		return nil
	}

	// Agregar entrada nueva
	f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("\n" + entry + "\n")
	return err
}
