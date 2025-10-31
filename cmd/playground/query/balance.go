package query

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	commoncmd "github.com/hanchon/hanchond/cmd/playground/common"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/sql"
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

		idNumber, err := strconv.ParseInt(nodeID, 10, 64)
		if err != nil {
			utils.ExitError(fmt.Errorf("failed to parse node id: %w", err))
		}

		node, err := queries.GetChainNode(context.Background(), idNumber)
		if err != nil {
			utils.ExitError(fmt.Errorf("failed to get chain node: %w", err))
		}

		ports := node.GetPorts()
		d, err := commoncmd.GetDaemonForNode(node.GetDaemonInfo(), &ports)
		if err != nil {
			utils.ExitError(fmt.Errorf("failed to get node daemon: %w", err))
		}

		balance, err := d.Balance(args[0])
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get the balance: %w", err))
		}

		fmt.Println(balance)
	},
}

func init() {
	QueryCmd.AddCommand(balanceCmd)
}
