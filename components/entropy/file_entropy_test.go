package entropy

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileEntropy(t *testing.T) {
	var data []FileEntropy
	var expected, actual FileEntropy

	data = []FileEntropy{
		{Filename: "file1.txt", Entropy: 1.56},
		{Filename: "file2.txt", Entropy: 2.34},
		{Filename: "file3.txt", Entropy: 1.89},
	}

	fmt.Println("Original:", data)

	// Sort by Filename
	SortByFilename(data)
	fmt.Println("Sorted by Filename:", data)
	expected = FileEntropy{Filename: "file1.txt", Entropy: 1.56}
	actual = data[0]
	assert.Equal(t, expected, actual)

	// Sort by Entropy
	// SortByEntropy(files)
	// fmt.Println("Sorted by Entropy:", files)

}
