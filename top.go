package pore

// Limited size heap, mutexes ignored for now
type Top struct {
	container
}

func (t *Top) Push(v interface{}) {
	for i := 0; i < len(t.xs); i++ {
		if t.fn(v, t.xs[i]) {
			v, t.xs[i] = t.xs[i], v
		}
	}
}

func (t *Top) Pop() interface{} {
	val := t.xs[0]
	t.xs = append(t.xs[1:], nil)
	return val
}

func NewTop(fn Comparator, size int) Top {
	return Top{
		container{fn, make([]interface{}, size)},
	}
}
