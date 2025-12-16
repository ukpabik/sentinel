package leakybucket

import (
	"fmt"
	"sync"
	"time"
)

type Token struct{}

type LeakyBucket struct {
	requestQueue  *ConcurrentQueue
	bucketSize    int
	outflowAmount int
	outflowRate   time.Duration

	mutex    sync.Mutex
	done     sync.Once
	stopChan chan struct{}
}

func Init(bucketSize, outflowAmount int, outflowRate time.Duration) (*LeakyBucket, error) {
	if bucketSize <= 0 || outflowAmount <= 0 || outflowRate <= 0 {
		return nil, fmt.Errorf("bucket inputs are <= 0, cannot create an empty bucket")
	}

	queue, err := ConcurrentQueueInitializer(bucketSize)
	if err != nil {
		return nil, fmt.Errorf("error initializing queue: %v", err)
	}

	lb := &LeakyBucket{
		requestQueue:  queue,
		bucketSize:    bucketSize,
		outflowAmount: outflowAmount,
		outflowRate:   outflowRate,
		stopChan:      make(chan struct{}),
	}

	go StartClock(lb)

	return lb, nil
}

func StartClock(lb *LeakyBucket) {
	ticker := time.NewTicker(lb.outflowRate)
	for {
		select {
		case <-ticker.C:
			lb.mutex.Lock()
			// pop from tail here and "process" request
			for range min(lb.outflowAmount, lb.requestQueue.Size()) {
				lb.requestQueue.Pop()
			}
			lb.mutex.Unlock()
		case <-lb.stopChan:
			return
		}
	}
}

func (lb *LeakyBucket) Allow() bool {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	err := lb.requestQueue.Append()
	return err == nil
}

func (lb *LeakyBucket) Stop() {
	if lb == nil {
		return
	}
	lb.done.Do(func() {
		close(lb.stopChan)
	})
}

func (lb *LeakyBucket) OutflowRate() time.Duration {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	return lb.outflowRate
}

func (lb *LeakyBucket) OutflowAmount() int {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	return lb.outflowAmount
}

func (lb *LeakyBucket) BucketSize() int {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	return lb.bucketSize
}

func (lb *LeakyBucket) QueueSize() int {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	return lb.requestQueue.size
}
