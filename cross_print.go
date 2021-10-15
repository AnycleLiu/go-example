package main

import (
	"fmt"
	"sync"
)

func main() {
	nc, cn := make(chan int), make(chan int)
	w := sync.WaitGroup{}

	w.Add(1)
	go func(w *sync.WaitGroup) {
		defer w.Done()

		for c := 'A'; c <= 'Z'; {
			<-nc
			fmt.Printf("%c", c)
			c++
			if c > 'Z' {
				close(cn)
				break
			}
			fmt.Printf("%c", c)
			c++
			if c > 'Z' {
				close(cn)
				break
			} else {
				cn <- 1
			}
		}
	}(&w)

	w.Add(1)
	go func(w *sync.WaitGroup) {
		defer w.Done()
		for c := 1; ; {
			if _, ok := <-cn; !ok {
				break
			}
			fmt.Printf("%d", c)
			c++
			fmt.Printf("%d", c)
			c++
			nc <- 1
		}
	}(&w)

	cn <- 1

	w.Wait()
}
