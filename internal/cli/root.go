package cli

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xarunoba/ccoco/internal/version"
	"github.com/xarunoba/ccoco/pkg/ccoco"
)

var app *ccoco.Ccoco

var cli = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "Change config on checkout",
	Long: `ccoco changes your config files based on your current branch.
Integrate with git hooks to automatically change config on checkout.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return nil
	},
	Version: version.Version,
}

func Execute() {
	if err := cli.Execute(); err != nil {
		log.Fatalf("Error executing ccoco: %v", err)
	}
}
