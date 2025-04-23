package filesmanager

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/hanchon/hanchond/playground/types"
)

func GitCloneHermesBranch(version string) error {
	return GitCloneBranch(version, GetBranchFolder(version), "https://github.com/informalsystems/hermes")
}

func GitCloneGitHubBranch(chainInfo types.ChainInfo, version string) error {
	return GitCloneBranch(version, GetBranchFolder(version), chainInfo.GetRepoURL())
}

func GitCloneBranch(version string, dstFolder string, repoURL string) error {
	if !strings.HasPrefix(repoURL, "https://github.com/") {
		panic("repoURL must start with 'https://github.com/'; got: " + repoURL)
	}

	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", version, repoURL, dstFolder)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone repository; error: %w: %s", err, string(out))
	}

	return nil
}
