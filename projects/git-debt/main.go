package main

import (
	"encoding/json"
	"fmt"
	"os"

	"techdebt/components/commitinfo"
	"techdebt/components/entropy"
	"techdebt/components/git"
	"techdebt/components/helpers"
)

// writeFileEntropies writes a slice of FileEntropy structs to a file in JSON format.
func writeFileEntropies(filename string, entropies []entropy.FileEntropy) error {
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

// aggregateByFile aggregates CommitInfo data into a map where the key is the filename
// and the value is a list of authors who committed to that file.
func aggregateByFile(commits []commitinfo.CommitInfo) map[string][]string {
	// Initialize the map to hold filenames as keys and list of authors as values
	aggregated := make(map[string][]string)

	for _, commit := range commits {
		// Check if the filename already exists in the map
		if authors, exists := aggregated[commit.Filename]; exists {
			// Add the author if they are not already in the list
			if !helpers.Contains(authors, commit.Author) {
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
	// Initialize the map
	// {
	//  filename1: {author1: 1, author2: 6},
	//  filename2: {author0: 5, author1: 3}
	// }
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

func transformMapCountsToArray(countmap map[string]int) []int {
	var res []int
	var i int

	res = make([]int, len(countmap))
	i = 0
	for _, count := range countmap {
		res[i] = count
		i++
	}

	return res
}

// aggregateByFile aggregates CommitInfo data into a map where the key is the filename
// and the value is a list of authors who committed to that file.
func aggregateAuthorCounts(commits []commitinfo.CommitInfo) map[string]int {
	// Initialize the map
	// {
	//  author0: 5,
	//  author1: 1, ...
	// }
	aggregated := make(map[string]int)

	for _, commit := range commits {
		if _, exists := aggregated[commit.Author]; !exists {
			aggregated[commit.Author] = 0
		}
		aggregated[commit.Author]++
	}

	return aggregated
}

func aggregateFileCounts(commits []commitinfo.CommitInfo) map[string]int {
	// Initialize the map
	// {
	//  author0: 5,
	//  author1: 1, ...
	// }
	aggregated := make(map[string]int)

	for _, commit := range commits {
		if _, exists := aggregated[commit.Filename]; !exists {
			aggregated[commit.Filename] = 0
		}
		aggregated[commit.Filename]++
	}

	return aggregated
}

func calcRepoEntropy(commits []commitinfo.CommitInfo) float64 {
	// entropy by author
	agg_counts_map := aggregateAuthorCounts(commits)
	agg_counts_arr := transformMapCountsToArray(agg_counts_map)
	entropy_author := float64(entropy.CalculateEntropyOfCounts(agg_counts_arr))
	fmt.Printf("author entropy: %.3f\n", entropy_author)

	// entropy by file
	agg_counts_map = aggregateFileCounts(commits)
	agg_counts_arr = transformMapCountsToArray(agg_counts_map)
	entropy_file := float64(entropy.CalculateEntropyOfCounts(agg_counts_arr))
	fmt.Printf("file entropy: %.3f\n", entropy_file)

	return (entropy_file + entropy_author) / 2.0
}

func calcEntroyDemo() {
	arr := []int{1, 1, 2, 2, 3, 3, 3}
	e := entropy.CalculateEntropyOfOccurances(arr)
	fmt.Printf("Entropy: %.4f\n", e)
}

func calcEntroyByFile(commits []commitinfo.CommitInfo) {

	aggregatedCounts := aggregateCountsByFile(commits)

	// fmt.Printf("aggregatedCounts %v\n", aggregatedCounts)
	fileEntropies := make([]entropy.FileEntropy, 0)
	for filename, authorCountMap := range aggregatedCounts {
		arr := []int{}
		for _, count := range authorCountMap {
			arr = append(arr, count)
		}
		// fmt.Printf("%s arr: %v\n", filename, arr)
		e := entropy.CalculateEntropyOfOccurances(arr)
		// fmt.Printf("entropy = %f\n", entropy)
		fileEntropies = append(fileEntropies, entropy.FileEntropy{Filename: filename, Entropy: float64(e)})
	}

	totalEntropy := 0.0
	for _, fe := range fileEntropies {
		totalEntropy = totalEntropy + fe.Entropy
	}
	avgEntropy := totalEntropy / float64(len(fileEntropies))

	fmt.Printf("Repo Entropy (average of all files):%f\n", avgEntropy)

	fmt.Println("file entropies:")
	entropy.PrintFileEntropySlice(fileEntropies)

}

func main() {
	var repoPath string
	if len(os.Args) == 2 {
		repoPath = os.Args[1]
	} else {
		repoPath = "." //"/Users/davidwright/Documents/javascript-dev/trebsirk.github.io"
		// Define local repo directory
	}

	var commits []commitinfo.CommitInfo = git.GetCommits(repoPath)
	var overallEntropy float64 = calcRepoEntropy(commits)
	fmt.Printf("overallEntropy = %f\n", overallEntropy)

	fmt.Println("commits")
	// Output the commit information
	// for _, commit := range commits {
	// 	commit.LogPrint()
	// }
	// fmt.Println()

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

	// fmt.Println()
	// // entropy.CalculateEntropy(commits)

	// // calcEntroyByFile(commits)
	// fmt.Println(commits)

	// fmt.Println("Optimizations")
	// arr := []int{5, 3, 1}
	// e := entropy.CalculateEntropyOfCounts(arr)
	// fmt.Printf("Entropy for %v: %.4f\n", arr, e)

	// arr = []int{6, 3, 1}
	// e = entropy.CalculateEntropyOfCounts(arr)
	// fmt.Printf("Entropy for %v: %.4f\n", arr, e)

	// arr = []int{5, 4, 1}
	// e = entropy.CalculateEntropyOfCounts(arr)
	// fmt.Printf("Entropy for %v: %.4f\n", arr, e)

	// arr = []int{5, 3, 2}
	// e = entropy.CalculateEntropyOfCounts(arr)
	// fmt.Printf("Entropy for %v: %.4f\n", arr, e)

}
