package cmd

import (
	"fmt"
	"os"

	"github.com/rtech91/justjump/pkg/config/global"
	"github.com/rtech91/justjump/pkg/config/local"
	"github.com/rtech91/justjump/pkg/util"
	"github.com/spf13/cobra"
)

var shellOutput string = ""

var rootCmd = &cobra.Command{
	Use:   "justjump",
	Short: "JustJump is a simple tool to help you jump between directories quickly.",
	Long: `JustJump is a simple tool to help you jump between directories quickly.
To use it simply run 'jj' in your terminal and select the directory you want to jump to.`,
	Run: func(cmd *cobra.Command, args []string) {
		if shellOutput != "" {
			performJump(shellOutput)
			return
		}
	},
}

func performJump(tmpFilePath string) {
	globalConfig, err := global.New()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	// check if the current directory contains a jump root
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

		prompt := util.PromptSelector(jumpPointPaths)

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			os.Exit(1)
		}

		// open tmpFilePath and write the selected jump point with command
		err = util.EchoCommand(tmpFilePath, jumpPointPaths[i]["fullPath"])
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&shellOutput, "shelloutput", "s", "", "Output the shell command to a temporary file")
	rootCmd.PersistentFlags().MarkHidden("shelloutput")
}
