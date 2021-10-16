package main

import (
	"log"
	"time"
)

type WorkerPool interface {
	Run()                // create workers (i.e. goroutines) and run them
	AddTask(task func()) // add task (which is a function to execute) to the queue
}

type workerPool struct {
	maxWorker int
	taskQueue chan func()
}

func (wp *workerPool) Run() {
	for i := 0; i < wp.maxWorker; i++ {
		go func(workerID int) {
			log.Printf("[workerPool] worker %v has been created", workerID)
			for task := range wp.taskQueue {
				log.Printf("[workerPool] worker %v picked up the task", workerID)
				task()
				log.Printf("[workerPool] worker %v finished the task", workerID)
			}
		}(i)
	}
}

func (wp *workerPool) AddTask(task func()) {
	wp.taskQueue <- task
}

func main() {
	// create worker pool and initialize the workers
	wp := &workerPool{maxWorker: 3, taskQueue: make(chan func())}
	wp.Run()

	resultChannel := make(chan int)

	go func() {
		for result := range resultChannel {
			log.Printf("[main] received result of task %v", result)
		}
	}()

	// add tasks to the queue
	for i := 0; i < 5; i++ {
		id := i
		log.Printf("adding task %v", id)
		wp.AddTask(func() {
			log.Printf("[main] starting task %v", id)
			time.Sleep(5 * time.Second)
			resultChannel <- id
		})
	}
}
