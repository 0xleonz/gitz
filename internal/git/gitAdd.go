package git

import (
	"bytes"
	"os/exec"
	"strings"
)

func ChangedFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-files", "--modified", "--others", "--exclude-standard")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}, nil
	}
	return lines, nil
}

func GitAddFile(file string) error {
	cmd := exec.Command("git", "add", file)
	return cmd.Run()
}

func GitAddAll() error {
	cmd := exec.Command("git", "add", ".")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
