package renderer

import (
	_ "embed"
	"engine/modules/camera"
	"engine/modules/genericrenderer"
	"engine/modules/groups"
	"engine/modules/render"
	"engine/modules/transform"
	"engine/services/assets"
	"engine/services/ecs"
	"engine/services/graphics/program"
	"engine/services/graphics/shader"
	"engine/services/graphics/texture"
	"engine/services/graphics/vao"
	"engine/services/graphics/vao/ebo"
	"engine/services/graphics/vao/vbo"
	"engine/services/logger"
	"engine/services/media/window"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
)

//go:embed s.vert
var vertSource string

//go:embed s.frag
var fragSource string

//

type locations struct {
	Mvp   int32 `uniform:"mvp"`
	Color int32 `uniform:"u_color"`
}

//

type textureKey struct {
	Texture render.TextureComponent
	Frame   int
}

type releasable struct {
	texturesImagesCount map[render.TextureComponent]int
	textures            map[textureKey]texture.Texture
	meshes              map[assets.AssetID]vao.VAO
	program             program.Program
	locations           locations
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
	world                ecs.World
	transformTransaction transform.Transaction
	groupsArray          ecs.ComponentsArray[groups.GroupsComponent]
	textureArray         ecs.ComponentsArray[render.TextureComponent]
	textureFrameArray    ecs.ComponentsArray[render.TextureFrameComponent]
	colorArray           ecs.ComponentsArray[render.ColorComponent]
	meshArray            ecs.ComponentsArray[render.MeshComponent]

	cameraArray ecs.ComponentsArray[camera.CameraComponent]

	window         window.Api
	assetsStorage  assets.AssetsStorage
	logger         logger.Logger
	vboFactory     vbo.VBOFactory[genericrenderer.Vertex]
	textureFactory texture.Factory
	camerasCtors   camera.Tool

	query ecs.LiveQuery

	releasable
}

func NewSystem(
	window window.Api,
	assetsStorage assets.AssetsStorage,
	logger logger.Logger,
	vboFactory vbo.VBOFactory[genericrenderer.Vertex],
	textureFactory texture.Factory,
	camerasCtors ecs.ToolFactory[camera.Tool],
	transformToolFactory ecs.ToolFactory[transform.Tool],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		vert, err := shader.NewShader(vertSource, shader.VertexShader)
		if err != nil {
			return err
		}
		defer vert.Release()

		frag, err := shader.NewShader(fragSource, shader.FragmentShader)
		if err != nil {
			return err
		}
		defer frag.Release()

		programID := gl.CreateProgram()
		gl.AttachShader(programID, vert.ID())
		gl.AttachShader(programID, frag.ID())

		p, err := program.NewProgram(programID, nil)
		if err != nil {
			return err
		}

		locations, err := program.GetProgramLocations[locations](p)
		if err != nil {
			return err
		}

		releasable := releasable{
			texturesImagesCount: make(map[render.TextureComponent]int),
			textures:            make(map[textureKey]texture.Texture),
			meshes:              make(map[assets.AssetID]vao.VAO),
			program:             p,
			locations:           locations,
		}

		w.SaveGlobal(releasable)

		transformTool := transformToolFactory.Build(w)
		system := &system{
			world:                w,
			transformTransaction: transformTool.Transaction(),
			groupsArray:          ecs.GetComponentsArray[groups.GroupsComponent](w),
			textureArray:         ecs.GetComponentsArray[render.TextureComponent](w),
			textureFrameArray:    ecs.GetComponentsArray[render.TextureFrameComponent](w),
			colorArray:           ecs.GetComponentsArray[render.ColorComponent](w),
			meshArray:            ecs.GetComponentsArray[render.MeshComponent](w),

			cameraArray: ecs.GetComponentsArray[camera.CameraComponent](w),

			window:         window,
			assetsStorage:  assetsStorage,
			logger:         logger,
			vboFactory:     vboFactory,
			textureFactory: textureFactory,
			camerasCtors:   camerasCtors.Build(w),

			query: transformTool.Query(w.Query()).
				Require(
					genericrenderer.PipelineComponent{},
					render.MeshComponent{},
					render.TextureComponent{},
				).
				Track(
					render.ColorComponent{},
					render.TextureFrameComponent{},
				).
				Build(),

			releasable: releasable,
		}

		events.ListenE(w.EventsBuilder(), system.Listen)
		return nil
	})

}

//

func (m *system) getTexture(entity ecs.EntityID) (texture.Texture, error) {
	textureComponent, err := m.textureArray.GetComponent(entity)
	if err != nil {
		return nil, err
	}
	imagesCount, okImagesCount := m.texturesImagesCount[textureComponent]
	textureFrameComponent, err := m.textureFrameArray.GetComponent(entity)
	if err != nil {
		textureFrameComponent = render.DefaultTextureFrameComponent()
	}
	var frame int
	if okImagesCount {
		frame = textureFrameComponent.GetFrame(imagesCount)
	}
	textureKey := textureKey{textureComponent, frame}
	if texture, ok := m.textures[textureKey]; ok {
		return texture, nil
	}

	textureAsset, err := assets.StorageGet[render.TextureAsset](m.assetsStorage, textureComponent.Asset)
	if err != nil {
		return nil, err
	}

	imagesCount = len(textureAsset.Images())
	if !okImagesCount {
		frame = textureFrameComponent.GetFrame(imagesCount)
	}

	image := textureAsset.Images()[frame]
	texture, err := m.textureFactory.New(image)
	if err != nil {
		return nil, err
	}
	m.textures[textureKey] = texture
	m.texturesImagesCount[textureComponent] = imagesCount
	return texture, nil
}

func (m *system) getMesh(asset assets.AssetID) (vao.VAO, error) {
	if mesh, ok := m.meshes[asset]; ok {
		return mesh, nil
	}
	meshAsset, err := assets.StorageGet[render.MeshAsset[genericrenderer.Vertex]](m.assetsStorage, asset)
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

func (m *system) Listen(render.RenderEvent) error {
	m.program.Use()

	for _, cameraEntity := range m.cameraArray.GetEntities() {
		cameraGroups, err := m.groupsArray.GetComponent(cameraEntity)
		if err != nil {
			cameraGroups = groups.DefaultGroups()
		}

		camera, err := m.camerasCtors.GetObject(cameraEntity)
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

			transform := m.transformTransaction.GetObject(entity)
			model := transform.Mat4()

			textureAsset, err := m.getTexture(entity)
			if err != nil {
				continue
			}

			colorComponent, err := m.colorArray.GetComponent(entity)
			if err != nil {
				colorComponent = render.DefaultColor()
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
			gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, meshAsset.EBO().ID())

			mvp := camera.Mat4().Mul4(model)
			gl.Viewport(camera.Viewport())
			gl.UniformMatrix4fv(m.locations.Mvp, 1, false, &mvp[0])
			gl.Uniform4fv(m.locations.Color, 1, &colorComponent.Color[0])

			gl.DrawElementsWithOffset(gl.TRIANGLES, int32(meshAsset.EBO().Len()), gl.UNSIGNED_INT, 0)
		}
	}

	return nil
}
