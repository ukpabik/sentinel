package leakybucket_test

import (
	"fmt"
	"testing"
	"time"

	leakybucket "github.com/ukpabik/sentinel/in_memory/leaky_bucket"
)

func TestInitializationHappyPath(t *testing.T) {
	limiter, err := leakybucket.Init(10, 5, time.Second)
	if err != nil {
		t.Error("init failed")
	}
	defer limiter.Stop()

	if limiter.BucketSize() != 10 || limiter.OutflowAmount() != 5 || limiter.OutflowRate() != time.Second {
		t.Error("the limiter values are not eq")
	}
}

func TestInitializationEmptyBucket(t *testing.T) {
	limiter, err := leakybucket.Init(0, 1, time.Second)
	if err == nil {
		t.Error("init succeeded when it should've failed")
	}
	defer limiter.Stop()
}

func TestRateLimitHappyPath(t *testing.T) {
	limiter, err := leakybucket.Init(10, 10, time.Second)
	if err != nil {
		t.Error("init failed")
	}
	defer limiter.Stop()

	for range 10 {
		if !limiter.Allow() {
			t.Error("bucket is empty before expected")
		}
	}
}

func TestRateLimitExceeded(t *testing.T) {
	limiter, err := leakybucket.Init(10, 2, time.Second)
	if err != nil {
		t.Error("init failed")
	}
	defer limiter.Stop()

	for range 10 {
		limiter.Allow()
	}
	if limiter.Allow() {
		t.Error("leaky bucket queue is not limiting requests as expected")
	}
}

func TestBucketRefill(t *testing.T) {
	limiter, err := leakybucket.Init(10, 5, time.Second)
	if err != nil {
		t.Error("init failed")
	}
	defer limiter.Stop()

	for range 10 {
		limiter.Allow()
	}
	time.Sleep(limiter.OutflowRate() + 100*time.Millisecond)
	fmt.Printf("Current queue size: %d\n", limiter.QueueSize())
	if limiter.QueueSize() != 5 || !limiter.Allow() {
		t.Error("queue was not emptied during outflow period")
	}
}
