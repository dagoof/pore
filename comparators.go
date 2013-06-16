package pore

type Comparator func(a, b interface{}) bool

func MaxInt(_a, _b interface{}) bool {
	if a, ok := _a.(int); ok {
		if b, ok := _b.(int); ok {
			return a > b
		}
	}
	return true
}

func MinInt(_a, _b interface{}) bool {
	if a, ok := _a.(int); ok {
		if b, ok := _b.(int); ok {
			return a < b
		}
	}
	return true
}
