package stack

//Stack struct represents actual Stack with Element and size as its properties
//Size increases or decreases with elements being pushed or popped
type Stack struct {
	top  *Node
	size int
}

//Node struct represents each node that forms a part of Stack with value of type Interface and prev representing previous node
type Node struct {
	value interface{}
	prev  *Node
}

//Size returns actual size of the Stack
func (s *Stack) Size() int {
	return s.size
}

//Empty returns bool if stack is empty
func (s *Stack) Empty() bool {
	if s.size == 0 {
		return true
	}
	return false
}

//Top returns value at top of the stack
func (s *Stack) Top() interface{} {
	return s.top.value
}

//Push pushes value on top of stack
func (s *Stack) Push(value interface{}) {
	//save top to prev
	prev := s.top
	s.top = &Node{value, prev}
	s.size++
}

//Pop pops value out of stack
func (s *Stack) Pop() (value interface{}) {
	if s.size > 0 {
		top := s.top
		//point new top to previous node
		s.top = top.prev
		s.size--
		return top.value
	}
	return nil
}
