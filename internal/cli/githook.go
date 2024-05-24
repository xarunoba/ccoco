package cli

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/xarunoba/ccoco/internal/config"
)

func init() {
	cli.AddCommand(githookCmd)
	githookCmd.Flags().BoolVarP(&skipGitHookExecute, "skip", "s", false, "Skip git hook execution")
}

// Generates the post-checkout script
func getPostCheckoutScript() string {
	script := `#!/bin/sh
# Skip ccoco if SKIP_CCOCO is set to 1
if [ "$SKIP_CCOCO" = "1" ]; then
	echo "SKIP_CCOCO is set to 1, skipping ccoco."
	exit 0
fi

# Run all preflight scripts
for file in ./` + config.PreflightsDir + `/*; do
  # Check if the file is executable
  if [ -x "$file" ]; then
    echo "Running $file"
    "$file"
  else
    echo "Cannot execute $file. Skipping."
  fi
done

# Run ccoco
`
	repository, err := git.PlainOpen(".")
	if err != nil {
		log.Fatalf("Error opening repository: %v", err)
	}
	worktree, err := repository.Worktree()
	if err != nil {
		log.Fatalf("Error getting worktree: %v", err)
	}

	relativePath, err := filepath.Rel(worktree.Filesystem.Root(), os.Args[0])
	if runtime.GOOS == "windows" {
		relativePath = filepath.ToSlash(relativePath)
	}
	if err != nil {
		log.Fatalf("Error getting relative path: %v", err)
	}
	script += relativePath + " run"

	return script
}

var githookCmd = &cobra.Command{
	Use:     "githook",
	Aliases: []string{"gh"},
	Short:   "Inject ccoco to git hooks",
	Long: `Injects ccoco to git hooks without depending on a git hook manager.
This will add a post-checkout hook to automatically change config on checkout.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		injectGitHook()
	},
}

// Injects ccoco to git hooks
func injectGitHook() {

	// Open repository to check if it exists
	_, err := git.PlainOpen(".")
	if err != nil {
		log.Fatalf("Error opening repository: %v", err)
	}

	script := getPostCheckoutScript()
	path := ".git/hooks/post-checkout"

	// Write the post-checkout hook script to the file
	if err := os.WriteFile(path, []byte(script), 0755); err != nil {
		log.Fatalf("Error writing post-checkout hook: %v", err)
	}

	// Execute the post-checkout hook when skipGitHookExecute is false
	if !skipGitHookExecute {
		executable := exec.Command("/bin/sh", path)
		executable.Stdout = os.Stdout
		executable.Stderr = os.Stderr
		err = executable.Run()
		if err != nil {
			log.Fatalf("Error executing post-checkout hook: %v", err)
		}
	} else {
		log.Println("Skipped post-checkout hook execution")
	}

	log.Println("Post-checkout hook injected")
}
