package cli

import (
	"log"

	"github.com/spf13/cobra"
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
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.GenerateConfigs(); err != nil {
			log.Fatalf("Error generating configs: %v", err)
		}
	},
}
