package filesmanager

import (
	"fmt"
	"os"
	"os/exec"
)

// NOTE: This requires that the version was already cloned
func BuildEvmosVersion(version string) error {
	return BuildEvmos(GetBranchFolder(version))
}

func BuildEvmos(path string) error {
	if err := os.Chdir(path); err != nil {
		return err
	}

	cmd := exec.Command("rm", "-rf", path+"/build") //nolint:gosec
	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}

	cmd = exec.Command("make", "build")
	_, err := cmd.CombinedOutput()
	return err
}

func SaveEvmosBuiltVersion(version string) error {
	// Ensure the path exists
	_ = CreateBuildsDir()
	return MoveFile(GetBranchFolder(version)+"/build/evmosd", GetEvmosdPath(version))
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
	command := exec.Command("bash", "-c", cmd)
	_, err := command.CombinedOutput()
	return err
}

func SaveHermesBuiltVersion(version string) error {
	return os.Rename(GetBranchFolder(version)+"/target/debug/hermes", GetHermesBinary())
}
