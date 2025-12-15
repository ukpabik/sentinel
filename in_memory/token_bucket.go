package in_memory

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucketLimiter struct {

	// Max number of tokens that can be inside of the bucket
	bucketSize int

	// Current amount of tokens that are available
	currentTokenAmount int

	// The number of tokens that are put back into the bucket
	refillAmount int

	// The time between each refill
	refillRate time.Duration

	mutex    sync.Mutex
	done     sync.Once
	stopChan chan struct{}
}

func Init(size int, refillAmount int, rate time.Duration) (*TokenBucketLimiter, error) {
	if size <= 0 || refillAmount <= 0 || rate <= 0 {
		return nil, fmt.Errorf("bucket inputs are <= 0, cannot create an empty bucket")
	}

	limiter := &TokenBucketLimiter{
		bucketSize:         size,
		currentTokenAmount: size,
		refillAmount:       refillAmount,
		refillRate:         rate,
		stopChan:           make(chan struct{}),
	}
	go startClock(limiter)
	return limiter, nil
}

func (tl *TokenBucketLimiter) Allow() bool {
	if tl == nil {
		return false
	}
	tl.mutex.Lock()
	defer tl.mutex.Unlock()

	if tl.currentTokenAmount == 0 {
		return false
	}

	tl.currentTokenAmount -= 1
	return true
}

func startClock(tl *TokenBucketLimiter) {
	if tl == nil {
		return
	}
	ticker := time.NewTicker(tl.refillRate)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			refill(tl)
		case <-tl.stopChan:
			return
		}
	}
}

func refill(limiter *TokenBucketLimiter) {
	if limiter == nil {
		return
	}

	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()
	limiter.currentTokenAmount = min(limiter.currentTokenAmount+limiter.refillAmount, limiter.bucketSize)
}

func (tl *TokenBucketLimiter) Stop() {
	if tl == nil {
		return
	}
	tl.done.Do(func() {
		close(tl.stopChan)
	})
}

func (tl *TokenBucketLimiter) BucketSize() int {
	tl.mutex.Lock()
	defer tl.mutex.Unlock()
	return tl.bucketSize
}

func (tl *TokenBucketLimiter) CurrentTokenAmount() int {
	tl.mutex.Lock()
	defer tl.mutex.Unlock()
	return tl.currentTokenAmount
}

func (tl *TokenBucketLimiter) RefillAmount() int {
	tl.mutex.Lock()
	defer tl.mutex.Unlock()
	return tl.refillAmount
}

func (tl *TokenBucketLimiter) RefillRate() time.Duration {
	tl.mutex.Lock()
	defer tl.mutex.Unlock()
	return tl.refillRate
}
