package ecs

type DirtySet interface {
	// get also clears
	Get() []EntityID
	Dirty(EntityID)
	Clear()

	Ok() bool
	Release()
}

type dirtySet struct {
	entities []EntityID
	set      []uint8
	ok       bool
}

func NewDirtySet() DirtySet {
	return &dirtySet{
		entities: nil,
		set:      nil,
		ok:       true,
	}
}

func (f *dirtySet) Get() []EntityID {
	if !f.ok {
		return nil
	}
	values := f.entities
	f.Clear()
	return values
}

func (f *dirtySet) Dirty(entity EntityID) {
	if !f.ok {
		return
	}
	byteIndex := int(entity / 8)
	entityMask := uint8(1 << (entity & 7)) // &7 == %8

	if byteIndex >= len(f.set) {
		newBytes := make([]uint8, byteIndex-len(f.set)+1)
		f.set = append(f.set, newBytes...)
	}

	if f.set[byteIndex]&entityMask != 0 {
		return
	}
	f.entities = append(f.entities, entity)
	f.set[byteIndex] |= entityMask
}

func (f *dirtySet) Clear() {
	if !f.ok {
		return
	}
	for _, entity := range f.entities {
		byteIndex := entity / 8
		f.set[byteIndex] = 0
	}
	f.entities = nil
}

func (f *dirtySet) Ok() bool {
	return f.ok
}
func (f *dirtySet) Release() {
	f.ok = false
	f.entities = nil
	f.set = nil
}
