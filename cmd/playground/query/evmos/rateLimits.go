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

// rateLimitsCmd represents the ibc-rate-limits command.
var rateLimitsCmd = &cobra.Command{
	Use:   "ibc-rate-limits",
	Short: "Get all active IBC rate limits",
	Run: func(cmd *cobra.Command, _ []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(errors.New("node not set"))
		}

		e := evmos.NewEvmosFromDB(queries, nodeID)
		client := requester.NewClient().WithUnsecureRestEndpoint(fmt.Sprintf("http://localhost:%d", e.Ports.P1317))
		rateLimits, err := client.GetIBCRateLimits()
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get the rateLimits: %w", err))
		}
		values, err := json.Marshal(rateLimits)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not marshal response: %w", err))
		}

		fmt.Println(string(values))
	},
}

func init() {
	EvmosCmd.AddCommand(rateLimitsCmd)
}
