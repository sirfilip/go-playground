package main

import (
	"fmt"
	"sync"
	"time"
)

type Job interface {
	Perform() error
}

type WorkerPool interface {
	Start() error
	Add(Job) error
	Shutdown() error
}

type WorkerPoolImpl struct {
	queue    chan Job
	wg       *sync.WaitGroup
	capacity int
}

func (self WorkerPoolImpl) Start() error {
	self.wg.Add(self.capacity)
	for i := 0; i < self.capacity; i++ {
		go func() {
			for job := range self.queue {
				job.Perform()
			}
			self.wg.Done()
		}()
	}
	return nil
}

func (self WorkerPoolImpl) Add(job Job) error {
	self.queue <- job
	return nil
}

func (self WorkerPoolImpl) Shutdown() error {
	close(self.queue)
	self.wg.Wait()
	return nil
}

func NewWorkerPool(capacity int) WorkerPoolImpl {
	queue := make(chan Job, capacity)
	wg := new(sync.WaitGroup)
	return WorkerPoolImpl{
		queue:    queue,
		wg:       wg,
		capacity: capacity,
	}
}

type JobImpl struct {
	id  int
	out chan int
}

func (self JobImpl) Perform() error {
	time.Sleep(1 * time.Second)
	self.out <- self.id
	return nil
}

func main() {
	fmt.Println("Worker Pool implementation")
	out := make(chan int)
	done := make(chan int)
	go func() {
		for i := range out {
			fmt.Println(i)
		}
		done <- 1
	}()
	pool := NewWorkerPool(10)
	pool.Start()
	for i := 0; i < 100; i++ {
		pool.Add(JobImpl{i, out})
	}
	fmt.Println("Shutting down the pool\n waitin for all of the jobs to finish")
	pool.Shutdown()
	fmt.Println("All done exiting")
	close(out)
	<-done
}
