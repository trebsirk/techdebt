package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntropyFromOccurances(t *testing.T) {
	var data []int
	var actual, expected []TypeProb

	data = []int{}
	actual = MakeProbabilitiesFromOccurances(data)
	expected = []TypeProb{}
	assert.Equal(t, expected, actual, "should be equal")

	data = []int{0}
	actual = MakeProbabilitiesFromOccurances(data)
	expected = []TypeProb{1.0}
	assert.Equal(t, expected, actual, "should be equal")

	data = []int{0, 1}
	actual = MakeProbabilitiesFromOccurances(data)
	expected = []TypeProb{0.5, 0.5}
	assert.Equal(t, expected, actual, "should be equal")

	data = []int{0, 0, 1}
	actual = MakeProbabilitiesFromOccurances(data)
	expected = []TypeProb{2.0 / 3.0, 1.0 / 3.0}
	assert.Equal(t, expected, actual, "should be equal")

	data = []int{0, 0, 1, 2}
	actual = MakeProbabilitiesFromOccurances(data)
	expected = []TypeProb{2.0 / 4.0, 1.0 / 4.0, 1.0 / 4.0}
	assert.Equal(t, expected, actual, "should be equal")

}

func TestTransformOccurancesToCounts(t *testing.T) {
	var data []int
	var actual, expected []int

	data = []int{}
	actual = transformOccurancesToCounts(data)
	expected = []int{}
	assert.Equal(t, expected, actual, "should be equal")

	data = []int{0}
	actual = transformOccurancesToCounts(data)
	expected = []int{1}
	assert.Equal(t, expected, actual, "should be equal")

	data = []int{1, 1, 3}
	actual = transformOccurancesToCounts(data)
	expected = []int{0, 2, 0, 1}
	assert.Equal(t, expected, actual, "should be equal")

	data = []int{1, 2, 3}
	actual = transformOccurancesToCounts(data)
	expected = []int{0, 1, 1, 1}
	assert.Equal(t, expected, actual, "should be equal")

}

func TestEntropyFromCounts(t *testing.T) {
	var data []int
	var actual, expected []TypeProb

	data = []int{}
	actual = MakeProbabilitiesFromCounts(data)
	expected = []TypeProb{}
	assert.Equal(t, expected, actual, "should be equal")

	data = []int{1}
	actual = MakeProbabilitiesFromCounts(data)
	expected = []TypeProb{1.0}
	assert.Equal(t, expected, actual, "should be equal")

	data = []int{1, 1}
	actual = MakeProbabilitiesFromCounts(data)
	expected = []TypeProb{1.0 / 2.0, 1.0 / 2.0}
	assert.Equal(t, expected, actual, "should be equal")

	data = []int{1, 1, 1}
	actual = MakeProbabilitiesFromCounts(data)
	expected = []TypeProb{1.0 / 3.0, 1.0 / 3.0, 1.0 / 3.0}
	assert.Equal(t, expected, actual, "should be equal")

	data = []int{1, 1, 2}
	actual = MakeProbabilitiesFromCounts(data)
	expected = []TypeProb{1.0 / 4.0, 1.0 / 4.0, 2.0 / 4.0}
	assert.Equal(t, expected, actual, "should be equal")

}
