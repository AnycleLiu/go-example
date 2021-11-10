package main

import "fmt"

func main() {
	arr := []int{3, 1, 2, 9, 10, 100, 1, 1, 4, 6, 5, 0}
	bubblesort(arr)

	fmt.Printf("%v", arr)
}

func bubblesort(arr []int) {
	for i := 0; i < len(arr)-1; i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}
