package solidity

import (
	"fmt"
	"path/filepath"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/filesmanager"
	"github.com/hanchon/hanchond/playground/solidity"
	"github.com/hanchon/hanchond/playground/sql"
	"github.com/spf13/cobra"
)

// compileContractCmd represents the compile command.
var compileContractCmd = &cobra.Command{
	Use:     "compile-contract [path_to_solidity_file]",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"c"},
	Short:   "Compile a solidity contract",
	Run: func(cmd *cobra.Command, args []string) {
		_ = sql.InitDBFromCmd(cmd)

		outputFolder, err := cmd.Flags().GetString("output-folder")
		if err != nil {
			utils.ExitError(errors.New("incorrect output folder"))
		}
		if outputFolder[len(outputFolder)-1] != '/' {
			outputFolder += "/"
		}

		// TODO: read from pragma the correct version and use it automatically
		solcVersion, err := cmd.Flags().GetString("solc-version")
		if err != nil {
			utils.ExitError(errors.New("incorrect solc version"))
		}

		pathToSolidityCode := args[0]

		if err := filesmanager.CleanUpTempFolder(); err != nil {
			utils.ExitError(fmt.Errorf("could not clean up temp folder: %w", err))
		}

		folderName := "compiler"
		if err := filesmanager.CreateTempFolder(folderName); err != nil {
			utils.ExitError(fmt.Errorf("could not create up temp folder: %w", err))
		}

		err = solidity.CompileWithSolc(solcVersion, pathToSolidityCode, filesmanager.GetBranchFolder(folderName))
		if err != nil {
			utils.ExitError(fmt.Errorf("could not compile the contract: %w", err))
		}

		if err := moveFiles(filesmanager.GetBranchFolder(folderName), outputFolder, "abi"); err != nil {
			utils.ExitError(fmt.Errorf("error copying the built files: %w", err))
		}

		if err := moveFiles(filesmanager.GetBranchFolder(folderName), outputFolder, "bin"); err != nil {
			utils.ExitError(fmt.Errorf("error copying the built files: %w", err))
		}

		if err := filesmanager.CleanUpTempFolder(); err != nil {
			utils.ExitError(fmt.Errorf("could not clean up temp folder: %w", err))
		}

		utils.Log("Contract compiled at %s", outputFolder)
		utils.ExitSuccess()
	},
}

func init() {
	SolidityCmd.AddCommand(compileContractCmd)
	compileContractCmd.Flags().StringP("output-folder", "o", "./", "Output folder where the compile code will be saved")
	compileContractCmd.Flags().StringP("solc-version", "v", "0.8.0", "Solc version used to compile the code")
}

func moveFiles(in, out, extension string) error {
	files, err := filepath.Glob(in + "/*." + extension)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return err
	}

	for _, v := range files {
		if err := filesmanager.CopyFile(
			v,
			out,
		); err != nil {
			return err
		}
	}

	return nil
}
