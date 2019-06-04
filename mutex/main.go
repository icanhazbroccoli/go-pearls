package main

import (
	"fmt"
	"sync"
)

type mutex struct {
	ch chan struct{}
}

func NewMutex() *mutex {
	return &mutex{
		ch: make(chan struct{}, 1),
	}
}

func (mx *mutex) Lock() {
	mx.ch <- struct{}{}
}

func (mx *mutex) Unlock() {
	<-mx.ch
}

type counter struct {
	mx  *mutex
	cnt int
}

func (cnt *counter) Inc() {
	// Intentional race condition
	v := cnt.cnt
	v += 1
	cnt.cnt = v
}

func (cnt *counter) SafeInc() {
	cnt.mx.Lock()
	defer cnt.mx.Unlock()
	v := cnt.cnt
	v += 1
	cnt.cnt = v
}

func (cnt *counter) Res() int {
	return cnt.cnt
}

func ExecCountLoop(safe bool) int {
	wg := sync.WaitGroup{}
	cnt := counter{}
	incFn := cnt.Inc
	if safe {
		cnt.mx = NewMutex()
		incFn = cnt.SafeInc
		defer close(cnt.mx.ch)
	}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			incFn()
			wg.Done()
		}()
	}

	wg.Wait()

	return cnt.Res()
}

func main() {
	unsafe, safe := ExecCountLoop(false), ExecCountLoop(true)
	fmt.Printf("Execution result: unsafe: %d, safe: %d\n", unsafe, safe)
}
