package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	file, err := os.OpenFile(tmpFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open temporary file: %w", err)
	}
	defer file.Close()

	// write the selected jump point with command
	_, err = file.WriteString("cd " + chosenFullPath)
	if err != nil {
		return fmt.Errorf("failed to write to temporary file: %w", err)
	}

	return nil
}

// FuzzyMatch returns true if the characters in the search string
// appear in the target string in the same order.
// Both search and target must be pre-lowercased for best performance.
func FuzzyMatch(search, target string) bool {
	if search == "" {
		return true
	}

	searchIdx := 0
	for targetIdx := 0; targetIdx < len(target); targetIdx++ {
		if target[targetIdx] == search[searchIdx] {
			searchIdx++
		}
		if searchIdx == len(search) {
			return true
		}
	}

	return false
}


// GetGitWorktrees runs 'git worktree list' and returns the paths of all associated worktrees.
func GetGitWorktrees() ([]map[string]string, error) {
	worktrees := make([]map[string]string, 0)

	// Run git worktree list --porcelain
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("not a git repository or git not found")
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "worktree ") {
			path := strings.TrimPrefix(line, "worktree ")
			name := filepath.Base(path)
			worktrees = append(worktrees, map[string]string{
				"jumpPoint": name,
				"fullPath":  path,
			})
		}
	}

	return worktrees, nil
}
