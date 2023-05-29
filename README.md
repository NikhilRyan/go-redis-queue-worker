# Redis Queue Implementation in Go

This repository contains an implementation of a Redis queue in Go. It includes separate files for the Redis client, worker, and an API to push data into the queue.

## Prerequisites

- Go 1.16 or higher installed
- Redis installed and running on your local machine

## Getting Started

1. Clone the repository:

   ```
   git clone <repository_url>
   ```

2. Install the required dependencies:

   ```
   go get github.com/go-redis/redis/v8
   ```

3. Start the Redis server on your local machine. If Redis is not installed, you can install it using Homebrew on macOS:

   ```
   brew install redis
   ```

4. Start the Redis server:

   ```
   brew services start redis
   ```

5. Build and run the Go program:

   ```
   go build
   ./<binary_file_name>
   ```

   The server and worker will start running. The worker will listen for tasks in the Redis queue, and the API server will be available on port 8080.

## Pushing Tasks to the Queue

To push tasks to the Redis queue, you can make HTTP requests to the `/push` endpoint of the API server. The tasks will be added to the Redis queue and processed by the worker.

Send a POST request to `http://localhost:8080/push?task=<your_task_here>` to push a task to the queue. Replace `<your_task_here>` with the task you want to add.

Example using cURL:

```
curl -X POST "http://localhost:8080/push?task=example_task"
```

## Stopping the Program

To stop the program gracefully, send a termination signal to the process. You can do this by pressing Ctrl+C in the terminal where the program is running.

## Customization

Feel free to customize the code according to your requirements. You can add error handling, authentication, or enhance the worker's processing logic.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
