package core

import (
	"time"
)

type PeriodicPushFunc func(queue map[string]interface{}) (interface{}, string, error)
type OnPushedFunc func(item interface{}, key string)
type PeriodicRemoveFunc func(queue map[string]interface{}) (string, error)
type OnRemovedFunc func(item interface{}, key string)
type OnErrorFunc func(err error)
type OnTickFunc func(queue map[string]interface{})

type QueueOptions struct {
	Capacity       int
	TickDelay      int64
	LoopDelay      int64
	OnInit         func()
	PeriodicPush   PeriodicPushFunc
	OnPushed       OnPushedFunc
	PeriodicRemove PeriodicRemoveFunc
	OnRemoved      OnRemovedFunc
	OnDrained      func()
	OnError        OnErrorFunc
	OnTick         OnTickFunc
}

type Queue struct {
	*QueueOptions
	// IsPaused  bool
	loopDelay time.Duration
	tickDelay time.Duration
	queue     map[string]interface{}
}

func NewQueue(opts QueueOptions) *Queue {
	if opts.Capacity < 1 {
		opts.Capacity = 1
	}

	if opts.OnInit == nil {
		opts.OnInit = func() {}
	}
	if opts.OnPushed == nil {
		opts.OnPushed = func(interface{}, string) {}
	}
	if opts.OnRemoved == nil {
		opts.OnRemoved = func(interface{}, string) {}
	}
	if opts.OnDrained == nil {
		opts.OnDrained = func() {}
	}
	if opts.OnError == nil {
		opts.OnError = func(error) {}
	}
	if opts.OnTick == nil {
		opts.OnTick = func(map[string]interface{}) {}
	}

	q := &Queue{
		QueueOptions: &opts,
	}

	q.loopDelay = time.Duration(q.LoopDelay) * time.Nanosecond
	q.tickDelay = time.Duration(q.TickDelay) * time.Millisecond
	q.queue = make(map[string]interface{}, q.Capacity)

	return q
}

func (q *Queue) Start() {
	var (
		readyRemoving chan struct{} = make(chan struct{})
		readyTicking  chan struct{} = make(chan struct{})
		readyPushing  chan struct{} = make(chan struct{})
	)

	defer close(readyRemoving)
	defer close(readyTicking)
	defer close(readyPushing)

	go func() {
		readyRemoving <- struct{}{}
		time.Sleep(q.loopDelay)
		readyTicking <- struct{}{}
		time.Sleep(q.loopDelay)
		readyPushing <- struct{}{}

		q.OnInit()
	}()

	go func() {
		for {
			q._watchRemoving()
			time.Sleep(q.loopDelay)
		}
	}()
	<-readyRemoving

	go func() {
		for {
			q._startTicking()
			time.Sleep(q.tickDelay)
		}
	}()
	<-readyTicking

	go func() {
		for {
			q._startPushing()
			time.Sleep(q.loopDelay)
		}
	}()
	<-readyPushing

	select {}
}

func (q *Queue) _startPushing() {
	var (
		err  error
		key  string
		item interface{}
	)

	if len(q.queue) >= q.Capacity {
		return
	}
	if item, key, err = q.PeriodicPush(q.queue); err != nil {
		q.OnError(err)
		return
	}
	if item == nil {
		return
	}

	q.queue[key] = item
	q.OnPushed(item, key)
}

func (q *Queue) _watchRemoving() {
	var (
		err error
		key string
	)
	if len(q.queue) == 0 {
		q.OnDrained()
		return
	}
	if key, err = q.PeriodicRemove(q.queue); err != nil {
		q.OnError(err)
		return
	}
	item := q.queue[key]
	if item == nil {
		return
	}
	delete(q.queue, key)
	q.OnRemoved(item, key)
}

func (q *Queue) _startTicking() {
	q.OnTick(q.queue)
}

func (q *Queue) Drain() {
	for key := range q.queue {
		delete(q.queue, key)
	}
}

func (q *Queue) Stop() {}
