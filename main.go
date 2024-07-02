package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	id int
}

func producer(id int, taskCh chan Task, rateLimit <-chan time.Time, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		<-rateLimit
		task := Task{10*id + i}
		fmt.Printf("Producer %d produced task %d\n", id, task.id)
		taskCh <- task
	}
}

func consumer(id int, taskCh chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskCh {
		fmt.Printf("Consumer %d consumed task %d\n", id, task.id)
	}
}

func main() {
	taskCh := make(chan Task, 10)
	rateLimit := time.Tick(200 * time.Millisecond)
	var wgp sync.WaitGroup
	var wgc sync.WaitGroup

	numProducers := 3
	numConsumers := 2

	for i := 0; i < numProducers; i++ {
		wgp.Add(1)
		go producer(i, taskCh, rateLimit, &wgp)
	}

	for i := 0; i < numConsumers; i++ {
		wgc.Add(1)
		go consumer(i, taskCh, &wgc)
	}

	wgp.Wait()
	close(taskCh)

	wgc.Wait()
	fmt.Println("All tasks processed.")
}
