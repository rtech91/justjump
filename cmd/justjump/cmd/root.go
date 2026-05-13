package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/rtech91/justjump/pkg/config/global"
	"github.com/rtech91/justjump/pkg/config/local"
	"github.com/rtech91/justjump/pkg/util"
	promptui_global "github.com/rtech91/justjump/pkg/util/promptui/global"
	promtui_local "github.com/rtech91/justjump/pkg/util/promptui/local"
	"github.com/spf13/cobra"
)

var shellOutput string = ""
var globalJump bool = false
var workspacesJump bool = false

var rootCmd = &cobra.Command{
	Use:   "justjump",
	Short: "JustJump is a simple tool to help you jump between directories quickly.",
	Long: `JustJump is a simple tool to help you jump between directories quickly.
To use it simply run 'jj' in your terminal and select the directory you want to jump to.

The --global or -G flag can be used not only to perform jumps across projects, but also as a modifier for other commands like add, verify, or remove.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if shellOutput != "" {

			if len(args) == 1 && args[0] == "-" {
				lastPath, err := util.ReadLastJump()
				if err != nil {
					if os.IsNotExist(err) {
						fmt.Println("No previous jump found")
					} else {
						fmt.Printf("Error reading jump history: %v\n", err)
					}
					os.Exit(1)
				}

				err = util.EchoCommand(shellOutput, lastPath)
				if err != nil {
					fmt.Printf("%v\n", err)
					os.Exit(1)
				}
				return
			}

			if workspacesJump {
				performWorkspaceJump(shellOutput, args)
				return
			}

			if globalJump {
				performGlobalJump(shellOutput, args)
				return
			}

			performLocalJump(shellOutput, args)
			return
		}
	},
}

func performGlobalJump(tmpFilePath string, args []string) {
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

	allPaths := util.BuildJumpRootPaths(jumpRoots)
	targetPaths := allPaths

	if len(args) > 0 {
		searchTerm := strings.ToLower(args[0])
		var filtered []map[string]string
		for _, p := range allPaths {
			if util.FuzzyMatch(searchTerm, strings.ToLower(p["jumpRoot"])) {
				filtered = append(filtered, p)
			}
		}

		if len(filtered) == 1 {
			err = util.EchoCommand(tmpFilePath, filtered[0]["fullPath"])
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
			return
		}

		if len(filtered) > 1 {
			targetPaths = filtered
		}
	}

	prompt := promptui_global.PromptSelector(targetPaths, "Select a jump root")

	i, _, err := prompt.Run()
	if err != nil {
		if (err == promptui.ErrInterrupt || err == promptui.ErrEOF) && len(targetPaths) < len(allPaths) {
			prompt = promptui_global.PromptSelector(allPaths, "Select a jump root")
			i, _, err = prompt.Run()
			if err != nil {
				os.Exit(0)
			}
			err = util.EchoCommand(tmpFilePath, allPaths[i]["fullPath"])
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
			return
		}
		os.Exit(0)
	}

	// open tmpFilePath and write the selected jump root with command
	err = util.EchoCommand(tmpFilePath, targetPaths[i]["fullPath"])
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func performLocalJump(tmpFilePath string, args []string) {
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

		allPaths := util.BuildJumpPointPaths(jumpRoot, jumpPoints)
		targetPaths := allPaths

		if len(args) > 0 {
			searchTerm := strings.ToLower(args[0])
			var filtered []map[string]string
			for _, p := range allPaths {
				if util.FuzzyMatch(searchTerm, strings.ToLower(p["jumpPoint"])) {
					filtered = append(filtered, p)
				}
			}

			if len(filtered) == 1 {
				err = util.EchoCommand(tmpFilePath, filtered[0]["fullPath"])
				if err != nil {
					fmt.Printf("%v\n", err)
					os.Exit(1)
				}
				return
			}

			if len(filtered) > 1 {
				targetPaths = filtered
			}
		}

		prompt := promtui_local.PromptSelector(targetPaths, "Select a jump point")

		i, _, err := prompt.Run()
		if err != nil {
			if (err == promptui.ErrInterrupt || err == promptui.ErrEOF) && len(targetPaths) < len(allPaths) {
				prompt = promtui_local.PromptSelector(allPaths, "Select a jump point")
				i, _, err = prompt.Run()
				if err != nil {
					os.Exit(0)
				}
				err = util.EchoCommand(tmpFilePath, allPaths[i]["fullPath"])
				if err != nil {
					fmt.Printf("%v\n", err)
					os.Exit(1)
				}
				return
			}
			os.Exit(0)
		}

		// open tmpFilePath and write the selected jump point with command
		err = util.EchoCommand(tmpFilePath, targetPaths[i]["fullPath"])
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Can't determine jump root for current directory")
		fmt.Println("Please run 'jj add -G' to add a global jump root.")
	}
}

func performWorkspaceJump(tmpFilePath string, args []string) {
	allPaths, err := util.GetGitWorktrees()
	if err != nil {
		fmt.Printf("Workspaces error: %v\n", err)
		os.Exit(1)
	}

	if len(allPaths) <= 1 {
		fmt.Println("No other git workspaces (worktrees) found for this repository")
		os.Exit(1)
	}

	targetPaths := allPaths

	if len(args) > 0 {
		searchTerm := strings.ToLower(args[0])
		var filtered []map[string]string
		for _, p := range allPaths {
			if util.FuzzyMatch(searchTerm, strings.ToLower(p["jumpPoint"])) {
				filtered = append(filtered, p)
			}
		}

		if len(filtered) == 1 {
			err = util.EchoCommand(tmpFilePath, filtered[0]["fullPath"])
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
			return
		}

		if len(filtered) > 1 {
			targetPaths = filtered
		}
	}

	prompt := promtui_local.PromptSelector(targetPaths, "Select a workspace/worktree")
	i, _, err := prompt.Run()
	if err != nil {
		if (err == promptui.ErrInterrupt || err == promptui.ErrEOF) && len(targetPaths) < len(allPaths) {
			prompt = promtui_local.PromptSelector(allPaths, "Select a workspace/worktree")
			i, _, err = prompt.Run()
			if err != nil {
				os.Exit(0)
			}
			err = util.EchoCommand(tmpFilePath, allPaths[i]["fullPath"])
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
			return
		}
		os.Exit(0)
	}

	err = util.EchoCommand(tmpFilePath, targetPaths[i]["fullPath"])
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
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

	rootCmd.PersistentFlags().BoolVarP(&globalJump, "global", "G", false, "Perform a global jump across registered projects or use as a modifier for other commands like add, verify, or remove")
	rootCmd.PersistentFlags().BoolVarP(&workspacesJump, "workspaces", "W", false, "Discovery: List all other git workspaces from current folder")
}
