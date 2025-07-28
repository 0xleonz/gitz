package inputs

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gitlab.com/0xleonz/gitz/internal/types"
	"gitlab.com/0xleonz/gitz/internal/utils"
)

func Subject(current string) string {
	reader := bufio.NewReader(os.Stdin)

	if current == "" {
		example := "Add login redirect after auth"
		fmt.Printf("📝 Subject (e.g. %s): ", utils.Colorize(example, utils.Yellow))
	} else {
		fmt.Printf("📝 Subject [%s]: ", utils.Colorize(current, utils.Red))
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "" {
		return input
	}
	return current
}

func Description(current []string) []string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("📖 Description (ENTER vacío para terminar):")

	if len(current) > 0 {
		fmt.Println(utils.Colorize("Descripción actual:", utils.Yellow))
		for _, d := range current {
			fmt.Println("•", d)
		}
	} else {
		example := []string{
			"Refactor login controller to improve error handling.",
			"Add coverage for edge cases in session middleware.",
			"Hacked the mainframe like pirate software",
		}
		fmt.Println(utils.Colorize("Ejemplo:", utils.Yellow))
		for _, e := range example {
			fmt.Println("•", e)
		}
	}

	var desc []string
	for {
		fmt.Print("- ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		desc = append(desc, line)
	}

	return desc
}

func Changes(current []types.Change) []types.Change {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🔧 Changes (formato tipo:resumen, ENTER vacío para terminar):")

	if len(current) > 0 {
		fmt.Println(utils.Colorize("Cambios actuales:", utils.Yellow))
		for _, ch := range current {
			fmt.Printf("• %s:%s\n", ch.Type, ch.Summary)
		}
	} else {
		fmt.Println(utils.Colorize("Ejemplos:", utils.Yellow))
		examples := []string{
			"fix:handle token expiry",
			"refactor:split handlers into separate files",
			"docs:add usage example for login endpoint",
		}
		for _, e := range examples {
			fmt.Println("•", e)
		}
	}

	var changes []types.Change
	for {
		fmt.Print("- ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			fmt.Println(utils.Colorize("❌ Formato inválido. Usa tipo:resumen", utils.Red))
			continue
		}
		changes = append(changes, types.Change{
			Type:    strings.TrimSpace(parts[0]),
			Summary: strings.TrimSpace(parts[1]),
		})
	}

	return changes
}

func Footer(current map[string]string) map[string]string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🧾 Footer (clave:valor, ENTER vacío para terminar):")
	if len(current) > 0 {
		fmt.Println(utils.Colorize("Footers actuales:", utils.Yellow))
		for k, v := range current {
			fmt.Printf("• %s: %s\n", k, v)
		}
	} else {
		fmt.Println(utils.Colorize("Ejemplos:", utils.Yellow))
		examples := []string{
			"Signed-off-by: Edgar M <edgar@example.com>",
			"Migration step: clear user sessions",
			"Reviewed-by: QA team",
		}
		for _, e := range examples {
			fmt.Println("•", e)
		}
	}

	result := make(map[string]string)
	for {
		fmt.Print("- ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			result[key] = val
		} else {
			fmt.Println(utils.Colorize("❌ Formato inválido. Usa clave:valor", utils.Red))
		}
	}

	return result
}

func Issue(current string) string {
	reader := bufio.NewReader(os.Stdin)

	exampleIssue := "JIRA-1912"
	label := current
	if label == "" {
		label = exampleIssue
	}

	fmt.Printf("📌 Issue ID [%s]: ", utils.Colorize(label, utils.Cyan))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "" {
		return input
	}
	if current == "" {
		return ""
	}
	return current
}
