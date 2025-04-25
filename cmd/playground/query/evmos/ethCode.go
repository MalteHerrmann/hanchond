package evmos

import (
	"fmt"

	"github.com/hanchon/hanchond/lib/requester"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/cosmosdaemon"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// ethCodeCmd represents the ethCode command
var ethCodeCmd = &cobra.Command{
	Use:   "eth-code [address]",
	Args:  cobra.ExactArgs(1),
	Short: "Get the smartcontract registered eth code",
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)

		height, err := cmd.Flags().GetString("height")
		if err != nil {
			utils.ExitError(fmt.Errorf("error getting the request height: %w", err))
		}

		endpoint, err := cosmosdaemon.GetWeb3Endpoint(queries, cmd)
		if err != nil {
			utils.ExitError(fmt.Errorf("error generting web3 endpoint: %w", err))
		}
		client := requester.NewClient().WithUnsecureWeb3Endpoint(endpoint)

		code, err := client.EthCode(args[0], height)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get the ethCode: %w", err))
		}

		fmt.Println(string(code))
	},
}

func init() {
	EvmosCmd.AddCommand(ethCodeCmd)
	ethCodeCmd.Flags().String("height", "latest", "Query at the given height.")
}
