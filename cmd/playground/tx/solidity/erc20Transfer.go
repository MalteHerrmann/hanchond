package solidity

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hanchon/hanchond/lib/converter"
	"github.com/hanchon/hanchond/lib/txbuilder"
	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/cosmosdaemon"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/solidity"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// erc20TransferCmd represents the erc20Transfer command
var erc20TransferCmd = &cobra.Command{
	Use:   "erc20-transfer [contract] [wallet] [amount]",
	Args:  cobra.ExactArgs(3),
	Short: "Transfer erc20 coins from the validator wallet",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		queries := sql.InitDBFromCmd(cmd)

		contract := strings.TrimSpace(args[0])
		wallet := strings.TrimSpace(args[1])
		amount := strings.TrimSpace(args[2])

		wallet, err = converter.NormalizeAddressToHex(wallet)
		if err != nil {
			utils.ExitError(fmt.Errorf("invalid wallet: %w", err))
		}

		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(fmt.Errorf("node not set"))
		}

		endpoint, err := cosmosdaemon.GetWeb3Endpoint(queries, cmd)
		if err != nil {
			utils.ExitError(fmt.Errorf("error generting web3 endpoint: %w", err))
		}

		e := evmos.NewEvmosFromDB(queries, nodeID)
		valWallet := txbuilder.NewSimpleWeb3WalletFromMnemonic(e.ValMnemonic, endpoint)

		callData, err := solidity.ERC20TransferCallData(wallet, amount)
		if err != nil {
			utils.ExitError(fmt.Errorf("error building the call data: %w", err))
		}
		to := common.HexToAddress(contract)
		txhash, err := valWallet.TxBuilder.SendTx(valWallet.Address, &to, big.NewInt(0), 200_000, callData, valWallet.PrivKey)
		if err != nil {
			utils.ExitError(fmt.Errorf("error sending the transaction: %w", err))
		}

		fmt.Println("{\"txhash\":\"" + txhash + "\"}")
		utils.ExitSuccess()
	},
}

func init() {
	SolidityCmd.AddCommand(erc20TransferCmd)
}
