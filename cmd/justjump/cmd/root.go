package cmd

import (
	"fmt"
	"os"

	"github.com/rtech91/justjump/pkg/config/global"
	"github.com/rtech91/justjump/pkg/config/local"
	"github.com/rtech91/justjump/pkg/util"
	promptui_global "github.com/rtech91/justjump/pkg/util/promptui/global"
	promtui_local "github.com/rtech91/justjump/pkg/util/promptui/local"
	"github.com/spf13/cobra"
)

var shellOutput string = ""
var globalJump bool = false

var rootCmd = &cobra.Command{
	Use:   "justjump",
	Short: "JustJump is a simple tool to help you jump between directories quickly.",
	Long: `JustJump is a simple tool to help you jump between directories quickly.
To use it simply run 'jj' in your terminal and select the directory you want to jump to.`,
	Run: func(cmd *cobra.Command, args []string) {
		if shellOutput != "" {

			if globalJump {
				performGlobalJump(shellOutput)
				return
			}

			performLocalJump(shellOutput)
			return
		}
	},
}

func performGlobalJump(tmpFilePath string) {
	globalConfig, err := global.New()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	jumpRoots := globalConfig.JumpRoots()
	if len(jumpRoots) == 0 {
		fmt.Println("No jump roots found")
		os.Exit(1)
	}

	jumpRootPaths := util.BuildJumpRootPaths(jumpRoots)

	prompt := promptui_global.PromptSelector(jumpRootPaths)

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	// open tmpFilePath and write the selected jump root with command
	err = util.EchoCommand(tmpFilePath, jumpRootPaths[i]["fullPath"])
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func performLocalJump(tmpFilePath string) {
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

		prompt := promtui_local.PromptSelector(jumpPointPaths)

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

	rootCmd.PersistentFlags().BoolVarP(&globalJump, "global", "G", false, "Perform a global jump accross registered projects")
}
