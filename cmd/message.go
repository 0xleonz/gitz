package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/0xleonz/gitz/internal/git"
	"gitlab.com/0xleonz/gitz/internal/types"
	"gitlab.com/0xleonz/gitz/internal/utils"
)

var (
	showFlag        bool
	rawFlag         bool
	verboseFlag     bool
	subjectFlag     bool
	changesFlag     bool
	addSubject      bool
	descriptionFlag bool
	issueFlag       bool
	fixFlag         bool
	refactorFlag    bool
	addChangeInput  string
	addChangesFlag  bool
)

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Edita o muestra el commitMessage.yml",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoRoot, err := git.FindRepoRoot()
		if err != nil {
			return fmt.Errorf("no se encontró la raíz del repo: %w", err)
		}
		path := filepath.Join(repoRoot, "commitMessage.yml")

		var msg types.CommitMessage
		if v := viper.Get("commitMessage"); v != nil {
			m, ok := utils.DecodeCommitMessage(v)
			if ok {
				msg = m
			}
		} else if !utils.IsEmptyOrMissing(path) {
			m, err := utils.LoadCommitMessage(path)
			if err == nil {
				msg = m
			}
		}

		// Mostrar si aplica
		if showFlag {
			printCommitMessage(msg)
			return nil
		}
		if rawFlag {
			fmt.Println(utils.FormatCommitMessage(msg))
			return nil
		}
		if subjectFlag {
			if msg.Subject != "" {
				fmt.Println(utils.Colorize("• Subject:", utils.Yellow), msg.Subject)
			}
			return nil
		}
		if changesFlag {
			if len(msg.Changes) > 0 {
				fmt.Println(utils.Colorize("🔧 Cambios:", utils.Cyan))
				for _, change := range msg.Changes {
					fmt.Printf("  - %s: %s\n", utils.Colorize(change.Type, utils.Cyan), change.Summary)
				}
			}
			return nil
		}

		r := bufio.NewReader(os.Stdin)

		if addSubject {
			msg.Subject = promptWithDefault("📝 Subject:", utils.Pink, msg.Subject, "feat: agregar validación de email")
			saveCommitMessage(path, msg)
			printCommitMessage(msg)
			return nil
		}
		if issueFlag {
			msg.Issue = promptWithDefault("🔗 Issue:", utils.Purple, msg.Issue, "#123")
			saveCommitMessage(path, msg)
			printCommitMessage(msg)
			return nil
		}
		if fixFlag {
			msg.Changes = append(msg.Changes, types.Change{Type: "fix", Summary: prompt("Resumen del fix (ej: arregla validación de token):")})
			saveCommitMessage(path, msg)
			printCommitMessage(msg)
			return nil
		}
		if refactorFlag {
			msg.Changes = append(msg.Changes, types.Change{Type: "refactor", Summary: prompt("Resumen del refactor (ej: simplifica lógica de parser):")})
			saveCommitMessage(path, msg)
			printCommitMessage(msg)
			return nil
		}
		if descriptionFlag {
			fmt.Println(utils.Colorize("📄 Description actual:", utils.Yellow))
			for _, d := range msg.Description {
				fmt.Println("  " + d)
			}
			fmt.Println(utils.Colorize("📄 Agrega nuevas líneas (ENTER vacío para terminar):", utils.Green))
			for {
				fmt.Print("  > ")
				line, _ := r.ReadString('\n')
				line = strings.TrimSpace(line)
				if line == "" {
					break
				}
				msg.Description = append(msg.Description, line)
			}
			saveCommitMessage(path, msg)
			printCommitMessage(msg)
			return nil
		}
		if addChangeInput != "" {
			parts := strings.SplitN(addChangeInput, ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("formato inválido para --add-change. Usa tipo:resumen")
			}
			msg.Changes = append(msg.Changes, types.Change{
				Type:    strings.TrimSpace(parts[0]),
				Summary: strings.TrimSpace(parts[1]),
			})
			saveCommitMessage(path, msg)
			printCommitMessage(msg)
			return nil
		}
		if addChangesFlag {
			fmt.Println(utils.Colorize("🔧 Cambios actuales:", utils.Yellow))
			for _, c := range msg.Changes {
				fmt.Printf("  - %s: %s\n", c.Type, c.Summary)
			}
			fmt.Println(utils.Colorize("🔧 Agrega cambios nuevos (tipo:resumen, ENTER vacío para terminar):", utils.Cyan))
			for {
				fmt.Print("  - ")
				line, _ := r.ReadString('\n')
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
			saveCommitMessage(path, msg)
			printCommitMessage(msg)
			return nil
		}

		// Modo completo interactivo
		msg.Subject = promptWithDefault("📝 Subject:", utils.Pink, msg.Subject, "feat: agregar validación")

		fmt.Println(utils.Colorize("📄 Description (ENTER vacío para terminar):", utils.Green))
		for _, d := range msg.Description {
			fmt.Println("  [actual] " + utils.Colorize(d, utils.Yellow))
		}
		desc := []string{}
		for {
			fmt.Print("  > ")
			line, _ := r.ReadString('\n')
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			desc = append(desc, line)
		}
		msg.Description = append(msg.Description, desc...)

		fmt.Println(utils.Colorize("🔧 Changes (tipo:resumen, ENTER vacío para terminar):", utils.Cyan))
		for _, c := range msg.Changes {
			fmt.Printf("  [actual] - %s: %s\n", utils.Colorize(c.Type, utils.Yellow), c.Summary)
		}
		for {
			fmt.Print("  - ")
			line, _ := r.ReadString('\n')
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
		for k, v := range msg.Footer {
			fmt.Printf("  [actual] %s: %s\n", k, v)
		}
		newFooter := msg.Footer
		for {
			fmt.Print("  - ")
			line, _ := r.ReadString('\n')
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				fmt.Println(utils.Colorize("  ⚠️  Formato inválido. Usa clave:valor", utils.Red))
				continue
			}
			newFooter[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
		msg.Footer = newFooter

		msg.Issue = promptWithDefault("🔗 Issue:", utils.Purple, msg.Issue, "#123")

		saveCommitMessage(path, msg)
		printCommitMessage(msg)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(messageCmd)
	messageCmd.Flags().BoolVarP(&showFlag, "show", "s", false, "Muestra el mensaje de commit formateado")
	messageCmd.Flags().BoolVar(&rawFlag, "raw", false, "Muestra el mensaje final plano")
	messageCmd.Flags().BoolVar(&verboseFlag, "verbose", false, "Imprime el mensaje después de editar")
	messageCmd.Flags().BoolVar(&subjectFlag, "subject", false, "Muestra solo el subject")
	messageCmd.Flags().BoolVar(&changesFlag, "changes", false, "Muestra solo los cambios")
	messageCmd.Flags().BoolVar(&addSubject, "add-subject", false, "Modifica solo el subject")
	messageCmd.Flags().BoolVar(&descriptionFlag, "description", false, "Edita la descripción")
	messageCmd.Flags().BoolVar(&issueFlag, "issue", false, "Edita solo el issue")
	messageCmd.Flags().BoolVar(&fixFlag, "fix", false, "Agrega un cambio de tipo fix")
	messageCmd.Flags().BoolVar(&refactorFlag, "refactor", false, "Agrega un cambio de tipo refactor")
	messageCmd.Flags().StringVar(&addChangeInput, "add-change", "", "Agrega un cambio en formato tipo:resumen")
	messageCmd.Flags().BoolVar(&addChangesFlag, "add-changes", false, "Agrega múltiples cambios de forma interactiva")
}

func prompt(label string) string {
	fmt.Print(utils.Colorize(label+" ", utils.Pink))
	r := bufio.NewReader(os.Stdin)
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}

func promptWithDefault(label, color, def, example string) string {
	if def != "" {
		fmt.Printf("%s %s %s\n", utils.Colorize(label, color), utils.Colorize("[actual]", utils.Yellow), def)
	} else if example != "" {
		fmt.Printf("%s %s\n", utils.Colorize(label, color), utils.Colorize("[ej: "+example+"]", utils.Blue))
	}
	fmt.Print("  > ")
	r := bufio.NewReader(os.Stdin)
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}

func saveCommitMessage(path string, msg types.CommitMessage) {
	data, err := utils.MarshalYAML(msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, utils.Colorize("❌ Error serializando commitMessage.yml: "+err.Error(), utils.Red))
		return
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		fmt.Fprintln(os.Stderr, utils.Colorize("❌ Error guardando commitMessage.yml: "+err.Error(), utils.Red))
		return
	}
	fmt.Println(utils.Colorize("✅ commitMessage.yml actualizado.", utils.Green))
}

func printCommitMessage(msg types.CommitMessage) {
	fmt.Println(utils.Colorize("┌────────────────────┐", utils.Purple))
	fmt.Println(utils.Colorize("│ Mensaje de Commit  │", utils.Purple))
	fmt.Println(utils.Colorize("└────────────────────┘", utils.Purple))

	if msg.Subject != "" {
		fmt.Println(utils.Colorize("• Subject:", utils.Yellow), msg.Subject)
	}
	if msg.Issue != "" {
		fmt.Println(utils.Colorize("• Issue:", utils.Cyan), msg.Issue)
	}
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
