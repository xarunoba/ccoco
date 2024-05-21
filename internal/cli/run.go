package cli

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/xarunoba/ccoco/internal/config"
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
		run()
	},
}

func run() {
	// Load config
	configData := config.GetFile()

	// Open repository
	repository, err := git.PlainOpen(".")
	if err != nil {
		log.Fatalf("Error opening repository: %v", err)
	}

	// Get current branch
	currentBranch, err := repository.Head()
	if err != nil {
		log.Fatalf("Error getting current branch: %v", err)
	}

	// Check if current branch is a sub-branch
	isSubBranch := false
	splitCurrentBranch := strings.Split(currentBranch.Name().Short(), "/")
	if len(splitCurrentBranch) > 1 {
		isSubBranch = true
		log.Printf("Current branch is a sub-branch: %s", currentBranch.Name().Short())
	}

	for _, file := range configData.Files {
		if isSubBranch {
			var subBranchPath string
			isSuccess := false

			// Iterate through sub-branches
			for i := len(splitCurrentBranch) - 1; i > 0; i-- {
				subBranchPath = strings.Join(splitCurrentBranch[:i], "/")
				path := filepath.Join(config.ConfigsDir, file, subBranchPath)

				// Check if current path is a directory
				info, err := os.Stat(path)
				if err != nil {
					log.Printf("Failed to stat current path: %v", err)
					continue
				}
				if info.IsDir() {
					log.Printf("Current path is a directory: %s", subBranchPath)
					continue
				}

				// Read data from current path
				data, err := os.ReadFile(path)
				if err != nil {
					log.Printf("Failed to read current path: %v", err)
					continue
				}

				// Write data to config file
				if err := os.WriteFile(filepath.Join(configData.FilesDir, file), data, 0755); err != nil {
					log.Printf("Failed to write current path: %v", err)
					continue
				}

				isSuccess = true
			}

			if !isSuccess {
				log.Printf("No config file found for current branch for %s: %s", file, currentBranch.Name().Short())
			}
		} else {
			path := filepath.Join(config.ConfigsDir, file, currentBranch.Name().Short())

			// Check if current path is a directory
			info, err := os.Stat(path)
			if err != nil {
				log.Printf("Failed to stat current path: %v", err)
				continue
			}
			if info.IsDir() {
				log.Printf("Current path is a directory, expected file: %s", currentBranch.Name().Short())
				continue
			}

			// Read data from current path
			data, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Failed to read current path: %v", err)
				continue
			}

			// Write data to config file
			if err := os.WriteFile(filepath.Join(configData.FilesDir, file), data, 0755); err != nil {
				log.Printf("Failed to write current path: %v", err)
				continue
			}

			log.Printf("Config file changed for branch %sw: %s", currentBranch.Name().Short(), file)
		}
	}
}
