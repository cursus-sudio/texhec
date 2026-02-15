package internal

import (
	"engine/modules/assets"
	"engine/services/datastructures"
	"fmt"

	"github.com/ogiusek/ioc/v2"
)

type index uint32
type dispatcher func(assets.Path) (any, error)
type asset struct {
	path       assets.Path
	dispatcher dispatcher
	cached     any
}

//

type service struct {
	Extensions assets.Extensions `inject:"1"`

	parentDirectory assets.Path

	assets []*asset
	ids    map[assets.Path]assets.ID
	//
	cached datastructures.SparseSet[index]
}

func NewService(c ioc.Dic, parentDirectory string) assets.Service {
	if len(parentDirectory) != 0 && parentDirectory[len(parentDirectory)-1] != '/' {
		parentDirectory += "/"
	}
	s := ioc.GetServices[*service](c)
	s.parentDirectory = assets.Path(parentDirectory)

	s.assets = make([]*asset, 0)
	s.ids = make(map[assets.Path]assets.ID)

	s.cached = datastructures.NewSparseSet[index]()
	return s
}

func (s *service) PathID(path assets.Path) (assets.ID, bool) {
	if id, ok := s.ids[path]; ok {
		return id, true
	}
	extension := s.Extensions.PathExntesion(path)
	dispatcher, ok := s.Extensions.ExtensionDispatcher(extension)
	if !ok {
		return 0, false
	}

	asset := &asset{
		s.parentDirectory + path,
		dispatcher,
		nil,
	}
	s.assets = append(s.assets, asset)

	index := index(len(s.assets) - 1)
	return s.indexId(index), true
}

func (s *service) RegisterExtension(
	/* shouldn't have dots and be after dots in asset */ extension string,
	dispatcher func(path assets.Path) (any, error),
) {
}

func (s *service) idIndex(id assets.ID) (index, bool) {
	if id == 0 {
		return 0, false
	}
	return index(id - 1), true
}
func (s *service) indexId(index index) assets.ID {
	return assets.ID(index + 1)
}

func (s *service) Get(id assets.ID) (any, error) {
	i, ok := s.idIndex(id)
	if !ok {
		return nil, fmt.Errorf("asset isn't registered")
	}
	asset := s.assets[i]
	if s.cached.Get(i) {
		return asset.cached, nil
	}
	obj, err := asset.dispatcher(asset.path)
	if err != nil {
		return nil, err
	}
	asset.cached = obj
	s.cached.Add(i)
	return obj, nil
}
func (s *service) release(index index) {
	if ok := s.cached.Get(index); !ok {
		return
	}
	asset := s.assets[index]
	if releasable, ok := asset.cached.(assets.Asset); ok {
		releasable.Release()
	}
	asset.cached = nil
	s.cached.Remove(index)
}
func (s *service) Release(ids ...assets.ID) {
	for _, id := range ids {
		index, ok := s.idIndex(id)
		if !ok {
			continue
		}
		s.release(index)
	}
}
func (s *service) ReleaseAll() {
	source := s.cached.GetIndices()
	dst := make([]index, len(source))
	copy(dst, source)
	for _, index := range dst {
		s.release(index)
	}
}
