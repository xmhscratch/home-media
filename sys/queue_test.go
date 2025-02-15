package sys

import (
	"strconv"
	"testing"
	"time"
)

type IFruit interface{}

type Fruit struct {
	QItem
	IFruit
}

func (ctx Fruit) Index() int {
	now := time.Now()
	return int(now.Unix())
}

func (ctx Fruit) Key() string {
	return strconv.FormatInt(int64(ctx.Index()), 5<<1)
}

func TestQueue(t *testing.T) {
	// q := NewQueue(QueueOptions{
	// 	Capacity:  2,
	// 	TickDelay: 1000 * 1,
	// 	LoopDelay: 500,
	// 	// OnInit:    func() {},
	// 	PeriodicPush: func(queue map[string]interface{}) (item interface{}, key string, err error) {
	// 		rd := Random(1, 2)
	// 		time.Sleep(time.Duration(rd) * time.Second)
	// 		return []string{"test", "test"}, strconv.Itoa(len(queue) + 1), err
	// 	},
	// 	OnPushed: func(item interface{}, key string) {
	// 		t.Log("item pushed", item)
	// 	},
	// 	PeriodicRemove: func(queue map[string]interface{}) (string, error) {
	// 		time.Sleep(time.Duration(Random(3, 4)) * time.Second)
	// 		return strconv.Itoa(len(queue)), nil
	// 	},
	// 	OnRemoved: func(item interface{}, key string) {
	// 		t.Log("item removed", item, key)
	// 	},
	// 	OnError: func(err error) {
	// 		t.Fatal(err)
	// 	},
	// 	OnTick: func(queue map[string]interface{}) {
	// 		t.Log(len(queue))
	// 	},
	// })

	q := NewQueue(QueueOptions[Fruit]{
		Capacity:  2,
		LoopDelay: 500,
		OnInit:    func() {},
		PeriodicPush: func(queue *QueueStack[Fruit]) (item Fruit, err error) {
			time.Sleep(time.Duration(1) * time.Second)
			return Fruit{}, err
		},
		OnPushed: func(item Fruit) {
			t.Log("item pushed", item.Key())
		},
		OnRemoved: func(item Fruit) {
			t.Log("item removed", item, item.Key())
		},
		OnError: func(err error) {
			t.Fatal(err)
		},
	})
	q.Start()
	// t.Log(q)
}
