package entropy

import "math"

// CalculateEntropy calculates the entropy of an array of integers.
func CalculateEntropy(arr []int) float64 {
	if len(arr) == 0 {
		return 0.0
	}

	// Frequency map to count occurrences of each integer
	frequency := make(map[int]int)
	for _, num := range arr {
		frequency[num]++
	}

	// Calculate the entropy
	var entropy float64
	total := float64(len(arr))

	for _, count := range frequency {
		probability := float64(count) / total
		entropy += -probability * math.Log2(probability)
	}

	return entropy
}
