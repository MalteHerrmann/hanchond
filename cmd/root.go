package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hanchon/hanchond/cmd/convert"
	"github.com/hanchon/hanchond/cmd/playground"
	"github.com/hanchon/hanchond/cmd/repo"
	"github.com/hanchon/hanchond/lib/utils"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "hanchond",
	Short: "Hanchon's toolkit",
	Long:  `Stop re-writing the same scripts one million times`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		utils.ExitError(err)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(VersionCmd)
	rootCmd.AddCommand(convert.ConvertCmd)
	rootCmd.AddCommand(playground.PlaygroundCmd)
	rootCmd.AddCommand(repo.RepoCmd)
}
