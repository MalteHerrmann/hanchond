package query

import (
	"fmt"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// represents the query command.
var balanceCmd = &cobra.Command{
	Use:   "balance [wallet]",
	Args:  cobra.ExactArgs(1),
	Short: "Get the wallet balance",
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(errors.New("node not set"))
		}

		wallet := args[0]

		// TODO: this should also check the node type of the running node and then use the correct daemon
		e := evmos.NewEvmosFromDB(queries, nodeID)
		balance, err := e.CheckBalance(wallet)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get the balance: %w", err))
		}

		fmt.Println(balance)
	},
}

func init() {
	QueryCmd.AddCommand(balanceCmd)
}
