package query

import (
	"fmt"

	"github.com/hanchon/hanchond/lib/requester"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// heightCmd represents the query height command
var heightCmd = &cobra.Command{
	Use:   "height",
	Short: "Get the current networkt height",
	Run: func(cmd *cobra.Command, _ []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(fmt.Errorf("node not set"))
		}

		e := evmos.NewEvmosFromDB(queries, nodeID)
		client := requester.NewClient().WithUnsecureTendermintEndpoint(fmt.Sprintf("http://localhost:%d", e.Ports.P26657))
		height, err := client.GetCurrentHeight()
		if err != nil {
			utils.ExitError(fmt.Errorf("could not query the current height: %w", err))
		}
		fmt.Println(height)
	},
}

func init() {
	QueryCmd.AddCommand(heightCmd)
}
