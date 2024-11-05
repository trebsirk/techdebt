package entropy

import (
	"math"
	"techdebt/components/helpers"
)

// CalculateEntropy calculates the entropy of an array of integers.
func CalculateEntropyOfProbabilities(ps []helpers.TypeProb) helpers.TypeProb {
	// ps is a list of probabilities that sum to one

	if len(ps) == 0 {
		return 0.0
	}

	var entropy helpers.TypeProb

	for _, p := range ps {
		if p != 0.0 {
			entropy += helpers.TypeProb(float64(p) * math.Log2(1.0/float64(p)))
		}
	}

	return entropy
}

// CalculateEntropy calculates the entropy of an array of integers.
func CalculateEntropyOfOccurances(data []int) helpers.TypeProb {
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

	return helpers.TypeProb(CalculateEntropyOfProbabilities(probs))
}

func CalculateEntropyOfCounts(data []int) helpers.TypeProb {
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

	return helpers.TypeProb(CalculateEntropyOfProbabilities(probs))
}
