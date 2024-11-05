package entropy

import (
	"fmt"
	"techdebt/components/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntropyCalc(t *testing.T) {
	var data []int
	var expected, actual helpers.TypeProb
	var probs []helpers.TypeProb

	data = []int{0, 1}
	probs = helpers.MakeProbabilitiesFromOccurances(data)
	fmt.Println(probs)
	actual = CalculateEntropyOfProbabilities(probs)
	expected = 1.0
	assert.EqualValues(t, expected, actual, "they should be equal")

	data = []int{0, 0, 0, 1}
	probs = helpers.MakeProbabilitiesFromOccurances(data)
	fmt.Println(probs)
	actual = CalculateEntropyOfProbabilities(probs)
	expected = 0.8112781
	assert.EqualValues(t, expected, actual, "they should be equal")

	data = []int{0, 0, 0, 0, 0, 0, 1}
	probs = helpers.MakeProbabilitiesFromOccurances(data)
	fmt.Println(probs)
	actual = CalculateEntropyOfProbabilities(probs)
	expected = 0.5916728
	assert.EqualValues(t, expected, actual, "they should be equal")

	data = []int{0, 0, 1, 1, 2, 2}
	probs = helpers.MakeProbabilitiesFromOccurances(data)
	fmt.Println(probs)
	actual = CalculateEntropyOfProbabilities(probs)
	expected = 1.5849626
	assert.EqualValues(t, expected, actual, "they should be equal")

	data = []int{0, 0, 1, 1, 2, 2, 3, 3, 3, 3, 3}
	probs = helpers.MakeProbabilitiesFromOccurances(data)
	fmt.Println(probs)
	actual = CalculateEntropyOfProbabilities(probs)
	expected = 1.8585552
	assert.EqualValues(t, expected, actual, "they should be equal")

	// from counts
	data = []int{0, 1, 10}
	probs = helpers.MakeProbabilitiesFromCounts(data)
	fmt.Println(probs)
	actual = CalculateEntropyOfProbabilities(probs)
	expected = 0.43949693
	assert.EqualValues(t, expected, actual, "they should be equal")

}
