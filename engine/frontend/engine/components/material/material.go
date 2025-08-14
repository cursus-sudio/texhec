package material

import (
	"errors"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/graphics/program"
)

var (
	ErrHaveToCallOnFrame error = errors.New("has to call on frame")
)

type Material struct {
	ID assets.AssetID
}

func NewMaterial(id assets.AssetID) Material {
	return Material{ID: id}
}

type MaterialStorageAsset interface {
	assets.StorageAsset
	Program() (program.Program, error)
	Render(world ecs.World, program program.Program, entities []ecs.EntityID) error
}

type materialStorageAsset struct {
	program func() (program.Program, error)
	render  func(world ecs.World, program program.Program, entities []ecs.EntityID) error
}

func NewMaterialStorageAsset(
	program func() (program.Program, error),
	render func(ecs.World, program.Program, []ecs.EntityID) error,
) MaterialStorageAsset {
	return &materialStorageAsset{
		program: program,
		render:  render,
	}
}

func (a *materialStorageAsset) Program() (program.Program, error) {
	return a.program()
}
func (a *materialStorageAsset) Render(world ecs.World, program program.Program, entities []ecs.EntityID) error {
	return a.render(world, program, entities)
}

func (a *materialStorageAsset) Cache() (assets.CachedAsset, error) {
	p, err := a.Program()
	if err != nil {
		return nil, err
	}
	var cached MaterialCachedAsset = &materialCachedAsset{
		program: p,
		render:  a.Render,
	}
	return cached, nil
}

type MaterialCachedAsset interface {
	assets.CachedAsset
	Render(world ecs.World, entities []ecs.EntityID) error
}

type materialCachedAsset struct {
	program program.Program
	render  func(world ecs.World, program program.Program, entities []ecs.EntityID) error
}

func (a *materialCachedAsset) Render(world ecs.World, entities []ecs.EntityID) error {
	a.program.Use()
	return a.render(world, a.program, entities)
}

func (a *materialCachedAsset) Release() {
	a.program.Release()
}
