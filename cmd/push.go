package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"gitlab.com/0xleonz/gitz/internal/config"
)

var (
	dryRun  bool
	verbose bool
	filterRemote string
	filterBranch string
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Empuja ramas del repo al remoto configurado en info.yml",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoRoot, err := findRepoRoot()
		if err != nil {
			return fmt.Errorf("no se pudo encontrar la raíz del repo: %w", err)
		}

		infoPath := filepath.Join(repoRoot, "info.yml")
		info, err := config.LoadInfo(infoPath)
		if err != nil {
			return fmt.Errorf("error cargando info.yml: %w", err)
		}

		pairs := buildPushPairs(info)

		// aplicar filtros si existen
		if filterRemote != "" || filterBranch != "" {
			pairs = filterPairs(pairs, filterRemote, filterBranch)
		}

		if len(pairs) == 0 {
			fmt.Println("No hay ramas/remotos para hacer push.")
			return nil
		}

		for _, pair := range pairs {
			fmt.Printf("[gitz] git push %s %s\n", pair.Remote, pair.Branch)
			if dryRun {
				continue
			}

			cmd := exec.Command("git", "push", pair.Remote, pair.Branch)
			cmd.Dir = repoRoot // asegura que se ejecuta desde la raíz
			out, err := cmd.CombinedOutput()
			if verbose || err != nil {
				fmt.Println(string(out))
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error al hacer push de %s a %s: %v\n", pair.Branch, pair.Remote, err)
			}
		}

		return nil
	},
}

type PushPair struct {
	Branch string
	Remote string
}

func buildPushPairs(info *config.Info) []PushPair {
	var pairs []PushPair

	// Si se define remote-branches, respétalo
	if len(info.RemoteBranches) > 0 {
		for _, rb := range info.RemoteBranches {
			pairs = append(pairs, PushPair{
				Branch: rb.Branch,
				Remote: rb.Remote,
			})
		}
	} else {
		// Default: todas las combinaciones rama x remoto
		for _, r := range info.Ramas {
			for _, b := range info.Branches {
				pairs = append(pairs, PushPair{
