package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/rtech91/justjump/pkg/config/global"
	"github.com/rtech91/justjump/pkg/config/local"
	"github.com/rtech91/justjump/pkg/util"
	"github.com/spf13/cobra"
)

var verifyGlobal bool
var verifyClean bool

// jrEntry represents a global jump root entry that was found to be missing.
type jrEntry struct {
	name string
	root string
}

// jpEntry represents a local jump point entry that was found to be missing.
type jpEntry struct {
	relPath  string
	fullPath string
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify and optionally clean up jump roots and points",
	Long: `Verify checks whether registered jump roots (global) or jump points (local)
still exist on disk. Non-existent entries are reported.

In an interactive terminal, you will be prompted to remove each invalid entry.
Use --clean / -c to remove all invalid entries automatically without prompting.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if verifyGlobal {
			verifyGlobalFolders()
		} else {
			verifyLocalFolders()
		}
	},
}

func verifyGlobalFolders() {
	globalConfig, err := global.New()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	jumpRoots := globalConfig.JumpRoots()
	if len(jumpRoots) == 0 {
		fmt.Println("No global jump roots found")
		os.Exit(1)
	}

	var missingRoots []jrEntry

	// Sort keys for deterministic output order
	names := make([]string, 0, len(jumpRoots))
	for name := range jumpRoots {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		jumpRoot := jumpRoots[name]
		if _, err := os.Stat(jumpRoot.Root); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				fmt.Printf("Warning: could not access jump root '%s' (%s): %v\n", name, jumpRoot.Root, err)
			}
			missingRoots = append(missingRoots, jrEntry{name: name, root: jumpRoot.Root})
		}
	}

	if len(missingRoots) == 0 {
		fmt.Println("All global jump roots are correct")
		os.Exit(0)
	}

	terminal := util.IsTerminal()
	var cleanedCount int
	hasUnresolvedIssues := false

	for _, jr := range missingRoots {
		fmt.Printf("Jump root '%s' (%s) does not exist\n", jr.name, jr.root)

		shouldRemove := false
		if verifyClean {
			shouldRemove = true
		} else if terminal {
			promptLabel := fmt.Sprintf("Remove global jump root '%s' from config", jr.name)
			shouldRemove = util.AskConfirmation(promptLabel)
		}

		if shouldRemove {
			err := globalConfig.DeleteJumpRoot(jr.name)
			if err != nil {
				fmt.Printf("Failed to remove global jump root '%s': %v\n", jr.name, err)
				hasUnresolvedIssues = true
			} else {
				fmt.Printf("Removed global jump root '%s'\n", jr.name)
				cleanedCount++
			}
		} else {
			hasUnresolvedIssues = true
		}
	}

	if cleanedCount > 0 {
		fmt.Printf("Successfully cleaned up %d global jump root(s)\n", cleanedCount)
	}

	if hasUnresolvedIssues {
		os.Exit(1)
	}
}

func verifyLocalFolders() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	globalConfig, err := global.New()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	jumpRoots := globalConfig.JumpRoots()
	jumpRootExists, jumpRoot := util.DetermineJumpRoot(currentDir, jumpRoots)

	if !jumpRootExists {
		fmt.Println("No local jump root found")
		os.Exit(1)
	}

	localConfig, err := local.New(jumpRoot)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	jumpPoints := localConfig.JumpPoints()
	if len(jumpPoints) == 0 {
		fmt.Println("No jump points found")
		os.Exit(1)
	}

	var missingPoints []jpEntry

	for _, relPath := range jumpPoints {
		fullPath := filepath.Join(jumpRoot, relPath)
		if _, err := os.Stat(fullPath); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				fmt.Printf("Warning: could not access jump point '%s': %v\n", fullPath, err)
			}
			missingPoints = append(missingPoints, jpEntry{relPath: relPath, fullPath: fullPath})
		}
	}

	if len(missingPoints) == 0 {
		fmt.Println("All local jump points are correct")
		os.Exit(0)
	}

	terminal := util.IsTerminal()
	var cleanedCount int
	hasUnresolvedIssues := false

	for _, jp := range missingPoints {
		fmt.Printf("Jump point '%s' does not exist\n", jp.fullPath)

		shouldRemove := false
		if verifyClean {
			shouldRemove = true
		} else if terminal {
			promptLabel := fmt.Sprintf("Remove local jump point '%s' from config", jp.relPath)
			shouldRemove = util.AskConfirmation(promptLabel)
		}

		if shouldRemove {
			err := localConfig.RemoveJumpPoint(jp.fullPath)
			if err != nil {
				fmt.Printf("Failed to remove local jump point '%s': %v\n", jp.relPath, err)
				hasUnresolvedIssues = true
			} else {
				fmt.Printf("Removed local jump point '%s'\n", jp.relPath)
				cleanedCount++
			}
		} else {
			hasUnresolvedIssues = true
		}
	}

	if cleanedCount > 0 {
		fmt.Printf("Successfully cleaned up %d local jump point(s)\n", cleanedCount)
	}

	if hasUnresolvedIssues {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().BoolVarP(&verifyGlobal, "global", "G", false, "Verify global folders")
	verifyCmd.Flags().BoolVarP(&verifyClean, "clean", "c", false, "Clean up non-existent folders")
}
