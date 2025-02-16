package sys

import (
	"log"
	"sync"
	"time"
)

type OnInitFunc[I QItem[I]] func(queue *QueueStack[I]) error
type PeriodicFunc[I QItem[I]] func(queue *QueueStack[I]) (*I, error)
type ConsumeFunc[I QItem[I]] func(queue *QueueStack[I], item *I) error
type OnPushedFunc[I QItem[I]] func(item *I)
type OnConsumedFunc[I QItem[I]] func(item *I)
type OnTickFunc[I QItem[I]] func(queue *QueueStack[I])
type OnErrorFunc func(err error)

type QueueOptions[I QItem[I]] struct {
	Capacity   int
	Throttle   int64
	OnInit     OnInitFunc[I]
	Periodic   PeriodicFunc[I]
	Consume    ConsumeFunc[I]
	OnPushed   OnPushedFunc[I]
	OnConsumed OnConsumedFunc[I]
	OnTick     OnTickFunc[I]
	OnError    OnErrorFunc
}

type Queue[I QItem[I]] struct {
	*QueueOptions[I]
	items *QueueStack[I]
	queue chan *I
}

func NewQueue[I QItem[I]](opts QueueOptions[I]) *Queue[I] {
	if opts.Capacity < 1 {
		opts.Capacity = 1
	}

	if opts.OnError == nil {
		opts.OnError = func(err error) { log.Panic(err) }
	}
	if opts.OnInit == nil {
		opts.OnInit = func(*QueueStack[I]) error { return nil }
	}
	if opts.OnPushed == nil {
		opts.OnPushed = func(*I) {}
	}
	if opts.OnConsumed == nil {
		opts.OnConsumed = func(*I) {}
	}
	if opts.OnTick == nil {
		opts.OnTick = func(*QueueStack[I]) {}
	}

	q := &Queue[I]{
		QueueOptions: &opts,
	}

	q.items = NewQueueStack[I]()
	q.queue = make(chan *I, q.Capacity)

	return q
}

func (q *Queue[I]) Start() {
	var (
		mu sync.Mutex
	)

	defer close(q.queue)

	q.OnInit(q.items)

	go func() {
	exitQueue:
		for {
			time.Sleep(time.Duration(q.Throttle) * time.Millisecond)

			if item, err := q.Periodic(q.items); err != nil {
				q.OnError(err)
				break exitQueue
			} else {
				if item == nil {
					continue
				}
				q.items.Push(item)
			}
		}
	}()

	for {
		time.Sleep(time.Duration(q.Throttle) * time.Millisecond)

		// 1-turn filling + 1-turn consume
		loopCircle := q.Capacity + q.Capacity

	exitFilling:
		for i := 0; i < loopCircle; i++ {
			time.Sleep(time.Duration(q.Throttle) * time.Millisecond)

			// continue popping item until queue reach its capacity
			if s := func() int {
				if q.items.Len() > 0 {
					item := q.items.Pop().(I)
					q.queue <- &item
					q.OnPushed(&item)
				}
				return len(q.queue)
			}(); s >= q.Capacity {
				go func() {
					mu.Lock()
					defer mu.Unlock()
					item := <-q.queue
					if err := q.Consume(q.items, item); err != nil {
						// put back item to the queue
						q.queue <- item
						q.OnError(err)
						return
					}
					go q.OnConsumed(item)
				}()
				// consume leftover items on the queue
				if q.items.Len() == 0 && s > 0 {
					continue
				}
				break exitFilling
			} else {
				// continue filling up the queue
				continue
			}
		}
	}
}

func (q *Queue[I]) Stop() {}
