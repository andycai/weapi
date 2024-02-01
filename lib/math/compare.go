package math

type Ordered interface {
	Integer | Float | ~string
}

// Integer
type Integer interface {
	Signed | Unsigned
}

// Float
type Float interface {
	~float32 | ~float64
}

// Signed
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Max 返回最大值
func Max[T Ordered](x, y T) T {
	if x > y {
		return x
	}

	return y
}

// Min 返回最小值
func Min[T Ordered](x, y T) T {
	if x < y {
		return x
	}

	return y
}
