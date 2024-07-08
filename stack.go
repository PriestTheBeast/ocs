package main

import "errors"

// StringStack represents a stack that holds a slice.
type StringStack struct {
	items []string
}

// Push adds an item to the top of the stack.
func (s *StringStack) Push(item string) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top item of the stack.
// It returns an error if the stack is empty.
func (s *StringStack) Pop() string {
	if s.IsEmpty() {
		panic(errors.New("stack is empty"))
	}
	topIndex := len(s.items) - 1
	item := s.items[topIndex]
	s.items = s.items[:topIndex]
	return item
}

// Peek returns the top item of the stack without removing it.
// It returns an error if the stack is empty.
func (s *StringStack) Peek() string {
	if s.IsEmpty() {
		panic(errors.New("stack is empty"))
	}
	return s.items[len(s.items)-1]
}

// IsEmpty checks if the stack is empty.
func (s *StringStack) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of items in the stack.
func (s *StringStack) Size() int {
	return len(s.items)
}
