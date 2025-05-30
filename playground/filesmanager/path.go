package filesmanager

import (
	"fmt"
	"os"
	"strings"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/types"
	"github.com/spf13/cobra"
)

var baseDir = "/tmp"

func SetBaseDir(path string) {
	baseDir = path
}

func SetHomeFolderFromCobraFlags(cmd *cobra.Command) string {
	home, err := cmd.Flags().GetString("home")
	if err != nil {
		utils.ExitError(err)
	}
	home, _ = strings.CutSuffix(home, "/")
	SetBaseDir(home)
	// Ensure that the folder exists
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.Mkdir(home, os.ModePerm); err != nil {
			// We panic here because if we can not create the folder we should inmediately stop
			panic(err)
		}
	}
	return home
}

func GetDatabaseFile() string {
	return fmt.Sprintf("%s/playground.db", GetBaseDir())
}

func GetDataFolder() string {
	return fmt.Sprintf("%s/data", GetBaseDir())
}

func getNodeHomePath(chainID int64, nodeID int64) string {
	return fmt.Sprintf("%s/%d-%d", GetDataFolder(), chainID, nodeID)
}

func GetNodeHomeFolder(chainID, nodeID int64) string {
	if _, err := os.Stat(GetDataFolder()); os.IsNotExist(err) {
		if err := os.Mkdir(GetDataFolder(), os.ModePerm); err != nil {
			// We panic here because if we can not create the folder we should inmediately stop
			panic(err)
		}
	}
	return getNodeHomePath(chainID, nodeID)
}

func IsNodeHomeFolderInitialized(chainID int64, nodeID int64) bool {
	return DoesFileExist(getNodeHomePath(chainID, nodeID))
}

func GetBaseDir() string {
	return baseDir
}

func GetBuildsDir() string {
	return baseDir + "/builds"
}

func GetDepsDir(name string) string {
	return baseDir + "/builds/deps/" + name
}

func GetTempDir() string {
	return baseDir + "/temp"
}

func GetBranchFolder(version string) string {
	return GetTempDir() + "/" + version
}

func GetDaemondPathWithVersion(ci types.ChainInfo, version string) string {
	return fmt.Sprintf("%s/%s", GetBuildsDir(), ci.GetVersionedBinaryName(version))
}

// TODO: also remove and turn ChainInfo into BinaryInfo to also support hermes here
func GetHermesBinary() string {
	return GetBuildsDir() + "/hermes"
}

func GetHermesPath() string {
	return GetDataFolder() + "/hermes"
}

func CreateBuildsDir() error {
	if _, err := os.Stat(GetBuildsDir()); os.IsNotExist(err) {
		return os.Mkdir(GetBuildsDir(), os.ModePerm)
	}
	return nil
}

func CreateDepsFolder() error {
	return os.MkdirAll(baseDir+"/builds/deps/", os.ModePerm)
}

func CreateTempFolder(version string) error {
	return os.MkdirAll(GetBranchFolder(version), os.ModePerm)
}

func CreateHermesFolder() error {
	return os.MkdirAll(GetHermesPath(), os.ModePerm)
}

func CleanUpTempFolder() error {
	return os.RemoveAll(GetTempDir())
}

func CleanUpData() error {
	_ = os.RemoveAll(GetDatabaseFile())
	return os.RemoveAll(GetDataFolder())
}

func GetSolcPath(version string) string {
	return GetBuildsDir() + "/solc" + version
}
