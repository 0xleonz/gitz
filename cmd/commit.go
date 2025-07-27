package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/0xleonz/gitz/internal/git"
	"gitlab.com/0xleonz/gitz/internal/types"
	"gitlab.com/0xleonz/gitz/internal/utils"
	"gopkg.in/yaml.v3"
)

var shortMessage string

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Hace commit usando commitMessage.yml o un mensaje corto",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoRoot, err := git.FindRepoRoot()
		if err != nil {
			return fmt.Errorf("no se encontró la raíz del repo: %w", err)
		}

		path := filepath.Join(repoRoot, "commitMessage.yml")

		if shortMessage != "" {
			return gitCommit(shortMessage)
		}

		if isEmptyOrMissing(path) {
			fmt.Println("📝 commitMessage.yml vacío o no existe. Ejecutando 'gitz message'...")
			if err := utils.Call("message"); err != nil {
				return fmt.Errorf("error ejecutando gitz message: %w", err)
			}
		}

		msg, err := loadCommitMessage(path)
		if err != nil {
			return fmt.Errorf("error leyendo commitMessage.yml: %w", err)
		}

		formatted := formatCommitMessage(msg)
		return gitCommit(formatted)
	},
}

func init() {
	commitCmd.Flags().StringVarP(&shortMessage, "short", "s", "", "Mensaje corto tipo git commit -m")
	rootCmd.AddCommand(commitCmd)
}

func gitCommit(msg string) error {
	cmd := exec.Command("git", "commit", "-m", msg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isEmptyOrMissing(path string) bool {
	data, err := os.ReadFile(path)
	return err != nil || strings.TrimSpace(string(data)) == ""
}

func loadCommitMessage(path string) (types.CommitMessage, error) {
	var msg types.CommitMessage
	data, err := os.ReadFile(path)
	if err != nil {
		return msg, err
	}
	err = yaml.Unmarshal(data, &msg)
	return msg, err
}

func formatCommitMessage(msg types.CommitMessage) string {
	var sb strings.Builder

	sb.WriteString(msg.Subject)

	if len(msg.Description) > 0 {
		sb.WriteString("\n\n" + strings.Join(msg.Description, "\n"))
	}

	if len(msg.Changes) > 0 {
		sb.WriteString("\n\nChanges:")
		for _, c := range msg.Changes {
			sb.WriteString(fmt.Sprintf("\n- %s: %s", c.Type, c.Summary))
		}
	}

	if len(msg.Footer) > 0 {
		sb.WriteString("\n")
		for k, v := range msg.Footer {
			sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
	}

	if msg.Issue != "" {
		sb.WriteString("\nIssue: " + msg.Issue)
	}

	return sb.String()
}
