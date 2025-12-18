package slidingwindowcounter_test

import (
	"testing"
	"time"

	slidingwindowcounter "github.com/ukpabik/sentinel/in_memory/sliding_window_counter"
)

func TestInitializationHappyPath(t *testing.T) {
	_, err := slidingwindowcounter.Init(10, time.Second)
	if err != nil {
		t.Error("initialization failed")
	}

	_, err = slidingwindowcounter.Init(0, time.Second)
	if err == nil {
		t.Error("initialization succeeded when it should've failed")
	}
}

func TestWindowAllowHappyPath(t *testing.T) {
	swc, err := slidingwindowcounter.Init(10, time.Second)
	if err != nil {
		t.Error("initialization failed")
	}

	for range 10 {
		swc.Allow()
	}
	if swc.Allow() {
		t.Error("window allowed a request when it was full")
	}
}

func TestWindowNoOverlap(t *testing.T) {
	swc, err := slidingwindowcounter.Init(10, time.Second)
	if err != nil {
		t.Error("initialization failed")
	}

	for range 10 {
		swc.Allow()
	}
	time.Sleep(50 * time.Millisecond)
	if swc.Allow() {
		t.Error("window allowed a request when it was full")
	}
}

func TestWindowWithOverlap(t *testing.T) {
	swc, err := slidingwindowcounter.Init(10, 100*time.Millisecond)
	if err != nil {
		t.Error("initialization failed")
	}

	for range 10 {
		swc.Allow()
	}
	time.Sleep(120 * time.Millisecond)
	if !swc.Allow() {
		t.Error("window should not be full")
	}
}
