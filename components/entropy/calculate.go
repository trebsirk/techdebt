package entropy

import (
	"math"
	"techdebt/components/helpers"
)

type prob_t float32

// CalculateEntropy calculates the entropy of an array of integers.
func CalculateEntropyOfProbabilities(ps []helpers.TypeProb) prob_t {
	// ps is a list of probabilities that sum to one

	if len(ps) == 0 {
		return 0.0
	}

	var entropy prob_t

	for _, p := range ps {
		entropy += prob_t(float64(p) * math.Log2(1.0/float64(p)))
	}

	return entropy
}

// CalculateEntropy calculates the entropy of an array of integers.
func CalculateEntropyOfOccurances(data []int) prob_t {
	// data is a list of probabilities that sum to one

	if len(data) == 0 {
		return 0.0
	}

	var probs []helpers.TypeProb
	var ps []helpers.TypeProb
	ps = helpers.MakeProbabilitiesFromOccurances(data)
	probs = make([]helpers.TypeProb, len(ps))
	for i, p := range ps {
		probs[i] = helpers.TypeProb(p)
	}

	return prob_t(CalculateEntropyOfProbabilities(probs))
}

func CalculateEntropyOfCounts(data []int) prob_t {
	// data is a list of probabilities that sum to one

	if len(data) == 0 {
		return 0.0
	}

	var probs []helpers.TypeProb
	var ps []helpers.TypeProb
	ps = helpers.MakeProbabilitiesFromCounts(data)
	probs = make([]helpers.TypeProb, len(ps))
	for i, p := range ps {
		probs[i] = helpers.TypeProb(p)
	}

	return prob_t(CalculateEntropyOfProbabilities(probs))
}
