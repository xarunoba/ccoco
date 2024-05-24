package cli

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/xarunoba/ccoco/internal/config"
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
		addFiles(args)
	},
}

func addFiles(files []string) {
	configFile := config.GetFile()

	configFile.Files = append(configFile.Files, files...)

	configFile.Files = dedupeSlice(configFile.Files)

	configData, err := json.MarshalIndent(configFile, "", "  ")
	if err != nil {
		log.Printf("Error marshalling default config: %v", err)
	}
	if err := os.WriteFile(config.FileName, configData, 0644); err != nil {
		log.Printf("Error creating file %s: %v", config.FileName, err)
	}
}
