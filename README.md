# Simple Go Rate Limiter

A thread-safe, zero-dependency Token Bucket rate limiter implemented in Go using standard library primitives (`sync.Mutex` and `time`).

## Features

- **Token Bucket Algorithm:** Supports bursts up to bucket capacity while enforcing a steady refill rate.
- **Lazy Evaluation:** Calculates token replenishment on-demand without background goroutine polling.
- **Thread-Safe:** Safe for concurrent use across multiple HTTP handler goroutines using `sync.Mutex`.
- **Zero External Dependencies:** Built entirely with the Go standard library.

## Project Structure

```text
simple-limiter/
├── go.mod             # Module definition
├── limiter/
│   └── limiter.go     # Core Token Bucket logic
└── main.go            # Example HTTP server demonstration