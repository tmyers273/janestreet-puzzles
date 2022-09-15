package june2018

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func printBin(num int) {
	fmt.Printf("%049b\n", num)
	s := fmt.Sprintf("%049b\n", num)
	for i := 0; i < 49; i += 7 {
		fmt.Printf("%d: %s\n", i/7, s[i:i+7])
	}
}

type tuple struct {
	row, col uint8
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

var oGridRotatedRight = [][]int{
	{0, 2, 4, 0, 0, 0, 0},
	{0, 0, 7, 0, 0, 0, 4},
	{0, 0, 0, 0, 0, 6, 0},
	{0, 7, 0, 4, 0, 3, 0},
	{0, 4, 0, 0, 0, 0, 0},
	{1, 0, 0, 0, 5, 0, 0},
	{0, 0, 0, 0, 5, 6, 0},
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
		grid:        oGrid,
		rowCount:    make([]int, 7),
		colCount:    make([]int, 7),
		digitCounts: make([]int, 7+1),
		rowSum:      make([]int, 7),
		colSum:      make([]int, 7),
	}

	b.update()

	return b
}

var rowCount = []int{1, 3, 2, 1, 2, 3, 1}
var colCount = []int{2, 2, 1, 3, 1, 2, 2}
var digitCounts = []int{0, 1, 1, 1, 4, 2, 2, 2}

