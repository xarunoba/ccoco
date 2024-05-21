package cli

import (
	"encoding/json"
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
		// Create directories
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

		// Create preflight script
		preflightFile := config.PreflightsDir + "/preflight"
		preflightScript := `#!/bin/sh
echo "Running preflight script"`
		if err := os.WriteFile(preflightFile, []byte(preflightScript), 0755); err != nil {
			log.Printf("Error creating file %s: %v", preflightFile, err)
		}

		// Create config file if it doesn't exist with default values
		if _, err := os.Stat(config.FileName); os.IsNotExist(err) {
			configData, err := json.Marshal(config.DefaultFile)
			if err != nil {
				log.Printf("Error marshalling default config: %v", err)
			}
			if err := os.WriteFile(config.FileName, configData, 0644); err != nil {
				log.Printf("Error creating file %s: %v", config.FileName, err)
			}
		}

		log.Println("Initialized ccoco-related files and directories")
	},
}
