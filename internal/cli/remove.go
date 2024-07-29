package cli

import (
	"github.com/spf13/cobra"
	"github.com/xarunoba/ccoco/pkg/ccoco"
)

func init() {
	cli.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Remove files to config",
	Long: `Remove files to config.
This will remove file/s from the config file for ccoco to generate.`,
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
		if err := app.RemoveFromFiles(args); err != nil {
			return err
		}
		return nil
	},
}
