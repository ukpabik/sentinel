package leakybucket_test

import (
	"fmt"
	"strconv"
	"sync"
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

	for i := range 10 {
		cq.Append(strconv.Itoa(i))
	}

	size := cq.Size()
	for i := size - 1; i >= 0; i-- {
		popped, err := cq.Pop()
		if err != nil {
			t.Error("error while popping")
		}
		expected := strconv.Itoa(i)
		if expected != popped {
			t.Errorf("expected value: %s, got value: %s", expected, popped)
		}
	}
}

func TestConcurrentQueueAppendLeftAndPopLeft(t *testing.T) {
	cq, err := leakybucket.ConcurrentQueueInitializer(10)
	if err != nil {
		t.Error("unable to initialize concurrent queue")
	}

	for i := range 10 {
		cq.AppendLeft(strconv.Itoa(i))
	}

	size := cq.Size()
	for i := size - 1; i >= 0; i-- {
		popped, err := cq.PopLeft()
		if err != nil {
			t.Error("error while popping")
		}
		expected := strconv.Itoa(i)
		if expected != popped {
			t.Errorf("expected value: %s, got value: %s", expected, popped)
		}
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

	wg.Add(N)
	for pid := 0; pid < N; pid++ {
		pid := pid
		go func() {
			defer wg.Done()
			<-start

			for i := 0; i < M; i++ {
				v := fmt.Sprintf("p%d-%d", pid, i)
				if err := cq.Append(v); err != nil {
					errCh <- err
					return
				}
			}
		}()
	}

	close(start)
	wg.Wait()
	close(errCh)

	for range errCh {
		t.Error("append failed")
	}

	seen := make(map[string]struct{}, N*M)
	for cq.Size() > 0 {
		v, err := cq.PopLeft()
		if err != nil {
			t.Error("pop failed")
		}

		if _, ok := seen[v]; ok {
			t.Error("duplicate value popped")
		}
		seen[v] = struct{}{}
	}

	if expected, got := len(seen), N*M; expected != got {
		t.Errorf("wrong total items: got %d, expected %d", got, expected)
	}
}
