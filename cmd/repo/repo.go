package repo

import (
	"github.com/spf13/cobra"

	"github.com/hanchon/hanchond/lib/utils"
)

// RepoCmd represents the repo command.
var RepoCmd = &cobra.Command{
	Use:     "repo",
	Aliases: []string{"r"},
	Short:   "Repo management utils",
	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
		utils.ExitSuccess()
	},
}

func init() {
	RepoCmd.AddCommand(BumpModuleVersionCmd)
}
