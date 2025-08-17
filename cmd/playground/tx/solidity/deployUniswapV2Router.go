package solidity

import (
	"encoding/hex"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/hanchon/hanchond/lib/smartcontract"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/solidity"
	"github.com/hanchon/hanchond/playground/sql"
)

// deployUniswapV2RouteryCmd represents the deploy command.
var deployUniswapV2RouteryCmd = &cobra.Command{
	Use:   "deploy-uniswap-v2-router [factory_address] [wrapped_coin_address]",
	Args:  cobra.ExactArgs(2),
	Short: "Deploy uniswap v2 router",
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(fmt.Errorf("node not set"))
		}

		gasLimit, err := cmd.Flags().GetInt("gas-limit")
		if err != nil {
			utils.ExitError(fmt.Errorf("incorrect gas limit"))
		}

		factoryAddress := args[0]
		wrappedCoinAddress := args[1]

		// TODO: allow mainnet as a valid endpoint
		e := evmos.NewEvmosFromDB(queries, nodeID)
		builder := e.NewTxBuilder(uint64(gasLimit))

		factoryCodeHash, err := e.NewRequester().EthCodeHash(factoryAddress, "latest")
		if err != nil {
			utils.ExitError(fmt.Errorf("failed to get the eth code: %w", err))
		}

		contractName := "/Router"
		// Clone v2-minified if needed
		path, err := solidity.DownloadUniswapV2Minified()
		if err != nil {
			utils.ExitError(fmt.Errorf("error downloading the uniswap v2 minified: %w", err))
		}

		// Keep working with the main contract
		path = path + "/contracts" + contractName + ".sol"
		libFile, err := filesmanager.ReadFile(path)
		if err != nil {
			utils.ExitError(fmt.Errorf("error opening the router file: %w", err))
		}

		regex := regexp.MustCompile(`hex".{3,}"`)
		libFile = regex.ReplaceAll(libFile, []byte(fmt.Sprintf("hex'%s'", factoryCodeHash)))
		if err := filesmanager.SaveFile(libFile, path); err != nil {
			utils.ExitError(fmt.Errorf("error saving the router file: %w", err))
		}

		// Set up temp folder
		if err := filesmanager.CleanUpTempFolder(); err != nil {
			utils.ExitError(fmt.Errorf("could not clean up the temp folder: %w", err))
		}

		folderName := "routerBuilder"
		if err := filesmanager.CreateTempFolder(folderName); err != nil {
			utils.ExitError(fmt.Errorf("could not create the temp folder: %w", err))
		}

		// Compile the contract
		err = solidity.CompileWithSolc("0.6.6", path, filesmanager.GetBranchFolder(folderName))
		if err != nil {
			utils.ExitError(fmt.Errorf("could not compile the erc20 contract: %w", err))
		}

		contractName = "/UniswapV2Router02"

		bytecode, err := filesmanager.ReadFile(
			filesmanager.GetBranchFolder(folderName) + contractName + ".bin",
		)
		if err != nil {
			utils.ExitError(fmt.Errorf("error reading the bytecode file: %w", err))
		}

		bytecode, err = hex.DecodeString(string(bytecode))
		if err != nil {
			utils.ExitError(fmt.Errorf("error converting bytecode to []byte: %w", err))
		}

		// Generate the constructor
		abiBytes, err := filesmanager.ReadFile(
			filesmanager.GetBranchFolder(folderName) + contractName + ".abi",
		)
		if err != nil {
			utils.ExitError(fmt.Errorf("error reading the abi file: %w", err))
		}

		// Get Params
		callArgs, err := smartcontract.StringsToABIArguments(
			[]string{
				fmt.Sprintf("a:%s", factoryAddress),
				fmt.Sprintf("a:%s", wrappedCoinAddress),
			},
		)
		if err != nil {
			utils.ExitError(fmt.Errorf("error converting arguments: %w", err))
		}

		callData, err := smartcontract.ABIPackRaw(abiBytes, "", callArgs...)
		if err != nil {
			utils.ExitError(fmt.Errorf("error converting arguments: %w", err))
		}
		bytecode = append(bytecode, callData...)

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

		fmt.Printf(
			"{\"contract_address\":\"%s\", \"code_hash\":\"%s\", \"tx_hash\":\"%s\"}\n",
			contractAddress,
			"0x"+codeHash,
			txHash,
		)

		// Clean up files
		if err := filesmanager.CleanUpTempFolder(); err != nil {
			utils.ExitError(fmt.Errorf("could not clean up the temp folder: %w", err))
		}
		utils.ExitSuccess()
	},
}

func init() {
	SolidityCmd.AddCommand(deployUniswapV2RouteryCmd)
	deployUniswapV2RouteryCmd.Flags().
		Int("gas-limit", 20_000_000, "GasLimit to be used to deploy the transaction")
}
