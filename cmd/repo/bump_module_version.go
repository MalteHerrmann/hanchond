package repo

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/spf13/cobra"
)

// BumpModuleVersionCmd represents the query command
var BumpModuleVersionCmd = &cobra.Command{
	Use:   "bump-module-version [path] [version]",
	Args:  cobra.ExactArgs(2),
	Short: "Bump the version of a go module, i.e., hanchond repo bump-version /tmp/repo v21",
	Run: func(_ *cobra.Command, args []string) {
		path := args[0]
		// Make sure that the path does not end with `/`
		path = strings.TrimSuffix(path, "/")

		version := args[1]
		// If the version arg is missing the `v` prefix, we add it
		if !strings.HasPrefix(version, "v") {
			version = fmt.Sprintf("v%s", version)
		}

		// Find the current version
		goModPath := fmt.Sprintf("%s/go.mod", path)
		utils.Log("using go.mod path as: %s", goModPath)
		goModFile, err := filesmanager.ReadFile(goModPath)
		if err != nil {
			utils.ExitError(fmt.Errorf("error reading the go.mod file: %w", err))
		}

		// Get the current version
		re := regexp.MustCompile(`(?m)^module\s+(\S+)$`)
		modules := re.FindAllStringSubmatch(string(goModFile), -1)
		if len(modules) == 0 {
			utils.ExitError(fmt.Errorf("the go.mod file does not define the module name"))
		}
		currentVersion := modules[0][1]
		utils.Log("the current version is: %s", currentVersion)

		// Create the new version with all the parts of the currentVersion but overwritting the last segment
		newVersion := ""
		parts := strings.Split(currentVersion, "/")
		for k, v := range parts {
			if k == len(parts)-1 {
				newVersion += version
				break
			}
			newVersion += v + "/"
		}

		utils.Log("the new version is: %s", newVersion)

		// Walk through the root directory recursively
		if err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				utils.ExitError(fmt.Errorf("error reading the directory %s: %w", path, err))
			}

			// Only process regular files
			if info.IsDir() {
				return nil
			}

			// Read the file
			content, err := filesmanager.ReadFile(path)
			if err != nil {
				utils.ExitError(fmt.Errorf("failed reading the file %s: %w", path, err))
			}

			fileContent := string(content)

			// Replace all occurrences
			re := regexp.MustCompile(regexp.QuoteMeta(currentVersion))
			updatedContent := re.ReplaceAllString(fileContent, newVersion)

			// Only write if the file was modified
			if updatedContent != fileContent {
				utils.Log("updating file: %s", path)
				err := filesmanager.SaveFileWithMode([]byte(updatedContent), path, info.Mode())
				if err != nil {
					utils.ExitError(fmt.Errorf("failed saving the file %s: %w", path, err))
				}
			}

			return nil
		}); err != nil {
			utils.ExitError(fmt.Errorf("error walking the directory: %w", err))
		}

		utils.ExitSuccess()
	},
}
