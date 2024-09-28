package util

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rtech91/justjump/pkg/config/global"
)

func DetermineJumpRoot(currentDir string, jumpRoots global.JumpRoots) (bool, string) {
	var exist bool = false
	var jumpRoot string = ""
	for _, jr := range jumpRoots {
		if strings.HasPrefix(currentDir, jr.Root) {
			exist = true
			jumpRoot = jr.Root
			break
		}
	}
	return exist, jumpRoot
}

func BuildJumpRootPaths(jumpRoots global.JumpRoots) []map[string]string {
	jumpRootPaths := make([]map[string]string, 0)

	for name, jr := range jumpRoots {

		if _, err := os.Stat(jr.Root); os.IsNotExist(err) {
			fmt.Printf("Can't add jump root to the list %s as it does not exist\n", jr.Root)
			continue
		}

		dict := map[string]string{
			"jumpRoot": name,
			"fullPath": jr.Root,
		}

		jumpRootPaths = append(jumpRootPaths, dict)
	}

	return jumpRootPaths
}

func BuildJumpPointPaths(jumpRoot string, jumpPoints []string) []map[string]string {
	jumpPointPaths := make([]map[string]string, 0)
	jumpPointPaths = append(jumpPointPaths, map[string]string{
		"jumpPoint": jumpRoot,
		"fullPath":  jumpRoot,
	})
	for _, jumpPoint := range jumpPoints {
		var fullPath string = jumpRoot + "/" + jumpPoint

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			fmt.Printf("Can't add jump point to the list %s as it does not exist\n", fullPath)
			continue
		}

		dict := map[string]string{
			"jumpPoint": jumpPoint,
			"fullPath":  jumpRoot + "/" + jumpPoint,
		}

		jumpPointPaths = append(jumpPointPaths, dict)
	}

	return jumpPointPaths
}

func EchoCommand(tmpFilePath string, chosenFullPath string) error {
	file, err := os.OpenFile(tmpFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("failed to open the temporary file")
	}

	// write the selected jump point with command
	_, err = file.WriteString("cd " + chosenFullPath)
	if err != nil {
		return errors.New("failed to write to the temporary file")
	}

	return nil
}
