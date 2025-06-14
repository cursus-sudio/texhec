package safemap

import "errors"

var (
	ErrDoNotExists error = errors.New("do not exists")
)

type SafeMap[Key any, Value any] interface {
	// can return error:
	// - ErrDoNotExists
	Get(Key) (Value, error)
	GetOrCreate(key Key, defaultValue Value) Value

	Set(Key, Value)
	Remove(Key)
}
