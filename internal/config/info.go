package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type RemoteBranch struct {
	Branch string `yaml:"branch"`
	Remote string `yaml:"remote"`
}

type Info struct {
	Description    string         `yaml:"description"`
	Ramas          []string       `yaml:"ramas"`
	Branches       []string       `yaml:"branches"`
	RemoteBranches []RemoteBranch `yaml:"remote-branches"`
}

func LoadInfo(path string) (*Info, error) {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("error leyendo info.yml: %w", err)
	}

	var info Info
	if err := yaml.Unmarshal(content, &info); err != nil {
		return nil, fmt.Errorf("error parseando YAML: %w", err)
	}

	return &info, nil
}

