# Sentinel: High-Performance Go Rate Limiters

**Sentinel** is a collection of high-performance rate limiter implementations written in Go. It showcases multiple strategies for controlling request rates, making it suitable for APIs, microservices, or distributed systems.

### 🚀 Key Features

* **Token Bucket Limiter:** Efficiently handles traffic bursts while enforcing a steady average rate.
* **Concurrent Safe:** Fully thread-safe and optimized for concurrent Go routines.
* **In-Memory Implementations:** Lightweight and zero-dependency, perfect for local testing or single-instance services.

---

### 📦 Usage Example

Below is an example of how to implement the **Token Bucket Limiter**.

```go
package main

import (
    "fmt"
    "time"
    "[github.com/ukpabik/sentinel/in_memory](https://github.com/ukpabik/sentinel/in_memory)"
)

func main() {
    // Initialize a limiter with:
    // - Capacity: 10 tokens
    // - Refill Rate: 2 tokens per second
    limiter, err := in_memory.Init(10, 2, time.Second)
    if err != nil {
        panic(err)
    }
    defer limiter.Stop()

    // Attempt to allow a request
    if limiter.Allow() {
        fmt.Println("Request allowed")
    } else {
        fmt.Println("Request denied")
    }

    // Inspect current state
    fmt.Printf("Current tokens: %d\n", limiter.CurrentTokenAmount())
}
```

### 🧪 Testing

To ensure stability, you can run the full suite of unit tests using the standard Go test tool:

```bash
go test ./...
```

### 🗺️ Future Plans & Roadmap

We are actively working on expanding the algorithms and capabilities of Sentinel:

* [ ] **Leaky Bucket Limiter:** Ensure requests are processed at a perfectly constant rate (smoothing traffic).
* [ ] **Fixed Window Limiter:** Simple rate limiting based on distinct time windows.
* [ ] **Sliding Window Limiter:** High-accuracy rate limiting over rolling time windows to prevent window-boundary spikes.
* [ ] **Distributed Support:** Redis-backed limiters for distributed systems.
* [ ] **Benchmarks:** Comprehensive performance comparisons between algorithms.
