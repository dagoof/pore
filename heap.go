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

	h.xs = append(h.xs, v)
	c := len(h.xs) - 1
	m := (c - 1) / 2
	for (c > 0) && h.fn(h.xs[c], h.xs[m]) {
		h.xs[c], h.xs[m] = h.xs[m], h.xs[c]
		c = m
		m = (m - 1) / 2
	}
}

// Blocks while heap is empty
func (h *Heap) Pop() interface{} {
	h.c.L.Lock()
	defer h.c.L.Unlock()

	// Wait for push signal
	for len(h.xs) < 1 {
		h.c.Wait()
	}

	val := h.xs[0]

	// Swap head and resize
	var a, b, c, ln int
	ln = len(h.xs) - 1
	h.xs[c] = h.xs[ln]
	h.xs = h.xs[:ln]

	for {
		// child addresses
		a = c*2 + 1
		b = a + 1

		if !(a < ln) {
			break
		}
		// Find child to swap with
		if b < ln && h.fn(h.xs[b], h.xs[a]) {
			a = b
		}
		// See if we should swap
		if !h.fn(h.xs[a], h.xs[c]) {
			break
		}
		// swap and continue
		h.xs[a], h.xs[c] = h.xs[c], h.xs[a]
		c = a
	}
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
