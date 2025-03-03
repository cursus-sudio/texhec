package arrays

func Select[TFrom any, TTo any](elements []TFrom, selector func(TFrom) TTo) []TTo {
	selected := make([]TTo, len(elements))
	for i, element := range elements {
		selected[i] = selector(element)
	}
	return selected
}
