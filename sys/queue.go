package sys

import (
	"fmt"
	"time"
)

type PeriodicPushFunc[I QItem] func(queue *QueueStack[I]) (I, error)
type OnPushedFunc[I QItem] func(item I)
type OnRemovedFunc[I QItem] func(item I)
type OnTickFunc[I QItem] func(queue *QueueStack[I])
type OnErrorFunc func(err error)

type QueueOptions[I QItem] struct {
	Capacity     int
	LoopDelay    int64
	OnInit       func()
	PeriodicPush PeriodicPushFunc[I]
	OnPushed     OnPushedFunc[I]
	OnRemoved    OnRemovedFunc[I]
	OnTick       OnTickFunc[I]
	OnError      OnErrorFunc
}

type Queue[I QItem] struct {
	*QueueOptions[I]
	loopDelay time.Duration
	items     *QueueStack[I]
	queue     chan I
}

func NewQueue[I QItem](opts QueueOptions[I]) *Queue[I] {
	if opts.Capacity < 1 {
		opts.Capacity = 1
	}

	if opts.OnInit == nil {
		opts.OnInit = func() {}
	}
	if opts.OnPushed == nil {
		opts.OnPushed = func(I) {}
	}
	if opts.OnRemoved == nil {
		opts.OnRemoved = func(I) {}
	}
	if opts.OnError == nil {
		opts.OnError = func(error) {}
	}
	if opts.OnTick == nil {
		opts.OnTick = func(*QueueStack[I]) {}
	}

	q := &Queue[I]{
		QueueOptions: &opts,
	}

	q.loopDelay = time.Duration(q.LoopDelay) * time.Nanosecond
	q.items = NewQueueStack[I]()
	q.queue = make(chan I, q.Capacity)

	return q
}

func (q *Queue[I]) Start() {
	var (
		waitRemoval chan struct{} = make(chan struct{})
		waitOnEmpty chan struct{} = make(chan struct{}, 1)
	)

	defer close(q.queue)
	defer close(waitRemoval)
	defer close(waitOnEmpty)

	go func() {
		for {
			var (
				err  error
				item I
			)

			if len(q.queue) >= q.Capacity {
				return
			}

			if item, err = q.PeriodicPush(q.items); err != nil {
				q.OnError(err)
				return
			}
			if q.items.Len() == 0 {
				waitOnEmpty <- struct{}{}
			}
			q.items.Push(item)

			time.Sleep(q.loopDelay)
		}
	}()

	fmt.Println("start queue...")

	for {
		time.Sleep(q.loopDelay)

		// lock on empty queue
		<-waitOnEmpty

	queueBlocking:
		for i := 0; i < q.Capacity; i++ {
			time.Sleep(q.loopDelay)

			if s := func() int {
				if q.items.Len() > 0 {
					item := q.items.Pop().(I)
					q.queue <- item
					q.OnPushed(item)
				}
				return len(q.queue)
			}(); s >= q.items.Len() {
				go func() {
					item := <-q.queue
					q.OnRemoved(item)
					waitRemoval <- struct{}{}
				}()

				if q.items.Len() == 0 && s > 0 {
					continue
				} else {
					break queueBlocking
				}
			} else {
				continue
			}
		}
		<-waitRemoval
		// litter.D("report:", q.queue, len(*q.items))
	}
}

func (q *Queue[I]) Stop() {}
