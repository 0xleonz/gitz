package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"gitlab.com/0xleonz/gitz/internal/config"
	"gitlab.com/0xleonz/gitz/internal/utils"
)

var (
	filterRemote string
	filterBranch string
	confirmPush  bool
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
			if confirmPush {
				fmt.Printf("¿Push a %s:%s? (y/n): ", pair.Remote, pair.Branch)
				if utils.Ask() != "y" {
					continue
				}
			}

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
					Branch: b,
					Remote: r,
				})
			}
		}
	}

	return pairs
}

func filterPairs(pairs []PushPair, remote, branch string) []PushPair {
	var filtered []PushPair
	for _, p := range pairs {
		if (remote == "" || p.Remote == remote) &&
			(branch == "" || p.Branch == branch) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("no se encontró .git en ningún directorio padre")
}

func init() {
	pushCmd.Flags().StringVar(&filterRemote, "remote", "", "Filtra por remoto")
	pushCmd.Flags().StringVar(&filterBranch, "branch", "", "Filtra por rama")
	pushCmd.Flags().BoolVarP(&confirmPush, "confirm", "c", false, "Pide confirmación antes de cada push")

	rootCmd.AddCommand(pushCmd)
}
