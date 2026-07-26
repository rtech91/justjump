package util

import (
	"os"

	"github.com/manifoldco/promptui"
)

// IsTerminal checks whether stdin is connected to an interactive terminal.
func IsTerminal() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// AskConfirmation prompts the user with a yes/no confirmation dialog.
// Returns true if the user confirmed, false otherwise.
func AskConfirmation(label string) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}
	_, err := prompt.Run()
	return err == nil
}
