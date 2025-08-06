package material

import (
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/graphics/program"
	"frontend/services/graphics/shader"
)

type Material struct {
	IDs []assets.AssetID
}

func NewMaterial(id ...assets.AssetID) Material {
	return Material{IDs: id}
}

type MaterialStorageAsset interface {
	assets.StorageAsset
	// shaders
	VertexShader() string
	FragmentShader() string
	// - space for more shaders (use for them *string because they are optional)

	// to set uniforms
	OnFrame() func(world ecs.World, program program.Program) error
	// - camera
	// - lights
	UseForEntity() func(world ecs.World, program program.Program, entityId ecs.EntityId) error
	// - position
	// - reflection

	// flags (parameters)
	Parameters() []program.Parameter
}

type materialStorageAsset struct {
	vertexShader   string
	fragmentShader string

	onFrame      func(world ecs.World, program program.Program) error
	useForEntity func(world ecs.World, program program.Program, entityId ecs.EntityId) error

	parameters []program.Parameter
}

func NewMaterialStorageAsset(
	vertexShader, fragmentShader string,
	onFrame func(world ecs.World, p program.Program) error,
	useForEntity func(world ecs.World, p program.Program, entityId ecs.EntityId) error,
	parameters []program.Parameter,
) MaterialStorageAsset {
	return &materialStorageAsset{
		vertexShader:   vertexShader,
		fragmentShader: fragmentShader,
		onFrame:        onFrame,
		useForEntity:   useForEntity,
		parameters:     parameters,
	}
}

func (a *materialStorageAsset) VertexShader() string   { return a.vertexShader }
func (a *materialStorageAsset) FragmentShader() string { return a.fragmentShader }
func (a *materialStorageAsset) OnFrame() func(world ecs.World, program program.Program) error {
	return a.onFrame
}
func (a *materialStorageAsset) UseForEntity() func(world ecs.World, program program.Program, entityId ecs.EntityId) error {
	return a.useForEntity
}
func (a *materialStorageAsset) Parameters() []program.Parameter { return a.parameters }

func (a *materialStorageAsset) Cache() (assets.CachedAsset, error) {
	vert, err := shader.NewShader(a.vertexShader, shader.VertexShader)
	if err != nil {
		return nil, err
	}
	frag, err := shader.NewShader(a.fragmentShader, shader.FragmentShader)
	if err != nil {
		return nil, err
	}
	p, err := program.NewProgram(vert, frag, a.parameters)
	if err != nil {
		vert.Release()
		frag.Release()
		return nil, err
	}
	vert.Release()
	frag.Release()
	var cached MaterialCachedAsset = &materialCachedAsset{
		program:    p,
		onFrame:    a.onFrame,
		drawEntity: a.useForEntity,
	}
	return cached, nil
}

type MaterialCachedAsset interface {
	assets.CachedAsset
	OnFrame(world ecs.World) error
	UseForEntity(world ecs.World, entityId ecs.EntityId) error
}

type materialCachedAsset struct {
	program    program.Program
	onFrame    func(world ecs.World, program program.Program) error
	drawEntity func(world ecs.World, program program.Program, entityId ecs.EntityId) error
}

func (a *materialCachedAsset) OnFrame(world ecs.World) error {
	a.program.Use()
	return a.onFrame(world, a.program)
}

func (a *materialCachedAsset) UseForEntity(world ecs.World, entityId ecs.EntityId) error {
	a.program.Use()
	return a.drawEntity(world, a.program, entityId)
}

func (a *materialCachedAsset) Release() {
	a.program.Release()
}
