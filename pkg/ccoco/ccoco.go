package ccoco

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/go-git/go-git/v5/plumbing"
)

const DefaultConfigFile = "ccoco.config.json"
const DefaultRootDirectory = ".ccoco"

const DefaultConfigDirectory = DefaultRootDirectory + "/configs"
const DefaultPreflightDirectory = DefaultRootDirectory + "/preflights"

type Ccoco struct {
	gitClient   *Git
	directories *Directories
	configFile  *File
}

func New() (*Ccoco, error) {
	gitClient, err := NewGitClient(".")
	if err != nil {
		return nil, err
	}
	directories := &Directories{
		Root:       DefaultRootDirectory,
		Configs:    DefaultConfigDirectory,
		Preflights: DefaultPreflightDirectory,
	}
	configFile := &File{
		Name: DefaultConfigFile,
		Content: &FileContent{
			Files: []string{".env"},
		},
	}

	instance, err := NewWithOptions(gitClient, directories, configFile)
	if err != nil {
		return nil, err
	}

	if instance.IsInitialized() {
		data, err := os.ReadFile(filepath.Join(gitClient.RootPathFromCwd, configFile.Name))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &configFile.Content)
		if err != nil {
			return nil, err
		}
		if err := instance.Load(LoadOptions{
			ConfigFile: configFile,
		}); err != nil {
			return nil, err
		}
	}

	return instance, nil
}

func NewWithOptions(gitClient *Git, directories *Directories, configFile *File) (*Ccoco, error) {
	ccoco := &Ccoco{
		gitClient:   gitClient,
		directories: directories,
		configFile:  configFile,
	}

	// Check state if valid
	if err := ccoco.CheckState(); err != nil {
		return nil, err
	}

	return ccoco, nil
}

func (c Ccoco) GitClient() *Git {
	return c.gitClient
}

func (c Ccoco) Directories() *Directories {
	return c.directories
}

func (c Ccoco) ConfigFile() *File {
	return c.configFile
}

func (c Ccoco) CheckState() error {
	if c.gitClient == nil {
		return errors.New("git client is nil")
	}
	if err := c.gitClient.CheckState(); err != nil {
		return err
	}

	if c.directories == nil {
		return errors.New("list of directories is nil")
	}
	if err := c.directories.CheckState(); err != nil {
		return err
	}

	if c.configFile == nil {
		return errors.New("config file is nil")
	}
	if err := c.configFile.CheckState(); err != nil {
		return err
	}
	return nil
}

type InitOptions struct {
	AddToGitIgnore bool
	AddToGitHooks  bool
	AddToGitHooksOptions
}

