package tx

import (
	"github.com/spf13/cobra"

	"github.com/hanchon/hanchond/cmd/playground/tx/solidity"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/filesmanager"
)

// TxCmd represents the tx command.
var TxCmd = &cobra.Command{
	Use:     "tx",
	Aliases: []string{"t"},
	Short:   "Send transactions",
	Run: func(cmd *cobra.Command, _ []string) {
		filesmanager.SetHomeFolderFromCobraFlags(cmd)
		_ = cmd.Help()
		utils.ExitSuccess()
	},
}

func init() {
	TxCmd.AddCommand(solidity.SolidityCmd)
	TxCmd.PersistentFlags().
		StringP("node", "n", "1", "Playground node that is sending the transaction")
}
