package june2018

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"strconv"
	"testing"
)

type Row byte
type Col byte
type RowCollection struct {
	rows  []Row
	count int
}

func (r RowCollection) String() string {
	out := "\n"
	for i := 0; i < len(r.rows); i++ {
		out += strconv.Itoa(i) + ": " + r.rows[i].String() + "\n"
	}

	return out
}

type tuple struct {
	row, col uint8
}

func (r RowCollection) PassesContinuity(grid [][]int) bool {
	visited := make(map[tuple]struct{})

	paddedGrid := make([][]int, len(grid)+2)
	paddedGrid[0] = make([]int, len(grid[0])+2)
	paddedGrid[len(grid)+1] = make([]int, len(grid[0])+2)
	for i := 0; i < len(grid); i++ {
		paddedGrid[i+1] = make([]int, len(grid[i])+2)
		for j := 0; j < len(grid[i]); j++ {
			paddedGrid[i+1][j+1] = grid[i][j]
		}
	}

	for i := 1; i < len(paddedGrid)-1; i++ {
		for j := 1; j < len(paddedGrid)-1; j++ {
			if paddedGrid[i][j] != 0 {
				cnt := r.getChunkCount(paddedGrid, i, j, visited)
				return cnt == r.count
			}
		}
	}

	return false
}

func (r RowCollection) getChunkCount(grid [][]int, i, j int, visited map[tuple]struct{}) int {
	t := tuple{uint8(i), uint8(j)}
	if _, ok := visited[t]; ok {
		return 0
	}

	var ok bool
	count := 1
	visited[t] = struct{}{}

	// Top
	if grid[i][j-1] != 0 {
		if _, ok = visited[tuple{uint8(i), uint8(j - 1)}]; !ok {
			count += r.getChunkCount(grid, i, j-1, visited)
		}
	}

	// Left
	if grid[i-1][j] != 0 {
		if _, ok = visited[tuple{uint8(i - 1), uint8(j)}]; !ok {
			count += r.getChunkCount(grid, i-1, j, visited)
		}
	}

	// Right
	if grid[i][j+1] != 0 {
		if _, ok = visited[tuple{uint8(i), uint8(j + 1)}]; !ok {
			count += r.getChunkCount(grid, i, j+1, visited)
		}
	}

	// Bottom
	if grid[i+1][j] != 0 {
		if _, ok = visited[tuple{uint8(i + 1), uint8(j)}]; !ok {
			count += r.getChunkCount(grid, i+1, j, visited)
		}
	}

	return count
}

func (r RowCollection) Passes2x2() bool {
	var sum uint8
	rows := r.rows

	sum += uint8(rows[0]) & uint8(rows[1]) & (3 << 0)
	sum += uint8(rows[1]) & uint8(rows[2]) & (3 << 1)
	sum += uint8(rows[2]) & uint8(rows[3]) & (3 << 2)
	sum += uint8(rows[3]) & uint8(rows[4]) & (3 << 3)
	sum += uint8(rows[4]) & uint8(rows[5]) & (3 << 4)
	sum += uint8(rows[5]) & uint8(rows[6]) & (3 << 5)

	return sum == 0
}

func (r Row) String() string {
	return fmt.Sprintf("%07b", r)
}

func NewRowRepresentation(grid [][]int) RowCollection {
	out := make([]Row, len(grid))
	count := 0
	for i := 0; i < len(grid); i++ {
		out[i] = NewRow(grid[i])

		for j := 0; j < 7; j++ {
			if grid[i][j] != 0 {
				count++
			}
		}
	}

	return RowCollection{
		rows:  out,
		count: count,
	}
}

func NewRow(row []int) Row {
	out := 0
	for i := 0; i < 7; i++ {
		if row[i] != 0 {
			out |= 1 << (6 - i)
		}
	}

	return Row(out)
}

func NewCol(col []int) Col {
	return Col(NewRow(col))
}

