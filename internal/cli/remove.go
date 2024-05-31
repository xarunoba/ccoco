package cli

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	cli.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Remove config files",
	Long:    `Remove config files`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.RemoveFromFiles(args); err != nil {
			log.Fatalf("Error removing files: %v", err)
		}
	},
}
