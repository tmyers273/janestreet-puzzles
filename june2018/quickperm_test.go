package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuickPerm(t *testing.T) {
	scenarios := []struct {
		Name     string
		Input    []int
		Expected int
	}{
		{
			Name:     "2",
			Input:    []int{1, 3, 2, 4},
			Expected: 24,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			// 1. Given
			count := 0

			// 2. Do this
			quickPerm(scenario.Input, func(in []int) error {
				count++
				return nil
			})

			// 3. Expect
			assert.Equal(t, scenario.Expected, count)
		})
	}
}
