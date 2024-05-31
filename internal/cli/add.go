package cli

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	cli.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:     "add [file1 file2 ...]",
	Aliases: []string{"a"},
	Short:   "Add file to config",
	Long: `Adds file/s to config.
This will add file/s to the ccoco.config.json for ccoco to generate.
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.AddToFiles(args); err != nil {
			log.Fatalf("Error adding files: %v", err)
		}
	},
}
