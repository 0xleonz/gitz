package utils

import (
	"fmt"
	"os"

	"gitlab.com/0xleonz/gitz/internal/types"
	"gopkg.in/yaml.v3"
)

func ResetCommitMessage(path string) error {
	empty := types.CommitMessage{
		Changes:     []types.Change{},
		Issue:       "",
		Subject:     "",
		Description: []string{},
		Footer:      map[string]string{},
	}

	data, err := yaml.Marshal(&empty)
	if err != nil {
		return fmt.Errorf("error serializando YAML vacío: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo limpio: %w", err)
	}

	return nil
}
