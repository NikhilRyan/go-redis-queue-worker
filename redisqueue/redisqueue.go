package redisqueue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// WorkerFunc is a function type for worker functions that process tasks
type WorkerFunc func(ctx context.Context, task string)

// RedisQueue represents a Redis queue
type RedisQueue struct {
	client  *redis.Client
	workers map[string]*Worker
}

// Worker represents a worker that processes tasks from a Redis queue
type Worker struct {
	queueName string
	function  WorkerFunc
}

// NewRedisQueue creates a new RedisQueue instance
func NewRedisQueue(addr, password string, db int) (*RedisQueue, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisQueue{
		client:  client,
		workers: make(map[string]*Worker),
	}, nil
}

// RegisterWorker registers a worker with a specific queue and worker function
func (rq *RedisQueue) RegisterWorker(queueName string, function WorkerFunc) {
	worker := &Worker{
		queueName: queueName,
		function:  function,
	}

	rq.workers[queueName] = worker
}

// Start starts all the workers
func (rq *RedisQueue) Start(ctx context.Context) {
	var wg sync.WaitGroup

	for _, worker := range rq.workers {
		wg.Add(1)
		go rq.startWorker(ctx, worker, &wg)
	}

	wg.Wait()
}

func (rq *RedisQueue) startWorker(ctx context.Context, worker *Worker, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Retrieve a task from the Redis queue
			task, err := rq.client.BLPop(ctx, 0, worker.queueName).Result()
			if err != nil {
				log.Printf("Failed to retrieve task from queue %s: %v\n", worker.queueName, err)
				continue
			}

			// Process the task
			worker.function(ctx, task[1])

			// Simulate some work
			time.Sleep(1 * time.Second)
		}
	}
}
