package tx

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hanchon/hanchond/lib/requester"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// upgradeProposalCmd represents the upgrade-proposal command.
var upgradeProposalCmd = &cobra.Command{
	Use:   "upgrade-proposal [version]",
	Args:  cobra.ExactArgs(1),
	Short: "Create an upgrade proposal. It defaults to 20 blocks in the future (1min).",
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(errors.New("node not set"))
		}
		version := strings.TrimSpace(args[0])

		e := evmos.NewEvmosFromDB(queries, nodeID)
		height, err := cmd.Flags().GetString("height")
		if err != nil || height == "" {
			diff, err := cmd.Flags().GetString("height-diff")
			if err != nil {
				utils.ExitError(fmt.Errorf("could not read any height related flag: %w", err))
			}
			diffInt, err := strconv.Atoi(diff)
			if err != nil {
				utils.ExitError(fmt.Errorf("could not convert diff to int: %w", err))
			}
			currentHeight, err := requester.
				NewClient().
				WithUnsecureTendermintEndpoint(
					fmt.Sprintf("http://localhost:%d", e.Ports.P26657),
				).GetCurrentHeight()
			if err != nil {
				utils.ExitError(fmt.Errorf("could not get current height: %w", err))
			}
			currentHeightInt, err := strconv.Atoi(currentHeight)
			if err != nil {
				utils.ExitError(fmt.Errorf("could convert height response to int: %w", err))
			}
			height = strconv.Itoa(currentHeightInt + diffInt)
		}

		txhash, err := e.CreateUpgradeProposal(version, height)
		if err != nil {
			utils.ExitError(fmt.Errorf("error sending the transaction: %w", err))
		}

		fmt.Printf("{\"txhash\":\"%s\", \"height\": %s}\n", txhash, height)
	},
}

func init() {
	TxCmd.AddCommand(upgradeProposalCmd)
	upgradeProposalCmd.Flags().String("height", "", "Upgrade height.")
	upgradeProposalCmd.Flags().String("height-diff", "20", "Blocks in the future when the upgrade is going to be executed.") //nolint:lll
}
