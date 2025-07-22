package program

import (
	"frontend/services/assets"
	"frontend/services/graphics/program"
	"frontend/services/graphics/shader"
)

type Program struct {
	ID assets.AssetID
}

func NewProgram(id assets.AssetID) Program {
	return Program{ID: id}
}

type ProgramStorageAsset[Locations any] interface {
	assets.StorageAsset
	VertexShader() string
	FragmentShader() string
	// other shaders can be added
	// tcs
	// tes
	// gs
	// cs (this shader is different and can be used to add ray tracing)
	//
	// and these shaders can return *string or something
}

type programStorageAsset[Locations any] struct {
	vertexShader   string
	fragmentShader string
}

func NewProgramStorageAsset[Locations any](
	vertexShader string,
	fragmentShader string,
) ProgramStorageAsset[Locations] {
	return &programStorageAsset[Locations]{
		vertexShader:   vertexShader,
		fragmentShader: fragmentShader,
	}
}

func (a *programStorageAsset[Locations]) VertexShader() string   { return a.vertexShader }
func (a *programStorageAsset[Locations]) FragmentShader() string { return a.fragmentShader }

func (a *programStorageAsset[Locations]) Cache() (assets.CachedAsset, error) {
	vert, err := shader.NewShader(a.vertexShader, shader.VertexShader)
	if err != nil {
		return nil, err
	}
	frag, err := shader.NewShader(a.fragmentShader, shader.FragmentShader)
	if err != nil {
		return nil, err
	}
	p, err := program.NewProgram[Locations](vert, frag)
	if err != nil {
		return nil, err
	}
	vert.Release()
	frag.Release()
	var asset ProgramCachedAsset[Locations] = &programCachedAsset[Locations]{program: p}
	return asset, nil
}

type ProgramCachedAsset[Locations any] interface {
	assets.CachedAsset
	Program() program.Program[Locations]
}

type programCachedAsset[Locations any] struct {
	program program.Program[Locations]
}

func (a *programCachedAsset[Locations]) Program() program.Program[Locations] { return a.program }
func (a *programCachedAsset[Locations]) Release()                            { a.program.Release() }
