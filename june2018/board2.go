package june2018

func NewBoard2FromKey(key uint64) Board2 {
	grid := NewBoardFromKey(key)
	rows := NewRowRepresentation(oGrid)

	counts := make([]int, 8)
	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			if grid[i][j] != 0 && grid[i][j] != 9 {
				counts[grid[i][j]]++
			}
		}
	}

	return Board2{
		grid:   grid,
		rows:   rows,
		counts: counts,
	}
}

type Board2 struct {
	grid   [][]int
	rows   RowCollection
	counts []int
}

func (b *Board2) Clone() Board2 {
	newGrid := make([][]int, 7)
	for i := 0; i < 7; i++ {
		newGrid[i] = make([]int, 7)
		copy(newGrid[i], b.grid[i])
	}

	newRows := RowCollection{
		rows: make([]Row, 7),
	}
	copy(newRows.rows, b.rows.rows)
	newRows.count = b.rows.count

	newCounts := make([]int, 8)
	copy(newCounts, b.counts)

	return Board2{
		grid:   newGrid,
		rows:   newRows,
		counts: newCounts,
	}
}

func (b *Board2) FillEasy() State {
	madeChanges := true
	changeCount := 0

	// Keep looping until we can't make any more "easy" fills
	for madeChanges {
		madeChanges = false

		// Loop through rows, looking for any with 3 numbers.
		//
		// If the sum of the 3 existing numbers is < 12, then
		// we know the board is invalid as it would require
		// an 8 to complete.
		//
		// If the sum is >= 13, then we know we can definitively
		// set the valid for that cell to be 20 minus the sum.
		for i := 0; i < 7; i++ {
			if b.rows.GetCountNumbers(i) == 3 {
				sum := 20
				missingIndex := 0
				for j := 0; j < 7; j++ {
					if b.grid[i][j] == 9 {
						missingIndex = j
					} else {
						sum -= b.grid[i][j]
					}
				}

				if sum > 7 {
					return StateInvalid
				}

				if !b.set(i, missingIndex, sum) {
					return StateInvalid
				}
				changeCount++
				madeChanges = true
			}
		}

		// Same thing, but looping through the columns
		for i := 0; i < 7; i++ {
			if b.rows.GetColCountNumbers(i) == 3 {
				sum := 20
				missingIndex := 0
				for j := 0; j < 7; j++ {
					if b.grid[j][i] == 9 {
						missingIndex = j
					} else {
						sum -= b.grid[j][i]
					}
				}

				if sum > 7 {
					return StateInvalid
				}

				if !b.set(missingIndex, i, sum) {
					return StateInvalid
				}
				changeCount++
				madeChanges = true
			}
		}
	}

	// Finally, if we have all 28 numbers, then the board is valid
	if b.rows.count == 28 {
		return StateValid
	}

	return StateUnknown
}

// set is a helper function. It sets the underlying grid's
// value, updates the bitwise row representation, and updates
// the counts for the set value.
//
// Returns false if the new value causes the count check to fail.
func (b *Board2) set(row, col, value int) bool {
	b.grid[row][col] = value
	b.rows.Set(uint8(row), uint8(col))
	b.counts[value]++

	return b.counts[value] <= value
}

func (b *Board2) passesSumChecks() bool {
	for i := 0; i < 7; i++ {
		if !b.passesRowSum(i) || !b.passesColSum(i) {
			return false
		}
	}

	return true
}

func (b *Board2) passesSumCheck(i, j int) bool {
	return b.passesRowSum(i) && b.passesColSum(j)
}

func (b *Board2) passesRowSum(i int) bool {
	sum := 0
	for j := 0; j < 7; j++ {
		sum += b.grid[i][j]
	}

	return sum == 20
}

func (b *Board2) passesColSum(i int) bool {
	sum := 0
	for j := 0; j < 7; j++ {
		sum += b.grid[j][i]
	}

	return sum == 20
}
