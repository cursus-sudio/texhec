package tacticalmap

import (
	"sync"
)

type impl struct {
	tiles map[Pos]Tile
	mutex sync.RWMutex

	createListenersMutex sync.RWMutex
	createListeners      []CreateListener

	destroyListenersMutex sync.RWMutex
	destroyListeners      []DestroyListener
}

func newTacticalMap() TacticalMap {
	return &impl{
		tiles:                 map[Pos]Tile{},
		mutex:                 sync.RWMutex{},
		createListenersMutex:  sync.RWMutex{},
		createListeners:       make([]CreateListener, 0),
		destroyListenersMutex: sync.RWMutex{},
		destroyListeners:      make([]DestroyListener, 0),
	}
}

func (tacticalMap *impl) Create(args CreateArgs) error {
	{
		tacticalMap.mutex.Lock()

		for _, tile := range args.Tiles {
			tacticalMap.tiles[tile.Pos] = tile
		}

		tacticalMap.mutex.Unlock()
	}

	{
		tacticalMap.createListenersMutex.RLock()

		for _, listener := range tacticalMap.createListeners {
			listener(args.Tiles)
		}

		tacticalMap.createListenersMutex.RUnlock()
	}

	return nil
}

func (tacticalMap *impl) Destroy(args DestroyArgs) error {
	{
		tacticalMap.mutex.Lock()

		for _, tile := range args.Tiles {
			delete(tacticalMap.tiles, tile.Pos)
		}

		tacticalMap.mutex.Unlock()
	}

	{
		tacticalMap.destroyListenersMutex.RLock()

		for _, listener := range tacticalMap.destroyListeners {
			listener(args.Tiles)
		}

		tacticalMap.destroyListenersMutex.RUnlock()
	}

	return nil
}

func (tacticalMap *impl) GetMap() ([]Tile, error) {
	tacticalMap.mutex.RLock()
	defer tacticalMap.mutex.RUnlock()
	tiles := make([]Tile, 0, len(tacticalMap.tiles))
	for _, tile := range tacticalMap.tiles {
		tiles = append(tiles, tile)
	}
	return tiles, nil
}

func (tacticalMap *impl) OnCreate(listener CreateListener) {
	tacticalMap.createListenersMutex.Lock()
	defer tacticalMap.createListenersMutex.Unlock()

	tacticalMap.createListeners = append(tacticalMap.createListeners, listener)
}

func (tacticalMap *impl) OnDestroy(listener DestroyListener) {
	tacticalMap.destroyListenersMutex.Lock()
	defer tacticalMap.destroyListenersMutex.Unlock()

	tacticalMap.destroyListeners = append(tacticalMap.destroyListeners, listener)
}

func TacticalMapImpl() TacticalMap {
	return &impl{
		tiles: map[Pos]Tile{},
		mutex: sync.RWMutex{},
	}
}
