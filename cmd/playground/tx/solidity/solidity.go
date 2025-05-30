package solidity

import (
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/spf13/cobra"
)

// SolidityCmd represents the solidity command
var SolidityCmd = &cobra.Command{
	Use:     "solidity",
	Aliases: []string{"s"},
	Short:   "Send transactions related to solidity contracts",
	Run: func(cmd *cobra.Command, _ []string) {
		filesmanager.SetHomeFolderFromCobraFlags(cmd)
		_ = cmd.Help()
		utils.ExitSuccess()
	},
}
