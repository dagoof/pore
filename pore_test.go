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

func TestHeap(t *testing.T) {
	h := NewHeap(MaxInt)
	ints := []int{3, 6, 2, 8, 13, 4}

	for _, n := range ints {
		h.Push(n)
	}

	if len(h.All()) != 6 {
		t.Fatal("heap should grow as items are pushed")
	}

	var a, b int
	a = 1000

	for i := 0; i < len(ints); i++ {
		b = h.Pop().(int)
		if b > a {
			t.Fatal("ints are not in order")
		}
		a = b
	}

	if len(h.All()) != 0 {
		t.Fatal("heap should shrink as items are popped", h)
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
