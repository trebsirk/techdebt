package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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

func (c *CommitInfo) print() {
	fmt.Printf("CommitInfo[Author: %s, Filename: %s, Timestamp: %s]\n",
		c.Author, filepath.Base(c.Filename),
		c.Timestamp.Format(time.RFC3339))
}

func (c *CommitInfo) logPrint() {
	fmt.Printf("%s %s %s\n", c.Author, filepath.Base(c.Filename),
		c.Timestamp.Format(time.RFC3339))
}

type FileEntropy struct {
	Filename string  `json:"filename"`
	Entropy  float64 `json:"entropy"`
}

func (fe *FileEntropy) printWithFormatting(colWidth int) {
	fmt.Printf("%-{colWidth}s\n", fe.Filename)
}

func printFileEntropySlice(entropies []FileEntropy) {
	// Find the maximum width for each column
	maxNameWidth := 0
	for _, fe := range entropies {
		if len(fe.Filename) > maxNameWidth {
			maxNameWidth = len(fe.Filename)
		}
	}

	// Print the header
	fmt.Printf("%-*s | %s\n", maxNameWidth, "Filename", "Score")
	fmt.Println("---------------------------")

	// Print the rows
	for i := 0; i < len(entropies); i++ {
		fmt.Printf("%-*s | %.4f\n", maxNameWidth,
			entropies[i].Filename,
			entropies[i].Entropy)
	}
}

// writeFileEntropies writes a slice of FileEntropy structs to a file in JSON format.
func writeFileEntropies(filename string, entropies []FileEntropy) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(entropies); err != nil {
		return fmt.Errorf("could not write JSON to file: %w", err)
	}

	return nil
}
func makeProbabilitiesFromOccurances(data []int) map[int]float64 {
	// data is occurances, eg [0 1 1 3 0 2 2 1 0]
	// where 0,1,2,3 are observation indices
	ps := make(map[int]float64)
	counts := make(map[int]int)
	total := len(data)
	for _, obj := range data {
		counts[obj]++
	}
	for obj, c := range counts {
		ps[obj] = float64(c) / float64(total)
	}
	return ps
}

func makeProbabilitiesFromCounts(data []int) []float64 {
	// data is occurances, eg [3 4 0]
	// where 3,4,0 are counts for each index
	ps := make([]float64, len(data))
	total := 0
	for _, count := range data {
		total = total + count
	}
	for i, c := range data {
		ps[i] = float64(c) / float64(total)
	}
	return ps
}

func calculateEntropy(data []int) float64 {
	if len(data) == 0 {
		return 0.0
	}

	ps := makeProbabilitiesFromCounts(data)

	// Calculate the entropy
	var entropy float64

	for obj, p := range ps {
		fmt.Printf("probability of %d = %f\n", obj, p)
		entropy += -p * math.Log2(p)
	}
	return entropy
}

// aggregateByFile aggregates CommitInfo data into a map where the key is the filename
// and the value is a list of authors who committed to that file.
func aggregateByFile(commits []CommitInfo) map[string][]string {
	// Initialize the map to hold filenames as keys and list of authors as values
	aggregated := make(map[string][]string)

	for _, commit := range commits {
		// Check if the filename already exists in the map
		if authors, exists := aggregated[commit.Filename]; exists {
			// Add the author if they are not already in the list
			if !contains(authors, commit.Author) {
				aggregated[commit.Filename] = append(authors, commit.Author)
			}
		} else {
			// Initialize the list with the current author
			aggregated[commit.Filename] = []string{commit.Author}
		}
	}

	return aggregated
}

func aggregateByFileWithRepeats(commits []CommitInfo) map[string][]string {
	// Initialize the map to hold filenames as keys and list of authors as values
	aggregated := make(map[string][]string)

	for _, commit := range commits {
		// Check if the filename already exists in the map
		if authors, exists := aggregated[commit.Filename]; exists {
			// Add the author if they are not already in the list
			aggregated[commit.Filename] = append(authors, commit.Author)
		} else {
			// Initialize the list with the current author
			aggregated[commit.Filename] = []string{commit.Author}
		}
	}

	return aggregated
}

// aggregateByFile aggregates CommitInfo data into a map where the key is the filename
// and the value is a list of authors who committed to that file.
func aggregateCountsByFile(commits []CommitInfo) map[string]map[string]int {
	// Initialize the map to hold filenames as keys and list of authors as values
	aggregated := make(map[string]map[string]int)

	for _, commit := range commits {
		// Check if the filename already exists in the map
		if _, exists := aggregated[commit.Filename]; !exists {
			aggregated[commit.Filename] = make(map[string]int)
		}

		// inc the author count
		if _, exists := aggregated[commit.Filename][commit.Author]; !exists {
			aggregated[commit.Filename][commit.Author] = 0
		}

		aggregated[commit.Filename][commit.Author]++

	}

	return aggregated
}

// Helper function to check if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func calcEntroyDemo() {
	arr := []int{1, 1, 2, 2, 3, 3, 3}
	entropy := calculateEntropy(arr)
	fmt.Printf("Entropy: %.4f\n", entropy)
}

func calcEntroyByFile(commits []CommitInfo) {

	aggregatedCounts := aggregateCountsByFile(commits)

	fmt.Printf("aggregatedCounts %v\n", aggregatedCounts)
	fileEntropy := make([]FileEntropy, 0)
	for filename, authorCountMap := range aggregatedCounts {
		arr := []int{}
		for _, count := range authorCountMap {
			arr = append(arr, count)
		}
		// fmt.Printf("%s arr: %v\n", filename, arr)
		entropy := calculateEntropy(arr)
		// fmt.Printf("entropy = %f\n", entropy)
		fileEntropy = append(fileEntropy, FileEntropy{filename, entropy})
	}

	fmt.Printf("Entropy:\n%v\n", fileEntropy)

	printFileEntropySlice(fileEntropy)

}

func main() {
	var repoPath string
	if len(os.Args) == 2 {
		repoPath = os.Args[1]
	} else {
		repoPath = "." //"/Users/davidwright/Documents/javascript-dev/trebsirk.github.io"
		// Define local repo directory
	}

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

	fmt.Println("commits")
	// Output the commit information
	for _, commit := range commits {
		//commit.print()
		commit.logPrint()
	}
	fmt.Println()

	// fmt.Println("aggregations")

	// // Aggregate the commit data by filename
	// aggregated := aggregateByFile(commits)

	// // Print the aggregated data
	// for filename, authors := range aggregated {
	// 	fmt.Printf("File: %s, Authors: %v\n", filename, authors)
	// }

	fmt.Println()

	// aggregatedCounts := aggregateCountsByFile(commits)
	// for filename, authorCountMap := range aggregatedCounts {
	// 	for author, count := range authorCountMap {
	// 		fmt.Printf("File: %s, Author: %s, Count: %d\n", filename, author, count)
	// 	}
	// }

	calcEntroyDemo()
	fmt.Println()
	calcEntroyByFile(commits)

}
