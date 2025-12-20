package ecs

type DirtySet interface {
	// get also clears
	Get() []EntityID
	Dirty(EntityID)
	Clear()
}

type dirtySet struct {
	entities []EntityID
	set      []uint8
}

func NewDirtySet() DirtySet {
	return &dirtySet{
		entities: nil,
		set:      nil,
	}
}

func (f *dirtySet) Get() []EntityID {
	values := f.entities
	f.Clear()
	return values
}

func (f *dirtySet) Dirty(entity EntityID) {
	byteIndex := entity / 8
	entityMask := uint8(1 << (entity % 8))

	// ensure index exists
	for len(f.set) <= int(byteIndex) {
		f.set = append(f.set, 0)
	}

	if f.set[byteIndex]&entityMask != 0 {
		return
	}
	f.entities = append(f.entities, entity)
	f.set[byteIndex] |= entityMask
}

func (f *dirtySet) Clear() {
	for _, entity := range f.entities {
		byteIndex := entity / 8
		f.set[byteIndex] = 0
	}
	f.entities = nil
}
