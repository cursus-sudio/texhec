package arrays

func MoveElement[T any](slice []T, from, to int) []T {
	// Check if indices are valid
	if from < 0 || from >= len(slice) || to < 0 || to >= len(slice) {
		panic("Invalid index")
	}

	// Remove the element at the 'from' index
	element := slice[from]
	slice = append(slice[:from], slice[from+1:]...)

	// Insert the element at the 'to' index
	slice = append(slice[:to], append([]T{element}, slice[to:]...)...)

	return slice
}
