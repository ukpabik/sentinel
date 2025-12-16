package leakybucket_test

import (
	"strconv"
	"testing"

	leakybucket "github.com/ukpabik/sentinel/in_memory/leaky_bucket"
)

func TestDLLInitialization(t *testing.T) {
	dll := leakybucket.DLLInitializer()

	if dll.Head == nil || dll.Tail == nil {
		t.Error("doubly linked list didn't initialize properly")
	}

	if dll.Head.Next != dll.Tail || dll.Tail.Prev != dll.Head {
		t.Error("head and tail nodes are not connected properly")
	}
}

func TestDLLAddAtHead(t *testing.T) {
	dll := leakybucket.DLLInitializer()

	dll.AddAtHead("1")

	if dll.Head.Next.Data != "1" || dll.Tail.Prev.Data != "1" {
		t.Error("adding at head didn't work")
	}

	dll.AddAtTail("2")
	if dll.Tail.Prev.Data != "2" {
		t.Error("adding at tail didn't work")
	}
}

func TestDLLAddNodes(t *testing.T) {
	dll := leakybucket.DLLInitializer()

	for i := range 10 {
		dll.AddAtTail(strconv.Itoa(i))
	}

	curr := dll.Head.Next
	counter := 0
	for curr != dll.Tail {
		if curr.Data != strconv.Itoa(counter) {
			t.Errorf("expected value %d, got value %s", counter, curr.Data)
		}
		counter += 1
		curr = curr.Next
	}

}

func TestDLLRemoveNodes(t *testing.T) {
	dll := leakybucket.DLLInitializer()

	for i := range 10 {
		dll.AddAtTail(strconv.Itoa(i))
	}

	removed := dll.RemoveFromTail()
	if removed.Data != "9" {
		t.Errorf("expected value: %d, got value %s", 9, removed.Data)
	}
}
