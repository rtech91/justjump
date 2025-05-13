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

var removeGlobalFlag bool

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a jump root (global or local)",
	Run: func(cmd *cobra.Command, args []string) {
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Failed to get current directory: %v\n", err)
			os.Exit(1)
		}

		if removeGlobalFlag {
			handleRemoveGlobal(currentDir)
		} else {
			handleRemoveLocal(currentDir)
		}
	},
}

func handleRemoveGlobal(currentDir string) {
	globalConfig, err := global.New()
	if err != nil {
		fmt.Printf("Failed to load global configuration: %v\n", err)
		os.Exit(1)
	}

	var targetName string
	for name, jumproot := range globalConfig.JumpRoots() {
		if filepath.HasPrefix(currentDir, jumproot.Root) {
			targetName = name
			break
		}
	}

	if targetName == "" {
		fmt.Println("No matching global jump root found for the current directory")
		os.Exit(1)
	}

	err = globalConfig.DeleteJumpRoot(targetName)
	if err != nil {
		fmt.Printf("Failed to remove global jump root: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Global jump root removed: %s\n", targetName)
}

func handleRemoveLocal(currentDir string) {
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

	// Remove the jump point
	err = localConfig.RemoveJumpPoint(currentDir)
	if err != nil {
		fmt.Printf("Failed to remove local jump root: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Local jump root removed: %s\n", currentDir)
}

func init() {
	removeCmd.Flags().BoolVarP(&removeGlobalFlag, "global", "G", false, "Remove as a global jump root")
	rootCmd.AddCommand(removeCmd)
}
