package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"simple-limiter/limiter"
)

// getIP extracts the IP address from the request's RemoteAddr (host:port).
func getIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func main() {
	// Each IP gets 3 burst tokens, refilling 1 token every 2 seconds
	ipLimiter := limiter.NewIPRateLimiter(3, 2*time.Second)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientIP := getIP(r)

		if !ipLimiter.Allow(clientIP) {
			http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
			return
		}

		fmt.Fprintf(w, "Request successful for IP: %s\n", clientIP)
	})

	fmt.Println("Server running on :8080 (Per-IP limit: 3 burst, 1 token/2s refill)")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}
