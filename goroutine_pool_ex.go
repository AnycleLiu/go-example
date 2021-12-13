package main

import (
	"fmt"
	"go-example/goroutinepool"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	pool := goroutinepool.NewGoroutinePoolWithDefault()
	var n int32
	var wg sync.WaitGroup

	for i := 0; i < 100000; i++ {
		wg.Add(1)
		pool.Put(func() {
			time.Sleep(time.Millisecond * 200)
			atomic.AddInt32(&n, 1)
			wg.Done()
		})
	}

	wg.Wait()
	fmt.Println("pool协程数量", pool.Running())
	fmt.Println(n)
}
