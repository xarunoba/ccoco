package ccoco

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

type Git struct {
	Repository      *git.Repository
	Worktree        *git.Worktree
	RootPathFromCwd string
}

func NewGitClient(path string) (*Git, error) {
	repository, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		log.Printf("Error opening repository: %v", err)
		return nil, err
	}

	// Get the worktree of the repository
	worktree, err := repository.Worktree()
	if err != nil {
		log.Printf("Error getting worktree: %v", err)
		return nil, err
	}

	// Check if worktree exists
	if worktree == nil {
		log.Println("No worktree found.")
		return nil, errors.New("Error getting worktree: No worktree found.")
	}

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %v", err)
		return nil, err
	}

	// Get the relative path from the current directory to the git worktree root
	rootPathFromCwd, err := filepath.Rel(cwd, worktree.Filesystem.Root())
	if err != nil {
		log.Printf("Error getting relative path: %v", err)
		return nil, err
	}

	return &Git{
		Repository:      repository,
		Worktree:        worktree,
		RootPathFromCwd: rootPathFromCwd}, nil
}

func (g *Git) CheckState() error {
	if g.Repository == nil {
		return errors.New("repository is nil")
	}
	if g.Worktree == nil {
		return errors.New("worktree is nil")
	}
	if g.RootPathFromCwd == "" {
		return errors.New("root path from current working directory is empty")
	}
	return nil
}
