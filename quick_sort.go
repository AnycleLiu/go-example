package main

import (
	"fmt"
	"math/rand"
)

func main() {

	arr := []int{3, 1, 2, 9, 10, 100, 1, 1, 4, 6, 5, 0, 10, 20, 11, 30}
	quicksort(arr)

	fmt.Printf("%v", arr)
}

func quicksort(arr []int) {
	if arr == nil || len(arr) == 0 {
		return
	}

	bi := int(rand.Int() % len(arr)) //随机选一个基准
	b := arr[bi]
	si := 0 //挡板
	//基准交换到数组最后位置
	arr[bi], arr[len(arr)-1] = arr[len(arr)-1], arr[bi]

	for i := 0; i < len(arr)-1; i++ {
		if arr[i] <= b {
			//如果小于或者等于基准
			arr[i], arr[si] = arr[si], arr[i]
			si++
		}
	}
	//基准和si交换
	arr[len(arr)-1], arr[si] = arr[si], arr[len(arr)-1]

	quicksort(arr[0:si])
	quicksort(arr[si+1:])
}
