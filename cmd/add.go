package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/0xleonz/gitz/internal/git"
	"gitlab.com/0xleonz/gitz/internal/utils"
)

var confirmAdd bool

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Agrega archivos al staging con confirmación opcional",
	RunE: func(cmd *cobra.Command, args []string) error {
		if confirmAdd {
			// obtener archivos modificados, sin stage
			files, err := git.ChangedFiles()
			if err != nil {
				return fmt.Errorf("no se pudo obtener lista de archivos: %w", err)
			}

			if len(files) == 0 {
				fmt.Println("✅ No hay archivos modificados para agregar.")
				return nil
			}

			for _, f := range files {
				fmt.Printf("¿Agregar %s? (y/n/a para ignorar): ", f)
				switch utils.Ask() {
				case "y":
					if err := git.GitAddFile(f); err != nil {
						fmt.Printf("❌ Error al agregar %s: %v\n", f, err)
					}
				case "a":
					if err := git.AppendToGitignore(f); err != nil {
						fmt.Printf("❌ Error al ignorar %s: %v\n", f, err)
					} else {
						fmt.Printf("🟡 %s agregado a .gitignore\n", f)
					}
				case "n":
					// skip
				default:
					fmt.Println("⚠️ Opción inválida. Usando 'n' por defecto.")
				}
			}
			return nil
		}

		// sin confirm, agregar todo
		return git.GitAddAll()
	},
}

func init() {
	addCmd.Flags().BoolVarP(&confirmAdd, "confirm", "c", false, "Pide confirmación por cada archivo")
	rootCmd.AddCommand(addCmd)
}
