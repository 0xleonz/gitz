package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.com/0xleonz/gitz/internal/git"
	"gitlab.com/0xleonz/gitz/internal/inputs"
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
	Short: "Edit o shows the commitMessage.yml",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoRoot, err := git.FindRepoRoot()
		if err != nil {
			return fmt.Errorf("no se encontró la raíz del repo: %w", err)
		}
		path := filepath.Join(repoRoot, "commitMessage.yml")

		var msg types.CommitMessage
		if v := viper.Get("commitMessage"); v != nil {
			if m, ok := utils.DecodeCommitMessage(v); ok {
				msg = m
			}
		} else if !utils.IsEmptyOrMissing(path) {
			if m, err := utils.LoadCommitMessage(path); err == nil {
				msg = m
			}
		}

		// Flags que solo muestran
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

		// Flags que modifican directamente
		if addSubject {
			msg.Subject = inputs.Subject(msg.Subject)
			saveCommitMessage(path, msg)
			printCommitMessage(msg)
			return nil
		}
		if issueFlag {
			msg.Issue = inputs.Issue(msg.Issue)
			saveCommitMessage(path, msg)
			printCommitMessage(msg)
			return nil
		}
		if fixFlag {
			inputs.Fix(&msg)
			saveCommitMessage(path, msg)
			printCommitMessage(msg)
			return nil
		}
		if refactorFlag {
			inputs.Refactor(&msg)
			saveCommitMessage(path, msg)
			printCommitMessage(msg)
			return nil
		}
		if descriptionFlag {
			msg.Description = inputs.Description(msg.Description)
			saveCommitMessage(path, msg)
			printCommitMessage(msg)
			return nil
		}
		if addChangeInput != "" {
			if err := inputs.AddChange(&msg, addChangeInput); err != nil {
				return err
			}
			saveCommitMessage(path, msg)
			printCommitMessage(msg)
			return nil
		}
		if addChangesFlag {
			inputs.Cambios(&msg)
			saveCommitMessage(path, msg)
			printCommitMessage(msg)
			return nil
		}

		// Modo interactivo completo
		msg.Subject = inputs.Subject(msg.Subject)
		msg.Description = inputs.Description(msg.Description)
		inputs.Cambios(&msg)
		msg.Footer = inputs.Footer(msg.Footer)
		msg.Issue = inputs.Issue(msg.Issue)

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
