package gameobject

import (
	"encoding/json"
	"reflect"
	"sync"

	"github.com/ogiusek/ioc"
)

// SERVICE
type MetadataErrors interface {
	AlreadyRegistered(service any) error
	AlreadySealed() error
	CannotUpdateUntilSealed(edited any) error
	NotFoundService(service any) error
}

type Metadata struct {
	// mutex here is used to ensure app runs synchronously
	mutex    *sync.Mutex
	metadata map[string][]byte
	sealed   bool
}

func NewMetadata() Metadata {
	return Metadata{
		mutex:    &sync.Mutex{},
		metadata: map[string][]byte{},
		sealed:   false,
	}
}

// json

func (m *Metadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.metadata)
}

func (m *Metadata) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &m.metadata); err != nil {
		return err
	}
	m.sealed = false
	m.mutex = &sync.Mutex{}
	return nil
}

// methods

func metadataKey[T any]() string {
	var t T
	return reflect.TypeOf(t).String()
}

func MetadataRegister[T any](c ioc.Dic, metadata *Metadata, service T) error {
	metadata.mutex.Lock()
	defer metadata.mutex.Unlock()

	errors := ioc.Get[MetadataErrors](c)
	if metadata.sealed {
		return errors.AlreadySealed()
	}

	key := metadataKey[T]()
	if _, ok := metadata.metadata[key]; ok {
		return errors.AlreadyRegistered(service)
	}

	bytes, err := json.Marshal(service)
	if err != nil {
		return err
	}

	metadata.metadata[key] = bytes
	return nil
}

func MetadateGet[T any](c ioc.Dic, metadata *Metadata) (T, error) {
	var service T
	errors := ioc.Get[MetadataErrors](c)
	key := metadataKey[T]()

	bytes, ok := metadata.metadata[key]
	if !ok {
		return service, errors.NotFoundService(service)
	}
	if err := json.Unmarshal(bytes, &service); err != nil {
		return service, err
	}
	return service, nil
}

func MetadataSeal(c ioc.Dic, metadata *Metadata) error {
	metadata.mutex.Lock()
	defer metadata.mutex.Unlock()
	if metadata.sealed {
		errors := ioc.Get[MetadataErrors](c)
		return errors.AlreadySealed()
	}

	metadata.sealed = true
	return nil
}

func MetadataUpdate[T any](c ioc.Dic, metadata *Metadata, service T) error {
	metadata.mutex.Lock()
	defer metadata.mutex.Unlock()
	errors := ioc.Get[MetadataErrors](c)
	if !metadata.sealed {
		return errors.CannotUpdateUntilSealed(service)
	}
	key := metadataKey[T]()
	if _, ok := metadata.metadata[key]; !ok {
		return errors.NotFoundService(service)
	}

	bytes, err := json.Marshal(service)
	if err != nil {
		return err
	}

	metadata.metadata[key] = bytes
	return nil
}
