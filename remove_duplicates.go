package main

import "fmt"

func main() {
	arr := []int{2, 0, 0, 4, 5, 6, 6, 6, 6, 10, 9, 8}

	fmt.Printf("remove dumplicates: %v\n", removeDumplicates(arr))

	arr = []int{2, 0, 0, 4, 5, 6, 6, 6, 6, 10, 9, 8}
	fmt.Printf("remove 6: %v\n", remove(arr, 6))
	arr = []int{2, 0, 0, 4, 5, 6, 6, 6, 6, 10, 9, 8}
	fmt.Printf("remove 4: %v\n", remove(arr, 4))
}

//原地移除重复元素，返回移除后的切片
func removeDumplicates(arr []int) []int {
	slow, fast := 0, 0

	for fast < len(arr) {
		if arr[fast] != arr[slow] {
			slow++
			arr[slow] = arr[fast]
		}
		fast++
	}

	return arr[:slow+1]
}

func remove(arr []int, target int) []int {
	slow, fast := 0, 0

	for fast < len(arr) {
		if arr[fast] != target {
			arr[slow] = arr[fast]
			slow++
		}
		fast++
	}
	return arr[:slow]
}
