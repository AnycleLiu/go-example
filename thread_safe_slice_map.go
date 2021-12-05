package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

//SyncMap： 基于读写锁实现的线程安全的map
type SyncMap struct {
	rw sync.RWMutex
	m  map[interface{}]interface{}
}

//SyncSlice：基于Mutex实现的线程安全的slice
type SyncSlice struct {
	m sync.Mutex
	s []interface{}
}

func NewSyncMap() *SyncMap {
	return &SyncMap{
		rw: sync.RWMutex{},
		m:  make(map[interface{}]interface{}),
	}
}
func NewSyncSlice() *SyncSlice {
	return &SyncSlice{
		m: sync.Mutex{},
		s: make([]interface{}, 0),
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
	fmt.Println("test sync map")
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer func() {
				wg.Done()
			}()
			for j := 0; j < 100; j++ {
				m.Put(j, j)
			}
		}()
	}
	wg.Wait()
	fmt.Println("sync map test done.")
}

func TestMap(m map[int]int) {
	fmt.Println("test map")
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			m[10] = 10
		}()
	}

	wg.Wait()
	fmt.Println("map test done")
}

func (s *SyncSlice) Append(items ...interface{}) {
	s.m.Lock()
	defer s.m.Unlock()

	s.s = append(s.s, items)
}
func (s *SyncSlice) Len() int {
	s.m.Lock()
	defer s.m.Unlock()
	return len(s.s)
}
func (s *SyncSlice) Range(i, j int) {
	s.m.Lock()
	defer s.m.Unlock()
	s.s = s.s[i:j]
}
func TestSyncSlice(s *SyncSlice) {
	fmt.Println("test sync slice")
	var wg sync.WaitGroup
	var n int32

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100000; j++ {
				s.Append(j)
				atomic.AddInt32(&n, 1)
			}
		}()
	}
	wg.Wait()

	fmt.Printf("test append, expect length is %d, actual %d\n", n, s.Len())

	fmt.Println("sync slice test done.")
}

func main() {
	//为什么slice和map不是线程安全的
	//1. map并发写入会引发fatal error: concurrent map writes，导致程序退出无法恢复
	//2. slice并发append会导致数据被覆盖

	sm := NewSyncMap()

	TestSyncMap(sm)

	//TestMap(make(map[int]int))

	ss := NewSyncSlice()
	TestSyncSlice(ss)

}
