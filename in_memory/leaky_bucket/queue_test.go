package leakybucket_test

import (
	"sync"
	"sync/atomic"
	"testing"

	leakybucket "github.com/ukpabik/sentinel/in_memory/leaky_bucket"
)

func TestConcurrentQueueInitialization(t *testing.T) {
	cq, err := leakybucket.ConcurrentQueueInitializer(10)
	if err != nil {
		t.Error("unable to initialize concurrent queue")
	}

	if cq.Size() != 0 {
		t.Errorf("expected size: %d, got value %d", 0, cq.Size())
	}

	if cq.Capacity() != 10 {
		t.Errorf("expected size: %d, got value %d", 10, cq.Capacity())
	}
}

func TestConcurrentQueueAppendAndPop(t *testing.T) {
	cq, err := leakybucket.ConcurrentQueueInitializer(10)
	if err != nil {
		t.Error("unable to initialize concurrent queue")
	}

	for range 10 {
		if err := cq.Append(); err != nil {
			t.Errorf("failed to append: %v", err)
		}
	}

	if cq.Size() != 10 {
		t.Errorf("expected size 10, got %d", cq.Size())
	}

	for i := 0; i < 10; i++ {
		_, err := cq.Pop()
		if err != nil {
			t.Error("error while popping")
		}
	}

	if cq.Size() != 0 {
		t.Errorf("expected size 0, got %d", cq.Size())
	}
}

func TestConcurrentQueueAppendLeftAndPopLeft(t *testing.T) {
	cq, err := leakybucket.ConcurrentQueueInitializer(10)
	if err != nil {
		t.Error("unable to initialize concurrent queue")
	}

	for range 10 {
		if err := cq.AppendLeft(); err != nil {
			t.Errorf("failed to append left: %v", err)
		}
	}

	if cq.Size() != 10 {
		t.Errorf("expected size 10, got %d", cq.Size())
	}

	for i := 0; i < 10; i++ {
		_, err := cq.PopLeft()
		if err != nil {
			t.Error("error while popping left")
		}
	}

	if cq.Size() != 0 {
		t.Errorf("expected size 0, got %d", cq.Size())
	}
}

func TestConcurrentQueueMultipleThreads(t *testing.T) {
	const (
		N = 20
		M = 100
	)

	cq, err := leakybucket.ConcurrentQueueInitializer(N * M)
	if err != nil {
		t.Error("unable to initialize concurrent queue")
	}

	start := make(chan struct{})
	var wg sync.WaitGroup
	errCh := make(chan error, N*M)

	var successCount int64

	wg.Add(N)
	for pid := 0; pid < N; pid++ {
		go func() {
			defer wg.Done()
			<-start

			for i := 0; i < M; i++ {
				if err := cq.Append(); err != nil {
					errCh <- err
					return
				}
				atomic.AddInt64(&successCount, 1)
			}
		}()
	}

	close(start)
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Errorf("append failed: %v", err)
	}

	if expected, got := int(successCount), cq.Size(); expected != got {
		t.Errorf("wrong queue size: got %d, expected %d", got, expected)
	}

	poppedCount := 0
	for cq.Size() > 0 {
		_, err := cq.PopLeft()
		if err != nil {
			t.Error("pop failed")
		}
		poppedCount++
	}

	if expected, got := int(successCount), poppedCount; expected != got {
		t.Errorf("wrong total items popped: got %d, expected %d", got, expected)
	}
}
