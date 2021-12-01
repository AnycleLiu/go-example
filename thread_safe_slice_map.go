package main

import (
	"fmt"
	"sync"
)

type SyncMap struct {
	rw sync.RWMutex
	m  map[interface{}]interface{}
}

func NewSyncMap() *SyncMap {
	return &SyncMap{
		rw: sync.RWMutex{},
		m:  make(map[interface{}]interface{}),
	}
}

func (m *SyncMap) Put(key interface{}, value interface{}) {
	m.rw.Lock()
	defer m.rw.Unlock()

	m.m[key] = value
}
func (m *SyncMap) Get(key interface{}) (value interface{}, found bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	value, found = m.m[key]
	return
}

func (m *SyncMap) GetOrStore(key interface{}, valueFactory func() interface{}) interface{} {
	if val, f := m.Get(key); f {
		return val
	}

	m.rw.Lock()
	defer m.rw.Unlock()
	val := valueFactory()

	m.m[key] = val

	return val
}

func (m *SyncMap) Remove(key interface{}) {
	m.rw.Lock()
	defer m.rw.Unlock()

	delete(m.m, key)
}
func (m *SyncMap) Keys() []interface{} {
	m.rw.RLock()
	defer m.rw.RUnlock()

	keys := make([]interface{}, 0, len(m.m))
	for k := range m.m {
		keys = append(keys, k)
	}
	return keys
}

func (m *SyncMap) Len() int {
	m.rw.RLock()
	defer m.rw.RUnlock()
	return len(m.m)
}

func TestSyncMap(m *SyncMap) {
	const (
		//10个协程并发put
		goroutines = 100
		//每个协程put 100000个
		numPerGoroutine = 100000
	)

	for i := 0; i < goroutines; i++ {
		start := numPerGoroutine * i
		end := start + numPerGoroutine

		for j := start; j < end; j++ {
			m.Put(j, j)
		}
	}

	fmt.Printf("预期SyncMap有%d个元素，实际上有%d\n", goroutines*numPerGoroutine, m.Len())
}

func TestMap(m map[int]int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("出错", err)
		}
	}()
	const (
		//10个协程并发put
		goroutines = 100
		//每个协程put 100000个
		numPerGoroutine = 100000
	)

	for i := 0; i < goroutines; i++ {
		start := numPerGoroutine * i
		end := start + numPerGoroutine

		for j := start; j < end; j++ {
			m[j] = j
		}
	}

	fmt.Printf("预期Map有%d个元素，实际上有%d\n", goroutines*numPerGoroutine, len(m))
}

func main() {
	sm := NewSyncMap()

	TestSyncMap(sm)

	TestMap(make(map[int]int))
}
