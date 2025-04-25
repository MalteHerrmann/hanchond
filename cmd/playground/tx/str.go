package tx

import (
	"fmt"
	"strings"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// strV1ProposalCmd represents the str-v1-proposal command
var strV1ProposalCmd = &cobra.Command{
	Use:     "str-v1-proposal [denom]",
	Aliases: []string{"strv1"},
	Args:    cobra.ExactArgs(1),
	Short:   "Creates a STRv1 proposal",
	Run: func(cmd *cobra.Command, args []string) {
		queries := sql.InitDBFromCmd(cmd)
		nodeID, err := cmd.Flags().GetString("node")
		if err != nil {
			utils.ExitError(fmt.Errorf("node not set"))
		}

		denom := args[0]

		exponent, err := cmd.Flags().GetInt("exponent")
		if err != nil {
			utils.ExitError(fmt.Errorf("exponent not set"))
		}

		alias, err := cmd.Flags().GetString("alias")
		if err != nil {
			utils.ExitError(fmt.Errorf("alias not set"))
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			utils.ExitError(fmt.Errorf("name not set"))
		}

		symbol, err := cmd.Flags().GetString("symbol")
		if err != nil {
			utils.ExitError(fmt.Errorf("symbol not set"))
		}

		e := evmos.NewEvmosFromDB(queries, nodeID)
		out, err := e.CreateSTRv1Proposal(evmos.STRv1{
			Denom:    denom,
			Exponent: exponent,
			Alias:    alias,
			Name:     name,
			Symbol:   symbol,
		})
		if err != nil {
			utils.ExitError(fmt.Errorf("error sending the transaction: %w", err))
		}

		if !strings.Contains(out, "code: 0") {
			utils.ExitError(fmt.Errorf("transaction failed! %s", out))
		}

		hash := strings.Split(out, "txhash: ")
		if len(hash) > 1 {
			hash[1] = strings.TrimSpace(hash[1])
		}

		utils.Log("transaction sent: %s", hash[1])
	},
}

func init() {
	TxCmd.AddCommand(strV1ProposalCmd)
	strV1ProposalCmd.Flags().IntP("exponent", "e", 18, "Exponents of the token")
	strV1ProposalCmd.Flags().StringP("alias", "a", "tokenalias", "Token alias")
	strV1ProposalCmd.Flags().String("name", "tokenname", "Token name")
	strV1ProposalCmd.Flags().StringP("symbol", "s", "tokensymbol", "Token symbol")
}
