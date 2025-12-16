# Sentinel: High-Performance Go Rate Limiters

**Sentinel** is a collection of high-performance rate limiter implementations written in Go. 

### 🚀 Key Features

* **Concurrent Safe:** Fully thread-safe and optimized for concurrent Go routines.
* **In-Memory Implementations:** Lightweight and zero-dependency, perfect for local testing or single-instance services.

---

### 🧪 Testing

To ensure stability, you can run the full suite of unit tests using the standard Go test tool:

```bash
go test ./...
```
Each of the test suites also shows how to use the specified implementation. I recommend looking over them.

### 🗺️ Future Plans & Roadmap

I'm currently working on implementing these algorithms below! I also want to create distributed versions for production-level rate limiting in the future.

* [x] **Token Bucket Limiter:** Efficiently handles traffic bursts while enforcing a steady average rate.
* [x] **Leaky Bucket Limiter:** Ensure requests are processed at a perfectly constant rate (smoothing traffic).
* [ ] **Fixed Window Limiter:** Simple rate limiting based on distinct time windows.
* [ ] **Sliding Window Limiter:** High-accuracy rate limiting over rolling time windows to prevent window-boundary spikes.
* [ ] **Distributed Support:** Redis-backed limiters for distributed systems.
* [ ] **Benchmarks:** Comprehensive performance comparisons between algorithms, involving throughput and speed.
