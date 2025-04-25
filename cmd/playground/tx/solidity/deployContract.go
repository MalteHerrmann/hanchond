package solidity

import (
	"encoding/hex"
	"fmt"

	"github.com/hanchon/hanchond/lib/smartcontract"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// deployContractCmd represents the deploy command
var deployContractCmd = &cobra.Command{
	Use:     "deploy-contract [path_to_bin_file]",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"d"},
	Short:   "Deploy a solidity contract",
	Long:    "The bytecode file must contain just the hex data",
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(fmt.Errorf("node not set"))
		}

		gasLimit, err := cmd.Flags().GetInt("gas-limit")
		if err != nil {
			utils.ExitError(fmt.Errorf("incorrect gas limit: %w", err))
		}

		pathToBytecode := args[0]

		e := evmos.NewEvmosFromDB(queries, nodeID)
		builder := e.NewTxBuilder(uint64(gasLimit))

		bytecode, err := filesmanager.ReadFile(pathToBytecode)
		if err != nil {
			utils.ExitError(fmt.Errorf("error reading the bytecode file: %w", err))
		}

		bytecode, err = hex.DecodeString(string(bytecode))
		if err != nil {
			utils.ExitError(fmt.Errorf("error converting bytecode to []byte: %w", err))
		}

		abiPath, err := cmd.Flags().GetString("abi")
		if err != nil {
			utils.ExitError(fmt.Errorf("could not read abi path: %w", err))
		}

		if abiPath != "" {
			// It requires a constructor
			abiBytes, err := filesmanager.ReadFile(abiPath)
			if err != nil {
				utils.ExitError(fmt.Errorf("error reading the abi file: %w", err))
			}
			// Get Params
			callArgs, err := smartcontract.StringsToABIArguments(params)
			if err != nil {
				utils.ExitError(fmt.Errorf("error converting arguments: %w", err))
			}

			callData, err := smartcontract.ABIPackRaw(abiBytes, "", callArgs...)
			if err != nil {
				utils.ExitError(fmt.Errorf("error converting arguments: %w", err))
			}
			bytecode = append(bytecode, callData...)
		}

		txHash, err := builder.DeployContract(0, bytecode, uint64(gasLimit))
		if err != nil {
			utils.ExitError(fmt.Errorf("error sending the transaction: %w", err))
		}

		contractAddress, err := e.NewRequester().GetContractAddress(txHash)
		if err != nil {
			utils.ExitError(fmt.Errorf("error getting the contract address: %w", err))
		}

		fmt.Printf("{\"contract_address\":\"%s\", \"tx_hash\":\"%s\"}\n", contractAddress, txHash)
		utils.ExitSuccess()
	},
}

func init() {
	SolidityCmd.AddCommand(deployContractCmd)
	deployContractCmd.Flags().Int("gas-limit", 2_000_000, "GasLimit to be used to deploy the transaction")
	deployContractCmd.Flags().String("abi", "", "ABI file if the contract has a contronstructor that needs params")
	deployContractCmd.Flags().StringSliceVarP(&params, "params", "p", []string{}, "A list of params. If the param is an address, prefix with `a:0x123...`")
}
