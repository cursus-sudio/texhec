package datastructures

import (
	"golang.org/x/exp/constraints"
)

type SparseSetReader[Index constraints.Integer] interface {
	Get(index Index) (ok bool)
	GetIndices() []Index
}

type SparseSet[Index constraints.Integer] interface {
	SparseSetReader[Index]
	Add(index Index) (added bool)
	Remove(index Index) (removed bool)
}

type sparseSet[Index constraints.Integer] struct {
	EmptyValue    Index
	valuesIndices []Index // here some indices have special meaning (read constants above)

	// both arrays below indices correspond
	indices []Index // here value means index in sparse array
}

func NewSparseSet[Index constraints.Integer]() SparseSet[Index] {
	var zero Index
	return &sparseSet[Index]{
		EmptyValue: ^zero,
	}
}

func (a *sparseSet[Index]) Get(index Index) bool {
	if int(index) >= len(a.valuesIndices) {
		return false
	}

	valueIndex := a.valuesIndices[index]
	if valueIndex == a.EmptyValue {
		return false
	}

	return true
}

func (a *sparseSet[Index]) GetIndices() []Index { return a.indices }

func (a *sparseSet[Index]) Add(index Index) bool {
	for int(index) >= len(a.valuesIndices) {
		a.valuesIndices = append(a.valuesIndices, a.EmptyValue)
	}

	valueIndex := a.valuesIndices[index]

	if valueIndex == a.EmptyValue {
		a.valuesIndices[index] = Index(len(a.indices))
		a.indices = append(a.indices, index)
		return true
	}

	a.indices[valueIndex] = index

	return false
}

func (a *sparseSet[Index]) Remove(index Index) bool {
	if int(index) >= len(a.valuesIndices) {
		return false
	}

	valueIndex := a.valuesIndices[index]
	if valueIndex == a.EmptyValue {
		return false
	}

	a.valuesIndices[index] = a.EmptyValue

	if len(a.indices)-1 != int(valueIndex) {
		movedIndex := a.indices[len(a.indices)-1]
		a.indices[valueIndex] = movedIndex

		a.valuesIndices[movedIndex] = valueIndex
	}

	a.indices = a.indices[:len(a.indices)-1]
	return true
}
