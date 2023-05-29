package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-redis/redis/v8"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func main() {
	// Create a new Redis client
	redisClient := NewRedisClient()

	// Create a context to handle cancellation
	ctx, cancel = context.WithCancel(context.Background())

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup

	// Create a worker
	worker := NewWorker(redisClient)

	// Start the worker goroutine
	wg.Add(1)
	go worker.Start(&wg)

	// Start the API server
	go startAPIServer(redisClient)

	fmt.Println("Server started.")

	// Wait for termination signal
	waitForTerminationSignal()

	// Stop the worker
	cancel()
	wg.Wait()

	fmt.Println("Server stopped.")
}

func startAPIServer(redisClient *redis.Client) {
	http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		task := r.URL.Query().Get("task")
		if task == "" {
			http.Error(w, "Missing task parameter", http.StatusBadRequest)
			return
		}

		// Push the task to the Redis queue
		err := redisClient.RPush(ctx, "myqueue", task).Err()
		if err != nil {
			log.Println("Failed to push task to Redis queue:", err)
			http.Error(w, "Failed to push task to Redis queue", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Task %s pushed to the Redis queue", task)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Failed to start API server:", err)
	}
}

func waitForTerminationSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
}
