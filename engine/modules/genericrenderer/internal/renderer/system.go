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
	"github.com/ogiusek/ioc/v2"
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

//

type system struct {
	EventsBuilder   events.Builder          `inject:"1"`
	World           ecs.World               `inject:"1"`
	GenericRenderer genericrenderer.Service `inject:"1"`
	Render          render.Service          `inject:"1"`
	Camera          camera.Service          `inject:"1"`
	Groups          groups.Service          `inject:"1"`
	Transform       transform.Service       `inject:"1"`

	Window         window.Api                             `inject:"1"`
	AssetsStorage  assets.AssetsStorage                   `inject:"1"`
	Logger         logger.Logger                          `inject:"1"`
	VboFactory     vbo.VBOFactory[genericrenderer.Vertex] `inject:"1"`
	TextureFactory gtexture.Factory                       `inject:"1"`

	texturesImagesCount map[render.TextureComponent]int
	textures            map[textureKey]gtexture.Texture
	meshes              map[assets.AssetID]vao.VAO
	program             program.Program
	locations           locations
}

func NewSystem(c ioc.Dic) genericrenderer.System {
	return ecs.NewSystemRegister(func() error {
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

		s := ioc.GetServices[*system](c)

		s.texturesImagesCount = make(map[render.TextureComponent]int)
		s.textures = make(map[textureKey]gtexture.Texture)
		s.meshes = make(map[assets.AssetID]vao.VAO)
		s.program = p
		s.locations = locations

		events.ListenE(s.EventsBuilder, s.Listen)
		return nil
	})
}

//

func (m *system) getTexture(entity ecs.EntityID) (gtexture.Texture, error) {
	textureComponent, ok := m.Render.Texture().Get(entity)
	if !ok {
		return nil, nil
	}
	imagesCount, okImagesCount := m.texturesImagesCount[textureComponent]
	textureFrameComponent, ok := m.Render.TextureFrame().Get(entity)
	if !ok {
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

	textureAsset, err := assets.StorageGet[render.TextureAsset](m.AssetsStorage, textureComponent.Asset)
	if err != nil {
		return nil, err
	}

	imagesCount = len(textureAsset.Images())
	if !okImagesCount {
		frame = textureFrameComponent.GetFrame(imagesCount)
	}

	image := textureAsset.Images()[frame]
	texture, err := m.TextureFactory.New(image)
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
	meshAsset, err := assets.StorageGet[render.MeshAsset[genericrenderer.Vertex]](m.AssetsStorage, asset)
	if err != nil {
		return nil, err
	}

	VBO := m.VboFactory()
	VBO.SetVertices(meshAsset.Vertices())
	EBO := ebo.NewEBO()
	EBO.SetIndices(meshAsset.Indices())
	VAO := vao.NewVAO(VBO, EBO)
	m.meshes[asset] = VAO
	return VAO, nil
}

func (m *system) Listen(render.RenderEvent) error {
	m.program.Use()

	for _, cameraEntity := range m.Camera.Component().GetEntities() {
		cameraGroups, ok := m.Groups.Component().Get(cameraEntity)
		if !ok {
			cameraGroups = groups.DefaultGroups()
		}

		for _, entity := range m.GenericRenderer.Pipeline().GetEntities() {
			entityGroups, ok := m.Groups.Component().Get(entity)
			if !ok {
				entityGroups = groups.DefaultGroups()
			}
			if !entityGroups.SharesAnyGroup(cameraGroups) {
				continue
			}

			model := m.Transform.Mat4(entity)

			textureAsset, err := m.getTexture(entity)
			if textureAsset == nil || err != nil {
				m.Logger.Warn(err)
				continue
			}

			colorComponent, ok := m.Render.Color().Get(entity)
			if !ok {
				colorComponent = render.DefaultColor()
			}

			meshComponent, ok := m.Render.Mesh().Get(entity)
			if !ok {
				continue
			}
			meshAsset, err := m.getMesh(meshComponent.ID)
			if err != nil {
				continue
			}

			textureAsset.Use()
			meshAsset.Use()
			gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, meshAsset.EBO().ID())

			mvp := m.Camera.Mat4(cameraEntity).Mul4(model)
			gl.Viewport(m.Camera.GetViewport(cameraEntity))
			gl.UniformMatrix4fv(m.locations.Mvp, 1, false, &mvp[0])
			gl.Uniform4fv(m.locations.Color, 1, &colorComponent.Color[0])

			gl.DrawElementsWithOffset(gl.TRIANGLES, int32(meshAsset.EBO().Len()), gl.UNSIGNED_INT, 0)
		}
	}

	return nil
}
