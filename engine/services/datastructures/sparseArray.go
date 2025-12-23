package datastructures

import (
	"golang.org/x/exp/constraints"
)

type SparseArray[Index constraints.Integer, Value any] interface {
	Get(index Index) (value Value, ok bool)
	GetValues() []Value
	GetIndices() []Index
	// if false then updated
	Set(index Index, value Value) (added bool)
	Remove(index Index) (removed bool)
}

type sparseArray[Index constraints.Integer, Value any] struct {
	EmptyValue    Index
	ValuesIndices []Index // here some indices have special meaning (read constants above)

	// both arrays below indices correspond
	Values  []Value
	Indices []Index // here value means index in sparse array
}

func NewSparseArray[Index constraints.Integer, Value any]() SparseArray[Index, Value] {
	var zero Index
	return &sparseArray[Index, Value]{
		EmptyValue: ^zero,
	}
}

func (a *sparseArray[Index, Value]) Get(index Index) (Value, bool) {
	if int(index) >= len(a.ValuesIndices) {
		var zero Value
		return zero, false
	}

	valueIndex := a.ValuesIndices[index]
	if valueIndex == a.EmptyValue {
		var zero Value
		return zero, false
	}

	value := a.Values[valueIndex]
	return value, true
}

func (a *sparseArray[Index, Value]) GetValues() []Value  { return a.Values }
func (a *sparseArray[Index, Value]) GetIndices() []Index { return a.Indices }

func (a *sparseArray[Index, Value]) Set(index Index, value Value) bool {
	for int(index) >= len(a.ValuesIndices) {
		a.ValuesIndices = append(a.ValuesIndices, a.EmptyValue)
	}

	valueIndex := a.ValuesIndices[index]

	if valueIndex == a.EmptyValue {
		a.ValuesIndices[index] = Index(len(a.Values))
		a.Values = append(a.Values, value)
		a.Indices = append(a.Indices, index)
		return true
	}

	a.Values[valueIndex] = value
	a.Indices[valueIndex] = index
	return false
}

func (a *sparseArray[Index, Value]) Remove(index Index) bool {
	if int(index) >= len(a.ValuesIndices) {
		return false
	}

	valueIndex := a.ValuesIndices[index]
	if valueIndex == a.EmptyValue {
		return false
	}

	a.ValuesIndices[index] = a.EmptyValue

	if len(a.Values)-1 != int(valueIndex) {
		movedValue := a.Values[len(a.Values)-1]
		movedIndex := a.Indices[len(a.Indices)-1]

		a.Values[valueIndex] = movedValue
		a.Indices[valueIndex] = movedIndex

		a.ValuesIndices[movedIndex] = valueIndex
	}

	a.Values = a.Values[:len(a.Values)-1]
	a.Indices = a.Indices[:len(a.Indices)-1]
	return true
}
