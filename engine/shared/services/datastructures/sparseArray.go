package datastructures

import (
	"golang.org/x/exp/constraints"
)

type SparseArray[Index constraints.Unsigned, Value any] interface {
	Get(index Index) (value Value, ok bool)
	GetValues() []Value
	GetIndices() []Index
	// if false then updated
	Set(index Index, value Value) (added bool)
	Remove(index Index) (removed bool)
}

type sparseArray[Index constraints.Unsigned, Value any] struct {
	EmptyValue    Index
	valuesIndices []Index // here some indices have special meaning (read constants above)

	// both arrays below indices correspond
	values  []Value
	indices []Index // here value means index in sparse array
}

func NewSparseArray[Index constraints.Unsigned, Value any]() SparseArray[Index, Value] {
	var zero Index
	return &sparseArray[Index, Value]{
		EmptyValue: ^zero,
	}
}

func (a *sparseArray[Index, Value]) Get(index Index) (Value, bool) {
	if int(index) >= len(a.valuesIndices) {
		var zero Value
		return zero, false
	}

	valueIndex := a.valuesIndices[index]
	if valueIndex == a.EmptyValue {
		var zero Value
		return zero, false
	}

	value := a.values[valueIndex]
	return value, true
}

func (a *sparseArray[Index, Value]) GetValues() []Value  { return a.values }
func (a *sparseArray[Index, Value]) GetIndices() []Index { return a.indices }

func (a *sparseArray[Index, Value]) Set(index Index, value Value) bool {
	for int(index) >= len(a.valuesIndices) {
		a.valuesIndices = append(a.valuesIndices, a.EmptyValue)
	}

	valueIndex := a.valuesIndices[index]

	if valueIndex == a.EmptyValue {
		a.valuesIndices[index] = Index(len(a.values))
		a.values = append(a.values, value)
		a.indices = append(a.indices, index)
		return true
	}

	a.values[valueIndex] = value
	a.indices[valueIndex] = index
	return false
}

func (a *sparseArray[Index, Value]) Remove(index Index) bool {
	if int(index) >= len(a.valuesIndices) {
		return false
	}

	valueIndex := a.valuesIndices[index]
	if valueIndex == a.EmptyValue {
		return false
	}

	a.valuesIndices[index] = a.EmptyValue

	if len(a.values)-1 != int(valueIndex) {
		movedValue := a.values[len(a.values)-1]
		movedIndex := a.indices[len(a.indices)-1]

		a.values[valueIndex] = movedValue
		a.indices[valueIndex] = movedIndex

		a.valuesIndices[movedIndex] = valueIndex
	}

	a.values = a.values[:len(a.values)-1]
	a.indices = a.indices[:len(a.indices)-1]
	return true
}
