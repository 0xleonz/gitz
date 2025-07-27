package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/0xleonz/gitz/internal/utils"
)

var quickConfirm bool
var quickShortMsg string

var quickCmd = &cobra.Command{
	Use:   "quick",
	Short: "check + add + commit + push en una sola línea",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Paso 1: check
		if err := utils.Call("check"); err != nil {
			fmt.Println("🟡 No hay cambios o check falló. Abortando.")
			return nil
		}

		// Paso 2: add
		addArgs := []string{"add"}
		if quickConfirm {
			addArgs = append(addArgs, "--confirm")
		}
		if err := utils.Call(addArgs...); err != nil {
			fmt.Println("🛑 gitz add falló o fue cancelado.")
			return nil
		}

		// Paso 3: commit
		commitArgs := []string{"commit"}
		if quickShortMsg != "" {
			commitArgs = append(commitArgs, "--short", quickShortMsg)
		}
		if err := utils.Call(commitArgs...); err != nil {
			fmt.Println("❌ Error al hacer commit:", err)
			return err
		}

		// Paso 4: push
		pushArgs := []string{"push"}
		if quickConfirm {
			pushArgs = append(pushArgs, "--confirm")
		}
		return utils.Call(pushArgs...)
	},
}

func init() {
	quickCmd.Flags().BoolVarP(&quickConfirm, "confirm", "c", false, "Pide confirmación antes de add/push")
	quickCmd.Flags().StringVarP(&quickShortMsg, "short", "s", "", "Mensaje corto tipo git commit -m")
	rootCmd.AddCommand(quickCmd)
}
