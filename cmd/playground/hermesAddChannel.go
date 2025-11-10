package playground

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/database"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/gaia"
	"github.com/hanchon/hanchond/playground/hermes"
	"github.com/hanchon/hanchond/playground/sagaos"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// hermesAddChannelCmd represents the hermesAddChannel command.
var hermesAddChannelCmd = &cobra.Command{
	Use:   "hermes-add-channel [chain_id] [chain_id]",
	Args:  cobra.ExactArgs(2),
	Short: "It uses the hermes client to open an IBC channel between two chains",
	Long:  `This command requires that Hermes was already built and at least one node for each chain running.`,
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)

		chainOne := args[0]
		chainOneID, err := strconv.Atoi(chainOne)
		if err != nil {
			utils.ExitError(errors.New("invalid chain id"))
		}
		chainTwo := args[1]
		chainTwoID, err := strconv.Atoi(chainTwo)
		if err != nil {
			utils.ExitError(errors.New("invalid chain id"))
		}

		chains := make([]database.GetAllChainNodesRow, 2)
		nodesChainOne, err := queries.GetAllNodesForChainID(context.Background(), int64(chainOneID))
		if err != nil {
			utils.ExitError(fmt.Errorf("could not find nodes for chain: %s", chainOne))
		}
		chains[0] = nodesChainOne[0]

		nodesChainTwo, err := queries.GetAllNodesForChainID(context.Background(), int64(chainTwoID))
		if err != nil {
			utils.ExitError(fmt.Errorf("could not find nodes for chain: %s", chainTwo))
		}
		chains[1] = nodesChainTwo[0]

		h := hermes.NewHermes()
		utils.Log("Relayer initialized")

		for i, v := range chains {
			if v.IsRunning != 1 {
				utils.ExitError(fmt.Errorf("node %d of chain %d is not running; start first", v.NodeID, v.ChainID))
			}

			chainInfo := v.MustParseChainInfo()
			binaryName := chainInfo.GetBinaryName()

			switch binaryName {
			case gaia.ChainInfo.GetBinaryName():
				utils.Log("Adding %s chain", binaryName)
				if err := h.AddCosmosChain(
					chainInfo,
					v.ChainID_2,
					hermes.LocalEndpoint(v.P26657),
					hermes.LocalEndpoint(v.P9090),
					v.ValidatorKeyName,
					v.ValidatorKey,
				); err != nil {
					utils.ExitError(fmt.Errorf("error adding chain %d to the relayer: %s", i, err.Error()))
				}
			case evmos.ChainInfo.GetBinaryName(), sagaos.ChainInfo.GetBinaryName():
				utils.Log("Adding chain %d: %s", i, binaryName)
				if err := h.AddEVMChain(
					chainInfo,
					v.ChainID_2,
					hermes.LocalEndpoint(v.P26657),
					hermes.LocalEndpoint(v.P9090),
					v.ValidatorKeyName,
					v.ValidatorKey,
				); err != nil {
					utils.ExitError(fmt.Errorf("error adding chain %d to the relayer: %s", i, err.Error()))
				}
			default:
				utils.ExitError(fmt.Errorf("incorrect binary name: %s", binaryName))
			}

		}

		utils.Log("Calling create channel")
		err = h.CreateChannel(chains[0].ChainID_2, chains[1].ChainID_2)
		if err != nil {
			utils.ExitError(fmt.Errorf("error creating channel: %w", err))
		}

		utils.Log("Channel created")
	},
}

func init() {
	PlaygroundCmd.AddCommand(hermesAddChannelCmd)
}