func (c Ccoco) Init(opts InitOptions) error {
	// Create directories
	if err := os.MkdirAll(filepath.Join(c.gitClient.RootPathFromCwd, c.directories.Root), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(c.gitClient.RootPathFromCwd, c.directories.Configs), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(c.gitClient.RootPathFromCwd, c.directories.Preflights), 0755); err != nil {
		return err
	}

	if opts.AddToGitIgnore {
		if err := c.AddToGitIgnore(); err != nil {
			return err
		}
	}

	if opts.AddToGitHooks {
		if err := c.AddToGitHooks(opts.AddToGitHooksOptions); err != nil {
			return err
		}
	}

	// Create config file if it doesn't exist with default values
	ccocoConfigFile := filepath.Join(c.gitClient.RootPathFromCwd, c.configFile.Name)
	if _, err := os.Stat(ccocoConfigFile); os.IsNotExist(err) {
		configData, err := json.MarshalIndent(&FileContent{
			Files: []string{".env"},
		}, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(ccocoConfigFile, configData, 0644); err != nil {
			return err
		}
	}

	log.Printf("Initialized ccoco-related files and directories at %s", c.gitClient.RootPathFromCwd)

	return nil
}

func (c Ccoco) IsInitialized() bool {
	if _, err := os.Stat(filepath.Join(c.gitClient.RootPathFromCwd, c.directories.Root)); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(filepath.Join(c.gitClient.RootPathFromCwd, c.directories.Configs)); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(filepath.Join(c.gitClient.RootPathFromCwd, c.directories.Preflights)); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(filepath.Join(c.gitClient.RootPathFromCwd, c.configFile.Name)); os.IsNotExist(err) {
		return false
	}

	return true
}

type LoadOptions struct {
	GitClient   *Git
	Directories *Directories
	ConfigFile  *File
}

func (c *Ccoco) Load(opts LoadOptions) error {
	if c == nil {
		return errors.New("ccoco is nil")
	}
	if opts.GitClient != nil {
		c.gitClient = opts.GitClient
	}
	if opts.Directories != nil {
		c.directories = opts.Directories
	}
	if opts.ConfigFile != nil {
		c.configFile = opts.ConfigFile
	}
	return nil
}

func (c Ccoco) AddToGitIgnore() error {

	// Add root directory to root .gitignore
	gitignoreFile := filepath.Join(c.gitClient.RootPathFromCwd, ".gitignore")
	gitignoreData := []byte("# ccoco directory\n" + c.directories.Root + "\n")

	// Create root .gitignore if it doesn't exist and write root directory in it
	if _, err := os.Stat(gitignoreFile); os.IsNotExist(err) {
		if err := os.WriteFile(gitignoreFile, gitignoreData, 0644); err != nil {
			return err
		}
	} else if err == nil {
		gitignoreFileData, err := os.ReadFile(gitignoreFile)
		if err != nil {
			log.Printf("Error reading file %s: %v", gitignoreFile, err)
		}
		// Check if root directory is already in .gitignore
		if !strings.Contains(string(gitignoreFileData), c.directories.Root) {
			// Append root directory to .gitignore
			gitignoreFileData = append(gitignoreFileData, []byte("\n# ccoco directory\n"+c.directories.Root+"\n")...)
			if err := os.WriteFile(gitignoreFile, gitignoreFileData, 0644); err != nil {
				return err
			}
		} else {
			log.Printf("Root .gitignore already contains %s", c.directories.Root)
		}
	}
	return nil
}

type AddToGitHooksOptions struct {
	SkipExecution bool
}

func (c Ccoco) AddToGitHooks(opts AddToGitHooksOptions) error {

	script := `#!/bin/sh
# Skip ccoco if SKIP_CCOCO is set to 1
if [ "$SKIP_CCOCO" = "1" ]; then
	echo "SKIP_CCOCO is set to 1, skipping ccoco."
	exit 0
fi
	
	# Run all preflight scripts
for file in ./` + c.directories.Preflights + `/*; do
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

	// Get the absolute path of the git worktree root
	absRootPath, err := filepath.Abs(c.gitClient.RootPathFromCwd)
	if err != nil {
		return err
	}
	// Get the relative path from the git worktree root to ccoco executable
	relativePath, err := filepath.Rel(absRootPath, os.Args[0])
	if err != nil {
		return err
	}
	// Convert Windows paths to Unix
	if runtime.GOOS == "windows" {
		relativePath = filepath.ToSlash(relativePath)
	}

	script += relativePath + " run"

	path := filepath.Join(c.gitClient.RootPathFromCwd, ".git/hooks/post-checkout")

	// Write the post-checkout hook script to the file
	if err := os.WriteFile(path, []byte(script), 0755); err != nil {
		return err
	}

	// Execute the post-checkout hook when SkipExecution is false
	if !opts.SkipExecution {
		executable := exec.Command("/bin/sh", path)
		executable.Stdout = os.Stdout
		executable.Stderr = os.Stderr
		err = executable.Run()
		if err != nil {
			return err
		}
	} else {
		log.Println("Skipped post-checkout hook execution")
	}

	log.Println("Post-checkout hook injected")

	return nil
}

type RunOptions struct {
	ForceToBranch *string
}

func (c Ccoco) Run(opts RunOptions) error {
	// Check if configs are initialized
	if !c.IsInitialized() {
		return errors.New("ccoco is not initialized properly. please reinitialize it")
	}
	// Get current branch from options
	currentBranch := ""
	if opts.ForceToBranch == nil || currentBranch == "" {
		// Get current branch from git
		currentBranchInfo, err := c.gitClient.Repository.Head()
		if err != nil {
			log.Fatalf("Error getting current branch: %v", err)
		}
		currentBranch = currentBranchInfo.Name().Short()
	} else {
		currentBranch = *opts.ForceToBranch
	}

	if strings.Contains(currentBranch, "/") {
		log.Printf("Current branch is a sub-branch: %s", currentBranch)

		// Split current branch path
		splitCurrentBranch := strings.Split(currentBranch, "/")

		// Initialize sub-branch path for looping
		subBranchPath := ""

		// Check if sub-branch exists
		isSuccess := false

		// Iterate through sub-branches from child to parent
		for i := len(splitCurrentBranch) - 1; i > 0; i-- {
			// Get current sub-branch path
			subBranchPath = strings.Join(splitCurrentBranch[:i], "/")
			path := filepath.Join(c.gitClient.RootPathFromCwd, c.directories.Configs, subBranchPath)

			// Check if path exists
			info, err := os.Stat(path)
			if err != nil {
				log.Printf("Failed to stat current path: %v", err)
				continue
			}

			// Check if path is a directory
			if !info.IsDir() {
				log.Printf("Current path is not a directory: %s", subBranchPath)
				continue
			}

			if err := c.ChangeConfigFiles(subBranchPath); err != nil {
				return err
			}

			isSuccess = true
		}
		if !isSuccess {
			return errors.New("failed to find any configs")
		}
	} else {
		path := filepath.Join(c.gitClient.RootPathFromCwd, c.directories.Configs, currentBranch)

		// Check if path exists
		info, err := os.Stat(path)
		if err != nil {
			return err
		}

		// Check if path is a directory
		if !info.IsDir() {
			return errors.New("current path is not a directory")
		}

		if err := c.ChangeConfigFiles(currentBranch); err != nil {
			return err
		}
	}
	return nil
}

func (c Ccoco) ChangeConfigFiles(currentBranch string) error {
	// Check if configs are initialized
	if !c.IsInitialized() {
		return errors.New("ccoco is not initialized properly. please reinitialize it")
	}

	for _, file := range c.configFile.Content.Files {
		// Encode file name to base58
		encodedFile := file
		if strings.Contains(file, "/") {
			encodedFile = filepath.Base(file) + "-" + base58.Encode([]byte(filepath.Dir(file)))
		}

		// Read data from current path
		data, err := os.ReadFile(filepath.Join(c.gitClient.RootPathFromCwd, c.directories.Configs, currentBranch, encodedFile))
		if err != nil {
			log.Printf("Failed to read current path: %v", err)
			continue
		}

		// Check if file is generated by ccoco
		if !strings.HasPrefix(string(data), "CCOCO GENERATED FILE - "+file+" - DO NOT REMOVE OR EDIT THIS LINE") {
			log.Printf("Malformed config file: %s", file)
			continue
		}

		// Remove first line from data
		if strings.Contains(string(data), "\n") {
			data = []byte(strings.Join(strings.Split(string(data), "\n")[1:], "\n"))
		} else {
			data = []byte("")
		}

		// Remove root path if it exists
		if err := os.RemoveAll(filepath.Join(c.gitClient.RootPathFromCwd, file)); err != nil {
			log.Printf("Failed to clear current path: %v", err)
			continue
		}

		// Write data to root path
		if err := os.WriteFile(filepath.Join(c.gitClient.RootPathFromCwd, file), data, 0644); err != nil {
			return err
		}
	}

	return nil
}

func (c Ccoco) GenerateConfigs() error {
	// Check if configs are initialized
	if !c.IsInitialized() {
		return errors.New("ccoco is not initialized properly. please reinitialize it")
	}

	headBranchInfo, err := c.gitClient.Repository.Head()
	if err != nil {
		log.Printf("Error getting current branch: %v", err)
		return err
	}
	headBranch := headBranchInfo.Name().Short()

	// Get all branches
	branches, err := c.gitClient.Repository.Branches()
	if err != nil {
		return err
	}

	// Generate per-branch config files
	if err := branches.ForEach(func(branch *plumbing.Reference) error {
		currentBranch := branch.Name().Short()

		// Create directory if it doesn't exist
		if err := os.MkdirAll(filepath.Join(c.gitClient.RootPathFromCwd, c.directories.Configs, branch.Name().Short()), 0755); err != nil {
			return err
		}
		for _, file := range c.configFile.Content.Files {
			// Add actual file name to the first line of the file
			data := []byte("CCOCO GENERATED FILE - " + file + " - DO NOT REMOVE OR EDIT THIS LINE\n")

			// Encode to base58 to flatten file
			encodedFile := file
			if strings.Contains(file, "/") {
				encodedFile = filepath.Base(file) + "-" + base58.Encode([]byte(filepath.Dir(file)))
			}

			if headBranch == currentBranch {
				// Read data from root path if it exists
				fileData, _ := os.ReadFile(filepath.Join(c.gitClient.RootPathFromCwd, file))
				if fileData != nil || len(fileData) > 0 {
					data = append(data, fileData...)
				}
			}

			if _, err := os.Stat(filepath.Join(c.gitClient.RootPathFromCwd, c.directories.Configs, currentBranch, encodedFile)); os.IsNotExist(err) {
				// Write data to current path
				if err := os.WriteFile(filepath.Join(c.gitClient.RootPathFromCwd, c.directories.Configs, currentBranch, encodedFile), data, 0644); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (c Ccoco) AddToFiles(files []string) error {
	// Check if configs are initialized
	if !c.IsInitialized() {
		return errors.New("ccoco is not initialized properly. please reinitialize it")
	}

	// Add files to config file
	filesMap := make(map[string]struct{})
	for _, file := range c.configFile.Content.Files {
		filesMap[filepath.ToSlash(file)] = struct{}{}
	}
	for _, f := range files {
		if _, exists := filesMap[f]; !exists {
			c.configFile.Content.Files = append(c.configFile.Content.Files, filepath.ToSlash(f))
		}
	}

	configData, err := json.MarshalIndent(c.configFile.Content, "", "  ")
	if err != nil {
		return err
	}

	if err := os.RemoveAll(filepath.Join(c.gitClient.RootPathFromCwd, c.configFile.Name)); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(c.gitClient.RootPathFromCwd, c.configFile.Name), configData, 0644); err != nil {
		return err
	}

	return nil
}

func (c Ccoco) RemoveFromFiles(files []string) error {
	// Check if configs are initialized
	if !c.IsInitialized() {
		return errors.New("ccoco is not initialized properly. please reinitialize it")
	}

	// Remove files from config file
	filesMap := make(map[string]struct{})
	for _, file := range files {
		filesMap[filepath.ToSlash(file)] = struct{}{}
	}
	newFiles := []string{}
	for _, f := range c.configFile.Content.Files {
		if _, exists := filesMap[f]; !exists {
			newFiles = append(newFiles, filepath.ToSlash(f))
		}
	}
	c.configFile.Content.Files = newFiles

	configData, err := json.MarshalIndent(c.configFile.Content, "", "  ")
	if err != nil {
		return err
	}

	if err := os.RemoveAll(filepath.Join(c.gitClient.RootPathFromCwd, c.configFile.Name)); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(c.gitClient.RootPathFromCwd, c.configFile.Name), configData, 0644); err != nil {
		return err
	}

	return nil
}
