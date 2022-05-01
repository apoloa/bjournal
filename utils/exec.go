package utils

import (
	"os"
	"os/exec"
)

func RunEditor(editor string, filePath string) error {
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
