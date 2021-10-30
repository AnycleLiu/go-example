package main

import (
	"fmt"
)

func main() {

	s := []byte{'a', 'b', ' ', ' ', 'c', ' ', 'd', 'e', 'f', ' '}

	n, r := trimspace(s)

	fmt.Printf("%d %s", n, string(r))
}

func trimspace(s []byte) (int, []byte) {
	n := 0
	i := 0
	j := 0
	c := byte(' ')
	for i < len(s)-n {
		if s[i] != c {
			i++
		} else {
			j = i + 1
			for j < len(s)-n && s[j] == c {
				j++
			}

			for k := j; k < len(s)-n; k++ {
				s[k-(j-i)] = s[k]
			}
			n = n + j - i
			i++
		}
	}
	return n, s[:len(s)-n]
}
