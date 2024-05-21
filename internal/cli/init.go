package cli

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/xarunoba/ccoco/internal/config"
)

func init() {
	cli.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize config file",
	Long:  `Initialize config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.MkdirAll(config.CcocoDir, 0755); err != nil {
			log.Printf("Error creating directory %s: %v", config.PreflightsDir, err)
		}
		if err := os.MkdirAll(config.CacheDir, 0755); err != nil {
			log.Printf("Error creating directory %s: %v", config.PreflightsDir, err)
		}
		if err := os.MkdirAll(config.ConfigsDir, 0755); err != nil {
			log.Printf("Error creating directory %s: %v", config.ConfigsDir, err)
		}
		if err := os.MkdirAll(config.PreflightsDir, 0755); err != nil {
			log.Printf("Error creating directory %s: %v", config.PreflightsDir, err)
		}

		preflightFile := config.PreflightsDir + "/preflight"
		preflightScript := `#!/bin/sh
echo "Running preflight script"`
		if err := os.WriteFile(preflightFile, []byte(preflightScript), 0755); err != nil {
			log.Printf("Error creating file %s: %v", preflightFile, err)
		}

		log.Println("Initialized ccoco-related files and directories")
	},
}
