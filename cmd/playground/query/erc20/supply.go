package erc20

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hanchon/hanchond/lib/requester"
	"github.com/hanchon/hanchond/lib/smartcontract/erc20"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/cosmosdaemon"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// supplyCmd represents the supply command.
var supplyCmd = &cobra.Command{
	Use:   "supply [contract]",
	Args:  cobra.ExactArgs(1),
	Short: "Get the wallet supply",
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)
		contract := strings.TrimSpace(args[0])

		endpoint, err := cosmosdaemon.GetWeb3Endpoint(queries, cmd)
		if err != nil {
			utils.ExitError(fmt.Errorf("error generting web3 endpoint: %w", err))
		}
		client := requester.NewClient().WithUnsecureWeb3Endpoint(endpoint)

		height, _ := cmd.Flags().GetString("height")
		heightInt := erc20.Latest
		if height != "latest" {
			temp, err := strconv.ParseInt(height, 10, 64)
			if err != nil {
				utils.ExitError(fmt.Errorf("invalid height: %w", err))
			}
			heightInt = int(temp)
		}

		supply, err := client.GetTotalSupply(contract, heightInt)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get the supply: %w", err))
		}

		fmt.Println(supply)
	},
}

func init() {
	ERC20Cmd.AddCommand(supplyCmd)
	supplyCmd.Flags().String("height", "latest", "Query at the given height.")
}
