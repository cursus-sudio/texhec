package arrays

func Where[T any](elements []T, where func(T, int) bool) []T {
	var newElements []T
	for i, element := range elements {
		if where(element, i) {
			newElements = append(newElements, element)
		}
	}
	return newElements
}
