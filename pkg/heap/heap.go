package heap

import (
	"container/heap"
	"fmt"
)

type Heap[T Lesser[T]] struct {
	_heap _heap[T]
}

type Lesser[T any] interface {
	Less(T) bool
}

func (h *Heap[T]) Len() int {
	return h._heap.Len()
}

func (h *Heap[T]) Push(v T) {
	heap.Push(&h._heap, v)
}

func (h *Heap[T]) Pop() T {
	return heap.Pop(&h._heap).(T)
}

func (h *Heap[T]) Peek() T {
	return h._heap.ess[0]
}

func (h Heap[T]) String() string {
	return fmt.Sprint(h._heap.ess)
}

type _heap[T Lesser[T]] struct {
	ess []T
}

func (h *_heap[T]) Len() int           { return len(h.ess) }
func (h *_heap[T]) Swap(i, j int)      { h.ess[i], h.ess[j] = h.ess[j], h.ess[i] }
func (h *_heap[T]) Less(i, j int) bool { return h.ess[i].Less(h.ess[j]) }

func (h *_heap[T]) Push(v any) {
	h.ess = append(h.ess, v.(T))
}

func (h *_heap[T]) Pop() any {
	last := h.ess[len(h.ess)-1]
	h.ess = h.ess[:len(h.ess)-1]
	return last
}
