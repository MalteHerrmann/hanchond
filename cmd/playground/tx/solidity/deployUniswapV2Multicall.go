package solidity

import (
	"encoding/hex"
	"fmt"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/solidity"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// deployUniswapV2MulticallyCmd represents the deploy command
var deployUniswapV2MulticallyCmd = &cobra.Command{
	Use:   "deploy-uniswap-v2-multicall",
	Args:  cobra.ExactArgs(0),
	Short: "Deploy uniswap v2 multicall",
	Run: func(cmd *cobra.Command, _ []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(fmt.Errorf("node not set"))
		}

		gasLimit, err := cmd.Flags().GetInt("gas-limit")
		if err != nil {
			utils.ExitError(fmt.Errorf("incorrect gas limit"))
		}

		// TODO: allow mainnet as a valid endpoint
		e := evmos.NewEvmosFromDB(queries, nodeID)
		builder := e.NewTxBuilder(uint64(gasLimit))

		contractName := "/Multicall"
		// Clone v2-minified if needed
		path, err := solidity.DownloadUniswapV2Minified()
		if err != nil {
			utils.ExitError(fmt.Errorf("error downloading the uniswap v2 minified: %w", err))
		}

		// Keep working with the main contract
		path = path + "/contracts" + contractName + ".sol"

		// Set up temp folder
		if err := filesmanager.CleanUpTempFolder(); err != nil {
			utils.ExitError(fmt.Errorf("could not clean up the temp folder: %w", err))
		}

		folderName := "multicallBuilder"
		if err := filesmanager.CreateTempFolder(folderName); err != nil {
			utils.ExitError(fmt.Errorf("could not create the temp folder: %w", err))
		}

		// Compile the contract
		err = solidity.CompileWithSolc("0.5.0", path, filesmanager.GetBranchFolder(folderName))
		if err != nil {
			utils.ExitError(fmt.Errorf("could not compile the erc20 contract: %w", err))
		}

		bytecode, err := filesmanager.ReadFile(filesmanager.GetBranchFolder(folderName) + contractName + ".bin")
		if err != nil {
			utils.ExitError(fmt.Errorf("error reading the bytecode file: %w", err))
		}

		bytecode, err = hex.DecodeString(string(bytecode))
		if err != nil {
			utils.ExitError(fmt.Errorf("error converting bytecode to []byte: %w", err))
		}

		txHash, err := builder.DeployContract(0, bytecode, uint64(gasLimit))
		if err != nil {
			utils.ExitError(fmt.Errorf("error sending the transaction: %w", err))
		}

		contractAddress, err := e.NewRequester().GetContractAddress(txHash)
		if err != nil {
			utils.ExitError(fmt.Errorf("error getting the contract address: %w", err))
		}

		codeHash, err := e.NewRequester().EthCodeHash(contractAddress, "latest")
		if err != nil {
			utils.ExitError(fmt.Errorf("failed to get the eth code: %w", err))
		}

		fmt.Printf("{\"contract_address\":\"%s\", \"code_hash\":\"%s\", \"tx_hash\":\"%s\"}\n", contractAddress, "0x"+codeHash, txHash)

		// Clean up files
		if err := filesmanager.CleanUpTempFolder(); err != nil {
			utils.ExitError(fmt.Errorf("could not clean up the temp folder: %w", err))
		}
		utils.ExitSuccess()
	},
}

func init() {
	SolidityCmd.AddCommand(deployUniswapV2MulticallyCmd)
	deployUniswapV2MulticallyCmd.Flags().Int("gas-limit", 20_000_000, "GasLimit to be used to deploy the transaction")
}
