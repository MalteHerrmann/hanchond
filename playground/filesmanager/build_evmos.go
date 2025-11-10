package filesmanager

import (
	"fmt"
	"os"
	"strings"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/types"
)

// BuildEVMChainVersion builds the downloaded binary of a Cosmos EVM based chain.
//
// NOTE: This requires that the version was already cloned.
func BuildEVMChainVersion(version string) error {
	return BuildEVMBinary(GetBranchFolder(version))
}

func BuildEVMBinary(path string) error {
	if err := os.Chdir(path); err != nil {
		return err
	}

	if _, err := utils.ExecCommand("rm", "-rf", path+"/build"); err != nil {
		return err
	}

	_, err := utils.ExecCommand("make", "build", "COSMOS_BUILD_OPTIONS=nooptimization,nostrip")
	if err != nil {
		return err
	}

	// TODO: is build output required here? Could be removed just as well...
	// utils.Log("build output: %s", out)
	return nil
}

func SaveBuiltVersion(chainInfo types.ChainInfo, version string) error {
	_ = CreateBuildsDir() // Make sure that the build dir exists

	return MoveFile(
		fmt.Sprintf("%s/build/%s", GetBranchFolder(version), chainInfo.GetBinaryName()),
		GetDaemondPathWithVersion(chainInfo, version),
	)
}

func MoveFile(origin string, destination string) error {
	if _, err := os.Stat(origin); os.IsNotExist(err) {
		return fmt.Errorf("file not found at %s", origin)
	}

	// check if folder of destination exists, assuming destination is a file
	parentDir := destination[:strings.LastIndex(destination, "/")]
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		return fmt.Errorf("destination parent folder not found at %s", parentDir)
	}

	return os.Rename(origin, destination)
}

func CopyFile(origin string, destination string) error {
	_, err := utils.ExecCommand("cp", origin, destination)

	return err
}

// NOTE: This requires that the version was already cloned
//
// TODO: avoid building from source and rather download from release page instead; depending on operating system get the correct binary.
func BuildHermes(version string) error {
	// Change directory to the cloned repository
	if err := os.Chdir(GetBranchFolder(version)); err != nil {
		return err
	}

	cmd := "CARGO_NET_GIT_FETCH_WITH_CLI=true cargo build"
	_, err := utils.ExecCommand("bash", "-c", cmd)

	return err
}

func SaveHermesBuiltVersion(version string) error {
	buildTarget := GetBranchFolder(version) + "/target/debug/hermes"

	_ = CreateBuildsDir() // Make sure that the build dir exists

	if err := MoveFile(buildTarget, GetHermesBinary()); err != nil {
		return fmt.Errorf("error moving hermes binary: %w", err)
	}

	return nil
}
