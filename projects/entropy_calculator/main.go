package main

import (
	"fmt"
	"log"

	"techdebt/components/entropy"
)

func main() {
	// Sample data
	entropies := []entropy.FileEntropy{
		{Filename: "file1.txt", Entropy: 1.56},
		{Filename: "file2.txt", Entropy: 2.34},
		{Filename: "file3.txt", Entropy: 1.89},
	}

	// Write to file
	if err := entropy.WriteFileEntropies("entropies.json", entropies); err != nil {
		log.Fatalf("Error writing to file: %v\n", err)
	}
	fmt.Println("Data successfully written to entropies.json")

	// Read from file
	fes, err := entropy.ReadFileEntropies("entropies.json")
	if err != nil {
		log.Fatalf("Error reading from file: %v\n", err)
	}

	// Print the read data
	fmt.Println("Data read from entropies.json:")
	entropy.PrintFileEntropySlice(fes)
}
