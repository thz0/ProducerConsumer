package main

import (
	"sync"
	"testing"
	"time"
)

func TestProducerConsumer(t *testing.T) {
	taskCh := make(chan Task, 10)
	rateLimit := time.Tick(200 * time.Millisecond)
	var wgp, wgc sync.WaitGroup

	numProducers := 1
	numConsumers := 1

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

	// 确保所有消费者完成
	wgc.Wait()

	remainingTasks := len(taskCh)
	if remainingTasks != 0 {
		t.Errorf("Expected 0 remaining tasks, but got %d", remainingTasks)
	}
}

func TestMultipleProducersConsumers(t *testing.T) {
	//taskCh := make(chan Task, 20)
	//rateLimit := time.Tick(200 * time.Millisecond)
	//var wgp, wgc sync.WaitGroup
	//
	//numProducers := 2
	//numConsumers := 2
	var tests = []struct {
		numP int
		numC int
	}{
		{2, 3},
		{3, 2},
		{3, 3},
	}

	for _, tt := range tests {
		o := tt

		t.Run("", func(t *testing.T) {
			t.Parallel()
			taskCh := make(chan Task, 20)
			rateLimit := time.Tick(200 * time.Millisecond)
			var wgp, wgc sync.WaitGroup

			numProducers := o.numP
			numConsumers := o.numC

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

			remainingTasks := len(taskCh)
			if remainingTasks != 0 {
				t.Errorf("Expected 0 remaining tasks, but got %d", remainingTasks)
			}
		})
	}
}
