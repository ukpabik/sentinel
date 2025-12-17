package fixedwindowcounter

import (
	"fmt"
	"sync"
	"time"
)

type FixedWindowCounter struct {
	counter    int
	capacity   int
	windowSize time.Duration
	lastWindow time.Time

	mutex  sync.Mutex
	done   sync.Once
	stopCh chan struct{}
}

func Init(capacity int, timeGap time.Duration) (*FixedWindowCounter, error) {
	if capacity <= 0 || timeGap <= 0 {
		return nil, fmt.Errorf("cannot create an empty window")
	}

	window := &FixedWindowCounter{
		capacity:   capacity,
		windowSize: timeGap,
		lastWindow: time.Now(),
		stopCh:     make(chan struct{}),
	}

	go startClock(window)
	return window, nil
}

func startClock(fwc *FixedWindowCounter) {
	if fwc == nil {
		return
	}
	ticker := time.NewTicker(fwc.WindowSize())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			refill(fwc)
		case <-fwc.stopCh:
			return
		}
	}
}

func refill(fwc *FixedWindowCounter) {
	if fwc == nil {
		return
	}
	fwc.mutex.Lock()
	defer fwc.mutex.Unlock()
	fwc.lastWindow = time.Now()
	fwc.counter = 0
}

func (fwc *FixedWindowCounter) Allow() bool {
	if fwc == nil {
		return false
	}

	fwc.mutex.Lock()
	defer fwc.mutex.Unlock()

	if fwc.counter == fwc.capacity {
		return false
	}

	fwc.counter += 1
	return true
}

func (fwc *FixedWindowCounter) Stop() {
	fwc.done.Do(func() {
		close(fwc.stopCh)
	})
}

func (fwc *FixedWindowCounter) WindowSize() time.Duration {
	fwc.mutex.Lock()
	defer fwc.mutex.Unlock()
	return fwc.windowSize
}

func (fwc *FixedWindowCounter) Counter() int {
	fwc.mutex.Lock()
	defer fwc.mutex.Unlock()
	return fwc.counter
}

func (fwc *FixedWindowCounter) Capacity() int {
	fwc.mutex.Lock()
	defer fwc.mutex.Unlock()
	return fwc.capacity
}

func (fwc *FixedWindowCounter) LastWindow() time.Time {
	fwc.mutex.Lock()
	defer fwc.mutex.Unlock()
	return fwc.lastWindow
}
