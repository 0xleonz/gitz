package utils

import (
	"fmt"
	"strings"

	"gitlab.com/0xleonz/gitz/internal/types"
)

func FormatCommitMessage(msg types.CommitMessage) string {
	var sb strings.Builder

	sb.WriteString(msg.Subject)

	if len(msg.Description) > 0 {
		sb.WriteString("\n\n" + strings.Join(msg.Description, "\n"))
	}

	if len(msg.Changes) > 0 {
		sb.WriteString("\n\nChanges:")
		for _, c := range msg.Changes {
			sb.WriteString(fmt.Sprintf("\n- %s: %s", c.Type, c.Summary))
		}
	}

	if len(msg.Footer) > 0 {
		sb.WriteString("\n")
		for k, v := range msg.Footer {
			sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
	}

	if msg.Issue != "" {
		sb.WriteString("\nIssue: " + msg.Issue)
	}

	return sb.String()
}
