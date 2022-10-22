package main

import "fmt"

func main() {
	s := []int{1, 2, 4, 4, 6, 7, 10}
	fmt.Println("search1(s, 11) ", search1(s, 11))
	fmt.Println("search2(s, 11) ", search2(s, 11))
	fmt.Println("search1(s, 4) ", search1(s, 4))
	fmt.Println("search2(s, 4) ", search2(s, 4))
	fmt.Println("search1(s, 10) ", search1(s, 10))
	fmt.Println("search2(s, 10) ", search2(s, 10))
	fmt.Println("search2(s, 0) ", search2(s, 0))
}

// 搜索>=n的最小值
func search1(s []int, n int) int {
	l, r := 0, len(s)-1

	for l <= r {
		mi := l + (r-l+1)/2
		if s[mi] >= n {
			r = mi - 1
		} else {
			l = mi + 1
		}
	}
	if l < len(s) {
		return l
	}
	return -1
}

// 搜索<=n的最大值
func search2(s []int, n int) int {
	l, r := 0, len(s)-1

	for l <= r {
		mi := l + (r-l+1)/2
		if s[mi] <= n {
			l = mi + 1
		} else {
			r = mi - 1
		}
	}
	if r > 0 {
		return l - 1
	}
	return -1
}
