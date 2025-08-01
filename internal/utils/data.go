package utils

import "gitlab.com/0xleonz/gitz/internal/types"

// Defaults para message
var DefaultDescription = []string{
	"Esta implementación añade validación para tokens expirados.",
	"Se usa la librería jwt-go con tolerancia de 1 minuto.",
}

var DefaultChanges = []types.Change{
	{Type: "feat", Summary: "agrega endpoint para exportar CSV"},
	{Type: "fix", Summary: "corrige error en validación de email"},
	{Type: "refactor", Summary: "renombra variables para mayor claridad"},
}

var DefaultFooter = map[string]string{
	"Issue":          "#123",
	"Signed-off-by":  "leonz <0xleonz@gmail.com>",
	"Co-authored-by": "otro-dev <dev@example.com>",
}

const (
	DefaultSubject           = "feat: agregar validación de email"
	DefaultIssue             = "#123"
	PromptDescriptionCurrent = "📄 Description actual:"
	PromptDescriptionAdd     = "📄 Agrega nuevas líneas (ENTER vacío para terminar):"
)
