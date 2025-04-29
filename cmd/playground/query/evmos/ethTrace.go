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

// ethTraceCmd represents the ethTrace command
var ethTraceCmd = &cobra.Command{
	Use:   "eth-trace [tx_hash]",
	Args:  cobra.ExactArgs(1),
	Short: "Get the trace for the given tx hash",
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)

		endpoint, err := cosmosdaemon.GetWeb3Endpoint(queries, cmd)
		if err != nil {
			utils.ExitError(fmt.Errorf("error generting web3 endpoint: %w", err))
		}
		client := requester.NewClient().WithUnsecureWeb3Endpoint(endpoint)

		receipt, err := client.GetTransactionTrace(args[0])
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get the ethTrace: %w", err))
		}

		val, err := json.Marshal(receipt.Result)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not process the ethTrace: %w", err))
		}

		fmt.Println(string(val))
	},
}

func init() {
	EvmosCmd.AddCommand(ethTraceCmd)
}
