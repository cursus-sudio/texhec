package genericrenderersys

import (
	_ "embed"
	"frontend/engine/components/camera"
	"frontend/engine/components/groups"
	meshcomponent "frontend/engine/components/mesh"
	texturecomponent "frontend/engine/components/texture"
	"frontend/engine/components/transform"
	"frontend/engine/systems/render"
	"frontend/engine/tools/cameras"
	"frontend/services/assets"
	"frontend/services/graphics/program"
	"frontend/services/graphics/shader"
	"frontend/services/graphics/texture"
	"frontend/services/graphics/vao"
	"frontend/services/graphics/vao/ebo"
	"frontend/services/graphics/vao/vbo"
	"frontend/services/media/window"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
)

//go:embed s.vert
var vertSource string

//go:embed s.frag
var fragSource string

//

type Vertex struct {
	Pos [3]float32
	// normal [3]float32
	TexturePos [2]float32
	// color [4]float32
	// vertexGroups (for animation) []VertexGroupWeight {Name string; weight float32} (weights should add up to one)
}

type PipelineComponent struct{}

type locations struct {
	Mvp int32 `uniform:"mvp"`
}

//

type releasable struct {
	textures  map[assets.AssetID]texture.Texture
	meshes    map[assets.AssetID]vao.VAO
	program   program.Program
	locations locations
}

func (r releasable) Release() {
	r.program.Release()
	for _, texture := range r.textures {
		texture.Release()
	}
	for _, mesh := range r.meshes {
		mesh.Release()
	}
}

//

type system struct {
	world          ecs.World
	transformArray ecs.ComponentsArray[transform.Transform]
	groupsArray    ecs.ComponentsArray[groups.Groups]
	textureArray   ecs.ComponentsArray[texturecomponent.Texture]
	meshArray      ecs.ComponentsArray[meshcomponent.Mesh]

	window        window.Api
	assetsStorage assets.AssetsStorage
	logger        logger.Logger
	vboFactory    vbo.VBOFactory[Vertex]
	camerasCtors  cameras.CameraConstructors

	query       ecs.LiveQuery
	cameraQuery ecs.LiveQuery

	releasable
}

func NewSystem(
	world ecs.World,
	window window.Api,
	assetsStorage assets.AssetsStorage,
	logger logger.Logger,
	vboFactory vbo.VBOFactory[Vertex],
	camerasCtors cameras.CameraConstructors,
	entitiesQueryAdditionalArguments []ecs.ComponentType,
) (ecs.SystemRegister, error) {
	vert, err := shader.NewShader(vertSource, shader.VertexShader)
	if err != nil {
		return nil, err
	}
	defer vert.Release()

	frag, err := shader.NewShader(fragSource, shader.FragmentShader)
	if err != nil {
		return nil, err
	}
	defer frag.Release()

	programID := gl.CreateProgram()
	gl.AttachShader(programID, vert.ID())
	gl.AttachShader(programID, frag.ID())

	p, err := program.NewProgram(programID, nil)
	if err != nil {
		return nil, err
	}

	locations, err := program.GetProgramLocations[locations](p)
	if err != nil {
		return nil, err
	}

	releasable := releasable{
		textures:  make(map[assets.AssetID]texture.Texture),
		meshes:    make(map[assets.AssetID]vao.VAO),
		program:   p,
		locations: locations,
	}

	world.SaveRegister(releasable)

	system := &system{
		world:          world,
		transformArray: ecs.GetComponentsArray[transform.Transform](world.Components()),
		groupsArray:    ecs.GetComponentsArray[groups.Groups](world.Components()),
		textureArray:   ecs.GetComponentsArray[texturecomponent.Texture](world.Components()),
		meshArray:      ecs.GetComponentsArray[meshcomponent.Mesh](world.Components()),

		window:        window,
		assetsStorage: assetsStorage,
		logger:        logger,
		vboFactory:    vboFactory,
		camerasCtors:  camerasCtors,

		query: world.QueryEntitiesWithComponents(
			append(
				entitiesQueryAdditionalArguments,
				ecs.GetComponentType(PipelineComponent{}),
				ecs.GetComponentType(transform.Transform{}),
				ecs.GetComponentType(meshcomponent.Mesh{}),
				ecs.GetComponentType(texturecomponent.Texture{}),
			)...,
		),
		cameraQuery: world.QueryEntitiesWithComponents(
			ecs.GetComponentType(camera.Camera{}),
		),

		releasable: releasable,
	}

	return system, nil
}

func (s *system) Register(b events.Builder) {
	events.ListenE(b, s.Listen)
}

//

func (m *system) getTexture(asset assets.AssetID) (texture.Texture, error) {
	if texture, ok := m.textures[asset]; ok {
		return texture, nil
	}
	textureAsset, err := assets.StorageGet[texturecomponent.TextureAsset](m.assetsStorage, asset)
	if err != nil {
		return nil, err
	}
	texture, err := texture.NewTexture(textureAsset.Image())
	if err != nil {
		return nil, err
	}
	m.textures[asset] = texture
	return texture, nil
}

func (m *system) getMesh(asset assets.AssetID) (vao.VAO, error) {
	if mesh, ok := m.meshes[asset]; ok {
		return mesh, nil
	}
	meshAsset, err := assets.StorageGet[meshcomponent.MeshAsset[Vertex]](m.assetsStorage, asset)
	if err != nil {
		return nil, err
	}

	VBO := m.vboFactory()
	VBO.SetVertices(meshAsset.Vertices())
	EBO := ebo.NewEBO()
	EBO.SetIndices(meshAsset.Indices())
	VAO := vao.NewVAO(VBO, EBO)
	m.meshes[asset] = VAO
	return VAO, nil
}

func (m *system) Listen(rendersys.RenderEvent) error {
	m.program.Use()

	for _, cameraEntity := range m.cameraQuery.Entities() {
		cameraGroups, err := m.groupsArray.GetComponent(cameraEntity)
		if err != nil {
			cameraGroups = groups.DefaultGroups()
		}

		camera, err := m.camerasCtors.Get(cameraEntity)
		if err != nil {
			continue
		}

		for _, entity := range m.query.Entities() {
			entityGroups, err := m.groupsArray.GetComponent(entity)
			if err != nil {
				entityGroups = groups.DefaultGroups()
			}
			if !entityGroups.SharesAnyGroup(cameraGroups) {
				continue
			}

			transformComponent, err := m.transformArray.GetComponent(entity)
			if err != nil {
				transformComponent = transform.NewTransform()
			}
			model := transformComponent.Mat4()

			textureComponent, err := m.textureArray.GetComponent(entity)
			if err != nil {
				continue
			}
			textureAsset, err := m.getTexture(textureComponent.ID)
			if err != nil {
				continue
			}

			meshComponent, err := m.meshArray.GetComponent(entity)
			if err != nil {
				continue
			}
			meshAsset, err := m.getMesh(meshComponent.ID)
			if err != nil {
				continue
			}

			textureAsset.Use()
			meshAsset.Use()

			mvp := camera.Mat4().Mul4(model)

			gl.UniformMatrix4fv(m.locations.Mvp, 1, false, &mvp[0])
			gl.DrawElementsWithOffset(gl.TRIANGLES, int32(meshAsset.EBO().Len()), gl.UNSIGNED_INT, 0)
		}
	}

	return nil
}
