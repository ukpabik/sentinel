package slidingwindowlog

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

type SlidingWindowLog struct {
	logQueue     *PriorityQueue
	allowedCount int
	rate         time.Duration
	counter      int
	mutex        sync.Mutex
}

func Init(allowedCount int, rate time.Duration) (*SlidingWindowLog, error) {
	if allowedCount <= 0 || rate <= 0 {
		return nil, fmt.Errorf("cannot create an empty sliding window log")
	}

	swl := &SlidingWindowLog{
		logQueue:     &PriorityQueue{},
		allowedCount: allowedCount,
		rate:         rate,
	}

	return swl, nil
}

func (swl *SlidingWindowLog) Allow() bool {
	if swl == nil {
		return false
	}

	swl.mutex.Lock()
	defer swl.mutex.Unlock()
	currLog := &Log{
		ID:        swl.counter,
		Timestamp: time.Now(),
	}
	swl.counter++

	timeToCheck := time.Now().Add(-swl.rate)
	for swl.logQueue.Len() > 0 && swl.logQueue.Top().Timestamp.Before(timeToCheck) {
		heap.Pop(swl.logQueue)
	}
	heap.Push(swl.logQueue, currLog)
	return swl.logQueue.Len() < swl.allowedCount
}

func (swl *SlidingWindowLog) AllowedCount() int {
	swl.mutex.Lock()
	defer swl.mutex.Unlock()
	return swl.allowedCount
}
