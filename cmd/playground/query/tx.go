package query

import (
	"fmt"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// txCmd represents the query tx command
var txCmd = &cobra.Command{
	Use:   "tx [txhash]",
	Args:  cobra.ExactArgs(1),
	Short: "Get the transaction info",
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(fmt.Errorf("node not set"))
		}

		txhash := args[0]

		e := evmos.NewEvmosFromDB(queries, nodeID)
		resp, err := e.GetTransaction(txhash)
		if err != nil {
			utils.ExitError(fmt.Errorf("error sending the request: %w", err))
		}
		fmt.Println(resp)
	},
}

func init() {
	QueryCmd.AddCommand(txCmd)
}
