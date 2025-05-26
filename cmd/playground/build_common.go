package playground

import (
	"fmt"
	"strings"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/config"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/spf13/cobra"
)

const LocalVersion = "local"

// TODO: this should already take the chain config as an argument
func BuildLocalEVMBinary(chainName string, path string) error {
	version := LocalVersion
	path = strings.TrimRight(path, "/")

	chainConfig, err := config.GetChainConfig(chainName)
	if err != nil {
		return fmt.Errorf("error getting chain config: %w", err)
	}

	utils.Log("Building %s...", chainConfig.BinaryName)
	if err := filesmanager.BuildEVMBinary(path); err != nil {
		return fmt.Errorf("error building %s: %w", chainConfig.BinaryName, err)
	}

	utils.Log("Moving built binary...")
	buildPath := fmt.Sprintf("%s/%s", path, chainConfig.Build.BinaryPath)
	if err := filesmanager.MoveFile(buildPath, filesmanager.GetDaemondPathWithVersion(chainConfig.ToChainInfo(), version)); err != nil {
		utils.Log("could not move the built binary: %s", err)
		return err
	}

	return nil
}

// TODO: there is duplication between the chain config stuff and the chain info that's stored per node,
// this should be refactored and passed as input here
// The config should be the source of truth.
func BuildEVMBinaryFromGitHub(chainName string, version string) error {
	chainConfig, err := config.GetChainConfig(chainName)
	if err != nil {
		return fmt.Errorf("error getting chain config: %w", err)
	}

	// Clone from github
	if err := filesmanager.CreateTempFolder(version); err != nil {
		return fmt.Errorf("could not create temp folder: %w", err)
	}

	utils.Log("Cloning %s version: %s", chainConfig.BinaryName, version)
	if err := filesmanager.GitCloneGitHubBranch(chainConfig.ToChainInfo(), version); err != nil {
		return fmt.Errorf("could not clone the %s version: %s", chainConfig.BinaryName, err)
	}

	utils.Log("Building %s...", chainConfig.BinaryName)
	if err := filesmanager.BuildEVMChainVersion(version); err != nil {
		return fmt.Errorf("error building %s: %w", chainConfig.BinaryName, err)
	}

	utils.Log("Moving built binary...")
	if err := filesmanager.SaveBuiltVersion(chainConfig.ToChainInfo(), version); err != nil {
		return fmt.Errorf("could not move the built binary: %w", err)
	}

	utils.Log("Cleaning up...")
	if err := filesmanager.CleanUpTempFolder(); err != nil {
		return fmt.Errorf("could not remove temp folder: %w", err)
	}

	return nil
}

func RunBuildEVMChainCmd(cmd *cobra.Command, chainName string, version string) error {
	_ = filesmanager.SetHomeFolderFromCobraFlags(cmd)

	// Create build folder if needed
	if err := filesmanager.CreateBuildsDir(); err != nil {
		return fmt.Errorf("could not create build folder: %w", err)
	}

	path, err := cmd.Flags().GetString("path")
	// Local build
	if err == nil && path != "" {
		if err = BuildLocalEVMBinary(chainName, path); err != nil {
			return fmt.Errorf("could not build binary: %w", err)
		}
		return nil
	}

	if err = BuildEVMBinaryFromGitHub(chainName, version); err != nil {
		return fmt.Errorf("could not build binary: %w", err)
	}

	return nil
}
