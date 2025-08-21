package tx

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hanchon/hanchond/cmd/playground/common"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/hanchon/hanchond/playground/types"
)

// ibcTransferCmd represents the ibc-transfer command.
var ibcTransferCmd = &cobra.Command{
	Use:     "ibc-transfer wallet amount",
	Args:    cobra.ExactArgs(2),
	Aliases: []string{"it"},
	Short:   "Sends an IBC transaction through the selected channel",
	Long: `Sends an IBC transfer from the validator wallet of the selected node ` +
		`to the given destination wallet on the receiving end of the selected IBC channel.`,
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(errors.New("node not set"))
		}

		channel, err := cmd.Flags().GetString("channel")
		if err != nil {
			utils.ExitError(errors.New("ibc channel not set"))
		}

		dstWallet := args[0]
		sentFunds, err := types.ParseCoin(args[1])
		if err != nil {
			utils.ExitError(fmt.Errorf("invalid coin amount; got: %w", err))
		}

		nID, err := strconv.ParseInt(nodeID, 10, 64)
		if err != nil {
			utils.ExitError(fmt.Errorf("invalid node ID: %s", nodeID))
		}

		node, err := queries.GetChainNode(context.Background(), nID)
		if err != nil {
			utils.ExitError(errors.New("node not found"))
		}

		ports := node.GetPorts()

		d, err := commoncmd.GetDaemonForNode(node.GetDaemonInfo(), &ports)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get the for node: %w", err))
		}

		memo, err := cmd.Flags().GetString("memo")
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get the memo: %w", err))
		}

		out, err := d.SendIBC("transfer", channel, dstWallet, sentFunds, memo)
		if err != nil {
			utils.ExitError(fmt.Errorf("error sending the transaction: %w", err))
		}

		if !strings.Contains(out, "code: 0") {
			utils.ExitError(fmt.Errorf("transaction failed: %s", out))
		}

		hash := strings.Split(out, "txhash: ")
		if len(hash) > 1 {
			hash[1] = strings.TrimSpace(hash[1])
		}

		utils.Log("transaction sent: %s", hash[1])
	},
}

func init() {
	TxCmd.AddCommand(ibcTransferCmd)
	ibcTransferCmd.Flags().StringP("channel", "c", "channel-0", "IBC channel")
	ibcTransferCmd.Flags().
		StringP("memo", "m", "", "Optional field to pass an IBC memo along with the transfer")
}
