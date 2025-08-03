package utils

import (
	"os"
	"strings"

	"gitlab.com/0xleonz/gitz/internal/types"
	"gopkg.in/yaml.v3"
)

func LoadCommitMessage(path string) (types.CommitMessage, error) {
	var msg types.CommitMessage
	data, err := os.ReadFile(path)
	if err != nil {
		return msg, err
	}
	err = yaml.Unmarshal(data, &msg)
	return msg, err
}

func IsEmptyOrMissing(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return true // archivo no existe o no se puede leer
	}

	var msg types.CommitMessage
	if err := yaml.Unmarshal(data, &msg); err != nil {
		return true // no se puede parsear => consideramos inválido
	}

	return len(msg.Changes) == 0 &&
		strings.TrimSpace(msg.Issue) == "" &&
		strings.TrimSpace(msg.Subject) == "" &&
		len(msg.Description) == 0 &&
		len(msg.Footer) == 0
}

func MarshalYAML(msg types.CommitMessage) ([]byte, error) {
	return yaml.Marshal(&msg)
}
