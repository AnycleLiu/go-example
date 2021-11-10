package main

import "fmt"

func main() {
	a, b, c := 3, 100, 2

	fmt.Println(middle(a, b, c))
}

func middle(a, b, c int) int {
	if a > b {
		a, b = b, a
	}

	if b > c {
		b, c = c, b
	}

	if a > b {
		a, b = b, a
	}

	return b
}
