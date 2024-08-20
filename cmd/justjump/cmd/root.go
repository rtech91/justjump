package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "justjump",
	Short: "JustJump is a simple tool to help you jump between directories quickly.",
	Long: `JustJump is a simple tool to help you jump between directories quickly.
You can use it as "justjump" or "jj" shortcut.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
