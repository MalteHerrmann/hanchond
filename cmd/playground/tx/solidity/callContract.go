package solidity

import (
	"fmt"
	"strings"

	"github.com/hanchon/hanchond/lib/requester"
	"github.com/hanchon/hanchond/lib/smartcontract"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/cosmosdaemon"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

var params []string

// callContractViewCmd represents the callContractView command.
var callContractViewCmd = &cobra.Command{
	Use:   "call-contract-view [contract] [abi_path] [method]",
	Args:  cobra.ExactArgs(3),
	Short: "Call a contract view with eth_call",
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)

		height, err := cmd.Flags().GetString("height")
		if err != nil {
			utils.ExitError(fmt.Errorf("could not read height value: %w", err))
		}

		contract := strings.TrimSpace(args[0])
		abiPath := strings.TrimSpace(args[1])
		method := strings.TrimSpace(args[2])

		abiBytes, err := filesmanager.ReadFile(abiPath)
		if err != nil {
			utils.ExitError(fmt.Errorf("error reading the abi file: %w", err))
		}

		endpoint, err := cosmosdaemon.GetWeb3Endpoint(queries, cmd)
		if err != nil {
			utils.ExitError(fmt.Errorf("error generting web3 endpoint: %w", err))
		}

		callArgs, err := smartcontract.StringsToABIArguments(params)
		if err != nil {
			utils.ExitError(fmt.Errorf("error converting arguments: %w", err))
		}

		client := requester.NewClient().WithUnsecureWeb3Endpoint(endpoint)

		callData, err := smartcontract.ABIPack(abiBytes, method, callArgs...)
		if err != nil {
			utils.ExitError(fmt.Errorf("error converting arguments: %w", err))
		}

		resp, err := client.EthCall(contract, callData, height)
		if err != nil {
			utils.ExitError(fmt.Errorf("error on eth call: %w", err))
		}
		fmt.Println(string(resp))
		utils.ExitSuccess()
	},
}

func init() {
	SolidityCmd.AddCommand(callContractViewCmd)
	callContractViewCmd.Flags().String("height", "latest", "Query at the given height.")
	callContractViewCmd.Flags().StringSliceVarP(&params, "params", "p", []string{}, "A list of params. If the param is an address, prefix with `a:0x123...`") //nolint:lll
}
