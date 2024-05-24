package cli

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/xarunoba/ccoco/internal/config"
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
		removeFiles(args)
	},
}

func removeFiles(files []string) {
	configFile := config.GetFile()

	for _, file := range files {
		for i, f := range configFile.Files {
			if f == file {
				configFile.Files = append(configFile.Files[:i], configFile.Files[i+1:]...)
			}
		}
	}

	configFile.Files = dedupeSlice(configFile.Files)

	configData, err := json.MarshalIndent(configFile, "", "  ")
	if err != nil {
		log.Printf("Error marshalling default config: %v", err)
	}
	if err := os.WriteFile(config.FileName, configData, 0644); err != nil {
		log.Printf("Error creating file %s: %v", config.FileName, err)
	}
}
