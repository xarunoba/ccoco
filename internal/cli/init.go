package cli

import (
	"log"

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
	Short:   "Initialize config file",
	Long:    `Initialize config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.Init(ccoco.InitOptions{
			AddToGitIgnore: addToGitIgnore,
			AddToGitHooks:  injectCcocoToGitHooks,
			AddToGitHooksOptions: ccoco.AddToGitHooksOptions{
				SkipExecution: skipGitHookExecute,
			},
		}); err != nil {
			log.Fatalf("Error initializing ccoco: %v", err)
		}
	},
}
