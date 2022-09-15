package june2018

import (
	"fmt"
	"strconv"
)

type Board struct {
	grid        [][]int
	rowCount    []int
	colCount    []int
	digitCounts []int
	rowSum      []int
	colSum      []int
}

func (b *Board) Set(row, col, value int) bool {
	row--
	col--
	if b.rowCount[row] > 4 || b.colCount[col] > 4 {
		return false
	}

	if b.rowSum[row]+value > 20 || b.colSum[col]+value > 20 {
		return false
	}

	if b.digitCounts[value]+1 > value {
		return false
	}

	tlCount := 0
	trCount := 0
	blCount := 0
	brCount := 0

	row++
	col++

	if b.grid[row-1][col-1] != 0 {
		tlCount++
	}
	if b.grid[row-1][col] != 0 {
		tlCount++
		trCount++
	}
	if b.grid[row-1][col+1] != 0 {
		trCount++
	}
	if b.grid[row][col-1] != 0 {
		tlCount++
		blCount++
	}
	if b.grid[row][col+1] != 0 {
		trCount++
		brCount++
	}
	if b.grid[row+1][col-1] != 0 {
		blCount++
	}
	if b.grid[row+1][col] != 0 {
		blCount++
		brCount++
	}
	if b.grid[row+1][col+1] != 0 {
		brCount++
	}

	if tlCount > 2 || trCount > 2 || blCount > 2 || brCount > 2 {
		return false
	}

	b.grid[row][col] = value
	// ..
	// .x

	// ..
	// x.

	// x.
	// ..

	// .x
	// ..
	return true
}

func (b *Board) String() string {
	out := "\n -------------\n"

	for i := 0; i < 7; i++ {
		out += "|"
		for j := 0; j < 7; j++ {
			if b.grid[i][j] == 0 {
				out += " |"
			} else {
				out += strconv.Itoa(b.grid[i][j]) + "|"
			}
		}
		if i == 6 {
			out += "\n ------------- \n"
		} else {
			out += "\n|-------------|\n"
		}
	}
	out += "\n"

	out += "  | Count   | Sum     | Digit\n"
	out += "i | Row Col | Row Col | Count\n"
	out += "-----------------------------\n"
	for i := 0; i < 7; i++ {
		out += fmt.Sprintf("%d | %3d %3d | %3d %3d |     %d\n", i+1, b.rowCount[i], b.colCount[i], b.rowSum[i], b.colSum[i], b.digitCounts[i+1])
	}
	//out += " --------------"

	return out
}

func (b *Board) update() {
	for i := 0; i < 7; i++ {
		b.getRowCount(i)
		b.getColCount(i)
		b.getRowSum(i)
		b.getColSum(i)
	}
	b.getDigitCounts()
}

func (b *Board) invalid() bool {
	for i := 0; i < 7; i++ {
		// Can only have 4 numbers in rows and cols
		if b.rowCount[i] > 4 || b.colCount[i] > 4 {
			return true
		}

		// Rows and cols must sum to exactly 20
		if b.rowSum[i] > 20 || b.colSum[i] > 20 {
			return true
		}

		// Can only have one 1, two 2's, etc.
		if b.digitCounts[i+1] > i+1 {
			return true
		}

		// @todo must form a connected region

	}

	// Every 2x2 subsquare must contain at least one empty cell
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			used := 0
			if b.grid[i][j] != 0 {
				used++
			}
			if b.grid[i+1][j] != 0 {
				used++
			}
			if b.grid[i][j+1] != 0 {
				used++
			}
			if b.grid[i+1][j+1] != 0 {
				used++
			}

			if used == 4 {
				return true
			}
		}
	}

	return false
}

func (b *Board) getDigitCounts() {
	for i := 1; i < 8; i++ {
		b.digitCounts[i] = 0
	}

	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			if b.grid[i][j] > 0 {
				b.digitCounts[b.grid[i][j]]++
			}
		}
	}
}

func (b *Board) getRowCount(i int) int {
	count := 0
	for _, v := range b.grid[i] {
		if v != 0 {
			count++
		}
	}

	b.rowCount[i] = count
	return count
}

func (b *Board) getColCount(i int) int {
	count := 0
	for j := 0; j < 7; j++ {
		if b.grid[j][i] != 0 {
			count++
		}
	}

	b.colCount[i] = count
	return count
}

func (b *Board) getRowSum(i int) int {
	sum := 0

	for _, v := range b.grid[i] {
		sum += v
	}

	b.rowSum[i] = sum
	return sum
}

func (b *Board) getColSum(i int) int {
	sum := 0

	for j := 0; j < 7; j++ {
		sum += b.grid[j][i]
	}

	b.colSum[i] = sum
	return sum
}
