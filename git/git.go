package git

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func CloneOrPull(repoUrl, branch, path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		//if repo does not exist clone it
		_, err := git.PlainClone(path, false, &git.CloneOptions{
			URL:           repoUrl,
			ReferenceName: plumbing.NewBranchReferenceName(branch), //plumbing convert our branches into git slang (like /ref/head/master)
			SingleBranch:  true,
			Depth:         1,
		})
		return err
	}

	//pull if cloned already
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	w, err := r.Worktree() //access of repo(commits , pr etc)
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	})
	if err == git.NoErrAlreadyUpToDate {
		return nil // Not an error, just skip applying
	}

	return err

}
