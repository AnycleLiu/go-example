/*
循环队列
*/

package main

import (
	"errors"
	"fmt"
)

func main() {
	q := NewQueue(3)

	q.Push(1)
	q.Push(2)
	q.Push(3)
	fmt.Println("queue length: ", q.Len())

	err, item := q.Pop()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("pop item: ", item)
	fmt.Println("queue length: ", q.Len())

	err, item = q.Pop()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("pop item: ", item)
	fmt.Println("queue length: ", q.Len())

	q.Push(4)
	q.Push(5)

	fmt.Println("queue length: ", q.Len())

	for q.Len() > 0 {
		err, item = q.Pop()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("pop item: ", item)
		fmt.Println("queue length: ", q.Len())
	}
	fmt.Println("queue length: ", q.Len())

}

type DQueue interface {
	Push(item interface{}) error
	Pop() (error, interface{})
	IsEmpty() bool
	IsFull() bool
	Len() int
}

type Queue struct {
	n     int           //队列容量
	head  int           //头指针
	tail  int           //尾指针
	count int           //队列长度
	buf   []interface{} //缓冲区
}

func NewQueue(capacity int) *Queue {
	return &Queue{
		n:   capacity,
		buf: make([]interface{}, capacity, capacity),
	}
}

func (q *Queue) Len() int {
	return q.count
}

func (q *Queue) Push(item interface{}) error {
	if q.IsFull() {
		return errors.New("queue is full")
	}

	q.buf[q.tail] = item
	q.tail++
	q.count++
	if q.tail >= q.n {
		q.tail = 0
	}
	return nil
}

func (q *Queue) Pop() (error, interface{}) {
	if q.IsEmpty() {
		return errors.New("queue is empty"), nil
	}
	item := q.buf[q.head]
	q.head++
	q.count--
	if q.head >= q.n {
		q.head = 0
	}
	return nil, item
}

func (q *Queue) IsEmpty() bool {
	return q.count == 0
}

func (q *Queue) IsFull() bool {
	return q.count == q.n
}
