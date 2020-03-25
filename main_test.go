package main

import (
	"testing"

	"github.com/YouDad/blockchain/types"
)

func TestQueue(t *testing.T) {
	q := types.NewQueue(5)
	t.Run("Push & Pop", func(t *testing.T) {
		q.Push(1)
		t.Log(q.Get())
		q.Push(2)
		t.Log(q.Get())
		q.Push(3)
		t.Log(q.Get())
		q.Push(4)
		t.Log(q.Get())
		q.Push(5)
		t.Log(q.Get())
		q.Pop()
		t.Log(q.Get())
		q.Pop()
		t.Log(q.Get())
		q.Pop()
		t.Log(q.Get())
		q.Pop()
		t.Log(q.Get())
		q.Pop()
		t.Log(q.Get())
		q.Push(1)
		t.Log(q.Get())
		q.Push(2)
		t.Log(q.Get())
		q.Push(3)
		t.Log(q.Get())
		q.Push(4)
		t.Log(q.Get())
		q.Push(5)
		t.Log(q.Get())
	})
}
