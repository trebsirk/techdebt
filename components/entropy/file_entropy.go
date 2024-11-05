package entropy

import (
	"sort"
)

type FileEntropy struct {
	Filename string  `json:"filename"`
	Entropy  float64 `json:"entropy"`
}

func SortByFilename(files []FileEntropy) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Filename < files[j].Filename
	})
}

// SortByEntropy sorts a slice of FileEntropy by Entropy in ascending order.
func SortByEntropy(files []FileEntropy) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Entropy < files[j].Entropy
	})
}
