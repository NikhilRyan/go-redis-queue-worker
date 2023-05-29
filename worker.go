package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// Worker represents a worker that processes tasks from the Redis queue
type Worker struct {
	client *redis.Client
}

// NewWorker creates a new Worker instance
func NewWorker(client *redis.Client) *Worker {
	return &Worker{
		client: client,
	}
}

// Start starts the worker and listens for tasks in the Redis queue
func (w *Worker) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Retrieve a task from the Redis queue
			task, err := w.client.BLPop(ctx, 0, "myqueue").Result()
			if err != nil {
				log.Println("Failed to retrieve task:", err)
				continue
			}

			// Process the task
			fmt.Println("Processing task:", task[1])

			// Simulate some work
			time.Sleep(1 * time.Second)
		}
	}
}
