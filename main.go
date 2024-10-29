package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type CommitInfo struct {
	Author    string
	Filename  string
	Timestamp time.Time
}

func main() {
	repoPath := "."                               //"/Users/davidwright/Documents/javascript-dev/trebsirk.github.io" // Define local repo directory
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

	var commits []CommitInfo

	// Iterate over each commit
	err = commitIter.ForEach(func(c *object.Commit) error {
		// Get the files modified in this commit
		files, err := c.Files()
		if err != nil {
			return err
		}

		// Add each file with commit metadata to commits slice
		files.ForEach(func(f *object.File) error {
			commitInfo := CommitInfo{
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

	// Output the commit information
	for _, commit := range commits {
		fmt.Printf("Author: %s, Filename: %s, Timestamp: %s\n",
			commit.Author, filepath.Base(commit.Filename), commit.Timestamp)
	}
}
