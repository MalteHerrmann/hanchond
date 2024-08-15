package erc20

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hanchon/hanchond/lib/requester"
	"github.com/hanchon/hanchond/lib/smartcontract/erc20"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// supplyCmd represents the supply command
var supplyCmd = &cobra.Command{
	Use:   "supply [contract]",
	Args:  cobra.ExactArgs(1),
	Short: "Get the wallet supply",
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			fmt.Println("node not set")
			os.Exit(1)
		}
		contract := strings.TrimSpace(args[0])
		e := evmos.NewEvmosFromDB(queries, nodeID)
		client := requester.NewClient().WithUnsecureWeb3Endpoint(fmt.Sprintf("http://localhost:%d", e.Ports.P8545))

		height, _ := cmd.Flags().GetString("height")
		heightInt := erc20.Latest
		if height != "latest" {
			temp, err := strconv.ParseInt(height, 10, 64)
			if err != nil {
				fmt.Printf("invalid height: %s\n", err.Error())
				os.Exit(1)
			}
			heightInt = int(temp)
		}

		supply, err := client.GetTotalSupply(contract, heightInt)
		if err != nil {
			fmt.Println("could not get the supply:", err.Error())
			os.Exit(1)
		}
		fmt.Println(supply)
	},
}

func init() {
	ERC20Cmd.AddCommand(supplyCmd)
	supplyCmd.Flags().String("height", "latest", "Query at the given height.")
}
