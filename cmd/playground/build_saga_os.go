package playground

import (
	"errors"

	"github.com/hanchon/hanchond/playground/sagaos"
	"github.com/spf13/cobra"
)

// buildSagaOSCmd represents the buildSagaOSCmd command.
var buildSagaOSCmd = &cobra.Command{
	Use:   "build-sagaos",
	Short: "Build an specific version of sagaosd (hanchond playground build-sagaos v0.8.0), it also supports local repositories (hanchond playground build-sagaos --path $HOME/sagaxyz/sagaos)", //nolint:lll
	Long:  `It downloads, builds and clean up temp files for any SagaOS tag. Using the --path flag will build you local repo`,                                                                   //nolint:lll
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("version is missing. Usage: hanchond playground build-sagaosd v0.8.0")
		}

		version := args[0]
		chainInfo := sagaos.ChainInfo

		return RunBuildEVMChainCmd(cmd, chainInfo, version)
	},
}

func init() {
	PlaygroundCmd.AddCommand(buildSagaOSCmd)
	buildSagaOSCmd.Flags().StringP("path", "p", "", "Path to you local clone of Evmos")
}
