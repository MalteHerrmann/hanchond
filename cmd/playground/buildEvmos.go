package playground

import (
	"errors"

	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/spf13/cobra"
)

// buildEvmosCmd represents the buildEvmos command.
var buildEvmosCmd = &cobra.Command{
	Use:   "build-evmos",
	Short: "Build an specific version of Evmos (hanchond playground build-evmos v18.0.0), it also supports local repositories (hanchond playground build-evmos --path /home/hanchon/evmos)", //nolint:lll
	Long:  `It downloads, builds and clean up temp files for any Evmos tag. Using the --path flag will build you local repo`,                                                                //nolint:lll
	RunE: func(cmd *cobra.Command, args []string) error {
		// Clone from github
		if len(args) == 0 {
			return errors.New("version is missing. Usage: hanchond playground build-evmosd v18.1.0")
		}

		version := args[0]
		chainInfo := evmos.ChainInfo

		return RunBuildEVMChainCmd(cmd, chainInfo, version)
	},
}

func init() {
	PlaygroundCmd.AddCommand(buildEvmosCmd)
	buildEvmosCmd.Flags().StringP("path", "p", "", "Path to you local clone of Evmos")
}
