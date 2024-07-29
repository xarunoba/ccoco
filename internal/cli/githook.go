package cli

import (
	"github.com/spf13/cobra"
	"github.com/xarunoba/ccoco/pkg/ccoco"
)

func init() {
	cli.AddCommand(githookCmd)
	githookCmd.Flags().BoolVarP(&skipGitHookExecute, "skip", "s", false, "Skip git hook execution")
}

var githookCmd = &cobra.Command{
	Use:     "githook",
	Aliases: []string{"gh"},
	Short:   "Inject ccoco to git hooks",
	Long: `Injects ccoco to git hooks without depending on a git hook manager.
This will add a post-checkout hook to automatically change config on checkout.
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		instance, err := ccoco.New()
		if err != nil {
			return err
		}
		app = instance
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := app.AddToGitHooks(ccoco.AddToGitHooksOptions{
			SkipExecution: skipGitHookExecute,
		}); err != nil {
			return err
		}
		return nil
	},
}
