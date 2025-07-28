package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gitlab.com/0xleonz/gitz/internal/git"
	"gitlab.com/0xleonz/gitz/internal/inputs"
	"gitlab.com/0xleonz/gitz/internal/types"
	"gitlab.com/0xleonz/gitz/internal/utils"
	"gopkg.in/yaml.v3"
)

type CommitMessage struct {
	Changes     []types.Change    `yaml:"changes"`
	Issue       string            `yaml:"issue"`
	Subject     string            `yaml:"subject"`
	Description []string          `yaml:"description"`
	Footer      map[string]string `yaml:"footer"`
}

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Crea y maneja mensajes de commit enriquecidos",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTextEditor()
	},
}

func runTextEditor() error {
	repoRoot, err := git.FindRepoRoot()
	if err != nil {
		return fmt.Errorf("no se encontró el repositorio git: %w", err)
	}
	path := filepath.Join(repoRoot, "commitMessage.yml")

	msg := CommitMessage{
		Changes:     []types.Change{},
		Description: []string{},
		Footer:      map[string]string{},
	}

	// Cargar YAML existente si hay
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err == nil {
			_ = yaml.Unmarshal(data, &msg)
		}
	}

	// Subject, Description
	msg.Subject = inputs.Subject(msg.Subject)
	msg.Description = inputs.Description(msg.Description)
	msg.Changes = inputs.Changes(msg.Changes)
	msg.Footer = inputs.Footer(msg.Footer)
	msg.Issue = inputs.Issue(msg.Issue)

	// Guardar archivo
	yamlData, err := yaml.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error serializando YAML: %w", err)
	}
	if err := os.WriteFile(path, yamlData, 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo: %w", err)
	}

	fmt.Println(utils.Colorize("\n✅ Mensaje de commit guardado", utils.Green))
	return nil
}

func init() {
	rootCmd.AddCommand(messageCmd)
}
