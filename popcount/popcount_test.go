package popcount

import (
	"math/bits"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func TestPopCount(t *testing.T) {
	var i uint64
	testingFuncs := []struct {
		name     string
		function func(uint64) int
	}{
		{"PopCount", PopCount},
		{"PopCountLoop", PopCountLoop},
		{"PopCountShift", PopCountShift},
		{"PopCountCleanup", PopCountCleanup},
	}
	for i = 0; i < 2^64; i++ {
		pc0 := bits.OnesCount64(i)
		for _, tf := range testingFuncs {
			pc := tf.function(i)
			if pc0 != pc {
				t.Errorf("%s returned a wrong result for i:%d: %d, want: %d",
					tf.name, i, pc, pc0)
			}
		}
	}
}

func BenchmarkPopCountLib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bits.OnesCount64(rand.Uint64())
	}
}

func BenchmarkPopCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCount(rand.Uint64())
	}
}

func BenchmarkPopCountLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountLoop(rand.Uint64())
	}
}

func BenchmarkPopCountShift(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountShift(rand.Uint64())
	}
}

func BenchmarkPopCountCleanup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountCleanup(rand.Uint64())
	}
}
