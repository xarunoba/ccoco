package cli

import (
	"log"

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
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.Run(ccoco.RunOptions{}); err != nil {
			log.Fatalf("Error running ccoco: %v", err)
		}
	},
}
