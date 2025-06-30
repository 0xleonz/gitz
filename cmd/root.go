package cmd

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
    Use:   "gitz",
    Short: "gitz: un wrapper ligero de git",
    Long:  "gitz es tu CLI para interactuar con git más rápido.",
    // Antes de cualquier subcomando, cargamos config:
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        initConfig()
    },
}

// Execute ejecuta el rootCmd.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize() // initConfig ya en PersistentPreRun
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/gitz/config.yaml)")
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        cobra.CheckErr(err)
        cfgDir := filepath.Join(home, ".config", "gitz")
        viper.AddConfigPath(cfgDir)
        viper.SetConfigName("config")
        viper.SetConfigType("yaml")
        viper.SetDefault("templates_dir", filepath.Join(home, ".config", "gitz", "templates"))
        viper.SetDefault("defaults.gitignore", "gitignore")
        viper.SetDefault("defaults.datos", "datos")
        // licencias
        viper.SetDefault("licenses", map[string]string{
            "GNU GPLv3": "licenses/GPL-3.0.txt",
            "MIT":       "licenses/MIT.txt",
        })
        viper.SetDefault("default_license", "GNU GPLv3")

        if err := viper.ReadInConfig(); err != nil {
            fmt.Fprintln(os.Stderr, "⚠️  no se encontró config, creando en", cfgDir)
            os.MkdirAll(cfgDir, 0o755)
            viper.WriteConfigAs(filepath.Join(cfgDir, "config.yaml"))
            return
        }
    }

    if f := viper.ConfigFileUsed(); f != "" {
        fmt.Fprintln(os.Stderr, "✅ usando config file:", f)
    }
    viper.AutomaticEnv()
}

