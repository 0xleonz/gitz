package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/0xleonz/gitz/internal/git"
	"gitlab.com/0xleonz/gitz/internal/utils"
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

		// --- Modo corto: git commit -m "<mensaje>"
		if shortMessage != "" {
			if !utils.IsEmptyOrMissing(path) {
				fmt.Println("⚠️ commitMessage.yml no está vacío. Podrías estar ignorando contenido valioso.")
				data, _ := os.ReadFile(path)
				fmt.Println("📝 Contenido actual:")
				fmt.Println(strings.Repeat("─", 40))
				fmt.Println(string(data))
				fmt.Println(strings.Repeat("─", 40))

				fmt.Print("¿Deseas limpiarlo ahora? (y/N): ")
				var resp string
				fmt.Scanln(&resp)
				resp = strings.ToLower(strings.TrimSpace(resp))

				if resp == "y" || resp == "s" {
					if err := utils.ResetCommitMessage(path); err != nil {
						fmt.Println("⚠️ No se pudo limpiar commitMessage.yml:", err)
					} else {
						fmt.Println("🧹 commitMessage.yml limpiado.")
					}
				} else {
					fmt.Println("⏭️ commitMessage.yml no fue limpiado.")
				}
			}
			return gitCommit(shortMessage)
		}

		// --- Modo normal: usar commitMessage.yml
		if utils.IsEmptyOrMissing(path) {
			fmt.Println("📝 commitMessage.yml vacío o no existe. Ejecutando 'gitz message'...")
			if err := utils.Call("message"); err != nil {
				return fmt.Errorf("error ejecutando gitz message: %w", err)
			}
		}

		msg, err := utils.LoadCommitMessage(path)
		if err != nil {
			return fmt.Errorf("error leyendo commitMessage.yml: %w", err)
		}

		formatted := utils.FormatCommitMessage(msg)
		if err := gitCommit(formatted); err != nil {
			return err
		}

		if err := utils.ResetCommitMessage(path); err != nil {
			fmt.Println("⚠️ No se pudo limpiar commitMessage.yml:", err)
		} else {
			fmt.Println("🧹 commitMessage.yml limpiado.")
		}

		return nil
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
