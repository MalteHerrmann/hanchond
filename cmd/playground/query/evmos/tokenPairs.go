package evmos

import (
	"encoding/json"
	"fmt"

	"github.com/hanchon/hanchond/lib/requester"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// tokenPairsCmd represents the tokenPairs command.
var tokenPairsCmd = &cobra.Command{
	Use:   "token-pairs",
	Short: "Get the network registered token-pairs",
	Run: func(cmd *cobra.Command, _ []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(errors.New("node not set"))
		}

		e := evmos.NewEvmosFromDB(queries, nodeID)
		client := requester.NewClient().WithUnsecureRestEndpoint(fmt.Sprintf("http://localhost:%d", e.Ports.P1317))
		pairs, err := client.GetEvmosERC20TokenPairs()
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get the tokenPairs: %w", err))
		}
		values, err := json.Marshal(pairs.TokenPairs)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not marshal response: %w", err))
		}

		fmt.Println(string(values))
	},
}

func init() {
	EvmosCmd.AddCommand(tokenPairsCmd)
}
