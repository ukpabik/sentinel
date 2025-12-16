package tokenbucket_test

import (
	"fmt"
	"testing"
	"time"

	token_bucket "github.com/ukpabik/sentinel/in_memory/token_bucket"
)

func TestInitializationHappyPath(t *testing.T) {
	limiter, err := token_bucket.Init(10, 1, time.Second)
	defer limiter.Stop()
	if err != nil {
		t.Error("init failed")
	}

	if limiter.BucketSize() != 10 || limiter.RefillAmount() != 1 || limiter.RefillRate() != time.Second {
		t.Error("the limiter values are not eq")
	}
}

func TestInitializationEmptyBucket(t *testing.T) {
	limiter, err := token_bucket.Init(0, 1, time.Second)
	defer limiter.Stop()
	if err == nil {
		t.Error("init succeeded when it should've failed")
	}

}

func TestRateLimitHappyPath(t *testing.T) {
	limiter, err := token_bucket.Init(10, 2, time.Second)
	defer limiter.Stop()
	if err != nil {
		t.Error("init failed")
	}

	for range 10 {
		if !limiter.Allow() {
			t.Error("bucket is empty before expected")
		}
	}
}

func TestRateLimitExceeded(t *testing.T) {
	limiter, err := token_bucket.Init(10, 2, time.Second)
	defer limiter.Stop()
	if err != nil {
		t.Error("init failed")
	}

	for range 10 {
		limiter.Allow()
	}
	if limiter.Allow() {
		t.Error("bucket tokens are not depleting as expected")
	}
}

func TestBucketRefill(t *testing.T) {
	limiter, err := token_bucket.Init(10, 2, time.Second)
	defer limiter.Stop()
	if err != nil {
		t.Error("init failed")
	}

	for range 10 {
		limiter.Allow()
	}
	time.Sleep(limiter.RefillRate() + 100*time.Millisecond)
	fmt.Printf("Current token amount: %d\n", limiter.CurrentTokenAmount())
	if limiter.CurrentTokenAmount() != 2 || !limiter.Allow() {
		t.Error("bucket tokens are not refilling as expected")
	}
}
