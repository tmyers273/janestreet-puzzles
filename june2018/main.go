// package main solves the puzzle!
//
// It works in a two pass approach. First, it generates a lsit of
// structurally possible boards. These are defined as boards that
// pass the 2x2 check, the 4 per row/col check, and the continuity
// check.
//
// Once we have a list of structurally possible boards, we then
// generate the actual numbers for each board.
package main

import (
	"errors"
	"fmt"
	"time"
)

// Notes:
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
// This was originally approached by taking a counter, and incrementing it by 1. Then taking the
// Hamming weight of that number and checking to see if it equals 15. If it does, then
// that is a valid "key" for the board. The lower 36 bits correspond to which of the 36
// cells should be filled. The loop exits when the counter reaches 1<<36.
//
// So for example, if we had a smaller board that had a key of 0b00100010 we would
// interpret that as setting the 2nd and 6th empty slots on the board.
//
// However, the bit twiddling below is ~2x as fast in benchmarks, so the Hamming weight
// based approach was scraped.
// https://graphics.stanford.edu/~seander/bithacks.html#NextBitPermutation
//
// These get plugged in, then we check if the board structure is valid. Board structure is
// just checking for 2x2 correctness, 4 per row/col correctness, and continuity correctness.
//
// Once we have generated a list of valid keys, we need to do a second pass to generate
// The actual numbers in the correct locations.
//
// ~63M checks / sec with a ~1.5min runtime
func getValidBoardStructures() []uint64 {
	var keys []uint64
	var iterations int
	start := time.Now()

	var key uint64 = 1<<15 - 1
	var rows RowCollection
	validCount := 0
	rows = NewRowRepresentation(oGrid)
	orig := NewRowRepresentation(oGrid)
	grid := make([][]int, 7)
	for i := 0; i < 7; i++ {
		grid[i] = make([]int, 7)
	}

	var x uint64

	for {
		rows.rows = orig.rows
		rows.count = orig.count

		// Find the lexicographically next combination that has 15 bits
		// ie 0 01111111 11111111
		// -> 0 10111111 11111111
		// -> 0 11011111 11111111
		x = (key | (key - 1)) + 1
		key = x | ((((x & -x) / (key & -key)) >> 1) - 1)
		if key > 1<<36 {
			return keys
		}

		// Convert that key into a mask, then OR it on the
		// current rows to set all the filled cells.
		rows.rows |= keyToMask(key)
		rows.count += 15

		// The order of checks is important for performance. We have the fastest
		// checks first, and the slowest checks last. Approx runtimes:
		//   - Passes2x2: 3ns/op
		//   - Passes4Check: 4ns/op
		//   - Transpose + Passes4Check: 27+4ns/op = 31ns/op
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
			//return nil // @todo just for dev / bench
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
	if len(keys) == 0 {
		panic("uhh oh - didn't get any valid boards")
	}

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
