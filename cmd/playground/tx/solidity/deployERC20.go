package solidity

import (
	"fmt"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/solidity"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// deployERC20Cmd represents the deploy command
var deployERC20Cmd = &cobra.Command{
	Use:   "deploy-erc20 [name] [symbol]",
	Args:  cobra.ExactArgs(2),
	Short: "Deploy an erc20 contract",
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(fmt.Errorf("node not set"))
		}

		gasLimit, err := cmd.Flags().GetUint64("gas-limit")
		if err != nil {
			utils.ExitError(fmt.Errorf("incorrect gas limit: %w", err))
		}

		initialAmount, err := cmd.Flags().GetString("initial-amount")
		if err != nil {
			utils.ExitError(fmt.Errorf("incorrect initial-amount: %w", err))
		}

		isWrapped, err := cmd.Flags().GetBool("is-wrapped-coin")
		if err != nil {
			utils.ExitError(fmt.Errorf("incorrect wrapped flag: %w", err))
		}

		name := args[0]
		symbol := args[1]

		// TODO: allow mainnet as a valid endpoint
		e := evmos.NewEvmosFromDB(queries, nodeID)
		builder := e.NewTxBuilder(gasLimit)

		txHash, err := solidity.BuildAndDeployERC20Contract(name, symbol, initialAmount, isWrapped, builder, gasLimit)
		if err != nil {
			utils.ExitError(fmt.Errorf("error building and deploying the erc20 contract: %w", err))
		}

		contractAddress, err := e.NewRequester().GetContractAddress(txHash)
		if err != nil {
			utils.ExitError(fmt.Errorf("error getting the contract address: %w", err))
		}

		fmt.Printf("{\"contract_address\":\"%s\", \"tx_hash\":\"%s\"}\n", contractAddress, txHash)

		// Clean up files
		if err := filesmanager.CleanUpTempFolder(); err != nil {
			utils.ExitError(fmt.Errorf("could not clean up the temp folder: %w", err))
		}
		utils.ExitSuccess()
	},
}

func init() {
	SolidityCmd.AddCommand(deployERC20Cmd)
	deployERC20Cmd.Flags().Uint64("gas-limit", 2_000_000, "GasLimit to be used to deploy the transaction")
	deployERC20Cmd.Flags().String("initial-amount", "1000000", "Initial amout of coins sent to the deployer address")
	deployERC20Cmd.Flags().Bool("is-wrapped-coin", false, "Flag used to indenfity if the contract is representing the base denom. It uses WETH9 instead of OpenZeppelin contracts")
}
