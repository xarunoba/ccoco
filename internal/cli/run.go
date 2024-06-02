package cli

import (
	"github.com/spf13/cobra"
	"github.com/xarunoba/ccoco/pkg/ccoco"
)

func init() {
	cli.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"r", "start"},
	Short:   "Run ccoco",
	Long: `Run ccoco. 
This will change config files based on your current branch.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		instance, err := ccoco.New()
		if err != nil {
			return err
		}
		app = instance
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := app.Run(ccoco.RunOptions{
			ForceToBranch: nil,
		}); err != nil {
			return err
		}
		return nil
	},
}
