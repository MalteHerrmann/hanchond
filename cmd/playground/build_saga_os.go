package playground

import (
	"fmt"
	"os"
	"strings"

	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/spf13/cobra"
)

// buildSagaOSCmd represents the buildSagaOSCmd command
var buildSagaOSCmd = &cobra.Command{
	Use:   "build-sagaos",
	Short: "Build an specific version of sagaosd (hanchond playground build-sagaos v0.8.0), it also supports local repositories (hanchond playground build-sagaos --path $HOME/sagaxyz/sagaos)",
	Long:  `It downloads, builds and clean up temp files for any SagaOS tag. Using the --path flag will build you local repo`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = filesmanager.SetHomeFolderFromCobraFlags(cmd)

		// Create build folder if needed
		if err := filesmanager.CreateBuildsDir(); err != nil {
			fmt.Println("could not create build folder:" + err.Error())
			os.Exit(1)
		}

		// TODO: refactor this to not have to manually replicate this cleanup logic etc. for every build
		// we should rather create an interface ChainBinaryInfo or smth. that serves unified Build() and Get... functions instead
		// and then force the implementation for every built binary.
		path, err := cmd.Flags().GetString("path")
		// Local build
		if err == nil && path != "" {
			version := LocalVersion
			path = strings.TrimRight(path, "/")

			fmt.Println("Building sagaos...")
			if err := filesmanager.BuildEvmos(path); err != nil {
				fmt.Println("error building sagaos:", err.Error())
				os.Exit(1)
			}
			fmt.Println("Moving built binary...")
			if err := filesmanager.MoveFile(path+"/build/sagaosd", filesmanager.GetSagaosdPath(version)); err != nil {
				fmt.Println("could not move the built binary:", err.Error())
				os.Exit(1)
			}
			os.Exit(0)
		}

		// Clone from github
		if len(args) == 0 {
			fmt.Println("version is missing. Usage: hanchond playground build-sagaosd v0.8.0")
			os.Exit(1)
		}
		version := args[0]
		if err := filesmanager.CreateTempFolder(version); err != nil {
			fmt.Println("could not create temp folder:" + err.Error())
			os.Exit(1)
		}
		fmt.Println("Cloning sagaos version:", version)
		if err := filesmanager.GitCloneSagaOSBranch(version); err != nil {
			fmt.Println("could not clone the sagaos version: ", err)
			os.Exit(1)
		}
		fmt.Println("Building sagaos...")
		if err := filesmanager.BuildEvmosVersion(version); err != nil {
			fmt.Println("error building sagaos:", err)
			os.Exit(1)
		}
		fmt.Println("Moving built binary...")
		if err := filesmanager.SaveSagaOSBuiltVersion(version); err != nil {
			fmt.Println("could not move the built binary:", err.Error())
			os.Exit(1)
		}
		fmt.Println("Cleaning up...")
		if err := filesmanager.CleanUpTempFolder(); err != nil {
			fmt.Println("could not remove temp folder", err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	PlaygroundCmd.AddCommand(buildSagaOSCmd)
	buildSagaOSCmd.Flags().StringP("path", "p", "", "Path to you local clone of Evmos")
}
