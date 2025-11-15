package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitClient struct {
	repoPath string
	repo     *git.Repository
	repoDir  string
}

// creates a new git client
func NewGitClient(workdir string) (*GitClient, error) {
	repoPath := filepath.Join(workdir, "repos")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create repo directory: %w", err)
	}

	return &GitClient{
		repoPath: repoPath,
	}, nil
}

// clones a repo
func (c *GitClient) Clone(repoURL, branch string) error {
	repoDir := filepath.Join(c.repoPath, sanitizeRepoName(repoURL))
	c.repoDir = repoDir

	if _, err := os.Stat(repoDir); err == nil {
		// Repository exists, try to pull
		repo, err := git.PlainOpen(repoDir)
		if err != nil {
			return fmt.Errorf("repository exists but failed to open: %w", err)
		}
		c.repo = repo

		// if conflict error, delete cache and re-clone
		if err := c.Pull(branch); err != nil {
			_ = os.RemoveAll(repoDir)
			return c.cloneFresh(repoURL, branch, repoDir)
		}

		return nil
	}

	return c.cloneFresh(repoURL, branch, repoDir)
}

// cloneFresh performs a fresh clone
func (c *GitClient) cloneFresh(repoURL, branch, repoDir string) error {
	repo, err := git.PlainClone(repoDir, false, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: getReferenceName(branch),
		SingleBranch:  true,
	})
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	c.repo = repo
	return nil
}

func (c *GitClient) Pull(branch string) error {
	if c.repo == nil {
		return fmt.Errorf("no repository opened")
	}
	// Get worktree
	worktree, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	err = worktree.Pull(&git.PullOptions{
		ReferenceName: getReferenceName(branch),
		SingleBranch:  true,
	})

	// git.NoErrAlreadyUpToDate means we're already up to date
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to pull: %w", err)
	}

	return nil
}

// returns all manifest files
func (c *GitClient) GetManifest(path string) ([]string, error) {
	if c.repo == nil {
		return nil, fmt.Errorf("no repository cloned")
	}

	worktree, err := c.repo.Worktree()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(worktree.Filesystem.Root(), path)

	var files []string
	err = filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// convert repo url to a dir name
func sanitizeRepoName(repoURL string) string {
	name := strings.TrimSuffix(repoURL, ".git")
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, "/", "-")
	return name
}

// converts branch  to git ref like
func getReferenceName(branch string) plumbing.ReferenceName {
	return plumbing.NewBranchReferenceName(branch)
}
