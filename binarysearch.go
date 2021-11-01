package main

import (
	"fmt"
)

func main() {
	arr := []int{1, 3, 6, 7, 10, 11, 13, 16, 20, 22, 30, 50, 100, 1000}

	fmt.Println(Search(arr, 4))
	fmt.Println(Search(arr, 1000))
	fmt.Println(Search(arr, 1))

	fmt.Println("great search 999", greatsearch(arr, 0, len(arr)-1, 999))
	fmt.Println("great search 1000", greatsearch(arr, 0, len(arr)-1, 1000))
	fmt.Println("great search 1", greatsearch(arr, 0, len(arr)-1, 1))
}

func Search(arr []int, target int) int {
	return binarysearch(arr, 0, len(arr)-1, target)
}

func binarysearch(arr []int, l, r, t int) int {
	if arr == nil || l > r {
		return -1
	}

	mi := l + (r-l)/2
	m := arr[mi]

	if m == t {
		return mi
	} else if m > t {
		return binarysearch(arr, l, mi-1, t)
	} else {
		return binarysearch(arr, mi+1, r, t)
	}
}

func greatsearch(arr []int, l, r, t int) int {
	if arr == nil || l > r {
		return -1
	}

	mi := l + (r-l)/2
	m := arr[mi]

	if m > t {
		li := greatsearch(arr, l, mi-1, t)
		if li != -1 {
			return li
		}
		return mi
	}

	return greatsearch(arr, mi+1, r, t)
}
