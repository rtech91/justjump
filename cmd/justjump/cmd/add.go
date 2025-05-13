package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rtech91/justjump/pkg/config/global"
	"github.com/rtech91/justjump/pkg/config/local"
	"github.com/rtech91/justjump/pkg/util"
	"github.com/spf13/cobra"
)

var addGlobalFlag bool

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a jump root (global or local)",
	Run: func(cmd *cobra.Command, args []string) {
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Failed to get current directory: %v\n", err)
			os.Exit(1)
		}

		if addGlobalFlag {
			handleAddGlobal(currentDir)
		} else {
			handleAddLocal(currentDir)
		}
	},
}

func handleAddGlobal(currentDir string) {
	globalConfig, err := global.New()
	if err != nil {
		fmt.Printf("Failed to load global configuration: %v\n", err)
		os.Exit(1)
	}

	name := filepath.Base(currentDir)
	err = globalConfig.RegisterJumpRoot(name, global.Jumproot{Name: name, Root: currentDir})
	if err != nil {
		fmt.Printf("Failed to add global jump root: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Global jump root added: %s -> %s\n", name, currentDir)
}

func handleAddLocal(currentDir string) {
	globalConfig, err := global.New()
	if err != nil {
		fmt.Printf("Failed to load global configuration: %v\n", err)
		os.Exit(1)
	}

	// Determine the jump root for the current directory
	jumpRoots := globalConfig.JumpRoots()
	jumpRootExists, jumpRoot := util.DetermineJumpRoot(currentDir, jumpRoots)

	if !jumpRootExists {
		fmt.Println("No matching local jump root found for the current directory")
		os.Exit(1)
	}

	// Load the local configuration using the determined jump root
	localConfig, err := local.New(jumpRoot)
	if err != nil {
		fmt.Printf("Failed to load local configuration: %v\n", err)
		os.Exit(1)
	}

	// Calculate the relative path
	relPath, err := filepath.Rel(jumpRoot, currentDir)
	if err != nil {
		fmt.Printf("Failed to calculate relative path: %v\n", err)
		os.Exit(1)
	}

	// Avoid adding a local jump point if it is the same folder as the global jump root
	if relPath == "." {
		fmt.Println("Cannot add a local jump point for the same folder as the global jump root")
		os.Exit(1)
	}

	// Add the relative path as a jump point
	err = localConfig.AddJumpPoint(relPath)
	if err != nil {
		fmt.Printf("Failed to add local jump root: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Local jump root added: %s\n", relPath)
}

func init() {
	addCmd.Flags().BoolVarP(&addGlobalFlag, "global", "G", false, "Add as a global jump root")
	rootCmd.AddCommand(addCmd)
}
