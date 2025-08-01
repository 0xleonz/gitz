package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/0xleonz/gitz/internal/git"
	"gitlab.com/0xleonz/gitz/internal/utils"
)

var (
	confirmAdd bool
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Agrega archivos al staging con confirmación, dry-run o modo detallado",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Caso: argumentos directos, por ejemplo: gitz add archivo.txt
		if len(args) > 0 {
			for _, file := range args {
				if dryRun {
					fmt.Printf("🔍 Agregando: agregar %s\n", file)
					continue
				}
				if err := git.GitAddFile(file); err != nil {
					return fmt.Errorf("error al agregar %s: %w", file, err)
				}
				if verbose {
					fmt.Println("🟢 Agregado:", file)
				}
			}
			return nil
		}

		// Modo confirmación
		if confirmAdd {
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
					if dryRun {
						fmt.Println("🔍 Agregado:", f)
					} else if err := git.GitAddFile(f); err != nil {
						fmt.Printf("❌ Error al agregar %s: %v\n", f, err)
					} else if verbose {
						fmt.Println("🟢 Agregado:", f)
					}
				case "a":
					if dryRun {
						fmt.Println("🔍 Simulado ignorar:", f)
					} else if err := git.AppendToGitignore(f); err != nil {
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

		// Por defecto, agregar todo
		if dryRun {
			files, err := git.ChangedFiles()
			if err != nil {
				return fmt.Errorf("no se pudo obtener lista de archivos: %w", err)
			}
			for _, f := range files {
				fmt.Println("🔍 Agregando:", f)
			}
			return nil
		}

		if err := git.GitAddAll(); err != nil {
			return fmt.Errorf("error al agregar todos los archivos: %w", err)
		}

		if verbose {
			fmt.Println("🟢 Todos los archivos agregados con éxito.")
		}

		return nil
	},
}

func init() {
	addCmd.Flags().BoolVarP(&confirmAdd, "confirm", "c", false, "Pide confirmación por cada archivo")
	rootCmd.AddCommand(addCmd)
}
