package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"techdebt/components/commitinfo"
	"techdebt/components/entropy"
	"techdebt/components/git"
)

type FileEntropy struct {
	Filename string  `json:"filename"`
	Entropy  float64 `json:"entropy"`
}

func (fe *FileEntropy) printWithFormatting(colWidth int) {
	fmt.Printf("%-{colWidth}s\n", fe.Filename)
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
	entropy.PrintFileEntropySlice(fileEntropy)

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
