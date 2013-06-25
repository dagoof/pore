package pore

type Comparator func(a, b interface{}) bool

func MaxInt(a, b interface{}) bool {
	if _a, ok := a.(int); ok {
		if _b, ok := b.(int); ok {
			return _a > _b
		}
	}
	return false
}

func MinInt(a, b interface{}) bool {
	if _a, ok := a.(int); ok {
		if _b, ok := b.(int); ok {
			return _a < _b
		}
	}
	return false
}
