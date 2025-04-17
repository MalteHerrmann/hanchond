package playground

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/hanchon/hanchond/playground/database"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/gaia"
	"github.com/hanchon/hanchond/playground/hermes"
	"github.com/hanchon/hanchond/playground/sagaos"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// hermesAddChannelCmd represents the hermesAddChannel command
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
			fmt.Println("invalid chain id")
			os.Exit(1)
		}
		chainTwo := args[1]
		chainTwoID, err := strconv.Atoi(chainTwo)
		if err != nil {
			fmt.Println("invalid chain id")
			os.Exit(1)
		}
		chains := make([]database.GetAllChainNodesRow, 2)
		nodesChainOne, err := queries.GetAllNodesForChainID(context.Background(), int64(chainOneID))
		if err != nil {
			fmt.Println("could not find nodes for chain:", chainOne)
			os.Exit(1)
		}
		chains[0] = nodesChainOne[0]

		nodesChainTwo, err := queries.GetAllNodesForChainID(context.Background(), int64(chainTwoID))
		if err != nil {
			fmt.Println("could not find nodes for chain:", chainTwo)
			os.Exit(1)
		}
		chains[1] = nodesChainTwo[0]

		h := hermes.NewHermes()
		fmt.Println("Relayer initialized")

		for i, v := range chains {
			if v.IsRunning != 1 {
				fmt.Printf("node %d of chain %d is not running; start first\n", v.NodeID, v.ChainID)
				os.Exit(1)
			}

			chainInfo := v.MustParseChainInfo()
			binaryName := chainInfo.GetBinaryName()

			switch binaryName {
			case gaia.ChainInfo.GetBinaryName():
				fmt.Printf("Adding %s chain\n", binaryName)
				if err := h.AddCosmosChain(
					chainInfo,
					v.ChainID_2,
					hermes.LocalEndpoint(v.P26657),
					hermes.LocalEndpoint(v.P9090),
					v.ValidatorKeyName,
					v.ValidatorKey,
				); err != nil {
					fmt.Printf("error adding chain %d to the relayer: %s\n", i, err.Error())
					os.Exit(1)
				}
			case evmos.ChainInfo.GetBinaryName(), sagaos.ChainInfo.GetBinaryName():
				fmt.Printf("Adding %s chain\n", binaryName)
				if err := h.AddEVMChain(
					chainInfo,
					v.ChainID_2,
					hermes.LocalEndpoint(v.P26657),
					hermes.LocalEndpoint(v.P9090),
					v.ValidatorKeyName,
					v.ValidatorKey,
				); err != nil {
					fmt.Printf("error adding chain %d to the relayer: %s\n", i, err.Error())
					os.Exit(1)
				}
			default:
				fmt.Println("incorrect binary name: ", binaryName)
				os.Exit(1)
			}

		}

		fmt.Println("Calling create channel")
		err = h.CreateChannel(chains[0].ChainID_2, chains[1].ChainID_2)
		if err != nil {
			fmt.Println("error creating channel", err.Error())
			os.Exit(1)
		}
		fmt.Println("Channel created")
	},
}

func init() {
	PlaygroundCmd.AddCommand(hermesAddChannelCmd)
}
