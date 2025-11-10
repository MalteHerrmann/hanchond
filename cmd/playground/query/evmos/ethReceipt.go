package evmos

import (
	"encoding/json"
	"fmt"

	"github.com/hanchon/hanchond/lib/requester"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/cosmosdaemon"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// ethReceiptCmd represents the ethReceipt command.
var ethReceiptCmd = &cobra.Command{
	Use:   "eth-receipt [tx_hash]",
	Args:  cobra.ExactArgs(1),
	Short: "Get the receipt for the given tx hash",
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)

		endpoint, err := cosmosdaemon.GetWeb3Endpoint(queries, cmd)
		if err != nil {
			utils.ExitError(fmt.Errorf("error generting web3 endpoint: %w", err))
		}
		client := requester.NewClient().WithUnsecureWeb3Endpoint(endpoint)

		receipt, err := client.GetTransactionReceipt(args[0])
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get the ethReceipt: %w", err))
		}

		val, err := json.Marshal(receipt.Result)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not process the ethReceipt: %w", err))
		}

		fmt.Println(string(val))
	},
}

func init() {
	EvmosCmd.AddCommand(ethReceiptCmd)
}
