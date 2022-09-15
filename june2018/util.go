package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type tuple struct {
	row, col uint8
}

var oGrid = [][]int{
	{0, 4, 0, 0, 0, 0, 0},
	{0, 0, 6, 3, 0, 0, 6},
	{0, 0, 0, 0, 0, 5, 5},
	{0, 0, 0, 4, 0, 0, 0},
	{4, 7, 0, 0, 0, 0, 0},
	{2, 0, 0, 7, 4, 0, 0},
	{0, 0, 0, 0, 0, 1, 0},
}

func NewGridFromKey(key uint64) [][]int {
	empties := generateEmpties()

	grid := make([][]int, 7)
	for i := 0; i < 7; i++ {
		grid[i] = make([]int, 7)
	}
	for i := 0; i < 7; i++ {
		copy(grid[i], oGrid[i])
	}
	for i := 0; i < 64; i++ {
		v := key & (1 << i)
		if v != 0 {
			grid[empties[i].row][empties[i].col] = 9
		}
	}

	return grid
}

func generateEmpties() []tuple {
	var empties []tuple
	var i, j uint8

	for i = 0; i < 7; i++ {
		for j = 0; j < 7; j++ {
			if oGrid[i][j] == 0 {
				empties = append(empties, tuple{i, j})
			}
		}
	}

	return empties
}

// printGrid is a helper function, it just pretty prints the grid.
func printGrid(grid [][]int) {
	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			if grid[i][j] == 0 {
				fmt.Printf("_ ")
			} else if grid[i][j] == 9 {
				fmt.Printf("* ")
			} else {
				fmt.Printf("%d ", grid[i][j])
			}
		}
		fmt.Println()
	}
}

// neq is a branchless version of a!=b. Go cannot natively cast
// a bool to an int val, but the compiler is smart enough to.
func neq(a, b uint64) int8 {
	if a != b {
		return 1
	}

	return 0
}

// eq is a branchless version of a==b. Go cannot natively cast
// a bool to an int val, but the compiler is smart enough to.
func eq(a, b int8) int {
	if a == b {
		return 1
	}
	return 0
}

func loadValidKeys() ([]uint64, error) {
	rawKeys, err := os.ReadFile("june2018/valid_keys.txt")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(rawKeys)), "\n")
	keys := make([]uint64, len(lines))
	for i, line := range lines {
		keys[i], _ = strconv.ParseUint(line, 10, 64)
	}

	return keys, nil
}
