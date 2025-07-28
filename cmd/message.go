package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/0xleonz/gitz/internal/git"
	"gitlab.com/0xleonz/gitz/internal/types"
	"gitlab.com/0xleonz/gitz/internal/utils"
)

var showFlag bool
var rawFlag bool

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Edita o muestra el commitMessage.yml",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showFlag && rawFlag {
			return fmt.Errorf("los flags --show y --raw no pueden usarse juntos")
		}

		repoRoot, err := git.FindRepoRoot()
		if err != nil {
			return fmt.Errorf("no se encontró la raíz del repo: %w", err)
		}
		path := filepath.Join(repoRoot, "commitMessage.yml")

		if showFlag || rawFlag {
			msg, err := utils.LoadCommitMessage(path)
			if err != nil {
				return fmt.Errorf("no se pudo cargar commitMessage.yml: %w", err)
			}

			if rawFlag {
				fmt.Println(utils.FormatCommitMessage(msg))
			} else {
				printCommitMessage(msg)
			}
			return nil
		}

		// Modo edición interactivo
		msg := types.CommitMessage{}
		if !utils.IsEmptyOrMissing(path) {
			existing, err := utils.LoadCommitMessage(path)
			if err == nil {
				msg = existing
			}
		}

		reader := bufio.NewReader(os.Stdin)

		fmt.Print(utils.Colorize("📝 Subject: ", utils.Pink))
		msg.Subject = readLineOrDefault(reader, msg.Subject)

		fmt.Println(utils.Colorize("📄 Description (ENTER vacío para terminar):", utils.Green))
		msg.Description = []string{}
		for {
			fmt.Print("  > ")
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			msg.Description = append(msg.Description, line)
		}

		fmt.Println(utils.Colorize("🔧 Changes (tipo:resumen, ENTER vacío para terminar):", utils.Cyan))
		msg.Changes = []types.Change{}
		for {
			fmt.Print("  - ")
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				fmt.Println(utils.Colorize("  ⚠️  Formato inválido. Usa tipo:resumen", utils.Red))
				continue
			}
			msg.Changes = append(msg.Changes, types.Change{
				Type:    strings.TrimSpace(parts[0]),
				Summary: strings.TrimSpace(parts[1]),
			})
		}

		fmt.Println(utils.Colorize("🔻 Footer (clave:valor, ENTER vacío para terminar):", utils.Yellow))
		msg.Footer = map[string]string{}
		for {
			fmt.Print("  - ")
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				fmt.Println(utils.Colorize("  ⚠️  Formato inválido. Usa clave:valor", utils.Red))
				continue
			}
			msg.Footer[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}

		fmt.Print(utils.Colorize("🔗 Issue (ID o referencia): ", utils.Purple))
		msg.Issue = readLineOrDefault(reader, msg.Issue)

		data, err := utils.MarshalYAML(msg)
		if err != nil {
			return fmt.Errorf("no se pudo serializar commitMessage.yml: %w", err)
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			return fmt.Errorf("no se pudo guardar commitMessage.yml: %w", err)
		}

		fmt.Println(utils.Colorize("✅ commitMessage.yml actualizado.", utils.Green))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(messageCmd)
	messageCmd.Flags().BoolVarP(&showFlag, "show", "s", false, "Muestra el mensaje de commit formateado")
	messageCmd.Flags().BoolVarP(&rawFlag, "raw", "r", false, "Muestra el mensaje final plano (estilo git commit -m)")
}

func readLineOrDefault(reader *bufio.Reader, def string) string {
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return def
	}
	return line
}

func printCommitMessage(msg types.CommitMessage) {
	fmt.Println(utils.Colorize("┌────────────────────┐", utils.Purple))
	fmt.Println(utils.Colorize("│ Mensaje de Commit │", utils.Purple))
	fmt.Println(utils.Colorize("└────────────────────┘", utils.Purple))

	fmt.Println(utils.Colorize("• Subject:", utils.Yellow), msg.Subject)
	fmt.Println(utils.Colorize("• Issue:", utils.Cyan), msg.Issue)

	if len(msg.Description) > 0 {
		fmt.Println(utils.Colorize("\n📄 Descripción:", utils.Green))
		for _, line := range msg.Description {
			fmt.Println("  " + line)
		}
	}

	if len(msg.Changes) > 0 {
		fmt.Println(utils.Colorize("\n🔧 Cambios:", utils.Pink))
		for _, change := range msg.Changes {
			fmt.Printf("  - %s: %s\n", utils.Colorize(change.Type, utils.Cyan), change.Summary)
		}
	}

	if len(msg.Footer) > 0 {
		fmt.Println(utils.Colorize("\n🔻 Footer:", utils.Yellow))
		for k, v := range msg.Footer {
			fmt.Printf("  %s: %s\n", utils.Colorize(k, utils.Green), v)
		}
	}
}
