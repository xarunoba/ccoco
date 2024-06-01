package cli

import (
	"github.com/spf13/cobra"
	"github.com/xarunoba/ccoco/pkg/ccoco"
)

func init() {
	cli.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate per-branch config files",
	Long: `Generates per-branch config files for the files specified in ccoco.config.json.
This will populate the branch configs folder based on the existing branches.
	`,
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		instance, err := ccoco.New()
		if err != nil {
			return err
		}
		app = instance
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := app.GenerateConfigs(); err != nil {
			return err
		}
		return nil
	},
}
