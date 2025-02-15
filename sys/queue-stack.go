package sys

import "container/heap"

type QItem interface {
	Index() int
	Key() string
}

type QueueStack[I QItem] ([]QItem)

func (h QueueStack[I]) Len() int {
	return len(h)
}
func (h QueueStack[I]) Less(i int, j int) bool { return h[i].Index() > h[j].Index() }
func (h QueueStack[I]) Swap(i int, j int)      { h[i], h[j] = h[j], h[i] }

func (h QueueStack[I]) FindKeyIndex(key string) int {
	for i, v := range h {
		if v.Key() != key {
			continue
		}
		return i
	}
	return -1
}

func (h *QueueStack[I]) Push(qItem any) {
	if h.FindKeyIndex(qItem.(QItem).Key()) != -1 {
		return
	}
	*h = append(*h, qItem.(QItem))
}

func (h *QueueStack[I]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func NewQueueStack[I QItem]() *QueueStack[I] {
	var ctx *QueueStack[I] = &QueueStack[I]{}
	heap.Init(ctx)

	return ctx
}
