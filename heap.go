package main

import (
	"fmt"
	"sort"
)

func main() {
	data := []int{100, 302, 10, 20, 0, 1, 3, 2, 6, 4, 0, 8, 10, 21, 45, 5, 1000, 101, 102, 111}
	h := &IntHeap{}

	Push(h, 0)
	for _, n := range data {
		Push(h, n)
	}

	for h.Len() > 0 {
		fmt.Println(Pop(h))
	}

}

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}
func (h *IntHeap) Pop() interface{} {
	l := len(*h)
	old := (*h)[l-1]
	*h = (*h)[:l-1]
	return old
}

type heap interface {
	sort.Interface
	Push(x interface{}) // add x as element Len()
	Pop() interface{}   // remove and return element Len() - 1.
}

func down(hp heap, i int) {
	now := i        //根
	nxt := now << 1 //左子节点
	l := hp.Len()

	for nxt <= l {
		if nxt != l && hp.Less(nxt, nxt+1) { //nxt为子节点较大的
			nxt = nxt + 1
		}

		if !hp.Less(nxt, now) { //交换直到当前节点大于子节点
			break
		}
		hp.Swap(nxt, now)
		now, nxt = nxt, nxt<<1
	}
}

func Push(hp heap, item interface{}) {
	hp.Push(item)
	now := hp.Len() - 1 //新节点位置
	nxt := now >> 1     //nxt为父节点

	for nxt > 0 {
		if !hp.Less(now, nxt) {
			break
		}

		hp.Swap(now, nxt)

		now, nxt = nxt, nxt>>1
	}
}

func Pop(hp heap) interface{} {
	//fmt.Println(hp)
	hp.Swap(1, hp.Len()-1)

	old := hp.Pop()

	down(hp, 1)

	return old
}

func Len(hp heap) int {
	return hp.Len() - 1
}
