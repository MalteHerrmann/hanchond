package relayer

import (
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/spf13/cobra"
)

// RelayerCmd represents the relayer command
var RelayerCmd = &cobra.Command{
	Use:     "relayer",
	Aliases: []string{"r"},
	Short:   "Relayer related functions",
	Run: func(cmd *cobra.Command, _ []string) {
		filesmanager.SetHomeFolderFromCobraFlags(cmd)
		_ = cmd.Help()
		utils.ExitSuccess()
	},
}
