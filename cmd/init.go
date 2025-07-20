package cmd

import (
  "fmt"
  "io"
  "os"
  "os/exec"
  "path/filepath"
  "strings"

  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/bubbles/list"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

// expandPath expande "~" a $HOME y variables de entorno.
func expandPath(p string) string {
  if strings.HasPrefix(p, "~"+string(os.PathSeparator)) {
    home, err := os.UserHomeDir()
    if err == nil {
      return filepath.Join(home, p[2:])
    }
  }
  return os.ExpandEnv(p)
}

const (
  listWidth  = 30
  listHeight = 6
)

// item implementa list.Item para Bubble Tea
type item struct{ title string }

func (i item) Title() string       { return i.title }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.title }

// model para Bubble Tea
type model struct{ list list.Model }

func (m model) Init() tea.Cmd                            { return nil }
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd
  m.list, cmd = m.list.Update(msg)
  if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
    return m, tea.Quit
  }
  return m, cmd
}
func (m model) View() string { return m.list.View() }

var force bool

var initCmd = &cobra.Command{
  Use:   "init",
  Short: "Inicializa un nuevo proyecto con plantilla básica",
  RunE:  runInit,
}

func init() {
  initCmd.Flags().BoolVarP(&force, "force", "f", false, "Sobrescribe archivos sin preguntar")
  rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
  // 0. Si ya hay .git, saltamos init
  if _, err := os.Stat(".git"); err == nil {
    fmt.Println("\n⚠️  Ya existe un repositorio Git en este directorio, se omite git init.")
  } else {
    g := exec.Command("git", "init")
    g.Stdout = os.Stdout
    g.Stderr = os.Stderr
    if err := g.Run(); err != nil {
      return fmt.Errorf("falló git init: %w", err)
    }
  }

  // 1. Leer templates_dir desde config
  tplDir := expandPath(viper.GetString("templates_dir"))

  // 2. Preparar lista de licencias
  licenses := viper.GetStringMapString("licenses")
  var items []list.Item
  for name := range licenses {
    items = append(items, item{title: name})
  }
  l := list.New(items, list.NewDefaultDelegate(), listWidth, listHeight)
  l.Title = "Selecciona la licencia (↑/↓ y Enter)"

  // 3. Ejecutar Bubble Tea
  p := tea.NewProgram(model{list: l})
  finalModel, err := p.Run()
  if err != nil {
    return fmt.Errorf("error en TUI de Bubble Tea: %w", err)
  }

  // 4. Obtener la selección
  choice := finalModel.(model).list.SelectedItem().(item).title

  // 5. Mapeo de archivos a copiar
  mappings := map[string]string{
    "LICENSE":    licenses[choice],
    ".gitignore": viper.GetString("defaults.gitignore"),
    "info.yml":   viper.GetString("defaults.info"),
  }

  // 6. Copiar cada plantilla con precaución
  for dst, tpl := range mappings {
    src := filepath.Join(tplDir, tpl)
    if err := copyFile(src, dst); err != nil {
      return fmt.Errorf("error copiando %s: %w", dst, err)
    }
  }

  // 7. Mensaje final
  fmt.Println("✅ Proyecto inicializado:")
  fmt.Printf("   • Licencia: %s\n", choice)
  fmt.Println("   • .gitignore")
  fmt.Println("   • info.yml")
  fmt.Println("   • LICENSE")
  return nil
}

// copyFile copia src a dst con confirmación si ya existe (a menos que --force)
func copyFile(src, dst string) error {
  if _, err := os.Stat(dst); err == nil && !force {
    fmt.Printf("⚠️  El archivo %s ya existe. ¿Deseas sobrescribirlo? [s/N]: ", dst)
    var resp string
    fmt.Scanln(&resp)
    resp = strings.ToLower(strings.TrimSpace(resp))
    if resp != "s" && resp != "y" {
      fmt.Printf("\n⏭️  Omitiendo %s\n", dst)
      return nil
    }
  }

  in, err := os.Open(src)
  if err != nil {
    return err
  }
  defer in.Close()

  out, err := os.Create(dst)
  if err != nil {
    return err
  }
  defer out.Close()

  if _, err := io.Copy(out, in); err != nil {
    return err
  }
  return out.Sync()
}
