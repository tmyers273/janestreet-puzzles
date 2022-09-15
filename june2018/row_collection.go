package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Row byte
type Col byte

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

type RowCollection struct {
	rows  []Row
	count int
}

func (r RowCollection) GetCountNumbers(i int) int {
	return hamming[r.rows[i]]
}

func (r RowCollection) GetColCountNumbers(i int) int {
	count := 0

	for j := 0; j < 7; j++ {
		if r.rows[j]&(1<<(6-i)) != 0 {
			count++
		}
	}

	return count
}

func (r RowCollection) Passes4Check() bool {
	for i := 0; i < len(r.rows); i++ {
		if hamming[r.rows[i]] > 4 {
			return false
		}
	}

	return true
}

// Passes2x2 is looking for any blocks of 2x2 where all the numbers are set.
// It does this by doing a bitwise AND on two adjacent rows, then checking
// two see if they have two adjacent bits set.
func (r RowCollection) Passes2x2() bool {
	var sum uint8
	rows := r.rows

	sum += (uint8(rows[0]) & uint8(rows[1])) & ((uint8(rows[0]) & uint8(rows[1])) >> 1)
	sum += (uint8(rows[1]) & uint8(rows[2])) & ((uint8(rows[1]) & uint8(rows[2])) >> 1)
	sum += (uint8(rows[2]) & uint8(rows[3])) & ((uint8(rows[2]) & uint8(rows[3])) >> 1)
	sum += (uint8(rows[3]) & uint8(rows[4])) & ((uint8(rows[3]) & uint8(rows[4])) >> 1)
	sum += (uint8(rows[4]) & uint8(rows[5])) & ((uint8(rows[4]) & uint8(rows[5])) >> 1)
	sum += (uint8(rows[5]) & uint8(rows[6])) & ((uint8(rows[5]) & uint8(rows[6])) >> 1)

	return sum == 0
}

func (r Row) String() string {
	return strings.Replace(fmt.Sprintf("%07b", r), "0", ".", -1)
}

func (r Row) ToSlice() []int {
	res := make([]int, 7)
	for i := 0; i < 7; i++ {
		res[i] = int(r >> (6 - i) & 1)
	}

	return res
}

func (r RowCollection) PassesContinuity(grid [][]int) bool {
	visited := make([]uint8, 49)

	// Find the first non-zero cell, then recursively count the number
	// of 4-way connected cells. The result of this should match the
	// total number of non-zero cells if the entire grid is connected.
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid); j++ {
			if grid[i][j] != 0 {
				cnt := r.getChunkCount(grid, i, j, visited)
				return cnt == r.count
			}
		}
	}

	return false
}

func (r RowCollection) String() string {
	out := "\n"
	out += "   0123456\n"
	for i := 0; i < len(r.rows); i++ {
		out += fmt.Sprintf("%d: %s (%d)\n", i, r.rows[i], hamming[r.rows[i]])
	}
	out += "   "
	sum := 0
	for i := 0; i < 7; i++ {
		sum += r.GetCountNumbers(i)
		out += fmt.Sprintf("%d", r.GetColCountNumbers(i))
	}
	out += " (" + strconv.Itoa(sum) + ")\n"
	out += "\n"

	return out
}

func (r RowCollection) Reset(grid [][]int) {
	r.count = 0

	for i := 0; i < len(grid); i++ {
		r.rows[i] = 0

		for j := 0; j < len(grid); j++ {
			if grid[i][j] != 0 {
				r.rows[i] |= 1 << (6 - i)
				r.count++
			}
		}
	}
}

func (r *RowCollection) Set(i uint8, j uint8) {
	r.rows[i] |= 1 << (6 - j)
	r.count++
}

func (r *RowCollection) SetV(i uint8, j uint8, v int8) {
	r.rows[i] |= Row(v) << (6 - j)
	r.count += eq(v, 1)
}

func (r RowCollection) Transpose() RowCollection {
	res := RowCollection{
		count: r.count,
		rows:  make([]Row, 7),
	}

	var v Row
	for i := 0; i < 7; i++ {
		v = r.rows[i] >> (6) & 1
		res.rows[6] |= v << (6 - i)

		v = r.rows[i] >> (5) & 1
		res.rows[5] |= v << (6 - i)

		v = r.rows[i] >> (4) & 1
		res.rows[4] |= v << (6 - i)

		v = r.rows[i] >> (3) & 1
		res.rows[3] |= v << (6 - i)

		v = r.rows[i] >> (2) & 1
		res.rows[2] |= v << (6 - i)

		v = r.rows[i] >> (1) & 1
		res.rows[1] |= v << (6 - i)

		v = r.rows[i] >> (0) & 1
		res.rows[0] |= v << (6 - i)
	}

	return res
}

func (r RowCollection) getChunkCount(grid [][]int, i, j int, visited []uint8) int {
	if visited[i*7+j] == 1 {
		return 0
	}

	count := 1
	visited[i*7+j] = 1

	// Top
	if j > 0 && grid[i][j-1] != 0 && visited[i*7+j-1] == 0 {
		count += r.getChunkCount(grid, i, j-1, visited)
	}

	// Left
	if i > 0 && grid[i-1][j] != 0 && visited[(i-1)*7+j] == 0 {
		count += r.getChunkCount(grid, i-1, j, visited)
	}

	// Right
	if j < 6 && grid[i][j+1] != 0 && visited[i*7+j+1] == 0 {
		count += r.getChunkCount(grid, i, j+1, visited)
	}

	// Bottom
	if i < 6 && grid[i+1][j] != 0 && visited[(i+1)*7+j] == 0 {
		count += r.getChunkCount(grid, i+1, j, visited)
	}

	return count
}
