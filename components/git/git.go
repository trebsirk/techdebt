package git

import (
	"log"
	"os"

	"github.com/go-git/go-git"
	"github.com/go-git/go-git/plumbing/object"

	"techdebt/components/commitinfo"
)

func GetCommits(repoPath string) []commitinfo.CommitInfo {
	repoURL := "https://github.com/user/repo.git" // Replace with your repo URL

	// Clone repo if it doesn't exist locally
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
			URL:      repoURL,
			Progress: os.Stdout,
		})
		if err != nil {
			log.Fatalf("Failed to clone repository: %v", err)
		}
	}

	// Open the local repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatalf("Failed to open repository: %v", err)
	}

	// Get the commit history
	ref, err := repo.Head()
	if err != nil {
		log.Fatalf("Failed to get HEAD reference: %v", err)
	}

	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Fatalf("Failed to get commit history: %v", err)
	}

	var commits []commitinfo.CommitInfo

	// Iterate over each commit
	err = commitIter.ForEach(func(c *object.Commit) error {
		// Get the files modified in this commit
		files, err := c.Files()
		if err != nil {
			return err
		}

		// Add each file with commit metadata to commits slice
		files.ForEach(func(f *object.File) error {
			commitInfo := commitinfo.CommitInfo{
				Author:    c.Author.Name,
				Filename:  f.Name,
				Timestamp: c.Author.When,
			}
			commits = append(commits, commitInfo)
			return nil
		})

		return nil
	})
	if err != nil {
		log.Fatalf("Failed to iterate commits: %v", err)
	}
}
