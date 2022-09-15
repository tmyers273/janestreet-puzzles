package main

import (
	"errors"
	"fmt"
	"time"
)

/*
 * The grid is incomplete. Place numbers in some of the empty cells below so that in total the grid contains
 * one 1, two 2’s, etc., up to seven 7’s. Furthermore, each row and column must contain exactly 4 numbers
 * which sum to 20. Finally, the numbered cells must form a connected region*, but every 2-by-2 subsquare
 * in the completed grid must contain at least one empty cell.
 *
 * The answer to this puzzle is the product of the areas of the connected groups of empty squares in the completed grid.
 */

// Should have 28 nums (7+6+5+4+3+2+1 = 28)
// Start with 13, need to add 15 more
// 49 spots, 13 filled, 36 left to place 15
// 36 pick 15 = 5,567,902,560 = 5.5B possible boards
// without looking at additional constraints.

func main() {
	keys := getValidBoardStructures()
	fmt.Printf("Found %d valid board structures\n", len(keys))
	validateBoardStructures(keys)
}

// getValidBoardStructures generates a list of valid board keys. Board keys are uint64s
// with bits set to indicate which empty cells should be filled to pass the structural
// checks. Structural checks include the 2x2 check, the 4 per row and col check, as
// well as the continuity check.
//
// We were given a board with 49 cells, 13 of which were filled. We know we need to have
// a total of 28 numbers at the end (7+6+5+4+3+2+1 = 28). So we need to fill in 15 more.
// We also start with 36 empty cells.
//
// This turns into a combinatorial problem, 36 pick 15 = ~5.5B possible board, without
// taking into account the other problem constraints. The tricky part here lies in how
// to iteratively generate the next combination, as the numbers are too large to do
// recursively without stack overflows.
//
// This was accomplished by taking a counter, and incrementing it by 1. Then taking the
// Hamming weight of that number and checking to see if it equals 15. If it does, then
// that is a valid "key" for the board. The lower 36 bits correspond to which of the 36
// cells should be filled. The loop exits when the counter reaches 1<<36.
//
// So for example, if we had a smaller board that had a key of 0b00100010 we would
// interpret that as setting the 2nd and 6th empty slots on the board.
//
// These get plugged in, then we check if the board structure is valid. Board structure is
// just checking for 2x2 correctness, 4 per row/col correctness, and continuity correctness.
//
// Once we have generated a list of valid keys, we need to do a second pass to generate
// The actual numbers in the correct locations.
//
// ~15M checks / sec with a ~6min runtime
func getValidBoardStructures() []uint64 {
	var keys []uint64
	empties := generateEmpties()

	var iterations int
	start := time.Now()

	var key uint64 = 0
	//var key uint64 = 7669990665
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

		// Find the next key where hammingWeight = 1
		hammingWeight = 0
		for hammingWeight != 15 {
			key++

			// This is the getHammingWeight64Fast function, but inlined
			x = key
			x -= (x >> 1) & m1              //put count of each 2 bits into those 2 bits
			x = (x & m2) + ((x >> 2) & m2)  //put count of each 4 bits into those 4 bits
			x = (x + (x >> 4)) & m4         //put count of each 8 bits into those 8 bits
			hammingWeight = (x * h01) >> 56 //returns left 8 bits of x + (x<<8) + (x<<16) + (x<<24) + ...
		}
		if key > 1<<36 {
			return keys
		}

		// This looks gross, but is quite simple. It is just an
		// unrolled loop performing the following:
		//
		// for i:=0;i <36; i++ {
		// 	 v := key & (1 << i)
		// 	 if v != 0 {
		// 	   rows.Set(empties[i].row,empties[i].col)
		// 	 }
		// }
		//
		// The neq(key&(1<<i)) is a branchless version of the if statement.
		// It will return 1 when not equal and 0 when they are equal.
		rows.SetV(empties[0].row, empties[0].col, neq(key&(1<<0), 0))
		rows.SetV(empties[1].row, empties[1].col, neq(key&(1<<1), 0))
		rows.SetV(empties[2].row, empties[2].col, neq(key&(1<<2), 0))
		rows.SetV(empties[3].row, empties[3].col, neq(key&(1<<3), 0))
		rows.SetV(empties[4].row, empties[4].col, neq(key&(1<<4), 0))
		rows.SetV(empties[5].row, empties[5].col, neq(key&(1<<5), 0))
		rows.SetV(empties[6].row, empties[6].col, neq(key&(1<<6), 0))
		rows.SetV(empties[7].row, empties[7].col, neq(key&(1<<7), 0))
		rows.SetV(empties[8].row, empties[8].col, neq(key&(1<<8), 0))
		rows.SetV(empties[9].row, empties[9].col, neq(key&(1<<9), 0))
		rows.SetV(empties[10].row, empties[10].col, neq(key&(1<<10), 0))
		rows.SetV(empties[11].row, empties[11].col, neq(key&(1<<11), 0))
		rows.SetV(empties[12].row, empties[12].col, neq(key&(1<<12), 0))
		rows.SetV(empties[13].row, empties[13].col, neq(key&(1<<13), 0))
		rows.SetV(empties[14].row, empties[14].col, neq(key&(1<<14), 0))
		rows.SetV(empties[15].row, empties[15].col, neq(key&(1<<15), 0))
		rows.SetV(empties[16].row, empties[16].col, neq(key&(1<<16), 0))
		rows.SetV(empties[17].row, empties[17].col, neq(key&(1<<17), 0))
		rows.SetV(empties[18].row, empties[18].col, neq(key&(1<<18), 0))
		rows.SetV(empties[19].row, empties[19].col, neq(key&(1<<19), 0))
		rows.SetV(empties[20].row, empties[20].col, neq(key&(1<<20), 0))
		rows.SetV(empties[21].row, empties[21].col, neq(key&(1<<21), 0))
		rows.SetV(empties[22].row, empties[22].col, neq(key&(1<<22), 0))
		rows.SetV(empties[23].row, empties[23].col, neq(key&(1<<23), 0))
		rows.SetV(empties[24].row, empties[24].col, neq(key&(1<<24), 0))
		rows.SetV(empties[25].row, empties[25].col, neq(key&(1<<25), 0))
		rows.SetV(empties[26].row, empties[26].col, neq(key&(1<<26), 0))
		rows.SetV(empties[27].row, empties[27].col, neq(key&(1<<27), 0))
		rows.SetV(empties[28].row, empties[28].col, neq(key&(1<<28), 0))
		rows.SetV(empties[29].row, empties[29].col, neq(key&(1<<29), 0))
		rows.SetV(empties[30].row, empties[30].col, neq(key&(1<<30), 0))
		rows.SetV(empties[31].row, empties[31].col, neq(key&(1<<31), 0))
		rows.SetV(empties[32].row, empties[32].col, neq(key&(1<<32), 0))
		rows.SetV(empties[33].row, empties[33].col, neq(key&(1<<33), 0))
		rows.SetV(empties[34].row, empties[34].col, neq(key&(1<<34), 0))
		rows.SetV(empties[35].row, empties[35].col, neq(key&(1<<35), 0))

		// The order of checks is important for performance. We have the fastest
		// checks first, and the slowest checks last. Approx runtimes:
		//   - Passes2x2: 3ns/op
		//   - Passes4Check: 6ns/op
		//   - Transpose + Passes4Check: 27+6ns/op = 33ns/op
		//   - PassesContinuity: 154ns/op + copy time
		if rows.Passes2x2() && rows.Passes4Check() && rows.Transpose().Passes4Check() {
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
				keys = append(keys, key)
				fmt.Printf("********** Found another valid board! key = %d, cnt = %d\n", key, validCount)
				validCount++
			}
		}

		iterations++
		if iterations%100000000 == 0 {
			fmt.Printf("Iterations: %dM after %s (%.2fM / sec)\n", iterations/1000000, time.Since(start), float64(iterations)/1000000/time.Since(start).Seconds())
			fmt.Printf("    Binary key: %64b %[1]d\n", key)
		}
	}
}

