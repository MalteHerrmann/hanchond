package relayer

import (
	"fmt"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/hermes"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/hanchon/hanchond/playground/types"
	"github.com/spf13/cobra"
)

// represents the addChainConfigCmd command
//
// TODO: this is probably broken but not high priority atm to fix this.
var addChainConfigCmd = &cobra.Command{
	Use:   "add-chain-config",
	Args:  cobra.ExactArgs(0),
	Short: "Add chain config to hermes, it is ignored if the chain id already exists",
	Run: func(cmd *cobra.Command, _ []string) {
		_ = sql.InitDBFromCmd(cmd)

		h := hermes.NewHermes()
		fmt.Println("Relayer initialized")

		chainID, err := cmd.Flags().GetString("chainid")
		if err != nil || chainID == "" {
			utils.ExitError(fmt.Errorf("missing chainid value"))
		}

		p26657, err := cmd.Flags().GetString("p26657")
		if err != nil || chainID == "" {
			utils.ExitError(fmt.Errorf("missing p26657 value"))
		}

		p9090, err := cmd.Flags().GetString("p9090")
		if err != nil || chainID == "" {
			utils.ExitError(fmt.Errorf("missing p9090 value"))
		}

		keyname, err := cmd.Flags().GetString("keyname")
		if err != nil || chainID == "" {
			utils.ExitError(fmt.Errorf("missing keyname value"))
		}

		keymnemonic, err := cmd.Flags().GetString("keymnemonic")
		if err != nil || chainID == "" {
			utils.ExitError(fmt.Errorf("missing keymnemonic value"))
		}

		prefix, err := cmd.Flags().GetString("prefix")
		if err != nil || chainID == "" {
			utils.ExitError(fmt.Errorf("missing prefix value"))
		}

		denom, err := cmd.Flags().GetString("denom")
		if err != nil || chainID == "" {
			utils.ExitError(fmt.Errorf("missing denom value"))
		}

		isEvm, err := cmd.Flags().GetBool("is-evm")
		if err != nil || chainID == "" {
			utils.ExitError(fmt.Errorf("missing is-evm value"))
		}

		hdPath := types.CosmosHDPath
		if isEvm {
			// TODO: maybe check here if that's the right hd path to use? could e.g. get user confirmation
			hdPath = types.EthHDPath
		}

		defaultChainInfo := types.NewChainInfo(
			prefix,
			"external",
			chainID,
			"external",
			denom,
			"external",
			hdPath,
			types.CosmosAlgo,
			types.GaiaSDK,
		)

		switch isEvm {
		case false:
			fmt.Println("Adding a NOT EVM chain")
			if err := h.AddCosmosChain(
				defaultChainInfo,
				chainID,
				p26657,
				p9090,
				keyname,
				keymnemonic,
			); err != nil {
				utils.ExitError(fmt.Errorf("error adding first chain to the relayer: %w", err))
			}
		case true:
			chainInfo := defaultChainInfo
			chainInfo.KeyAlgo = types.EthAlgo
			chainInfo.SdkVersion = types.EvmosSDK

			fmt.Println("Adding a EVM chain")
			if err := h.AddEVMChain(
				chainInfo,
				chainID,
				p26657,
				p9090,
				keyname,
				keymnemonic,
			); err != nil {
				utils.ExitError(fmt.Errorf("error adding first chain to the relayer: %w", err))
			}
		}
	},
}

func init() {
	RelayerCmd.AddCommand(addChainConfigCmd)
	addChainConfigCmd.Flags().String("chainid", "", "Chain-id, i.e., evmos_9001-2")
	addChainConfigCmd.Flags().String("p26657", "", "Endpoint where the port 26657 is exposed, i.e., http://localhost:26657")
	addChainConfigCmd.Flags().String("p9090", "", "Endpoint where the port 9090 is exposed, i.e., http://localhost:9090")
	addChainConfigCmd.Flags().String("keyname", "", "Key name, it's used to identify the files inside hermes, i.e., relayerkey")
	addChainConfigCmd.Flags().String("keymnemonic", "", "Key mnemonic, mnemonic for the wallet")
	addChainConfigCmd.Flags().String("prefix", "", "Prefix for the bech32 address, i.e, osmo")
	addChainConfigCmd.Flags().String("denom", "", "Denom of the base token, i.e, aevmos")
	addChainConfigCmd.Flags().Bool("is-evm", false, "If the chain is evm compatible, this is used to determinate the type of wallet.")
}
