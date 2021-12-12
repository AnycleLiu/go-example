package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
	// GiB // 1073741824
	// TiB // 1099511627776             (超过了int32的范围)
	// PiB // 1125899906842624
	// EiB // 1152921504606846976
	// ZiB // 1180591620717411303424    (超过了int64的范围)
	// YiB // 1208925819614629174706176
)

func main() {
	pool := NewGoroutinePoolWithDefault()
	var n int32
	var wg sync.WaitGroup

	for i := 0; i < 100000; i++ {
		wg.Add(1)
		pool.Put(func() {
			time.Sleep(time.Microsecond * 1000)
			atomic.AddInt32(&n, 1)
			wg.Done()
		})
	}

	wg.Wait()
	fmt.Println(n)
}

var curMem uint64

const (
	n     = 100000
	param = 100
)

func demoFunc(args interface{}) {
	n := args.(int)
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func TestDetaultPool(t *testing.T) {
	var wg sync.WaitGroup
	pool := NewGoroutinePoolWithDefault()

	for i := 0; i < n; i++ {
		wg.Add(1)
		pool.Put(func() {
			demoFunc(param)
			wg.Done()
		})
	}
	wg.Wait()
	t.Logf("pool, running workers number:%d", pool.Running())
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}
func TestNoPool(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			demoFunc(param)
			wg.Done()
		}()
	}

	wg.Wait()
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

type Task func()

type Worker struct {
	initTask Task
	pool     *GoroutinePool
	exit     chan struct{}
}

func (w *Worker) getTask() Task {
	if w.initTask != nil {
		t := w.initTask
		w.initTask = nil
		return t
	}
	var timeout bool

	for {
		timed := false
		for {
			workerCount := atomic.LoadInt32(&w.pool.workerCount)
			//如果worker数量比最小数量要大，timed就为true，需要进行回收worker
			timed := workerCount > w.pool.min

			if workerCount <= w.pool.max && !(timed && timeout) {
				break
			}
			if atomic.CompareAndSwapInt32(&w.pool.workerCount, workerCount, workerCount-1) {
				return nil
			}
		}

		if timed {
			select {
			case task := <-w.pool.tasks: //拿到任务
				return task
			case <-time.After(w.pool.maxIdle): //拿不到任务，超时
				timeout = true
			}
		} else {
			select {
			case task := <-w.pool.tasks: //拿到任务
				return task
			}
		}
	}
}

func (w *Worker) Work() {
	var task Task

	for {
		task = w.getTask()
		if task == nil {
			return
		}
		safeRun(task)
	}
	w.processExit()
}
func (w *Worker) processExit() {
	for {
		workerCount := atomic.LoadInt32(&w.pool.workerCount)

		//如果worker数量小于最小worker数量，重新开启一个worker
		if workerCount < w.pool.min {
			if atomic.CompareAndSwapInt32(&w.pool.workerCount, workerCount, workerCount+1) {
				w.pool.addWorker(nil)
				return
			} else {
				continue
			}
		} else {
			return
		}
	}
}

func safeRun(f func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	f()
}

type GoroutinePool struct {
	min, max     int32         //最小/大协程数量
	workerCount  int32         //协程数量
	maxTaskCount int           //任务队列长度
	maxIdle      time.Duration //最大空闲时间

	tasks chan Task //任务队列
}

func NewGoroutinePoolWithDefault() *GoroutinePool {
	return NewGoroutinePool(int32(runtime.NumCPU())*2, 64, time.Second*60, 1000)
}

func NewGoroutinePool(minWorkerCount, maxWorkerCount int32, maxIdle time.Duration, maxTaskCount int) *GoroutinePool {
	p := &GoroutinePool{
		min:          minWorkerCount,
		max:          maxWorkerCount,
		maxIdle:      maxIdle,
		tasks:        make(chan Task, maxTaskCount),
		maxTaskCount: maxTaskCount,
	}

	return p
}
func (p *GoroutinePool) Running() int {
	return int(atomic.LoadInt32(&p.workerCount))
}

func (p *GoroutinePool) addWorker(task Task) {
	worker := &Worker{
		pool:     p,
		exit:     make(chan struct{}),
		initTask: task,
	}
	go worker.Work()
}

func (p *GoroutinePool) Put(t Task) {
	/*
		1. 如果workerCount < min，则创建并启动一个协程来执行新提交的任务。
		2. 如果workerCount >= min，且池内的队列未满，则将任务添加到该队列中。
		3. 如果workerCount >= min && workerCount < max，且池内的队列已满，则创建并启动一个协程来执行新提交的任务。
		4. 如果workerCount >= max，并且池内的队列已满, 阻塞等待写入。
	*/
retry:
	for {
		workerCount := atomic.LoadInt32(&p.workerCount)

		if workerCount < p.min {
			if atomic.CompareAndSwapInt32(&p.workerCount, workerCount, workerCount+1) {
				p.addWorker(t)
				return
			} else {
				continue retry //状态改变，更新失败，重试
			}
		} else {
			taskCount := len(p.tasks)
			if taskCount < p.maxTaskCount {
				p.tasks <- t
				return
			} else {
				if workerCount >= p.max {
					p.tasks <- t
					return
				} else {
					if atomic.CompareAndSwapInt32(&p.workerCount, workerCount, workerCount+1) {
						p.addWorker(t)
						return
					} else {
						continue retry //状态改变，更新失败，重试
					}
				}
			}
		}
	}
}
