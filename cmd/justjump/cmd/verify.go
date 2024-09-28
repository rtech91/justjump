package cmd

import (
	"fmt"
	"os"

	"github.com/rtech91/justjump/pkg/config/global"
	"github.com/rtech91/justjump/pkg/config/local"
	"github.com/rtech91/justjump/pkg/util"
	"github.com/spf13/cobra"
)

var verifyGlobal bool = false

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify if local or global folders exist",
	Run: func(cmd *cobra.Command, args []string) {
		if verifyGlobal {
			verifyGlobalFolders()
		} else {
			verifyLocalFolders()
		}
	},
}

func verifyGlobalFolders() {
	var totalStatus bool = true
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

	for _, jumpRoot := range jumpRoots {
		if _, err := os.Stat(jumpRoot.Root); os.IsNotExist(err) {
			fmt.Printf("Jump root %s does not exist\n", jumpRoot.Root)
			totalStatus = false
		}
	}

	if totalStatus {
		fmt.Println("All global jump roots are correct")
	}
}

func verifyLocalFolders() {
	var totalStatus bool = true
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

	if jumpRootExists {
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

		jumpPointPaths := util.BuildJumpPointPaths(jumpRoot, jumpPoints)

		for _, jumpPoint := range jumpPointPaths {
			if _, err := os.Stat(jumpPoint["fullPath"]); os.IsNotExist(err) {
				fmt.Printf("Jump point %s does not exist\n", jumpPoint["fullPath"])
				totalStatus = false
			}
		}

		if totalStatus {
			fmt.Println("All local jump points are correct")
		}
	} else {
		fmt.Println("No local jump root found")
	}
}

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().BoolVarP(&verifyGlobal, "global", "G", false, "Verify global folders")
}
