/*
Alternative imagining of heap API. Rather than defining repetitive methods
for each type, create a heap with one function that acts as your heap.Less
method, and interface with it through typical Push/Pop methods, or channels.
pore.Heap is concurrency-safe and Pop blocks while waiting for Push
*/
package pore

type container struct {
	fn Comparator
	xs []interface{}
}

func (c *container) All() []interface{} {
	return c.xs
}

type Container interface {
	Push(interface{})
	Pop() interface{}
}
