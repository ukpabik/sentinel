package leakybucket

import (
	"fmt"
	"sync"
)

type ConcurrentQueue struct {
	nodes    *DLL
	capacity int
	size     int
	mutex    sync.Mutex
}

func ConcurrentQueueInitializer(capacity int) (*ConcurrentQueue, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("cannot create a queue of capacity <= 0")
	}

	return &ConcurrentQueue{
		nodes:    DLLInitializer(),
		capacity: capacity,
	}, nil
}

func (cq *ConcurrentQueue) PopLeft() (string, error) {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	if cq.size == 0 {
		return "", fmt.Errorf("queue is empty")
	}

	removed := cq.nodes.RemoveFromHead()
	cq.size -= 1
	return removed.Data, nil
}
func (cq *ConcurrentQueue) Pop() (string, error) {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	if cq.size == 0 {
		return "", fmt.Errorf("queue is empty")
	}

	removed := cq.nodes.RemoveFromTail()
	cq.size -= 1
	return removed.Data, nil
}
func (cq *ConcurrentQueue) AppendLeft(value string) error {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	if cq.size == cq.capacity {
		return fmt.Errorf("queue is full")
	}
	cq.size += 1
	cq.nodes.AddAtHead(value)
	return nil
}
func (cq *ConcurrentQueue) Append(value string) error {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	if cq.size == cq.capacity {
		return fmt.Errorf("queue is full")
	}
	cq.size += 1
	cq.nodes.AddAtTail(value)
	return nil
}
func (cq *ConcurrentQueue) Size() int {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	return cq.size
}
func (cq *ConcurrentQueue) Capacity() int {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	return cq.capacity
}
