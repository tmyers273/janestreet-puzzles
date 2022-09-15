package main

import (
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

// 4.149 ns/op
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

// ~126M / sec
func BenchmarkHammingSpeed(b *testing.B) {
	var key uint64 = 1<<15 - 1
	var hammingWeight uint64
	//var x uint64

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hammingWeight = 0
		for hammingWeight != 15 {
			key++

			hammingWeight = key
			hammingWeight -= (hammingWeight >> 1) & m1                         //put count of each 2 bits into those 2 bits
			hammingWeight = (hammingWeight & m2) + ((hammingWeight >> 2) & m2) //put count of each 4 bits into those 4 bits
			hammingWeight = (hammingWeight + (hammingWeight >> 4)) & m4        //put count of each 8 bits into those 8 bits
			hammingWeight = (hammingWeight * h01) >> 56                        //returns left 8 bits of x + (x<<8) + (x<<16) + (x<<24) + ...
		}
	}
	b.StopTimer()
	dur := time.Since(start)

	b.ReportMetric(float64(b.N)/1000000/dur.Seconds(), "M/sec")
}

// ~204M / sec
func BenchmarkLexicographicallyNext(b *testing.B) {
	var key uint64 = 1<<15 - 1
	var t uint64

	start := time.Now()
	b.ResetTimer()
	//fmt.Printf("%64b\n", key)
	for i := 0; i < b.N; i++ {
		t = (key | (key - 1)) + 1
		key = t | ((((t & -t) / (key & -key)) >> 1) - 1)
		//fmt.Printf("%64b\n", key)
		//if i > 1000 {
		//	break
		//}
	}
	b.StopTimer()
	dur := time.Since(start)

	b.ReportMetric(float64(b.N)/1000000/dur.Seconds(), "M/sec")
}

func Benchmark_getValidBoardStructures(b *testing.B) {
	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getValidBoardStructures()
	}
	b.StopTimer()
	dur := time.Since(start)

	b.ReportMetric(float64(b.N*100000000)/1000000/dur.Seconds(), "Mops/sec")
}
