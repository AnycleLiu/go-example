package main

import (
	"fmt"
)

func main() {
	arr := []int{1, 3, 6, 7, 10, 11, 13, 16, 20, 22, 30, 50, 100, 1000}

	fmt.Printf("数据：%v\n", arr)

	fmt.Println("查找4", BinarySearch(arr, 4))
	fmt.Println("查找1000", BinarySearch(arr, 1000))
	fmt.Println("查找1", BinarySearch(arr, 1))

	fmt.Println("great search 999", GreatSearch(arr, 999))
	fmt.Println("great search 1000", GreatSearch(arr, 1000))
	fmt.Println("great search 1", GreatSearch(arr, 1))
	fmt.Println("great search 10", GreatSearch(arr, 10))
	fmt.Println("great search 12", GreatSearch(arr, 12))

	arr = []int{1, 3, 3, 3, 4, 5, 5}

	fmt.Printf("arr: %v\n", arr)
	fmt.Println("search 3 first idx", SearchFirst(arr, 3))
	fmt.Println("search 1 first idx", SearchFirst(arr, 1))
	fmt.Println("search 5 first idx", SearchFirst(arr, 5))
	fmt.Println("search 0 first idx", SearchFirst(arr, 0))

	fmt.Println("search 3 latest idx", SearchLatest(arr, 3))
	fmt.Println("search 1 latest idx", SearchLatest(arr, 1))
	fmt.Println("search 5 latest idx", SearchLatest(arr, 5))
	fmt.Println("search 0 latest idx", SearchLatest(arr, 0))
}

/*
//二分查找递归写法
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
}*/

//二分查找
func BinarySearch(arr []int, target int) int {
	l, r := 0, len(arr)-1

	for l <= r {
		mid := l + (r-l)/2
		if arr[mid] == target {
			return mid
		} else if arr[mid] > target {
			r = mid - 1
		} else if arr[mid] < target {
			l = mid + 1
		}
	}

	return -1
}

//在有序数组arr中查找第一个等于target的索引
func SearchFirst(arr []int, target int) int {
	l, r := 0, len(arr)-1

	for l <= r {
		mid := l + (r-l)/2
		if arr[mid] == target {
			r = r - 1
		} else if arr[mid] > target {
			r = r - 1
		} else if arr[mid] < target {
			l = l + 1
		}
	}

	if l == len(arr) || arr[l] != target {
		return -1
	}
	return l
}

//在有序数组arr中查找最后一个等于target的索引
func SearchLatest(arr []int, target int) int {
	l, r := 0, len(arr)-1

	for l <= r {
		mid := l + (r-l)/2
		if arr[mid] == target {
			l = l + 1
		} else if arr[mid] > target {
			r = r - 1
		} else if arr[mid] < target {
			l = l + 1
		}
	}
	if r < 0 || arr[r] != target {
		return -1
	}
	return r
}

//在有序数组中查找第一个大于target的值的索引
func GreatSearch(arr []int, target int) int {
	l, r := 0, len(arr)-1

	for l <= r {
		mid := l + (r-l)/2
		if arr[mid] > target {
			r = mid - 1
		} else if arr[mid] <= target {
			l = mid + 1
		}
	}
	if l >= len(arr) || arr[l] <= target {
		return -1
	}
	return l
}

/*
递归写法
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
*/
