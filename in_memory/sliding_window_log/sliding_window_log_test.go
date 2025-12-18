package slidingwindowlog_test

import (
	"testing"
	"time"

	slidingwindowlog "github.com/ukpabik/sentinel/in_memory/sliding_window_log"
)

func TestInitializationHappyPath(t *testing.T) {
	window, err := slidingwindowlog.Init(10, time.Second)
	if err != nil {
		t.Error("initialization failure")
	}

	if window.AllowedCount() != 10 {
		t.Errorf("expected: %d, got: %d", 10, window.AllowedCount())
	}
}

func TestInitializationEmpty(t *testing.T) {
	_, err := slidingwindowlog.Init(0, time.Second)
	if err == nil {
		t.Error("initialization succeeded in unhappy path")
	}
}

func TestWindowResetsAfterTimePeriod(t *testing.T) {
	window, err := slidingwindowlog.Init(10, 100*time.Millisecond)
	if err != nil {
		t.Error("initialization failure")
	}

	for range 10 {
		window.Allow()
	}

	if window.Allow() {
		t.Error("window should be full")
	}

	time.Sleep(150 * time.Millisecond)

	if !window.Allow() {
		t.Error("window should have reset")
	}
}
