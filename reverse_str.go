package main

import "fmt"

func reverse(s string) (string, bool) {
	if len(s) > 5000 {
		return s, false
	}
	cs := []rune(s)
	for i, j := 0, len(cs)-1; i < j; i, j = i+1, j-1 {
		cs[i], cs[j] = cs[j], cs[i]
	}

	return string(cs), true
}

func main() {
	s := "hello world. 你好，中国."
	rs, ok := reverse(s)
	if !ok {
		fmt.Println("reverse fail. string limit 5000")
	} else {
		fmt.Println(rs)
	}
}
