package sys

import (
	"testing"
	"time"
)

type TFruitName (string)

const (
	_                 = iota
	Apple  TFruitName = "Apple"
	Banana TFruitName = "Banana"
	Cherry TFruitName = "Cherry"
	Grape  TFruitName = "Grape"
	Mango  TFruitName = "Mango"
	Orange TFruitName = "Orange"
)

func (f TFruitName) String() string {
	return map[TFruitName]string{
		Apple:  "Apple",
		Banana: "Banana",
		Cherry: "Cherry",
		Grape:  "Grape",
		Mango:  "Mango",
		Orange: "Orange",
	}[f]
}

type Fruit struct {
	QItem[Fruit]
	T TFruitName
}

func (ctx Fruit) Index() int {
	return int(time.Now().Unix())
}

func (ctx Fruit) Key() string {
	return ctx.T.String()
}

func TestQueue(t *testing.T) {
	q := NewQueue(QueueOptions[Fruit]{
		Capacity: 3,
		Throttle: 1,
		OnInit: func(queue *QueueStack[Fruit]) error {
			queue.Push(&Fruit{T: "Apple"})
			queue.Push(&Fruit{T: "Banana"})
			queue.Push(&Fruit{T: "Cherry"})
			queue.Push(&Fruit{T: "Grape"})
			queue.Push(&Fruit{T: "Mango"})
			queue.Push(&Fruit{T: "Orange"})

			return nil
		},
		Periodic: func(queue *QueueStack[Fruit]) (item *Fruit, err error) {
			queue.Push(&Fruit{T: "Banana"})
			queue.Push(&Fruit{T: "Mango"})

			return nil, err
		},
		Consume: func(queue *QueueStack[Fruit], item *Fruit) error {
			time.Sleep(time.Duration(2) * time.Second)
			return nil
		},
		OnPushed: func(item *Fruit) {
			t.Log("Item pushed", item.Key())
		},
		OnConsumed: func(item *Fruit) {
			t.Log("Item consumed", item.Key())
		},
		OnError: func(err error) {
			t.Fatal(err)
		},
	})
	q.Start()
	// t.Log(q)
}