func loadValidKeys() ([]uint64, error) {
	rawKeys, err := os.ReadFile("valid_keys.txt")
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

func TestValidBoardsPass2(t *testing.T) {
	keys, err := loadValidKeys()
	if err != nil {
		panic(err)
	}

	for _, key := range keys {
		madeChanges := true
		grid := NewBoardFromKey(key)
		rows := NewRowRepresentation(oGrid)
		changeCount := 0
		//cols := NewRowRepresentation(oGrid).Transpose()

		for madeChanges {
			madeChanges = false

			for i := 0; i < 7; i++ {
				//fmt.Printf("count for row i: %d is %d\n", i, rows.GetCountNumbers(i))
				if rows.GetCountNumbers(i) == 3 {
					sum := 20
					missingIndex := 0
					for j := 0; j < 7; j++ {
						if grid[i][j] == 9 {
							missingIndex = j
						} else {
							sum -= grid[i][j]
						}
					}

					if sum > 7 {
						continue
					}
					fmt.Printf("Found 3 in row %d, so we can conclusively set (%d,%d) to %d\n", i, i, missingIndex, sum)

					grid[i][missingIndex] = sum
					rows.Set(uint8(i), uint8(missingIndex))
					changeCount++
					//cols = rows.Transpose()
					madeChanges = true

					//fmt.Println("After change:")
					//printGrid(grid)
					//fmt.Printf("rows: %s\n", rows)
					//fmt.Println()
					//fmt.Printf("rows: %s\n", rows)
					//fmt.Printf("cols: %s\n", cols)
				}
			}
			for i := 0; i < 7; i++ {
				if rows.GetColCountNumbers(i) == 3 {
					sum := 20
					missingIndex := 0
					for j := 0; j < 7; j++ {
						if grid[j][i] == 9 {
							missingIndex = j
						} else {
							sum -= grid[j][i]
						}
					}

					if sum > 7 {
						continue
					}
					fmt.Printf("Found 3 in col %d, so we can conclusively set (%d,%d) to %d\n", i, missingIndex, i, sum)

					grid[missingIndex][i] = sum
					rows.Set(uint8(missingIndex), uint8(i))
					changeCount++
					//cols = rows.Transpose()

					//fmt.Println("After change")
					//printGrid(grid)
					//fmt.Printf("rows: %s\n", rows)
					//fmt.Println()
					//fmt.Printf("cols: %s\n", cols)

					madeChanges = true
				}
			}
		}

		fmt.Printf("Finished for key %d. Have %d nums after %d changes\n", key, rows.count, changeCount)
		printGrid(grid)
		fmt.Println(rows)
		if rows.count == 28 {
			//fmt.Printf("cols: %s\n", cols)
			fmt.Printf("rows: %s\n", rows)
			//fmt.Println(rows.rows[i].ToSlice())

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
			fmt.Print("\n\n")

			return
		}
	}
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

	//fmt.Printf("Board info for key %d:\n", key)

	//fmt.Printf("\n\n")
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

// TestGetBoards generates a list of valid board keys. We need a way to generate a list
// of potential places to put the 15 unused numbers. This is done by looping through and
// iterating a counter. The Hamming weight of the iterator is used and when it reaches
// 15, we use the index. The set bits in this index are used to reference the unused
// slots on the board. These get plugged in, then we check if the board structure is
// valid. Board structure is just checking for 2x2 correctness, 4 per row/col correctness,
// and continuity correctness.
//
// Once we have generated a list of valid keys, we need to do a second pass to generate
// The actual numbers in the correct locations. Still a todo.
func TestGetBoards(t *testing.T) {
	empties := generateEmpties()

	var iterations int
	start := time.Now()

	var key uint64 = 7669990665
	var rows RowCollection
	validCount := 0
	rows = NewRowRepresentation(oGrid)
	orig := NewRowRepresentation(oGrid)
	grid := make([][]int, 7)
	for i := 0; i < 7; i++ {
		grid[i] = make([]int, 7)
	}
	var hammingWeight uint64
	var x uint64
	for {
		copy(rows.rows, orig.rows)
		rows.count = orig.count

		hammingWeight = 0
		for hammingWeight != 15 {
			key++

			x = key
			x -= (x >> 1) & m1              //put count of each 2 bits into those 2 bits
			x = (x & m2) + ((x >> 2) & m2)  //put count of each 4 bits into those 4 bits
			x = (x + (x >> 4)) & m4         //put count of each 8 bits into those 8 bits
			hammingWeight = (x * h01) >> 56 //returns left 8 bits of x + (x<<8) + (x<<16) + (x<<24) + ...
			//hammingWeight = getHammingWeight64Fast(key)
		}

		for i := 0; i < 64; i++ {
			v := key & (1 << i)
			if v != 0 {
				rows.Set(empties[i].row, empties[i].col)
			}
		}

		if rows.Passes2x2() && rows.Passes4Check() {
			if rows.Transpose().Passes4Check() {
				for i := 0; i < 7; i++ {
					copy(grid[i], oGrid[i])
				}
				for i := 0; i < 64; i++ {
					v := key & (1 << i)
					if v != 0 {
						grid[empties[i].row][empties[i].col] = 9
					}
				}

				if rows.PassesContinuity(grid) {
					fmt.Printf("********** Found another valid board! key = %d, cnt = %d\n", key, validCount)
					validCount++

					//return
				}
			}
		}

		iterations++
		if iterations%100000000 == 0 {
			fmt.Printf("Iterations: %dM after %s (%.2fM / sec)\n", iterations/1000000, time.Since(start), float64(iterations)/1000000/time.Since(start).Seconds())
			fmt.Printf("    Binary key: %64b %[1]d\n", key)
			//return
		}
		//rows.Set()
	}
	//spew.Dump(pickCount, l)
	//spew.Dump(picked)
	//spew.Dump(empties)
}

// Should have 28 nums (7+6+5+4+3+2+1 = 28)
// Start with 13, need to add 15 more
// 49 spots, 13 filled, 36 left to place 15
// 36 pick 15 = 5,567,902,560 = 5.5B possible boards
// without looking at additional constraints.

// 3, 1, 2, 3, 2, 1, 3 = 20
func TestItWorks(t *testing.T) {
	b := NewBoard()
	spew.Dump(b)
	spew.Dump(!b.invalid())
	b.update()
	os.Exit(1)

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
