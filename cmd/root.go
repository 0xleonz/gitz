package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/0xleonz/gitz/internal/git"
	"gopkg.in/yaml.v3"
)

var (
	cfgFile       string
	infoFile      string
	commitMsgFile string
	dryRun        bool
	verbose       bool
)

var rootCmd = &cobra.Command{
	Use:   "gitz",
	Short: "gitz: un wrapper ligero de git",
	Long:  "gitz es un CLI para interactuar con git más rápido.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/gitz/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Simula acciones sin aplicar cambios")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Muestra operaciones realizadas")
}

func initConfig() {
	// Config global en ~/.config/gitz
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		cfgDir := filepath.Join(home, ".config", "gitz")
		viper.AddConfigPath(cfgDir)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.SetDefault("templates_dir", filepath.Join(cfgDir, "templates"))
		viper.SetDefault("defaults.gitignore", "gitignore")
		viper.SetDefault("defaults.info.yml", "info")
		viper.SetDefault("licenses", map[string]string{
			"GNU GPLv3": "licenses/GPL-3.0.txt",
			"MIT":       "licenses/MIT.txt",
		})
		viper.SetDefault("default_license", "GNU GPLv3")

		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintln(os.Stderr, "⚠️  no se encontró config, creando en", cfgDir)
			os.MkdirAll(cfgDir, 0o755)
			_ = viper.WriteConfigAs(filepath.Join(cfgDir, "config.yaml"))
		}
	}

	viper.AutomaticEnv()

	// Cargar info.yml y commitMessage.yml desde la raíz del repo
	repoRoot, err := git.FindRepoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, "⚠️  No estás en un repo git:", err)
		return
	}

	// Ruta a info.yml
	infoFile = filepath.Join(repoRoot, "info.yml")
	if _, err := os.Stat(infoFile); err == nil {
		var infoData map[string]interface{}
		content, err := os.ReadFile(infoFile)
		if err == nil && yaml.Unmarshal(content, &infoData) == nil {
			viper.Set("info", infoData)
		} else {
			fmt.Fprintln(os.Stderr, "⚠️  No se pudo leer info.yml:", err)
		}
	}

	// Ruta a commitMessage.yml
	commitMsgFile = filepath.Join(repoRoot, "commitMessage.yml")
	if _, err := os.Stat(commitMsgFile); err == nil {
		var msgData map[string]interface{}
		content, err := os.ReadFile(commitMsgFile)
		if err == nil && yaml.Unmarshal(content, &msgData) == nil {
			viper.Set("commitMessage", msgData)
		} else {
			fmt.Fprintln(os.Stderr, "⚠️  No se pudo leer commitMessage.yml:", err)
		}
	}
}
