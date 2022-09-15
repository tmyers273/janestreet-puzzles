package main

// hamming is a lookup table of the first 127 hamming weights.
// Hamming weights are the number of bits set to 1 in a given number.
// This gets used to quickly check the number of items in a row or col.
var hamming []int

func init() {
	max := 1<<7 - 1
	hamming = make([]int, max+1)

	for i := 0; i <= max; i++ {
		hamming[i] = getHammingWeight(uint8(i))
	}
}

// getHammingWeight returns the number of bits set to 1 in a given number.
// A pretty naive algorithm, but it's fast enough for our purposes since
// we are only running it 127 times
func getHammingWeight(i uint8) int {
	count := 0
	for count = 0; i > 0; count++ {
		i &= i - 1
	}
	return count
}

// Pulled from https://en.wikipedia.org/wiki/Hamming_weight
// And inlined in FirstPass
const m1 uint64 = 0x5555555555555555  //binary: 0101...
const m2 uint64 = 0x3333333333333333  //binary: 00110011..
const m4 uint64 = 0x0f0f0f0f0f0f0f0f  //binary:  4 zeros,  4 ones ...
const m8 uint64 = 0x00ff00ff00ff00ff  //binary:  8 zeros,  8 ones ...
const m16 uint64 = 0x0000ffff0000ffff //binary: 16 zeros, 16 ones ...
const m32 uint64 = 0x00000000ffffffff //binary: 32 zeros, 32 ones
const h01 uint64 = 0x0101010101010101 //the sum of 256 to the power of 0,1,2,3...
func getHammingWeight64Fast(x uint64) int {
	x -= (x >> 1) & m1             //put count of each 2 bits into those 2 bits
	x = (x & m2) + ((x >> 2) & m2) //put count of each 4 bits into those 4 bits
	x = (x + (x >> 4)) & m4        //put count of each 8 bits into those 8 bits
	return int((x * h01) >> 56)    //returns left 8 bits of x + (x<<8) + (x<<16) + (x<<24) + ...
}
