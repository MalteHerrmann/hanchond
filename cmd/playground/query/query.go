package query

import (
	"github.com/spf13/cobra"

	"github.com/hanchon/hanchond/cmd/playground/query/erc20"
	"github.com/hanchon/hanchond/cmd/playground/query/evmos"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/filesmanager"
)

// QueryCmd represents the query command.
var QueryCmd = &cobra.Command{
	Use:     "query",
	Aliases: []string{"q"},
	Short:   "Query the blockchain data",
	Run: func(cmd *cobra.Command, _ []string) {
		filesmanager.SetHomeFolderFromCobraFlags(cmd)
		_ = cmd.Help()
		utils.ExitSuccess()
	},
}

func init() {
	QueryCmd.AddCommand(erc20.ERC20Cmd)
	QueryCmd.AddCommand(evmos.EvmosCmd)
	QueryCmd.PersistentFlags().
		StringP("node", "n", "1", "Playground node used to get the information")
}
