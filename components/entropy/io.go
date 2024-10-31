package entropy

import (
	"encoding/json"
	"fmt"
	"os"
)

// WriteFileEntropies writes a slice of FileEntropy structs to a file in JSON format.
func WriteFileEntropies(filename string, entropies []FileEntropy) error {
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

// ReadFileEntropies reads a JSON file and decodes it into a slice of FileEntropy structs.
func ReadFileEntropies(filename string) ([]FileEntropy, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	var entropies []FileEntropy
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&entropies); err != nil {
		return nil, fmt.Errorf("could not decode JSON: %w", err)
	}

	return entropies, nil
}

func PrintFileEntropySlice(entropies []FileEntropy) {
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
