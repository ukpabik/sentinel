package leakybucket

type Node struct {
	Data string
	Next *Node
	Prev *Node
}

type DLL struct {
	Head *Node
	Tail *Node
}

func DLLInitializer() *DLL {
	head := &Node{}
	tail := &Node{}
	head.Next = tail
	tail.Prev = head

	doubly := &DLL{
		Head: head,
		Tail: tail,
	}

	return doubly
}

func (dll *DLL) AddAtHead(value string) {
	if dll == nil {
		return
	}

	newNode := &Node{
		Data: value,
	}

	next := dll.Head.Next
	newNode.Prev = dll.Head
	dll.Head.Next = newNode
	next.Prev = newNode
	newNode.Next = next
}

func (dll *DLL) AddAtTail(value string) {
	if dll == nil {
		return
	}

	newNode := &Node{
		Data: value,
	}

	prev := dll.Tail.Prev
	dll.Tail.Prev = newNode
	prev.Next = newNode
	newNode.Next = dll.Tail
	newNode.Prev = prev
}

func (dll *DLL) RemoveNode(node *Node) {
	if dll == nil {
		return
	}

	if node == dll.Head || node == dll.Tail {
		return
	}

	prev := node.Prev
	next := node.Next

	if prev != nil {
		prev.Next = next
	}
	if next != nil {
		next.Prev = prev
	}

	next.Prev = prev
	prev.Next = next

	node.Prev = nil
	node.Next = nil
}

func (dll *DLL) RemoveFromTail() *Node {
	if dll == nil {
		return nil
	}

	node := dll.Tail.Prev
	if node == dll.Head {
		return nil
	}

	dll.RemoveNode(node)
	return node
}

func (dll *DLL) RemoveFromHead() *Node {
	if dll == nil {
		return nil
	}

	node := dll.Head.Next
	if node == dll.Tail {
		return nil
	}

	dll.RemoveNode(node)
	return node
}
