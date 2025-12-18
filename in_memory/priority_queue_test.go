package inmemory_test

import (
	"container/heap"
	"testing"
	"time"

	inmemory "github.com/ukpabik/sentinel/in_memory"
)

func TestPriorityQueueOrdersByTimestamp(t *testing.T) {
	q := &inmemory.PriorityQueue{}
	heap.Init(q)

	base := time.Unix(1_700_000_000, 0)

	heap.Push(q, &inmemory.Log{ID: 2, Timestamp: base.Add(2 * time.Second)})
	heap.Push(q, &inmemory.Log{ID: 0, Timestamp: base.Add(0 * time.Second)})
	heap.Push(q, &inmemory.Log{ID: 1, Timestamp: base.Add(1 * time.Second)})

	for wantID := 0; wantID < 3; wantID++ {
		got := heap.Pop(q).(*inmemory.Log)
		if got.ID != wantID {
			t.Fatalf("wrong pop order: got ID=%d want ID=%d", got.ID, wantID)
		}
	}

	if q.Len() != 0 {
		t.Fatalf("expected empty queue, got len=%d", q.Len())
	}
}
