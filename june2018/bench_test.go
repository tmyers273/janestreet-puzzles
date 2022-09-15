package june2018

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// 26.52 ns/op
func BenchmarkRowCollection_Transpose(b *testing.B) {
	r := NewRowRepresentation(oGrid)

	var res RowCollection
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res = r.Transpose()
	}
	b.StopTimer()

	expected := NewRowRepresentation(oGridRotatedRight)
	assert.Equal(b, expected, res)
}

// 154 ns/op
func BenchmarkPassesContinuity(b *testing.B) {
	rows1 := NewRowRepresentation(oGrid)
	rows2 := NewRowRepresentation(contiguous1)

	for i := 0; i < b.N; i++ {
		rows1.PassesContinuity(oGrid)
		rows2.PassesContinuity(contiguous1)
	}
}

// 6.043 ns/op
func BenchmarkPasses4Check(b *testing.B) {
	rows1 := NewRowRepresentation(oGrid)
	rows2 := NewRowRepresentation(contiguous1)

	for i := 0; i < b.N; i++ {
		rows1.Passes4Check()
		rows2.Passes4Check()
	}
}

// 2.921 ns/op
func BenchmarkPasses2x2(b *testing.B) {
	oGrid[5][1] = 1
	rows := NewRowRepresentation(oGrid)

	var r bool
	for i := 0; i < b.N; i++ {
		r = rows.Passes2x2()
	}
	assert.Equal(b, false, r)
}

// ~15M / sec
func BenchmarkHammingSpeed(b *testing.B) {
	var key uint64
	var hammingWeight uint64

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for hammingWeight != 15 {
			key++

			var x uint64
			x = key
			x -= (x >> 1) & m1              //put count of each 2 bits into those 2 bits
			x = (x & m2) + ((x >> 2) & m2)  //put count of each 4 bits into those 4 bits
			x = (x + (x >> 4)) & m4         //put count of each 8 bits into those 8 bits
			hammingWeight = (x * h01) >> 56 //returns left 8 bits of x + (x<<8) + (x<<16) + (x<<24) + ...
			//hammingWeight = getHammingWeight64Fast(key)
		}
	}
	b.StopTimer()
	dur := time.Since(start)

	fmt.Printf("Took %s (%.2fM / sec)\n", dur, float64(b.N)/1000000/dur.Seconds())
}
