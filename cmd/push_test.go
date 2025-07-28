package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushCmd(t *testing.T) {
	// buffer para capturar la salida
	var out bytes.Buffer
	pushCmd.SetOut(&out)

	// Simulamos argumentos
	pushCmd.SetArgs([]string{"--dry-run"})

	// Ejecutamos y checanmos salida
	err := pushCmd.Execute()
	assert.NoError(t, err)

	output := out.String()
	assert.Contains(t, output, "dry-run")
}
