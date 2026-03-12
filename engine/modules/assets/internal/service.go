package internal

import (
	"engine/modules/assets"
	"engine/modules/registry"
	"engine/services/ecs"
	"fmt"

	"github.com/ogiusek/ioc/v2"
)

//

type assetsService struct {
	Registry registry.Service `inject:"1"`
	World    ecs.World        `inject:"1"`

	*extensions
	path  ecs.ComponentsArray[assets.PathComponent]
	cache ecs.ComponentsArray[assets.CacheComponent]
}

func NewService(c ioc.Dic) assets.Service {
	s := ioc.GetServices[*assetsService](c)

	s.extensions = NewExtensions(c)
	s.path = ecs.GetComponentsArray[assets.PathComponent](s.World)
	s.cache = ecs.GetComponentsArray[assets.CacheComponent](s.World)

	return s
}

func (s *assetsService) Path() ecs.ComponentsArray[assets.PathComponent]   { return s.path }
func (s *assetsService) Cache() ecs.ComponentsArray[assets.CacheComponent] { return s.cache }

func (s *assetsService) Get(entity ecs.EntityID) (assets.Asset, error) {
	if cache, ok := s.cache.Get(entity); ok {
		return cache.Cache, nil
	}

	path, ok := s.path.Get(entity)
	if !ok {
		return nil, assets.ErrAssetNotFound
	}

	ext := s.PathExntesion(path)
	dispatcher, ok := s.ExtensionDispatcher(ext)
	if !ok {
		fmt.Printf("\"%v\" path.\n", path)
		return nil, assets.ErrAssetNotFound
	}
	asset, err := dispatcher(path)
	if err != nil {
		return nil, err
	}
	s.cache.Set(entity, assets.NewCache(asset))
	return asset, nil
}

func (s *assetsService) Release(entities ...ecs.EntityID) {
	for _, entity := range entities {
		if cache, ok := s.cache.Get(entity); ok {
			cache.Cache.Release()
			s.cache.Remove(entity)
		}
	}
}
func (s *assetsService) ReleaseAll() {
	src := s.cache.GetEntities()
	dst := make([]ecs.EntityID, len(src))
	copy(dst, src)
	for _, entity := range dst {
		if cache, ok := s.cache.Get(entity); ok {
			cache.Cache.Release()
			s.cache.Remove(entity)
		}
	}
}
