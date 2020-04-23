package types

import (
	"errors"
	"fmt"
)

type Queue struct {
	array      []interface{}
	capability int
	front      int
	back       int
}

var (
	ErrInit = errors.New("Queue need initize")
)

func NewQueue(cap int) Queue {
	cap += 1
	return Queue{
		make([]interface{}, cap),
		cap, 0, 0,
	}
}

func (q *Queue) String() string {
	if q.array == nil {
		panic(ErrInit)
	}
	return fmt.Sprintf("Queue[%d,%d/%d]:%+v", q.front, q.back, q.capability, q.array)
}

func (q *Queue) Push(data interface{}) {
	if q.array == nil {
		panic(ErrInit)
	}

	if (q.back+1)%q.capability == q.front {
		q.Pop()
	}

	q.array[q.back] = data
	q.back = (q.back + 1) % q.capability
}

func (q *Queue) Pop() interface{} {
	if q.array == nil {
		panic(ErrInit)
	}
	ret := q.array[q.front]
	q.front = (q.front + 1) % q.capability
	return ret
}

func (q *Queue) Len() int {
	if q.array == nil {
		panic(ErrInit)
	}
	return len(q.array)
}

func (q *Queue) Get() []interface{} {
	if q.back >= q.front {
		return q.array[q.front:q.back]
	} else {
		return append(q.array[q.front:], q.array[:q.back]...)
	}
}
