package triangle

import (
	"bytes"
	"embed"
	"frontend/components/mesh"
	"frontend/components/program"
	"frontend/components/texture"
	"frontend/components/transform"
	"frontend/services/assets"
	"frontend/services/frames"
	"frontend/services/graphics/vao/ebo"
	"frontend/services/graphics/vao/vbo"
	"frontend/services/media/window"
	"path/filepath"
	appruntime "shared/services/runtime"
	"time"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

//go:embed square.png
var textureSource []byte

var shadersDir string = "shaders_old/texture"

//go:embed shaders_old/*
var shaders embed.FS

type FrontendPkg struct{}

func FrontendPackage() FrontendPkg {
	return FrontendPkg{}
}

const (
	MeshAssetID    assets.AssetID = "vao_asset"
	TextureAssetID assets.AssetID = "texture_asset"
	ProgramAssetID assets.AssetID = "program_asset"
	SceneAssetID   assets.AssetID = "scene_asset"
)

type Locations struct {
	Resolution int32 `uniform:"resolution"`
	Model      int32 `uniform:"model"`
	Camera     int32 `uniform:"camera"`
	Projection int32 `uniform:"projection"`
}

func (FrontendPkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b assets.AssetsStorageBuilder) assets.AssetsStorageBuilder {
		b.RegisterAsset(MeshAssetID, func() (assets.StorageAsset, error) {
			verticies := []vbo.Vertex{
				{Pos: [3]float32{0, 0, 0}, TexturePos: [2]float32{0, 0}},
				{Pos: [3]float32{100, 0, 0}, TexturePos: [2]float32{1, 0}},
				{Pos: [3]float32{0, 100, 0}, TexturePos: [2]float32{0, 1}},
				{Pos: [3]float32{100, 100, 0}, TexturePos: [2]float32{1, 1}},
			}
			indicies := []ebo.Index{
				0, 1, 3,
				0, 2, 3,
			}
			asset := mesh.NewMeshStorageAsset(transform.NewSize(100, 100, 0), verticies, indicies)
			return asset, nil
		})
		return b
	})

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b assets.AssetsStorageBuilder) assets.AssetsStorageBuilder {
		b.RegisterAsset(TextureAssetID, func() (assets.StorageAsset, error) {
			asset := texture.NewTextureStorageAsset(bytes.NewBuffer(textureSource))
			return asset, nil
		})
		return b
	})

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b assets.AssetsStorageBuilder) assets.AssetsStorageBuilder {
		b.RegisterAsset(ProgramAssetID, func() (assets.StorageAsset, error) {
			vertSource, err := shaders.ReadFile(filepath.Join(shadersDir, "s.vert"))
			if err != nil {
				return nil, err
			}
			fragSource, err := shaders.ReadFile(filepath.Join(shadersDir, "s.frag"))
			if err != nil {
				return nil, err
			}
			asset := program.NewProgramStorageAsset[Locations](string(vertSource), string(fragSource))
			return asset, nil
		})
		return b
	})

	var t time.Duration

	ioc.WrapService(b, frames.Draw, func(c ioc.Dic, b events.Builder) events.Builder {
		window := ioc.Get[window.Api](c).Window()
		assets := ioc.Get[assets.Assets](c)
		events.Listen(b, func(e frames.FrameEvent) {
			rawMeshAsset, err := assets.Get(MeshAssetID)
			if err != nil {
				panic(err)
			}
			meshAsset, ok := rawMeshAsset.(mesh.MeshCachedAsset)
			if !ok {
				panic("not an expected asset")
			}

			//

			rawTextureAsset, err := assets.Get(TextureAssetID)
			if err != nil {
				panic(err)
			}
			textureAsset, ok := rawTextureAsset.(texture.TextureCachedAsset)
			if !ok {
				panic("not an expected asset")
			}

			//

			rawProgramAsset, err := assets.Get(ProgramAssetID)
			if err != nil {
				panic(err)
			}
			programAsset, ok := rawProgramAsset.(program.ProgramCachedAsset[Locations])
			if !ok {
				panic("not an expected asset")
			}

			//

			t += e.Delta

			transformSize := [3]float32{100, 100, 0}

			// before setting uniforms
			programAsset.Program().Use()
			width, height := window.GetSize()
			{
				// gl.Uniform3f(tools.Locations.Resolution, float32(width), float32(height), 1)
				gl.Uniform3f(programAsset.Program().Locations().Resolution, float32(width), float32(height), 1)
			}
			{
				transformSize := [3]float32{
					transformSize[0],
					transformSize[1] * (1 + float32(t.Seconds())),
					transformSize[2],
				}
				scale := [3]float32{
					transformSize[0] / max(0.1, meshAsset.Size().X),
					transformSize[1] / max(0.1, meshAsset.Size().Y),
				}
				radians := mgl32.DegToRad(float32(t.Seconds()) * 100)
				rotation := mgl32.QuatIdent().
					Mul(mgl32.QuatRotate(radians, mgl32.Vec3{0, 0, 1}))
				matrices := []mgl32.Mat4{
					rotation.Mat4(),
					mgl32.Translate3D(
						0-transformSize[0]/2,
						0-transformSize[1]/2,
						0-transformSize[2]/2),
					mgl32.Scale3D(scale[0], scale[1], scale[2]),
				}
				var model mgl32.Mat4
				for i, matrix := range matrices {
					if i == 0 {
						model = matrix
						continue
					}
					model = model.Mul4(matrix)
				}
				// gl.UniformMatrix4fv(tools.Locations.Model, 1, false, &model[0])
				gl.UniformMatrix4fv(programAsset.Program().Locations().Model, 1, false, &model[0])
			}
			{
				camera := mgl32.Translate3D(0, 0, 0)
				// gl.UniformMatrix4fv(tools.Locations.Camera, 1, false, &camera[0])
				gl.UniformMatrix4fv(programAsset.Program().Locations().Camera, 1, false, &camera[0])
			}
			{
				projection := mgl32.Ortho2D(
					-float32(width)/2,
					float32(width)/2,
					-float32(height)/2,
					float32(height)/2,
				)
				// gl.UniformMatrix4fv(tools.Locations.Projection, 1, false, &projection[0])
				gl.UniformMatrix4fv(programAsset.Program().Locations().Projection, 1, false, &projection[0])
			}

			textureAsset.Texture().Use()
			meshAsset.VAO().Draw()
		})
		return b
	})

	ioc.WrapService(b, appruntime.OrderCleanUp, func(c ioc.Dic, b appruntime.Builder) appruntime.Builder {
		assets := ioc.Get[assets.Assets](c)
		b.OnStop(func(r appruntime.Runtime) {
			assets.Release(
				MeshAssetID,
				TextureAssetID,
				ProgramAssetID,
			)
		})
		return b
	})
}
