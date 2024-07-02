package main

import (
	"fmt"
	"sync"
	"time"
)

// Buffer represents the shared buffer between producer and consumer
type Buffer struct {
	items []int
	mutex sync.Mutex
	cond  *sync.Cond
	size  int
}

// NewBuffer creates a new Buffer with the given size
func NewBuffer(size int) *Buffer {
	b := &Buffer{
		items: make([]int, 0, size),
		size:  size,
	}
	b.cond = sync.NewCond(&b.mutex)
	return b
}

// Produce adds an item to the buffer
func (b *Buffer) Produce(item int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Wait if buffer is full
	for len(b.items) == b.size {
		b.cond.Wait()
	}

	// Add item to buffer
	b.items = append(b.items, item)
	fmt.Printf("Produced: %d\n", item)

	// Signal that a new item is available
	b.cond.Signal()
}

// Consume removes and returns an item from the buffer
func (b *Buffer) Consume() int {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Wait if buffer is empty
	for len(b.items) == 0 {
		b.cond.Wait()
	}

	// Remove item from buffer
	item := b.items[0]
	b.items = b.items[1:]
	fmt.Printf("Consumed: %d\n", item)

	// Signal that there is space available
	b.cond.Signal()

	return item
}

// Producer function
func Producer(b *Buffer, rate time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		time.Sleep(rate)
		b.Produce(i)
	}
}

// Consumer function
func Consumer(b *Buffer, rate time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		time.Sleep(rate)
		b.Consume()
	}
}

func main() {
	bufferSize := 5
	buffer := NewBuffer(bufferSize)

	var wg sync.WaitGroup
	producerRate := time.Second
	consumerRate := 2 * time.Second

	wg.Add(2)
	go Producer(buffer, producerRate, &wg)
	go Consumer(buffer, consumerRate, &wg)

	wg.Wait()
	fmt.Println("Done")
}
