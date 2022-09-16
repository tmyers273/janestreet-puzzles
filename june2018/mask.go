package main

var maskList [][]uint64
var empties []tuple

// Here we will generate a lookup table mapping all the key bytes to
// possible set bits in the final mask. We will use this to generate
// the mask for a given key by looping through the keys bytes and
// OR'ing the results. Thing final mask can then be OR'ed on the
// rows to generate all the set bits.
func init() {
	empties = generateEmpties()
	emptyOffsets := make([]uint8, len(empties))
	for i := 0; i < len(empties); i++ {
		emptyOffsets[i] = 7*(6-empties[i].row) + 6 - empties[i].col
	}

	maskList = make([][]uint64, 8)
	for i := 0; i < 8; i++ {
		maskList[i] = make([]uint64, 256)
	}

	for b := 0; b < 8; b++ {
		for k := 0; k < 256; k++ {
			for j := 7; j >= 0; j-- {
				originalIndex := 8*b + j
				if originalIndex >= len(emptyOffsets) {
					continue
				}

				if (k)&(1<<j) != 0 {
					maskList[b][k] |= 1 << emptyOffsets[originalIndex]
				}
			}
		}
	}
}

// keyToMask takes a key and returns the mask for that key.
// This mask can then be OR'ed on the current rows to set
// all the filled slots.
func keyToMask(key uint64) uint64 {
	var final uint64

	final |= maskList[0][byte(key>>(8*0))] // 0-8
	final |= maskList[1][byte(key>>(8*1))] // 8-16
	final |= maskList[2][byte(key>>(8*2))] // 16-24
	final |= maskList[3][byte(key>>(8*3))] // 24-32
	final |= maskList[4][byte(key>>(8*4))] // 32-40
	final |= maskList[5][byte(key>>(8*5))] // 40-48
	final |= maskList[6][byte(key>>(8*6))] // 48-56

	// Last byte can be skipped, since we will only
	// ever have 49 cells.
	//final |= maskList[7][byte(key>>(8*7))]

	return final
}
