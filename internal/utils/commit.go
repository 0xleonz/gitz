package utils

import (
	"gitlab.com/0xleonz/gitz/internal/types"
	"gopkg.in/yaml.v3"
)

// DecodeCommitMessage convierte un objeto genérico a types.CommitMessage
func DecodeCommitMessage(v interface{}) (types.CommitMessage, bool) {
	raw, ok := v.(map[string]interface{})
	if !ok {
		return types.CommitMessage{}, false
	}

	data, err := yaml.Marshal(raw)
	if err != nil {
		return types.CommitMessage{}, false
	}

	var msg types.CommitMessage
	if err := yaml.Unmarshal(data, &msg); err != nil {
		return types.CommitMessage{}, false
	}

	return msg, true
}
