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

// txCmd represents the query tx command.
var txCmd = &cobra.Command{
	Use:   "tx [txhash]",
	Args:  cobra.ExactArgs(1),
	Short: "Get the transaction info",
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(errors.New("node flag not set"))
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
			utils.ExitError(fmt.Errorf("failed to get daemon: %w", err))
		}

		resp, err := d.Tx(args[0])
		if err != nil {
			utils.ExitError(fmt.Errorf("error sending the request: %w", err))
		}
		fmt.Println(resp)
	},
}

func init() {
	QueryCmd.AddCommand(txCmd)
}
