package main

import (
	"fmt"
	"sync"
)

func main() {
	var m sync.Mutex
	var wg sync.WaitGroup
	n := 50

	wg.Add(2)
	go func() {
		fmt.Println("每次拿2个")
		defer wg.Done()
		//2个每次
		for {
			m.Lock()
			if n < 2 {
				m.Unlock()
				break
			} else {
				n -= 2
				fmt.Println("2")
				m.Unlock()
			}
		}
	}()

	go func() {
		fmt.Println("每次拿3个")
		defer wg.Done()
		//3个每次
		for {
			m.Lock()
			if n < 3 {
				m.Unlock()
				break
			} else {
				n -= 3
				fmt.Println("3")
				m.Unlock()
			}
		}
	}()

	wg.Wait()
	fmt.Println(n)
}
