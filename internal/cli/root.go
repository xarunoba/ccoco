package cli

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var cli = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "Change config on checkout",
	Long: `ccoco changes your config files based on your current branch.
Integrate with git hooks to automatically change config on checkout.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			log.Fatalf("Error executing ccoco: %v", err)
		}
	},
}

func Execute() {
	if err := cli.Execute(); err != nil {
		log.Fatalf("Error executing ccoco: %v", err)
	}
}
