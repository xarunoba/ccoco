package cli

import (
	"github.com/spf13/cobra"
	"github.com/xarunoba/ccoco/pkg/ccoco"
)

func init() {
	cli.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&addToGitIgnore, "gitignore", "i", true, "Add to .gitignore")
	initCmd.Flags().BoolVarP(&injectCcocoToGitHooks, "githook", "g", false, "Inject ccoco to .git/hooks/post-checkout")
	initCmd.Flags().BoolVarP(&skipGitHookExecute, "skip", "s", true, "Skip git hook execution when used with --githook")
}

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Initialize ccoco",
	Long:    `Initialize ccoco in the current git repository.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		instance, err := ccoco.New()
		if err != nil {
			return err
		}
		app = instance
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := app.Init(ccoco.InitOptions{
			AddToGitIgnore: addToGitIgnore,
			AddToGitHooks:  injectCcocoToGitHooks,
			AddToGitHooksOptions: ccoco.AddToGitHooksOptions{
				SkipExecution: skipGitHookExecute,
			},
		}); err != nil {
			return err
		}
		return nil
	},
}
