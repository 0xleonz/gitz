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

// ReadGitignoreLines devuelve un map con las líneas actuales del .gitignore
func ReadGitignoreLines() (map[string]bool, error) {
	lines := make(map[string]bool)

	f, err := os.Open(".gitignore")
	if err != nil {
		return lines, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines[line] = true
		}
	}
	return lines, nil
}

func MergeGitignoreWithTemplate(templatePath string) error {
	existing, _ := ReadGitignoreLines()

	tplData, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("no se pudo leer el template .gitignore: %w", err)
	}

	tplLines := strings.Split(string(tplData), "\n")
	var toAdd []string

	for _, line := range tplLines {
		line = strings.TrimSpace(line)
		if line == "" || existing[line] {
			continue
		}
		toAdd = append(toAdd, line)
	}

	if len(toAdd) == 0 {
		fmt.Println("✅ .gitignore ya está completo.")
		return nil
	}

	f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("no se pudo abrir .gitignore para escribir: %w", err)
	}
	defer f.Close()

	// Agrupamos con comentario si hay más de 1 línea
	if len(toAdd) > 0 {
		f.WriteString("\n# Añadido automáticamente desde template por gitz init\n")
		for _, line := range toAdd {
			f.WriteString(line + "\n")
		}
	}

	fmt.Printf("📝 %d nuevas entradas agregadas a .gitignore\n", len(toAdd))
	return nil
}