/*
 * The grid is incomplete. Place numbers in some of the empty cells below so that in total the grid contains
 * one 1, two 2’s, etc., up to seven 7’s. Furthermore, each row and column must contain exactly 4 numbers
 * which sum to 20. Finally, the numbered cells must form a connected region*, but every 2-by-2 subsquare
 * in the completed grid must contain at least one empty cell.
 *
 * The answer to this puzzle is the product of the areas of the connected groups of empty squares in the completed grid.
 */

var oGrid = [][]int{
	{0, 4, 0, 0, 0, 0, 0},
	{0, 0, 6, 3, 0, 0, 6},
	{0, 0, 0, 0, 0, 5, 5},
	{0, 0, 0, 4, 0, 0, 0},
	{4, 7, 0, 0, 0, 0, 0},
	{2, 0, 0, 7, 4, 0, 0},
	{0, 0, 0, 0, 0, 1, 0},
}

var grid = [][]int{
	// One row and col "border" of 0s
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 4, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 6, 3, 0, 0, 6, 0},
	{0, 0, 0, 0, 0, 0, 5, 5, 0},
	{0, 0, 0, 0, 4, 0, 0, 0, 0},
	{0, 4, 7, 0, 0, 0, 0, 0, 0},
	{0, 2, 0, 0, 7, 4, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
}

func NewBoard() Board {
	b := Board{
		grid:        grid,
		rowCount:    make([]int, 7),
		colCount:    make([]int, 7),
		digitCounts: make([]int, 7+1),
		rowSum:      make([]int, 7),
		colSum:      make([]int, 7),
	}

	b.update()

	return b
}

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

	for i := 1; i < 8; i++ {
		out += "|"
		for j := 1; j < 8; j++ {
			if b.grid[i][j] == 0 {
				out += " |"
			} else {
				out += strconv.Itoa(b.grid[i][j]) + "|"
			}
		}
		if i == 7 {
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

var rowCount = []int{1, 3, 2, 1, 2, 3, 1}
var colCount = []int{2, 2, 1, 3, 1, 2, 2}
var digitCounts = []int{0, 1, 1, 1, 4, 2, 2, 2}

func TestItWorks(t *testing.T) {
	b := NewBoard()
	spew.Dump(b)
	spew.Dump(!b.invalid())

	b.update()
	r := check(b)
	spew.Dump(r, b)
	//i := 0
	//j := 0
	//for k := 0; k < 7; k++ {
	//	b.grid[i][j] = k + 1
	//	b.update()
	//	fmt.Printf("invalid after setting %d,%d to %d? %t\n", i, j, k+1, b.invalid())
	//}
	//
	//i = 0
	//j = 1
	//for k := 0; k < 7; k++ {
	//	b.grid[i][j] = k + 1
	//	b.update()
	//	fmt.Printf("invalid after setting %d,%d to %d? %t\n", i, j, k+1, b.invalid())
	//}
	//
	//i = 1
	//j = 1
	//for k := 0; k < 7; k++ {
	//	b.grid[i][j] = k + 1
	//	b.update()
	//	fmt.Printf("invalid after setting %d,%d to %d? %t\n", i, j, k+1, b.invalid())
	//}
}

func check(b Board) bool {
	if b.invalid() {
		//fmt.Printf("board is invalid\n")
		return false
	}

	for i := 1; i < 8; i++ {
		for j := 1; j < 8; j++ {
			if b.grid[i][j] != 0 {
				continue
			}

			//hasValid := false
			for k := 0; k < 7; k++ {
				//fmt.Printf("setting %d,%d to %d\n", i, j, k+1)
				ok := b.Set(i, j, k+1)
				if !ok {
					b.grid[i][j] = 0
					continue
				}

				//b.grid[i][j] = k + 1
				//b.update()

				if check(b) {
					//fmt.Println("found a potential valid")
					//fmt.Println()
					//hasValid = true
				}
			}
			//if !hasValid {
			//	//fmt.Printf("no valid combinations for %d,%d - setting to 0\n", i, j)
			//	b.grid[i][j] = 0
			//	b.update()
			//}
		}
	}

	return !b.invalid()
}
