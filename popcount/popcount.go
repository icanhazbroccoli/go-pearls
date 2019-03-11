package popcount

var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func PopCount(x uint64) int {
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

func PopCountLoop(x uint64) int {
	var res byte
	var i uint
	for i = 0; i < 8; i++ {
		res += pc[byte(x>>(i*8))]
	}
	return int(res)
}

func PopCountShift(x uint64) int {
	var res int
	x0 := x
	for i := 0; i < 64; i++ {
		res += int(x0 & 1)
		x0 = x0 >> 1
	}
	return res
}

func PopCountCleanup(x uint64) int {
	x0 := x
	var res int
	for x0 > 0 {
		res++
		x0 &= (x0 - 1)
	}
	return res
}
