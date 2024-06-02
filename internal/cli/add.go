package cli

import (
	"github.com/spf13/cobra"
	"github.com/xarunoba/ccoco/pkg/ccoco"
)

func init() {
	cli.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:     "add [file1 file2 ...]",
	Aliases: []string{"a"},
	Short:   "Add file/s to config",
	Long: `Adds file/s to config.
This will add file/s to the config file for ccoco to generate.
`,
	Args: cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		instance, err := ccoco.New()
		if err != nil {
			return err
		}
		app = instance
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := app.AddToFiles(args); err != nil {
			return err
		}
		return nil
	},
}
