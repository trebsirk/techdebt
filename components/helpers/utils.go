package helpers

// Contains checks if a slice contains a specific string.
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func transformOccurancesToCounts(data []int) []int {
	/*
		transforms [0 1 1 3] into [1 2 0 1]
		this works for dense data
		TODO: implement sparse data solution using map[int]int
	*/
	if len(data) == 0 {
		return data
	}

	var maxElem int
	var counts []int

	maxElem = data[0]
	for _, obj := range data {
		maxElem = max(obj, maxElem)
	}

	counts = make([]int, maxElem+1)
	for _, obj := range data {
		counts[obj]++
	}
	return counts
}

func MakeProbabilitiesFromOccurances(data []int) []float64 {
	// data is occurances, eg [0 1 1 3 0 2 2 1 0]
	// where 0,1,2,3 are observation indices
	var counts []int
	var probs []float64

	counts = transformOccurancesToCounts(data)
	probs = MakeProbabilitiesFromCounts(counts)
	return probs
}

func MakeProbabilitiesFromCounts(data []int) []float64 {
	// data is occurances, eg [3 4 0]
	// where 3,4,0 are counts for each index
	var total int
	var probs []float64

	probs = make([]float64, len(data))
	total = 0

	for _, count := range data {
		total = total + count
	}

	for i, c := range data {
		probs[i] = float64(c) / float64(total)
	}

	return probs
}