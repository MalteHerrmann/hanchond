package tx

import (
	"fmt"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// vote represents the vote command.
var voteCmd = &cobra.Command{
	Use:   "vote",
	Short: "Vote on all the active proposals",
	Run: func(cmd *cobra.Command, _ []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(errors.New("node not set"))
		}

		option, err := cmd.Flags().GetString("option")
		if err != nil {
			utils.ExitError(errors.New("option not set"))
		}

		e := evmos.NewEvmosFromDB(queries, nodeID)
		txhashes, err := e.VoteOnAllTheProposals(option)
		if err != nil {
			utils.ExitError(fmt.Errorf("error sending the transaction: %w", err))
		}
		for _, v := range txhashes {
			utils.Log("vote sent in tx: %s", v)
		}
	},
}

func init() {
	TxCmd.AddCommand(voteCmd)
	voteCmd.Flags().StringP("option", "o", "yes", "Vote option")
}
