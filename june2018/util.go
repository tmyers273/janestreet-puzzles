package june2018

import "fmt"

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

func NewBoardFromKey(key uint64) [][]int {
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
