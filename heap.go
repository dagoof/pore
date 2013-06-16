package pore

import (
	"sync"
)

type Heap struct {
	container
	c *sync.Cond
}

func (h *Heap) Push(v interface{}) {
	defer h.c.Signal()

	h.c.L.Lock()
	defer h.c.L.Unlock()

	for i := 0; i < len(h.xs); i++ {
		if h.fn(v, h.xs[i]) {
			v, h.xs[i] = h.xs[i], v
		}
	}
	h.xs = append(h.xs, v)
}

// Blocks while heap is empty
func (h *Heap) Pop() interface{} {
	h.c.L.Lock()
	defer h.c.L.Unlock()

	for len(h.xs) < 1 {
		h.c.Wait()
	}

	val := h.xs[0]
	h.xs = h.xs[1:]
	return val
}

// Push into the heap using a given channel
func (h *Heap) In(c chan interface{}) {
	for {
		h.Push(<-c)
	}
}

// Pop out of heap into a given channel
func (h *Heap) Out(c chan interface{}) {
	for {
		c <- h.Pop()
	}
}

// Push into heap on the returned channel
func (h *Heap) InC() chan interface{} {
	c := make(chan interface{})
	go h.In(c)
	return c
}

// Pop out of heap on the returned channel
func (h *Heap) OutC() chan interface{} {
	c := make(chan interface{})
	go h.Out(c)
	return c
}

// Create a new heap based on a comparator function
func NewHeap(fn Comparator) *Heap {
	return &Heap{
		container{fn: fn},
		sync.NewCond(&sync.Mutex{}),
	}
}
