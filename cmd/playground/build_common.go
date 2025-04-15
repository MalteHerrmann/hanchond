package playground

import (
	"fmt"
	"strings"

	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/sagaos"
	"github.com/hanchon/hanchond/playground/types"
	"github.com/spf13/cobra"
)

const LocalVersion = "local"

func BuildLocalEVMBinary(chainInfo types.ChainInfo, path string) error {
	version := LocalVersion
	path = strings.TrimRight(path, "/")

	fmt.Printf("Building %s...\n", chainInfo.GetBinaryName())
	if err := filesmanager.BuildEVMBinary(path); err != nil {
		return fmt.Errorf("error building %s: %w", chainInfo.GetBinaryName(), err)
	}

	fmt.Println("Moving built binary...")
	buildPath := fmt.Sprintf("%s/build/%s", path, chainInfo.GetBinaryName())
	if err := filesmanager.MoveFile(buildPath, filesmanager.GetDaemondPathWithVersion(chainInfo, version)); err != nil {
		fmt.Println("could not move the built binary:", err.Error())
		return err
	}

	return nil
}

func BuildEVMBinaryFromGitHub(chainInfo types.ChainInfo, version string) error {
	// Clone from github
	if err := filesmanager.CreateTempFolder(version); err != nil {
		return fmt.Errorf("could not create temp folder: %w\n", err)
	}

	fmt.Printf("Cloning %s version: %s\n", chainInfo.GetBinaryName(), version)
	if err := filesmanager.GitCloneGitHubBranch(chainInfo, version); err != nil {
		return fmt.Errorf("could not clone the %s version: %s\n", chainInfo.GetBinaryName(), err)
	}

	fmt.Printf("Building %s...\n", chainInfo.GetBinaryName())
	if err := filesmanager.BuildEVMChainVersion(version); err != nil {
		return fmt.Errorf("error building %s: %w\n", chainInfo.GetBinaryName(), err)
	}

	fmt.Println("Moving built binary...")
	if err := filesmanager.SaveBuiltVersion(chainInfo, version); err != nil {
		return fmt.Errorf("could not move the built binary: %w", err)
	}

	fmt.Println("Cleaning up...")
	if err := filesmanager.CleanUpTempFolder(); err != nil {
		return fmt.Errorf("could not remove temp folder: %w", err)
	}

	return nil
}

func RunBuildEVMChainCmd(cmd *cobra.Command, chainInfo types.ChainInfo, version string) error {
	_ = filesmanager.SetHomeFolderFromCobraFlags(cmd)

	// Create build folder if needed
	if err := filesmanager.CreateBuildsDir(); err != nil {
		return fmt.Errorf("could not create build folder: %w", err)
	}

	path, err := cmd.Flags().GetString("path")
	// Local build
	if err == nil && path != "" {
		if err = BuildLocalEVMBinary(sagaos.ChainInfo, path); err != nil {
			return fmt.Errorf("could not build %s binary: %w", chainInfo.GetBinaryName(), err)
		}

		return nil
	}

	if err = BuildEVMBinaryFromGitHub(sagaos.ChainInfo, version); err != nil {
		return fmt.Errorf("could not build %s binary: %w", chainInfo.GetBinaryName(), err)
	}

	return nil
}
