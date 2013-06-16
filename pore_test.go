package pore

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkHeap(b *testing.B) {
	h := NewHeap(MaxInt)
	n := b.N

	for i := 0; i < n; i++ {
		h.Push(rand.Int())
	}
	for i := 0; i < n; i++ {
		h.Pop()
	}
}

func TestTop(t *testing.T) {
	top_five := NewTop(MaxInt, 5)

	top_five.Push(3)
	top_five.Push(6)
	top_five.Push(2)
	top_five.Push(8)
	top_five.Push(13)
	top_five.Push(4)

	if len(top_five.All()) != 5 {
		t.Fatal("Size should be limited to 5")
	}

	if top_five.Pop().(int) != 13 {
		t.Fatal("top int should be 13")
	}

	top_five.Pop()
	top_five.Pop()
	top_five.Pop()
	top_five.Pop()
	if len(top_five.All()) != 5 {
		t.Fatal("Size should be maintained despite pops")
	}
}

func TestHeap(t *testing.T) {
	top_ints := NewHeap(MaxInt)

	top_ints.Push(3)
	top_ints.Push(6)
	top_ints.Push(2)
	top_ints.Push(8)
	top_ints.Push(13)
	top_ints.Push(4)

	if len(top_ints.All()) != 6 {
		t.Fatal("heap should grow as items are pushed")
	}

	if top_ints.Pop().(int) != 13 {
		t.Fatal("first int should be 13")
	}

	top_ints.Pop()

	if len(top_ints.All()) != 4 {
		t.Fatal("heap should shrink as items are popped")
	}

}

func TestHeapIO(t *testing.T) {
	h := NewHeap(MaxInt)
	in := make(chan interface{})
	out := make(chan interface{})

	go h.In(in)
	in <- 11
	in <- 5
	in <- 115

	go h.Out(out)
	if (<-out).(int) != 115 {
		t.Fatal("heap not working")
	}
	// Clear heap
	<-out
	<-out

	in <- 11
	in <- 5
	in <- 115

	if (<-out).(int) != 11 {
		t.Fatal("11 should be the first placed in channel. Highest at time of pop from heap")
	}

	if (<-out).(int) != 115 {
		t.Fatal("Now heap should work correctly")
	}

}

func TestHeapPopWait(t *testing.T) {
	h := NewHeap(MaxInt)
	c := make(chan bool)

	go func() {
		if h.Pop().(int) == 5 {
			c <- true
		}
	}()

	time.Sleep(time.Second / 4)
	h.Push(5)

	select {
	case <-time.After(time.Second):
		t.Fatal("heap Pop never fired")
	case <-c:
		return
	}

}

func TestHeapSync(t *testing.T) {
	h := NewHeap(MaxInt)
	out := h.OutC()

	PushN := func(h *Heap, n int) {
		for i := 0; i < n; i++ {
			h.Push(rand.Int())
		}
	}

	Pop := func(h *Heap) {
		for {
			h.Pop()
		}
	}

	go Pop(h)
	go Pop(h)

	go PushN(h, 500)
	go PushN(h, 500)

	go Pop(h)
	go Pop(h)

	go PushN(h, 500)
	go PushN(h, 500)

	for {
		select {
		case <-time.After(time.Second / 4):
			if len(h.All()) > 0 {
				t.Fatal("Heap was not drained correctly")
			}
			return
		case <-out:
			continue
		}
	}

}
