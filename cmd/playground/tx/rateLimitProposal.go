package tx

import (
	"fmt"
	"strings"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// rateLimitProposalCmd represents the rateLimit-proposal command
var rateLimitProposalCmd = &cobra.Command{
	Use:   "rate-limit-proposal [denom]",
	Args:  cobra.ExactArgs(1),
	Short: "Create an rate-limit propsal",
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(fmt.Errorf("node not set"))
		}

		channel, err := cmd.Flags().GetString("channel")
		if err != nil {
			utils.ExitError(fmt.Errorf("channel not set"))
		}

		duration, err := cmd.Flags().GetString("duration")
		if err != nil {
			utils.ExitError(fmt.Errorf("duration not set"))
		}

		maxSend, err := cmd.Flags().GetString("max-send")
		if err != nil {
			utils.ExitError(fmt.Errorf("max-send not set"))
		}

		maxRecv, err := cmd.Flags().GetString("max-recv")
		if err != nil {
			utils.ExitError(fmt.Errorf("max-recv not set"))
		}

		e := evmos.NewEvmosFromDB(queries, nodeID)
		txhash, err := e.CreateRateLimitProposal(evmos.RateLimitParams{
			Channel:  channel,
			Denom:    strings.TrimSpace(args[0]),
			MaxSend:  maxSend,
			MaxRecv:  maxRecv,
			Duration: duration,
		})
		if err != nil {
			utils.ExitError(fmt.Errorf("error sending the transaction: %w", err))
		}

		fmt.Println(txhash)
	},
}

func init() {
	TxCmd.AddCommand(rateLimitProposalCmd)
	rateLimitProposalCmd.Flags().StringP("channel", "c", "channel-0", "IBC channel")
	rateLimitProposalCmd.Flags().String("max-send", "10", "Max send rate")
	rateLimitProposalCmd.Flags().String("max-recv", "10", "Max recv rate")
	rateLimitProposalCmd.Flags().String("duration", "24", "Duration in hours")
}
