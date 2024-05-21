package cli

import (
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/xarunoba/ccoco/internal/config"
)

func init() {
	cli.AddCommand(githookCmd)
}

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
    echo "Skipping $file (not executable)"
  fi
done

# Run ccoco
./` + config.CcocoExecutable

	if runtime.GOOS == "windows" {
		script = script + `.exe`
	}
	script = script + ` run`

	return script
}

var githookCmd = &cobra.Command{
	Use:   "githook",
	Short: "Inject ccoco to git hooks",
	Long: `Injects ccoco to git hooks without depending on a git hook manager.
This will add a post-checkout hook to automatically change config on checkout.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		injectGitHook()
	},
}

func injectGitHook() {

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

	// Execute the post-checkout hook
	executable := exec.Command("bash", path)
	executable.Stdout = os.Stdout
	executable.Stderr = os.Stderr
	err = executable.Run()
	if err != nil {
		log.Fatalf("Error executing post-checkout hook: %v", err)
	}

	log.Println("Post-checkout hook injected")
}
