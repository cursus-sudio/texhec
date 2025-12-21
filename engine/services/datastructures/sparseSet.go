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
	ValuesIndices []Index // here some indices have special meaning (read constants above)

	// both arrays below Indices correspond
	Indices []Index // here value means index in sparse array
}

func NewSparseSet[Index constraints.Integer]() SparseSet[Index] {
	var zero Index
	return &sparseSet[Index]{
		EmptyValue: ^zero,
	}
}

func (a *sparseSet[Index]) Get(index Index) bool {
	if int(index) >= len(a.ValuesIndices) {
		return false
	}

	valueIndex := a.ValuesIndices[index]
	if valueIndex == a.EmptyValue {
		return false
	}

	return true
}

func (a *sparseSet[Index]) GetIndices() []Index { return a.Indices }

func (a *sparseSet[Index]) Add(index Index) bool {
	for int(index) >= len(a.ValuesIndices) {
		a.ValuesIndices = append(a.ValuesIndices, a.EmptyValue)
	}

	valueIndex := a.ValuesIndices[index]

	if valueIndex == a.EmptyValue {
		a.ValuesIndices[index] = Index(len(a.Indices))
		a.Indices = append(a.Indices, index)
		return true
	}

	a.Indices[valueIndex] = index

	return false
}

func (a *sparseSet[Index]) Remove(index Index) bool {
	if int(index) >= len(a.ValuesIndices) {
		return false
	}

	valueIndex := a.ValuesIndices[index]
	if valueIndex == a.EmptyValue {
		return false
	}

	a.ValuesIndices[index] = a.EmptyValue

	if len(a.Indices)-1 != int(valueIndex) {
		movedIndex := a.Indices[len(a.Indices)-1]
		a.Indices[valueIndex] = movedIndex

		a.ValuesIndices[movedIndex] = valueIndex
	}

	a.Indices = a.Indices[:len(a.Indices)-1]
	return true
}
