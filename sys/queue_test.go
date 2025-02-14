package sys

import (
	"strconv"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	q := NewQueue(QueueOptions{
		Capacity:  2,
		TickDelay: 1000 * 1,
		LoopDelay: 500,
		// OnInit:    func() {},
		PeriodicPush: func(queue map[string]interface{}) (item interface{}, key string, err error) {
			rd := Random(1, 2)
			time.Sleep(time.Duration(rd) * time.Second)
			return []string{"test", "test"}, strconv.Itoa(len(queue) + 1), err
		},
		OnPushed: func(item interface{}, key string) {
			t.Log("item pushed", item)
		},
		PeriodicRemove: func(queue map[string]interface{}) (string, error) {
			time.Sleep(time.Duration(Random(3, 4)) * time.Second)
			return strconv.Itoa(len(queue)), nil
		},
		OnRemoved: func(item interface{}, key string) {
			t.Log("item removed", item, key)
		},
		OnDrained: func() {
			// t.Log("queue empty")
		},
		OnError: func(err error) {
			t.Fatal(err)
		},
		OnTick: func(queue map[string]interface{}) {
			t.Log(len(queue))
		},
	})
	q.Start()
	// t.Log(q)
}
