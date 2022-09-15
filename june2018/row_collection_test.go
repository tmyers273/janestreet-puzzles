package june2018

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRowCollection_Transpose(t *testing.T) {
	scenarios := []struct {
		Name     string
		Input    [][]int
		Expected [][]int
	}{
		{
			Name:     "it works",
			Input:    oGrid,
			Expected: oGridRotatedRight,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			// 1. Given
			orig := NewRowRepresentation(scenario.Input)

			// 2. Do this
			r := orig.Transpose()

			// 3. Expect
			expected := NewRowRepresentation(scenario.Expected)
			assert.Equal(t, expected, r)
		})
	}
}
