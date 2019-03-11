package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(42)
	t0 := time.Now()
	for i := 0; i < 10000000; i++ {
		bits.OnesCount(rand.Intn(10000))
	}
	t1 := time.Now()
	elapsed := t1.Sub(t0)
	fmt.Printf("Benchmark time: %d\n", elapsed)
}
