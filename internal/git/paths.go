package git

import (
	"fmt"
	"os"
	"path/filepath"
)

// just going up to find fir .git
func FindRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		gitPath := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			return dir, nil
		}
		if parent := filepath.Dir(dir); parent == dir {
			break
		} else {
			dir = parent
		}
	}
	return "", fmt.Errorf("no se encontró .git en ningún directorio padre")
}
