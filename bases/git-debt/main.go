package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"techdebt/components/commitinfo"
	"techdebt/components/git"
	"techdebt/components/helpers"
)

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

func calculateEntropy(data []int) float64 {
	if len(data) == 0 {
		return 0.0
	}

	ps := helpers.MakeProbabilitiesFromCounts(data)

	// Calculate the entropy
	var entropy float64

	for _, p := range ps {
		entropy += -p * math.Log2(p)
	}
	return entropy
}

// aggregateByFile aggregates CommitInfo data into a map where the key is the filename
// and the value is a list of authors who committed to that file.
func aggregateByFile(commits []commitinfo.CommitInfo) map[string][]string {
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

func aggregateByFileWithRepeats(commits []commitinfo.CommitInfo) map[string][]string {
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
func aggregateCountsByFile(commits []commitinfo.CommitInfo) map[string]map[string]int {
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

func calcEntroyDemo() {
	arr := []int{1, 1, 2, 2, 3, 3, 3}
	entropy := calculateEntropy(arr)
	fmt.Printf("Entropy: %.4f\n", entropy)
}

func calcEntroyByFile(commits []CommitInfo) {

	aggregatedCounts := aggregateCountsByFile(commits)

	// fmt.Printf("aggregatedCounts %v\n", aggregatedCounts)
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

	totalEntropy := 0.0
	for _, fe := range fileEntropy {
		totalEntropy = totalEntropy + fe.Entropy
	}
	avgEntropy := totalEntropy / float64(len(fileEntropy))

	fmt.Printf("Repo Entropy (average of all files):%f\n", avgEntropy)

	fmt.Println("file entropies:")
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

	var commits []commitinfo.CommitInfo = git.GetCommits()

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

	fmt.Println()
	calcEntroyByFile(commits)

	fmt.Println("Optimizations")
	arr := []int{5, 3, 1}
	entropy := calculateEntropy(arr)
	fmt.Printf("Entropy: %.4f\n", entropy)

	arr = []int{6, 3, 1}
	entropy = calculateEntropy(arr)
	fmt.Printf("Entropy: %.4f\n", entropy)

	arr = []int{5, 4, 1}
	entropy = calculateEntropy(arr)
	fmt.Printf("Entropy: %.4f\n", entropy)

	arr = []int{5, 3, 2}
	entropy = calculateEntropy(arr)
	fmt.Printf("Entropy: %.4f\n", entropy)

}
