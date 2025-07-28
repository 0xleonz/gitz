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
	return err != nil || strings.TrimSpace(string(data)) == ""
}

func MarshalYAML(msg types.CommitMessage) ([]byte, error) {
	return yaml.Marshal(&msg)
}
