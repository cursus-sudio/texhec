package ecs

import (
	"engine/services/datastructures"
)

type DirtySet interface {
	// get also clears
	Get() []EntityID
	Dirty(EntityID)
	Clear()
}

type dirtySet struct {
	set datastructures.SparseSet[EntityID]
}

func NewDirtySet() DirtySet {
	return &dirtySet{
		set: datastructures.NewSparseSet[EntityID](),
	}
}

func (f *dirtySet) Get() []EntityID {
	original := f.set.GetIndices()
	values := make([]EntityID, len(original))
	copy(values, original)

	for _, entity := range original {
		f.set.Remove(entity)
	}

	return values
}

func (f *dirtySet) Dirty(entity EntityID) {
	f.set.Add(entity)
}

func (f *dirtySet) Clear() {
	for _, entity := range f.set.GetIndices() {
		f.set.Remove(entity)
	}
}
