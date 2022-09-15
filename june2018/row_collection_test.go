package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRowRepresentation(t *testing.T) {
	scenarios := []struct {
		Name     string
		Input    [][]int
		Expected RowCollection
	}{
		{
			Name:  "",
			Input: oGrid,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			// 1. Given

			// 2. Do this
			r := NewRowRepresentation(scenario.Input)

			fmt.Println(r)

			// 3. Expect
		})
	}
}

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
			fmt.Println(orig)
			fmt.Println(r)
			fmt.Println(expected)
			assert.Equal(t, expected, r)
		})
	}
}

var oGridRotatedRight = [][]int{
	{0, 2, 4, 0, 0, 0, 0},
	{0, 0, 7, 0, 0, 0, 4},
	{0, 0, 0, 0, 0, 6, 0},
	{0, 7, 0, 4, 0, 3, 0},
	{0, 4, 0, 0, 0, 0, 0},
	{1, 0, 0, 0, 5, 0, 0},
	{0, 0, 0, 0, 5, 6, 0},
}
