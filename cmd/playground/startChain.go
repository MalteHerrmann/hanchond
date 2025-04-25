package playground

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/database"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/gaia"
	"github.com/hanchon/hanchond/playground/sagaos"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// startChainCmd represents the startChainCmd
//
// TODO: refactor all of the available commands to more layers for easier contextualization.
// e.g. use `h p chain start` and `h p node info` instead of `h p start-chain` and `h p get-node`
var startChainCmd = &cobra.Command{
	Use:   "start-chain [chain_id]",
	Args:  cobra.ExactArgs(1),
	Short: "Start all the validators of the chain",
	Long:  `Start all the required processes to run the chain`,
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)

		chainNumber, err := strconv.Atoi(strings.TrimSpace(args[0]))
		if err != nil {
			utils.ExitError(fmt.Errorf("invalid chain id: %w", err))
		}
		nodes, err := queries.GetAllNodesForChainID(context.Background(), int64(chainNumber))
		if err != nil {
			utils.ExitError(fmt.Errorf("could not find the chain: %w", err))
		}

		for _, v := range nodes {
			binaryName := v.MustParseChainInfo().GetBinaryName()

			var pID int
			var err error
			switch binaryName {
			case gaia.ChainInfo.GetBinaryName():
				d := gaia.NewGaia(v.Moniker, v.ConfigFolder, v.ChainID_2, v.ValidatorKeyName)
				pID, err = d.Start()
			case evmos.ChainInfo.GetBinaryName():
				d := evmos.NewEvmos(v.Moniker, v.Version, v.ConfigFolder, v.ChainID_2, v.ValidatorKeyName)
				pID, err = d.Start()
			case sagaos.ChainInfo.GetBinaryName():
				d := sagaos.NewSagaOS(v.Moniker, v.Version, v.ConfigFolder, v.ChainID_2, v.ValidatorKeyName)
				pID, err = d.Start()
			default:
				utils.ExitError(fmt.Errorf("binary %s not configured", binaryName))
			}

			if err != nil {
				utils.ExitError(fmt.Errorf("could not start the node: %w", err))
			}

			fmt.Println("Node is running with pID:", pID)
			err = queries.SetProcessID(context.Background(), database.SetProcessIDParams{
				ProcessID: int64(pID),
				IsRunning: 1,
				ID:        v.ID,
			})
			if err != nil {
				utils.ExitError(fmt.Errorf("could not save the process ID to the db: %w", err))
			}
		}
	},
}

func init() {
	PlaygroundCmd.AddCommand(startChainCmd)
}
