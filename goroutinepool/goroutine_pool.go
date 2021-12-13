package goroutinepool

import (
	"fmt"
	"sync/atomic"
	"time"
)

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
	return <-w.pool.tasks //拿到任务
	/*
		for {
			timed := false

			workerCount := atomic.LoadInt32(&w.pool.workerCount)
			//如果worker数量比最小数量要大，timed就为true，需要进行回收worker
			timed = workerCount > w.pool.min

			if timed {
				delay := time.NewTimer(w.pool.maxIdle)
				select {
				case task := <-w.pool.tasks: //拿到任务
					if !delay.Stop() {
						<-delay.C
					}
					return task
				case <-delay.C: //拿不到任务，超时
					if atomic.CompareAndSwapInt32(&w.pool.workerCount, workerCount, workerCount-1) {
						return nil
					}
				}
			} else {
				return <-w.pool.tasks //拿到任务
			}
		}*/
}

func (w *Worker) Work() {
	var task Task

	for {
		task = w.getTask()
		if task == nil {
			break
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
	return NewGoroutinePool(100, 10000, time.Second*60, 1024)
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
						p.addWorker(nil)
						p.tasks <- t
						return
					} else {
						continue retry //状态改变，更新失败，重试
					}
				}
			}
		}
	}
}
