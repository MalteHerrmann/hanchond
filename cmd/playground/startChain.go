package playground

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/cosmosdaemon"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// startChainCmd represents the startChainCmd
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
			chainInfo := v.MustParseChainInfo()

			daemon := cosmosdaemon.NewDameon(
				chainInfo,
				v.Moniker,
				v.Version,
				v.ConfigFolder,
				v.ChainID_2,
				v.ValidatorKeyName,
			)

			pid, err := daemon.StartNodeAndStoreInfo(queries, v.ID)
			if err != nil {
				utils.ExitError(fmt.Errorf("could not start the node: %w", err))
			}

			utils.Log("Node is running with pID: %d", pid)
		}
	},
}

func init() {
	PlaygroundCmd.AddCommand(startChainCmd)
}
