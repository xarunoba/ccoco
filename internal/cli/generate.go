package cli

import (
	"os"
	"path/filepath"

	"log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"

	"github.com/xarunoba/ccoco/internal/config"
)

func init() {
	cli.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate [FILE1 FILE2 ...]",
	Short: "Generate per-branch config files",
	Long: `Generates per-branch config files for a specific config file.
This will populate the branch configs folder based on the existing branches.
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		generate(args)
	},
}

func generate(args []string) {
	for _, configFileAsDir := range args {
		path := filepath.Join(config.ConfigsDir, configFileAsDir)

		// Create directory if it doesn't exist
		if _, err := os.Stat(path); err != nil {
			if err := os.MkdirAll(path, 0755); err != nil {
				log.Printf("Error creating directory %s: %v", path, err)
				continue
			}
		}

		// Open repository to check if it exists
		repository, err := git.PlainOpen(".")
		if err != nil {
			log.Fatalf("Error opening repository: %v", err)
		}

		// Get current branch
		currentBranch, err := repository.Head()
		if err != nil {
			log.Fatalf("Error getting current branch: %v", err)
		}

		// Get all branches
		branches, err := repository.Branches()
		if err != nil {
			log.Fatalf("Error getting branches: %v", err)
		}

		// Generate per-branch config files
		if err := branches.ForEach(func(branch *plumbing.Reference) error {
			configFile := filepath.Join(".", configFileAsDir)
			branchFile := filepath.Join(path, branch.Name().Short())

			// Check if branch file already exists and skip if it does
			if _, err := os.Stat(branchFile); err == nil {
				log.Printf("Config file %s already exists", branchFile)
				return nil
			}

			// If current branch, read the config file and save it to the branch file
			data := []byte("")
			if currentBranch.Name().Short() == branch.Name().Short() {
				data, _ = os.ReadFile(configFile)
			}

			// Create parent directory if it doesn't exist
			if err := os.MkdirAll(filepath.Dir(branchFile), 0755); err != nil {
				log.Printf("Error creating config file %s: %v", branchFile, err)
				return nil
			}

			// Write config file
			if err := os.WriteFile(branchFile, data, 0644); err != nil {
				log.Printf("Error creating config file %s: %v", branchFile, err)
				return nil
			}

			log.Printf("Created config file %s", branchFile)

			return nil
		}); err != nil {
			log.Fatalf("Error iterating over branches: %v", err)
		}
	}
}
