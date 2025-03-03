package arrays

// returns -1 if not found
func FirstIndexWhere[T any](elements []T, where func(T) bool) int {
	for i := 0; i < len(elements); i++ {
		if where(elements[i]) {
			return i
		}
	}
	return -1
}

// returns -1 if not found
func LastIndexWhere[T any](elements []T, where func(T) bool) int {
	for i := len(elements) - 1; i >= 0; i-- {
		if where(elements[i]) {
			return i
		}
	}
	return -1
}

// WhereIndex
