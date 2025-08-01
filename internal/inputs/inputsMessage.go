package inputs

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gitlab.com/0xleonz/gitz/internal/types"
	"gitlab.com/0xleonz/gitz/internal/utils"
)

// =========================
// === CAMPOS PRINCIPALES ===
// =========================

func Subject(current string) string {
	val := promptField("📝 Subject:", utils.Pink, current, utils.DefaultSubject)
	if val == "" {
		return current
	}
	return val
}

func Issue(issue string) string {
	return promptField("🔗 Issue:", utils.Purple, issue, utils.DefaultIssue)
}

func Description(desc []string) []string {
	r := bufio.NewReader(os.Stdin)

	if len(desc) > 0 {
		fmt.Println(utils.Colorize(utils.PromptDescriptionCurrent, utils.Yellow))
		for _, d := range desc {
			fmt.Println("  " + utils.Colorize(d, utils.Yellow))
		}
	} else {
		fmt.Println(utils.Colorize("📄 No hay descripción actual. Usando ejemplo por defecto:", utils.Cyan))
		for _, d := range utils.DefaultDescription {
			fmt.Println("  " + utils.Colorize("[ejemplo] "+d, utils.Blue))
		}
	}

	fmt.Println(utils.Colorize(utils.PromptDescriptionAdd, utils.Green))
	for {
		fmt.Print("  > ")
		line, _ := r.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		desc = append(desc, line)
	}

	return desc
}

// =======================
// === CAMBIOS / CHANGES ===
// =======================

func Cambios(msg *types.CommitMessage) {
	if len(msg.Changes) > 0 {
		fmt.Println(utils.Colorize("🔧 Cambios actuales:", utils.Yellow))
		for _, c := range msg.Changes {
			fmt.Printf("  - %s: %s\n", utils.Colorize(c.Type, utils.Yellow), c.Summary)
		}
	} else {
		fmt.Println(utils.Colorize("🔧 No hay cambios actuales. Ejemplos:", utils.Cyan))
		for _, c := range utils.DefaultChanges {
			fmt.Printf("  [ejemplo] - %s: %s\n", utils.Colorize(c.Type, utils.Blue), c.Summary)
		}
	}

	fmt.Println(utils.Colorize("🔧 Agrega cambios nuevos (tipo:resumen, ENTER vacío para terminar):", utils.Cyan))
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("  - ")
		line, _ := r.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			fmt.Println(utils.Colorize("  ⚠️  Formato inválido. Usa tipo:resumen", utils.Red))
			continue
		}
		msg.Changes = append(msg.Changes, types.Change{
			Type:    strings.TrimSpace(parts[0]),
			Summary: strings.TrimSpace(parts[1]),
		})
	}
}

func AddChange(msg *types.CommitMessage, input string) error {
	parts := strings.SplitN(input, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("formato inválido para --add-change. Usa tipo:resumen")
	}
	msg.Changes = append(msg.Changes, types.Change{
		Type:    strings.TrimSpace(parts[0]),
		Summary: strings.TrimSpace(parts[1]),
	})
	return nil
}

func Fix(msg *types.CommitMessage) {
	summary := prompt("Resumen del fix (ej: arregla validación de token):")
	msg.Changes = append(msg.Changes, types.Change{Type: "fix", Summary: summary})
}

func Refactor(msg *types.CommitMessage) {
	summary := prompt("Resumen del refactor (ej: simplifica lógica de parser):")
	msg.Changes = append(msg.Changes, types.Change{Type: "refactor", Summary: summary})
}

// ====================
// === FOOTER SECTION ===
// ====================

func Footer(current map[string]string) map[string]string {
	newFooter := current
	r := bufio.NewReader(os.Stdin)

	if len(current) > 0 {
		fmt.Println(utils.Colorize("🔻 Footer actual:", utils.Yellow))
		for k, v := range current {
			fmt.Printf("  [actual] %s: %s\n", k, v)
		}
	} else {
		fmt.Println(utils.Colorize("🔻 Footer vacío. Ejemplos:", utils.Cyan))
		for k, v := range utils.DefaultFooter {
			fmt.Printf("  [ejemplo] %s: %s\n", utils.Colorize(k, utils.Blue), v)
		}
	}

	fmt.Println(utils.Colorize("🔻 Agrega footers (clave:valor, ENTER vacío para terminar):", utils.Green))
	for {
		fmt.Print("  - ")
		line, _ := r.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			fmt.Println(utils.Colorize("  ⚠️  Formato inválido. Usa clave:valor", utils.Red))
			continue
		}
		newFooter[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return newFooter
}

// ==================
// === HELPERS ======
// ==================

func prompt(label string) string {
	fmt.Print(utils.Colorize(label+" ", utils.Pink))
	r := bufio.NewReader(os.Stdin)
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}

func promptField(label, color, current, example string) string {
	if current != "" {
		fmt.Printf("%s %s \n", utils.Colorize(label, color), current)
	} else if example != "" {
		fmt.Printf("%s %s\n", utils.Colorize(label, color), utils.Colorize("[ej: "+example+"]", utils.Blue))
	}
	fmt.Print("  > ")
	r := bufio.NewReader(os.Stdin)
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}
