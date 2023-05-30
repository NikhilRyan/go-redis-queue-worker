package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"go-redis-queue/redisqueue"
)

func main() {
	// Create a new Redis queue
	rq, err := redisqueue.NewRedisQueue("localhost:6379", "", 0)
	if err != nil {
		log.Fatal("Failed to create Redis queue:", err)
	}

	// Register worker functions for different queues
	rq.RegisterWorker("queue1", workerFunction1)
	rq.RegisterWorker("queue2", workerFunction2)

	// Start the Redis queue and workers
	go rq.Start(context.Background())

	// Set up API routes
	http.HandleFunc("/push", handlePushRequest)

	// Start the API server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePushRequest(w http.ResponseWriter, r *http.Request) {
	// Parse the query parameter "task" from the URL
	task := r.URL.Query().Get("task")
	if task == "" {
		http.Error(w, "Missing task parameter", http.StatusBadRequest)
		return
	}

	// Create a Redis client to push the task into the queue
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Push the task into the queue
	err := client.RPush(context.Background(), "queue1", task).Err()
	if err != nil {
		http.Error(w, "Failed to push task into the queue", http.StatusInternalServerError)
		return
	}

	// Return a success response
	response := map[string]interface{}{
		"status":  "success",
		"message": "Task added to the queue",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func workerFunction1(ctx context.Context, task string) {
	// Implement the worker function for queue1
	fmt.Printf("Worker1 processing task from queue1: %s\n", task)
	// Perform the required task processing logic for queue1
}

func workerFunction2(ctx context.Context, task string) {
	// Implement the worker function for queue2
	fmt.Printf("Worker2 processing task from queue2: %s\n", task)
	// Perform the required task processing logic for queue2
}
