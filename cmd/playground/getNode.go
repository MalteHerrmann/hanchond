package playground

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/hanchon/hanchond/lib/converter"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/sql"
)

// CLI flags.
var (
	getBinary, getChainID, getHome, getVal bool
	retrievedPort                          uint16
)

// getNodeCmd represents the getNode command.
var getNodeCmd = &cobra.Command{
	Use:   "get-node [id]",
	Args:  cobra.ExactArgs(1),
	Short: "Get the node configuration",
	Long:  `It will return the node configuration that is stored in the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)

		id := args[0]
		idNumber, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not parse the ID: %w", err))
		}

		ports, err := queries.GetNodePorts(context.Background(), idNumber)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get the ports: %w", err))
		}

		// This means the port was specified
		if retrievedPort != 0 {
			fmt.Println(ports.Get(retrievedPort))
			utils.ExitSuccess()
		}

		node, err := queries.GetNode(context.Background(), idNumber)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get the node: %w", err))
		}

		if getHome {
			fmt.Println(node.ConfigFolder)
			utils.ExitSuccess()
		}

		chain, err := queries.GetChain(context.Background(), node.ChainID)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not get the chain: %w", err))
		}

		// retrieve only binary
		if getBinary {
			fmt.Println(
				filesmanager.GetDaemondPathWithVersion(chain.MustParseChainInfo(), node.Version),
			)
			utils.ExitSuccess()
		}

		if getVal {
			fmt.Println(node.ValidatorWallet)
			utils.ExitSuccess()
		}

		if getChainID {
			fmt.Println(chain.ChainID)
			utils.ExitSuccess()
		}

		hexWallet, err := converter.Bech32ToHex(node.ValidatorWallet)
		if err != nil {
			utils.ExitError(fmt.Errorf("could not convert validator wallet to eth: %w", err))
		}

		fmt.Printf(`Node: %d
General Configuration:
    - Binary: %s
    - ChainID: %s
    - Home Folder: %s
Process:
    - IsRunning: %d
    - ProcessID: %d
Keys:
    - KeyName: %s
    - Mnemonic: %s
    - Wallet: %s
    - Wallet(hex): %s
Ports:
    - 8545(web3): %d
    - 26657(cli/tendermint): %d
    - 1317(cosmos rest): %d
    - 9090(grpc): %d
`,
			idNumber,
			chain.MustParseChainInfo().GetVersionedBinaryName(node.Version),
			chain.ChainID,
			node.ConfigFolder,
			node.IsRunning,
			node.ProcessID,
			node.ValidatorKeyName,
			node.ValidatorKey,
			node.ValidatorWallet,
			hexWallet,
			ports.P8545,
			ports.P26657,
			ports.P1317,
			ports.P9090,
		)
	},
}

func init() {
	getNodeCmd.Flags().BoolVarP(&getBinary, "bin", "b", false, "Get the node's running binary path")
	getNodeCmd.Flags().
		BoolVarP(&getChainID, "chain-id", "c", false, "Get the chain ID of the node's network")
	getNodeCmd.Flags().BoolVarP(&getHome, "node-home", "", false, "Get the node's home folder")
	getNodeCmd.Flags().BoolVarP(&getVal, "val", "v", false, "Get the node's validator address")
	getNodeCmd.Flags().Uint16VarP(&retrievedPort, "port", "p", 0, "Get the node's remapped port")

	PlaygroundCmd.AddCommand(getNodeCmd)
}
