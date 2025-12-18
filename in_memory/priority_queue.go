package inmemory

import "time"

type Log struct {
	ID        int
	Timestamp time.Time
}

type PriorityQueue []*Log

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Timestamp.Before(pq[j].Timestamp)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq PriorityQueue) Top() *Log {
	return pq[pq.Len()-1]
}

func (pq *PriorityQueue) Push(x any) {
	log := x.(*Log)
	*pq = append(*pq, log)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	log := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return log
}
