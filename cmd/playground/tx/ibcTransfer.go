package tx

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/sql"
)

// ibcTransferCmd represents the ibc-transfer command.
var ibcTransferCmd = &cobra.Command{
	Use:     "ibc-transfer wallet amount",
	Args:    cobra.ExactArgs(2),
	Aliases: []string{"it"},
	Short:   "Sends an ibc transaction",
	Long:    `It sends an IBC transfer from the validator wallet to the destination wallet`,
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(fmt.Errorf("node not set"))
		}

		channel, err := cmd.Flags().GetString("channel")
		if err != nil {
			utils.ExitError(fmt.Errorf("ibc channel not set"))
		}

		dstWallet := args[0]
		amount := args[1]

		e := evmos.NewEvmosFromDB(queries, nodeID)
		denom, err := cmd.Flags().GetString("denom")
		if err != nil {
			utils.ExitError(fmt.Errorf("denom not set"))
		}

		if denom == "" {
			denom = e.BaseDenom
		}

		out, err := e.SendIBC("transfer", channel, dstWallet, amount+denom)
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
		StringP("denom", "d", "", "Denom that you are sending, it defaults to the base denom of the chain")
}
