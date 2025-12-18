package slidingwindowcounter

import (
	"fmt"
	"sync"
	"time"

	inmemory "github.com/ukpabik/sentinel/in_memory"
)

type SlidingWindowCounter struct {
	logQueue        *inmemory.PriorityQueue
	allowedCount    int
	rate            time.Duration
	prevWindowCount int
	prevWindowStart time.Time
	currWindowCount int
	currWindowStart time.Time

	mutex sync.Mutex
}

func Init(allowedCount int, rate time.Duration) (*SlidingWindowCounter, error) {
	if allowedCount <= 0 || rate <= 0 {
		return nil, fmt.Errorf("cannot create an empty sliding window counter")
	}

	swc := &SlidingWindowCounter{
		logQueue:        &inmemory.PriorityQueue{},
		allowedCount:    allowedCount,
		rate:            rate,
		currWindowStart: time.Now(),
	}

	return swc, nil
}

func (swc *SlidingWindowCounter) Allow() bool {
	if swc == nil {
		return false
	}
	swc.mutex.Lock()
	defer swc.mutex.Unlock()
	currTime := time.Now()
	if swc.currWindowStart.Add(swc.rate).Before(currTime) {
		swc.prevWindowStart = swc.currWindowStart
		swc.prevWindowCount = swc.currWindowCount
		swc.currWindowStart = currTime
		swc.currWindowCount = 0
	}

	start := currTime.Add(-swc.rate)

	prevWindowEnd := swc.prevWindowStart.Add(swc.rate)

	overlapDuration := prevWindowEnd.Sub(start)
	overlapPercentage := max(0, min(1, float64(overlapDuration)/float64(swc.rate)))

	// current window reqs + (prev window reqs * overlap %)
	requestsInOverlappedWindow := int(float64(swc.currWindowCount) + (float64(swc.prevWindowCount) * (overlapPercentage)))
	if requestsInOverlappedWindow >= swc.allowedCount {
		return false
	}
	swc.currWindowCount += 1

	return true
}

func (swc *SlidingWindowCounter) CurrentWindowCount() int {
	swc.mutex.Lock()
	defer swc.mutex.Unlock()
	return swc.currWindowCount
}

func (swc *SlidingWindowCounter) PreviousWindowCount() int {
	swc.mutex.Lock()
	defer swc.mutex.Unlock()
	return swc.prevWindowCount
}
