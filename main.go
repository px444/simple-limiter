package main

import (
	"fmt"
	"net/http"
	"time"

	"simple-limiter/limiter"
)

func main() {
	// Bucket holds up to 3 tokens, and refills 1 token every 1 second
	rl := limiter.New(3, 1*time.Second)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !rl.Allow() {
			http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
			return
		}
		fmt.Fprintln(w, "Request successful!")
	})

	fmt.Println("Server running on :8080 (Limit: 3 burst, 1 req/sec refill)")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}
