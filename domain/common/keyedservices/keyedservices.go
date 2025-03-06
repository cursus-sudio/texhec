package keyedservices

import (
	"github.com/ogiusek/ioc"
	"github.com/ogiusek/null"
)

type KeyedServicesErrors interface {
	ServiceAlreadyExists(key any) error
}

type KeyedServices[TKey comparable, TService any] interface {
	// panics when there is no
	Get(key TKey) null.Nullable[TService]
	Add(key TKey, service TService) error
}

type keyedServicesImpl[TKey comparable, TService any] struct {
	c        ioc.Dic
	services map[TKey]TService
}

func (services *keyedServicesImpl[TKey, TService]) Get(key TKey) null.Nullable[TService] {
	service, ok := services.services[key]
	if !ok {
		return null.Null[TService]()
	}
	return null.New(service)
}

func (services *keyedServicesImpl[TKey, TService]) Add(key TKey, service TService) error {
	_, ok := services.services[key]
	if !ok {
		services.services[key] = service
		return nil
	}
	return ioc.Get[KeyedServicesErrors](services.c).ServiceAlreadyExists(key)
}

func NewKeyedService[TKey comparable, TService any](c ioc.Dic) KeyedServices[TKey, TService] {
	return &keyedServicesImpl[TKey, TService]{
		c:        c,
		services: make(map[TKey]TService),
	}
}
