package filesmanager

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/types"
)

// BuildEVMChainVersion builds the downloaded binary of a Cosmos EVM based chain.
//
// NOTE: This requires that the version was already cloned
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

	out, err := utils.ExecCommand("make", "build", "COSMOS_BUILD_OPTIONS=nooptimization,nostrip")
	if err != nil {
		return err
	}

	// TODO: is build output required here? Could be removed just as well...
	utils.Log("build output: %s", out)
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
	return os.Rename(origin, destination)
}

func CopyFile(origin string, destination string) error {
	cmd := exec.Command("cp", origin, destination)
	if out, err := cmd.CombinedOutput(); err != nil {
		err = fmt.Errorf("error %s: %s", err.Error(), string(out))
		return err
	}
	return nil
}

// NOTE: This requires that the version was already cloned
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
	return os.Rename(GetBranchFolder(version)+"/target/debug/hermes", GetHermesBinary())
}
