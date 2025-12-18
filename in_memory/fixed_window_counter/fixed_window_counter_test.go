package fixedwindowcounter_test

import (
	"testing"
	"time"

	fixedwindowcounter "github.com/ukpabik/sentinel/in_memory/fixed_window_counter"
)

func TestInitializationHappyPath(t *testing.T) {
	window, err := fixedwindowcounter.Init(10, 5*time.Second)
	if err != nil {
		t.Error("initialization failed")
	}
	defer window.Stop()

	if window.WindowSize() != 5*time.Second || window.Capacity() != 10 {
		t.Errorf("expected size: %d, got %d. expected capacity: %d, got %d.", 5*time.Second, window.WindowSize(), 10, window.Capacity())
	}
}

func TestInitializationEmpty(t *testing.T) {
	window, err := fixedwindowcounter.Init(0, 5*time.Second)
	if err == nil {
		window.Stop()
		t.Error("initialization succeeded when it should've failed")
	}
}

func TestWindowNoRateLimit(t *testing.T) {
	window, err := fixedwindowcounter.Init(10, time.Second)
	if err != nil {
		t.Error("initialization failed")
	}
	defer window.Stop()

	for range 3 {
		for range 10 {
			if !window.Allow() {
				t.Error("rate limit exceeded during happy path")
			}
		}
		time.Sleep(time.Second + time.Millisecond*100)
	}

}
func TestWindowRateLimitExceeded(t *testing.T) {
	window, err := fixedwindowcounter.Init(10, time.Second)
	if err != nil {
		t.Error("initialization failed")
	}
	defer window.Stop()

	for range 10 {
		window.Allow()
	}

	if window.Allow() {
		t.Error("rate limit didn't exceed during unhappy path")
	}
}
