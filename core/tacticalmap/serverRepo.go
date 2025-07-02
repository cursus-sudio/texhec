package tacticalmap

import (
	"shared/utils/httperrors"
	"sync"
)

type Pos struct{ X, Y int }
type Tile struct{ Pos Pos }

type CreateArgs struct{ Tiles []Tile }
type DestroyArgs struct{ Tiles []Tile }

type CreateListener func(tiles []Tile)
type DestroyListener func(tiles []Tile)

type TacticalMap interface {
	Create(CreateArgs) error
	Destroy(DestroyArgs) error
	GetMap() ([]Tile, error)
	OnCreate(CreateListener)
	OnDestroy(DestroyListener)
}

// impl

type tacticalMap struct {
	tiles map[Pos]Tile
	mutex sync.RWMutex

	createListenersMutex sync.RWMutex
	createListeners      []CreateListener

	destroyListenersMutex sync.RWMutex
	destroyListeners      []DestroyListener
}

func newTacticalMap() TacticalMap {
	return &tacticalMap{
		tiles:                 map[Pos]Tile{},
		mutex:                 sync.RWMutex{},
		createListenersMutex:  sync.RWMutex{},
		createListeners:       make([]CreateListener, 0),
		destroyListenersMutex: sync.RWMutex{},
		destroyListeners:      make([]DestroyListener, 0),
	}
}

func (tacticalMap *tacticalMap) Create(args CreateArgs) error {
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

func (tacticalMap *tacticalMap) Destroy(args DestroyArgs) error {
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

func (tacticalMap *tacticalMap) GetMap() ([]Tile, error) {
	tacticalMap.mutex.RLock()
	defer tacticalMap.mutex.RUnlock()
	tiles := make([]Tile, 0, len(tacticalMap.tiles))
	for _, tile := range tacticalMap.tiles {
		tiles = append(tiles, tile)
	}
	if len(tiles) == 0 {
		return nil, httperrors.Err404
	}
	return tiles, nil
}

func (tacticalMap *tacticalMap) OnCreate(listener CreateListener) {
	tacticalMap.createListenersMutex.Lock()
	defer tacticalMap.createListenersMutex.Unlock()

	tacticalMap.createListeners = append(tacticalMap.createListeners, listener)
}

func (tacticalMap *tacticalMap) OnDestroy(listener DestroyListener) {
	tacticalMap.destroyListenersMutex.Lock()
	defer tacticalMap.destroyListenersMutex.Unlock()

	tacticalMap.destroyListeners = append(tacticalMap.destroyListeners, listener)
}

func TacticalMapImpl() TacticalMap {
	return &tacticalMap{
		tiles: map[Pos]Tile{},
		mutex: sync.RWMutex{},
	}
}
