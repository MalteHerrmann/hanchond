package playground

import (
	"context"
	dbsql "database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/config"
	"github.com/hanchon/hanchond/playground/cosmosdaemon"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// initChainCmd represents the initChain command
var initChainCmd = &cobra.Command{
	Use:   "init-chain [amount_of_validators]",
	Args:  cobra.ExactArgs(1),
	Short: "Init the genesis and configurations files for a new chain",
	Long:  `Set up the validators nodes for the new chain.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := cmd.Flags().GetString("client")
		if err != nil {
			utils.ExitError(fmt.Errorf("client flag was not set"))
		}
		version, err := cmd.Flags().GetString("version")
		if err != nil {
			utils.ExitError(fmt.Errorf("version flag was not set"))
		}
		chainID, err := cmd.Flags().GetString("chainid")
		if err != nil {
			utils.ExitError(fmt.Errorf("chainid flag was not set"))
		}

		amountOfValidators, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			utils.ExitError(fmt.Errorf("invalid amount of validators"))
		}

		queries := sql.InitDBFromCmd(cmd)

		latestChain, err := queries.GetLatestChain(context.Background())
		chainNumber := 1
		if err == nil {
			chainNumber = int(latestChain.ID) + 1
		} else if !errors.Is(err, dbsql.ErrNoRows) { // NOTE: no rows can be expected for an empty db
			utils.ExitError(fmt.Errorf("could not get the chains info from db: %w", err))
		}

		chainConfig, err := config.GetChainConfig(strings.ToLower(client))
		if err != nil {
			utils.ExitError(fmt.Errorf("error getting chain config: %w", err))
		}

		chainInfo := chainConfig.ToChainInfo()
		nodes := make([]*cosmosdaemon.Daemon, amountOfValidators)

		if chainID == "" {
			chainID = fmt.Sprintf("%s%d", chainInfo.GetChainIDBase(), chainNumber)
		}

		for k := range nodes {
			if filesmanager.IsNodeHomeFolderInitialized(int64(chainNumber), int64(k)) {
				utils.ExitError(fmt.Errorf("the home folder already exists: %d-%d", chainNumber, k))
			}

			path := filesmanager.GetNodeHomeFolder(int64(chainNumber), int64(k))
			nodes[k] = cosmosdaemon.NewDameon(
				chainInfo,
				fmt.Sprintf("moniker-%d-%d", chainNumber, k),
				version,
				path,
				chainID,
				fmt.Sprintf("validator-key-%d-%d", chainNumber, k),
			)
		}

		dbID, err := cosmosdaemon.InitMultiNodeChain(nodes, queries)
		if err != nil {
			utils.ExitError(fmt.Errorf("error: %w", err))
		}

		utils.Log("New chain created with id: %d", dbID)
	},
}

func init() {
	PlaygroundCmd.AddCommand(initChainCmd)
	initChainCmd.Flags().String("client", "evmos", "Client that you want to use. Options: evmos, sagaos, gaia")
	initChainCmd.Flags().StringP("version", "v", "local", "Version of the Evmos node that you want to use, defaults to local. Tag names are supported. If selected node is gaia, the flag is ignored.")
	initChainCmd.Flags().StringP("chainid", "c", "", "Chain-ID to be used when creating the genesis file, it defaults to `evmos_9001-X` or `cosmoshub-X`, depending on the client.")
}
