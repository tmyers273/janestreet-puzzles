package june2018

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRowRepresentation(t *testing.T) {
	// 1. Given
	rows := NewRowRepresentation(oGrid)

	// 2. Do this
	r := rows.Passes2x2()
	assert.Equal(t, true, r)

	// 3. Expect
	// Set a number so that it fails the 2x2 test
	oGrid[5][1] = 1
	rows = NewRowRepresentation(oGrid)
	r = rows.Passes2x2()
	assert.Equal(t, false, r)
}

var contiguous1 = [][]int{
	{0, 4, 1, 0, 0, 0, 0},
	{0, 0, 6, 3, 1, 1, 6},
	{0, 0, 0, 0, 0, 5, 5},
	{0, 0, 1, 4, 1, 1, 0},
	{4, 7, 1, 0, 0, 0, 1},
	{2, 0, 1, 7, 4, 1, 1},
	{0, 0, 0, 0, 0, 1, 0},
}

// 2585 ns / op
// 2014 ns/op
// 154 ns/op
func BenchmarkPassesContinuity(b *testing.B) {
	rows1 := NewRowRepresentation(oGrid)
	rows2 := NewRowRepresentation(contiguous1)

	for i := 0; i < b.N; i++ {
		rows1.PassesContinuity(oGrid)
		rows2.PassesContinuity(contiguous1)
	}
}

func TestPassesContinuity(t *testing.T) {
	scenarios := []struct {
		Name     string
		Grid     [][]int
		Expected bool
	}{
		{
			Name:     "original",
			Grid:     oGrid,
			Expected: false,
		}, {
			Name:     "contiguous1",
			Grid:     contiguous1,
			Expected: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			// 1. Given
			rows := NewRowRepresentation(scenario.Grid)

			// 2. Do this
			r := rows.PassesContinuity(scenario.Grid)

			// 3. Expect
			assert.Equal(t, scenario.Expected, r)
		})
	}
}

// 1.87 ns/op
func BenchmarkPasses2x2(b *testing.B) {
	oGrid[5][1] = 1
	rows := NewRowRepresentation(oGrid)

	var r bool
	for i := 0; i < b.N; i++ {
		r = rows.Passes2x2()
	}
	assert.Equal(b, false, r)
}

func TestRow(t *testing.T) {
	scenarios := []struct {
		In       []int
		Expected string
	}{
		{
			In:       []int{1, 1, 1, 1, 1, 1, 1},
			Expected: "01111111",
		},
		{
			In:       []int{0, 0, 0, 0, 0, 0, 1},
			Expected: "00000001",
		},
		{
			In:       []int{1, 0, 0, 0, 0, 0, 1},
			Expected: "01000001",
		},
	}

	for _, scenario := range scenarios {
		t.Run("", func(t *testing.T) {
			// 1. Given

			// 2. Do this
			r := NewRow(scenario.In)
			s := fmt.Sprintf("%08b", r)

			// 3. Expect
			assert.Equal(t, scenario.Expected, s)
		})
	}
}

func TestBoard_Invalid(t *testing.T) {
	scenarios := []struct {
		Name string
	}{
		{
			Name: "",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			// 1. Given

			// 2. Do this

			// 3. Expect
		})
	}
}
