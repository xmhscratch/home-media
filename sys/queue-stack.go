package sys

import (
	"container/heap"
)

type QItem[I any] interface {
	Index() int
	Key() string
}

type QueueStack[I QItem[I]] ([]I)

func (h QueueStack[I]) Len() int {
	return len(h)
}
func (h QueueStack[I]) Less(i int, j int) bool { return h[i].Index() < h[j].Index() }
func (h QueueStack[I]) Swap(i int, j int)      { h[i], h[j] = h[j], h[i] }

func (h QueueStack[I]) findKeyIndex(key string) int {
	for j, v := range h {
		if I.Key(v) != key {
			continue
		}
		return j
	}
	return -1
}

func (h *QueueStack[I]) Push(i any) {
	qItem := *i.(*I)
	if h.findKeyIndex(I.Key(qItem)) != -1 {
		return
	}
	*h = append(*h, qItem)
}

func (h *QueueStack[I]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func NewQueueStack[I QItem[I]]() *QueueStack[I] {
	var ctx *QueueStack[I] = &QueueStack[I]{}
	heap.Init(ctx)

	return ctx
}