// validateBoardStructures loops through a given set of keys and looks for any
// boards that can be filled with a set of numbers that validate the remaining
// checks.
//
// It does this in two passes:
//  1. The "easy" fills - looking for rows or cols with 3 items
//  2. Looping through all remaining permutations
//
// that map to a board that does not pass the numeric checks. If it finds a
// valid board, it prints it out to count the last step manually.
func validateBoardStructures(keys []uint64) {
	valid := true
	var cp Board

	for _, key := range keys {
		b := NewBoardFromKey(key)
		state := b.FillEasy()

		if state == StateInvalid {
			continue
		}
		fmt.Printf("Key: %d, State: %s passed initial screening\n", key, state)

		// Generate a list of the remaining numbers and their positions
		remaining, empties := generateRemainingAndEmptyTuples(b)

		// Then generate all the possible permutations for those 10 numbers.
		// Step through each permutation and see if it's valid.
		quickPerm(remaining, func(nums []int) error {
			cp = b.Clone()
			valid = true

			// Nums will be a permutation of the remaining numbers.
			// We'll loop through these and plug them into the board,
			// while checking for validity.
			for i := 0; i < len(nums); i++ {
				ok := cp.set(int(empties[i].row), int(empties[i].col), nums[i])
				if !ok {
					valid = false
					break
				}
			}

			// If we found a valid permutation, exit!
			if valid && cp.passesSumChecks() {
				return errors.New("done")
			}

			return nil
		})

		// If this key has a valid combo, we don't have to look for anymore.
		if valid {
			break
		}
	}

	// The answer to this puzzle is the product of the areas of the
	// connected groups of empty squares in the completed grid.
	printGrid(cp.grid)
	fmt.Println("Found answer!")

	/*
		7 4 3 _ 6 _ _
		_ _ 6 3 5 _ 6
		_ _ 5 _ 5 5 5
		_ 3 6 4 _ _ 7
		4 7 _ _ _ 7 2
		2 _ _ 7 4 7 _
		7 6 _ 6 _ 1 _
	*/

	// Good 'ole counting leads us to the final answer here
	// 1*3*1*5**8*1*2=240
	// 49 boxes - 21 missing = 28 filled

	return
}

func generateRemainingAndEmptyTuples(b Board) ([]int, []tuple) {
	var remaining []int
	for i := 1; i <= 7; i++ {
		for j := 0; j < i-b.counts[i]; j++ {
			remaining = append(remaining, int(i))
		}
	}

	empties := make([]tuple, len(remaining))
	cnt := 0
	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			if b.grid[i][j] == 9 {
				empties[cnt] = tuple{uint8(i), uint8(j)}
				cnt++
			}
		}
	}

	return remaining, empties
}
