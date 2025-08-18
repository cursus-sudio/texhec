package material

import (
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/graphics/program"
)

type Material struct {
	IDs []assets.AssetID
}

func NewMaterial(ids ...assets.AssetID) Material {
	return Material{IDs: ids}
}

type MaterialStorageAsset interface {
	assets.CachableAsset
	Program() (program.Program, error)
	Render(world ecs.World, program program.Program) error
}

type materialStorageAsset struct {
	program func() (program.Program, error)
	render  func(world ecs.World, program program.Program) error
}

func NewMaterialStorageAsset(
	program func() (program.Program, error),
	render func(ecs.World, program.Program) error,
) MaterialStorageAsset {
	return &materialStorageAsset{
		program: program,
		render:  render,
	}
}

func (a *materialStorageAsset) Program() (program.Program, error) {
	return a.program()
}
func (a *materialStorageAsset) Render(world ecs.World, program program.Program) error {
	return a.render(world, program)
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
	Render(world ecs.World) error
}

type materialCachedAsset struct {
	program program.Program
	render  func(world ecs.World, program program.Program) error
}

func (a *materialCachedAsset) Render(world ecs.World) error {
	a.program.Use()
	return a.render(world, a.program)
}

func (a *materialCachedAsset) Release() {
	a.program.Release()
}
