package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xarunoba/ccoco/internal/version"
)

func init() {
	cli.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print the version of ccoco",
	Long:    `Print the version of ccoco`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Version)
	},
}
