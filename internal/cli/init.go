package cli

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xarunoba/ccoco/internal/config"
)

func init() {
	cli.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&addToGitIgnore, "gitignore", "i", false, "Add to .gitignore")
	initCmd.Flags().BoolVarP(&injectCcocoToGitHooks, "githook", "g", false, "Inject ccoco to .git/hooks/post-checkout")
	initCmd.Flags().BoolVarP(&skipGitHookExecute, "skip", "s", false, "Skip git hook execution when used with --githook")
}

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Initialize config file",
	Long:    `Initialize config file.`,
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

		if addToGitIgnore {
			// Add .ccoco to .gitignore
			gitignoreFile := ".gitignore"
			gitignoreData := []byte("\n.ccoco\n")

			// Create .gitignore if it doesn't exist and write .ccoco to it
			if _, err := os.Stat(gitignoreFile); os.IsNotExist(err) {
				if err := os.WriteFile(gitignoreFile, gitignoreData, 0644); err != nil {
					log.Printf("Error creating file %s: %v", gitignoreFile, err)
				}
			} else if err == nil {
				gitignoreData, err := os.ReadFile(gitignoreFile)
				if err != nil {
					log.Printf("Error reading file %s: %v", gitignoreFile, err)
				}
				// Check if .ccoco is already in .gitignore
				if !strings.Contains(string(gitignoreData), ".ccoco") {
					// Append .ccoco to .gitignore
					gitignoreData = append(gitignoreData, gitignoreData...)
					if err := os.WriteFile(gitignoreFile, gitignoreData, 0644); err != nil {
						log.Printf("Error creating file %s: %v", gitignoreFile, err)
					}
				}
			}
		}

		// Ijnect ccoco to git hooks
		if injectCcocoToGitHooks {
			injectGitHook()
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
