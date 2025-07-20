package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"gitlab.com/0xleonz/gitz/internal/utils"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"

  "gitlab.com/0xleonz/gitz/internal/git"
)

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Crea y maneja mensajes de commit enriquecidos",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(utils.Colorize("🧠 Iniciando editor interactivo...", utils.Cyan))
		return runInteractiveEditor()
	},
}

type Change struct {
	Type    string `yaml:"type"`
	Summary string `yaml:"summary"`
}

type CommitMessage struct {
	Changes     []Change          `yaml:"changes"`
	Issue       string            `yaml:"issue"`
	Subject     string            `yaml:"subject"`
	Description []string          `yaml:"description"`
	Footer      map[string]string `yaml:"footer"`
}

type messageModel struct {
	step int

	subjectInput textinput.Model
	descInput    textinput.Model
	typeInput    textinput.Model
	footerInput  textinput.Model
	issueInput   textinput.Model

	commit CommitMessage
	done   bool
	error  string
}

func initialModel() messageModel {
	subj := textinput.New()
	subj.Placeholder = "Subject (<= 50 chars)"
	subj.Focus()
	subj.CharLimit = 50

	desc := textinput.New()
	desc.Placeholder = "Description (ENTER para agregar, vacío para continuar)"
	desc.CharLimit = 200

	typeI := textinput.New()
	typeI.Placeholder = "Tipo:Resumen (ej. feat:Añade login)"

	footer := textinput.New()
	footer.Placeholder = "Footer (clave:valor)"

	issue := textinput.New()
	issue.Placeholder = "Issue ID (opcional)"

	return messageModel{
		step:         0,
		subjectInput: subj,
		descInput:    desc,
		typeInput:    typeI,
		footerInput:  footer,
		issueInput:   issue,
		commit: CommitMessage{
			Changes:     []Change{},
			Footer:      make(map[string]string),
			Description: []string{},
		},
	}
}

func (m messageModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m messageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.step {
		case 0:
			if msg.Type == tea.KeyEnter {
				m.commit.Subject = m.subjectInput.Value()
				m.step++
				return m, textinput.Blink
			}
			var cmd tea.Cmd
			m.subjectInput, cmd = m.subjectInput.Update(msg)
			return m, cmd

		case 1:
			if msg.Type == tea.KeyEnter {
				val := strings.TrimSpace(m.descInput.Value())
				if val == "" {
					m.step++
				} else {
					m.commit.Description = append(m.commit.Description, val)
					m.descInput.SetValue("")
				}
				return m, textinput.Blink
			}
			m.descInput, _ = m.descInput.Update(msg)
			return m, nil

		case 2:
			if msg.Type == tea.KeyEnter {
				val := strings.TrimSpace(m.typeInput.Value())
				if val == "" {
					m.step++
					return m, textinput.Blink
				}
				parts := strings.SplitN(val, ":", 2)
				if len(parts) == 2 {
					m.commit.Changes = append(m.commit.Changes, Change{Type: parts[0], Summary: parts[1]})
					m.typeInput.SetValue("")
				}
				return m, nil
			}
			m.typeInput, _ = m.typeInput.Update(msg)
			return m, nil

		case 3:
			if msg.Type == tea.KeyEnter {
				val := strings.TrimSpace(m.footerInput.Value())
				if val == "" {
					m.step++
					return m, textinput.Blink
				}
				parts := strings.SplitN(val, ":", 2)
				if len(parts) == 2 {
					m.commit.Footer[parts[0]] = parts[1]
					m.footerInput.SetValue("")
				}
				return m, nil
			}
			m.footerInput, _ = m.footerInput.Update(msg)
			return m, nil

		case 4:
			if msg.Type == tea.KeyEnter {
				m.commit.Issue = m.issueInput.Value()
				m.done = true
				return m, tea.Quit
			}
			m.issueInput, _ = m.issueInput.Update(msg)
			return m, nil
		}
	default:
		return m, nil
	}
	return m, nil
}

func (m messageModel) View() string {
	if m.done {
		repoRoot, _ := git.FindRepoRoot()
		path := filepath.Join(repoRoot, "commitMessage.yml")
		b, _ := yaml.Marshal(m.commit)
		_ = os.WriteFile(path, b, 0644)
		return utils.Colorize("\n✅ Mensaje de commit guardado en ./commitMessage.yml\n", utils.Green)
	}
	switch m.step {
	case 0:
		return "\nSubject:\n📝 " + m.subjectInput.View()
	case 1:
		return "\nDescription:\n📖 " + m.descInput.View() + "\nPresiona ENTER en vacío para continuar."
	case 2:
		return "\nChanges:\n🔧 " + m.typeInput.View() + "\nTipo:Resumen o ENTER vacío para continuar."
	case 3:
		return "\nFooter:\n🧾 " + m.footerInput.View() + "\nFooter: clave:valor o ENTER vacío para continuar."
	case 4:
		return "\nIssue ID:\n📌 " + m.issueInput.View() + "\nID de issue (opcional). ENTER para terminar."
	default:
		return ""
	}
}

func runInteractiveEditor() error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func init() {
	rootCmd.AddCommand(messageCmd)
}
