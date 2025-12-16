package leakybucket_test

import (
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

	dll.AddAtHead(&leakybucket.Token{})

	if dll.Head.Next == dll.Tail || dll.Tail.Prev == dll.Head {
		t.Error("adding nodes didn't work")
	}
}

func TestDLLAddNodes(t *testing.T) {
	dll := leakybucket.DLLInitializer()

	for range 10 {
		dll.AddAtTail(&leakybucket.Token{})
	}

	curr := dll.Head.Next
	counter := 0
	for curr != dll.Tail {
		counter += 1
		curr = curr.Next
	}

	if counter != 10 {
		t.Errorf("expected size: %d, got: %d", 10, counter)
	}
}

func TestDLLRemoveNodes(t *testing.T) {
	dll := leakybucket.DLLInitializer()

	for range 10 {
		dll.AddAtTail(&leakybucket.Token{})
	}

	removed := dll.RemoveFromTail()
	if removed.Next != nil || removed.Prev != nil {
		t.Error("node disconnection failed")
	}
}
