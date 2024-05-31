package cli

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xarunoba/ccoco/pkg/ccoco"
)

var app *ccoco.Ccoco

func init() {
	initialized, err := ccoco.IsInitialized(ccoco.DefaultRootDirectory, ccoco.DefaultConfigFile, ccoco.DefaultConfigDirectory, ccoco.DefaultPreflightDirectory)
	if err != nil {
		log.Fatalf("Error initializing ccoco: %v", err)
	}
	instance, err := ccoco.New()
	if err != nil {
		log.Fatalf("Error initializing ccoco: %v", err)
	}

	if *initialized {
		if err := instance.Load(nil, nil, &ccoco.File{
			Name: ccoco.DefaultConfigFile,
			Content: &ccoco.FileContent{
				Files: []string{".env"},
			},
		}); err != nil {
			log.Fatalf("Error initializing ccoco: %v", err)
		}
	}

	app = instance

}

var cli = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "Change config on checkout",
	Long: `ccoco changes your config files based on your current branch.
Integrate with git hooks to automatically change config on checkout.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			log.Fatalf("Error executing ccoco: %v", err)
		}
	},
}

func Execute() {
	if err := cli.Execute(); err != nil {
		log.Fatalf("Error executing ccoco: %v", err)
	}
}
