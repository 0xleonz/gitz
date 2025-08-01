package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// ChangedFiles devuelve una lista de archivos modificados o no trackeados
// (excluyendo los ignorados por .gitignore), con rutas absolutas basadas en
// la raíz del repositorio.
func ChangedFiles() ([]string, error) {
	root, err := FindRepoRoot()
	if err != nil {
		return nil, fmt.Errorf("no estás en un repositorio Git: %w", err)
	}

	cmd := exec.Command("git", "ls-files", "--modified", "--others", "--exclude-standard")
	cmd.Dir = root

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error ejecutando git ls-files: %w", err)
	}

	output := strings.TrimSpace(out.String())
	if output == "" {
		return []string{}, nil
	}

	lines := strings.Split(output, "\n")
	for i := range lines {
		lines[i] = filepath.Join(root, lines[i]) // rutas absolutas
	}

	return lines, nil
}

// GitAddFile ejecuta `git add <file>` desde la raíz del repositorio.
func GitAddFile(file string) error {
	root, err := FindRepoRoot()
	if err != nil {
		return fmt.Errorf("no estás en un repositorio Git: %w", err)
	}

	cmd := exec.Command("git", "add", file)
	cmd.Dir = root

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error al hacer git add %s: %w", file, err)
	}
	return nil
}

// GitAddAll ejecuta `git add .` desde la raíz del repositorio.
func GitAddAll() error {
	root, err := FindRepoRoot()
	if err != nil {
		return fmt.Errorf("no estás en un repositorio Git: %w", err)
	}

	cmd := exec.Command("git", "add", ".")
	cmd.Dir = root

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error al hacer git add .: %w", err)
	}
	return nil
}
