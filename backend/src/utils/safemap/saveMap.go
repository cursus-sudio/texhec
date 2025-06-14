package safemap

import "sync"

type safeMap[Key any, Value any] struct {
	syncMap sync.Map
}

func (safeMap *safeMap[Key, Value]) Get(key Key) (Value, error) {
	val, ok := safeMap.syncMap.Load(key)
	if !ok {
		var val Value
		return val, ErrDoNotExists
	}

	return val.(Value), nil
}

func (safeMap *safeMap[Key, Value]) GetOrCreate(key Key, defaultValue Value) Value {
	value, exists := safeMap.syncMap.LoadOrStore(key, defaultValue)
	if !exists {
		return defaultValue
	}
	return value.(Value)
}

func (safeMap *safeMap[Key, Value]) Set(key Key, value Value) {
	safeMap.Set(key, value)
}

func (safeMap *safeMap[Key, Value]) Remove(key Key) {
	safeMap.syncMap.Delete(key)
}
