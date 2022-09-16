package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Row byte
type Col byte

var mask uint8 = 1<<7 - 1

func NewRowRepresentation(grid [][]int) RowCollection {
	var row uint64
	count := 0
	//printGrid(grid)
	for i := uint64(0); i < uint64(len(grid)); i++ {
		var r uint64
		for j := 0; j < 7; j++ {
			if grid[i][j] != 0 {
				r |= 1 << (6 - j)
				count++
			}
		}
		//fmt.Printf("row = %d %07b\n", i, r)
		row |= r << (7 * (6 - i))
	}

	return RowCollection{
		rows:  row,
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
	rows  uint64
	count int
}

func (r RowCollection) GetCountNumbers(i int) int {
	return hamming[uint8(r.rows>>(7*(6-i)))&mask]
}

func (r RowCollection) GetColCountNumbers(i int) int {
	count := 0

	for j := 0; j < 7; j++ {
		if r.rows>>(7*(6-j))&(1<<(6-i)) != 0 {
			count++
		}
	}

	return count
}

func gt(a, b uint8) uint8 {
	if a > b {
		return 1
	}

	return 0
}

func (r RowCollection) Passes4Check() bool {
	if hamming[uint8(r.rows>>(7*(6-0)))&mask] > 4 ||
		hamming[uint8(r.rows>>(7*(6-1)))&mask] > 4 ||
		hamming[uint8(r.rows>>(7*(6-2)))&mask] > 4 ||
		hamming[uint8(r.rows>>(7*(6-3)))&mask] > 4 ||
		hamming[uint8(r.rows>>(7*(6-4)))&mask] > 4 ||
		hamming[uint8(r.rows>>(7*(6-5)))&mask] > 4 ||
		hamming[uint8(r.rows>>(7*(6-6)))&mask] > 4 {
		return false
	}

	return true
}

// Passes2x2 is looking for any blocks of 2x2 where all the numbers are set.
// It does this by doing a bitwise AND on two adjacent rows, then checking
// two see if they have two adjacent bits set.
func (r RowCollection) Passes2x2() bool {
	a := uint8(r.rows>>(7*(6-0))) & mask & uint8(r.rows>>(7*(6-1)))
	b := uint8(r.rows>>(7*(6-1))) & mask & uint8(r.rows>>(7*(6-2)))
	c := uint8(r.rows>>(7*(6-2))) & mask & uint8(r.rows>>(7*(6-3)))
	d := uint8(r.rows>>(7*(6-3))) & mask & uint8(r.rows>>(7*(6-4)))
	e := uint8(r.rows>>(7*(6-4))) & mask & uint8(r.rows>>(7*(6-5)))
	f := uint8(r.rows>>(7*(6-5))) & mask & uint8(r.rows>>(7*(6-6)))

	return a&(a>>1)+b&(b>>1)+c&(c>>1)+d&(d>>1)+e&(e>>1)+f&(f>>1) == 0
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
	for i := 0; i < 7; i++ {
		v := uint8(r.rows>>(7*(6-i))) & mask
		out += fmt.Sprintf("%d: %07b (%d)\n", i, v, hamming[v])
	}
	out += "   "
	sum := 0
	for i := 0; i < 7; i++ {
		sum += r.GetCountNumbers(i)
		out += fmt.Sprintf("%d", r.GetColCountNumbers(i)) //@todo
	}
	out += " (" + strconv.Itoa(sum) + ")\n"
	out += "\n"

	return out
}

func (r RowCollection) Reset(grid [][]int) {
	r.count = 0
	r.rows = 0

	for i := 0; i < len(grid); i++ {
		var tmp uint64
		for j := 0; j < 7; j++ {
			if grid[i][j] != 0 {
				tmp |= 1 << (6 - i)
				r.count++
			}
		}
		r.rows |= tmp << (6 - i)
	}
}

func (r *RowCollection) Set(i uint8, j uint8) {
	r.rows |= 1 << (7*(6-i) + 6 - j)
	r.count++
}

func (r *RowCollection) SetV(i uint8, j uint8, v int8) {
	r.rows |= uint64(v) << (7*(6-i) + 6 - j)
	//r.rows[i] |= Row(v) << (6 - j)
	////r.count += eq(v, 1) // Removing the count gets us ~28M / sec,
	//// but the wrong answer (242 valid boards with panic) 242 without... no change
}

func (r RowCollection) Transpose() RowCollection {
	res := RowCollection{
		count: r.count,
		rows:  0,
	}

	var v uint8
	for i := 0; i < 7; i++ {
		v = uint8(r.rows>>(7*(6-i)+6)) & mask & 1
		res.rows |= (uint64(v) << (6 - i)) << (7 * (6 - 6))

		v = uint8(r.rows>>(7*(6-i)+5)) & mask & 1
		res.rows |= (uint64(v) << (6 - i)) << (7 * (6 - 5))

		v = uint8(r.rows>>(7*(6-i)+4)) & mask & 1
		res.rows |= (uint64(v) << (6 - i)) << (7 * (6 - 4))

		v = uint8(r.rows>>(7*(6-i)+3)) & mask & 1
		res.rows |= (uint64(v) << (6 - i)) << (7 * (6 - 3))

		v = uint8(r.rows>>(7*(6-i)+2)) & mask & 1
		res.rows |= (uint64(v) << (6 - i)) << (7 * (6 - 2))

		v = uint8(r.rows>>(7*(6-i)+1)) & mask & 1
		res.rows |= (uint64(v) << (6 - i)) << (7 * (6 - 1))

		v = uint8(r.rows>>(7*(6-i)+0)) & mask & 1
		res.rows |= (uint64(v) << (6 - i)) << (7 * (6 - 0))
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
