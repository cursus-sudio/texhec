package gamescene

import (
	"core/modules/construct"
	"core/modules/generation"
	"core/modules/registry"
	"core/modules/settings"
	"core/modules/ui"
	gamescenes "core/scenes"
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/grid"
	"engine/modules/groups"
	"engine/modules/inputs"
	"engine/modules/render"
	"engine/modules/seed"
	"engine/modules/transform"
	"engine/modules/uuid"
	"engine/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

const (
	UiGroup groups.Group = iota + 1
	GameGroup
	BgGroup
)

func addScene(
	world gamescenes.World,
	sceneParent ecs.EntityID,
	isServer bool,
) {
	// biggest maps on mods in rusted warfare 2560x1440
	// - all tiles are rendered at once
	// - strategic map is used at some point
	// biggest zoom out in factorio is 448x256 (in 4k)
	rows := 1000
	cols := 1000

	{
		uiCamera := world.NewEntity()
		world.Hierarchy.SetParent(uiCamera, sceneParent)
		world.Camera.Priority().Set(uiCamera, camera.NewPriority(1))
		world.Camera.Ortho().Set(uiCamera, camera.NewOrtho(-1000, +1000))
		world.Groups.Component().Set(uiCamera, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())
		world.Ui.UiCamera().Set(uiCamera, ui.UiCameraComponent{})
		world.Ui.CursorCamera().Set(uiCamera, ui.CursorCameraComponent{})

		settingsEntity := world.NewEntity()
		world.Hierarchy.SetParent(settingsEntity, uiCamera)
		world.Transform.Pos().Set(settingsEntity, transform.NewPos(10, -10, 0))
		world.Transform.Size().Set(settingsEntity, transform.NewSize(50, 50, 1))
		world.Transform.PivotPoint().Set(settingsEntity, transform.NewPivotPoint(0, 1, .5))
		world.Transform.Parent().Set(settingsEntity, transform.NewParent(transform.RelativePos))
		world.Transform.ParentPivotPoint().Set(settingsEntity, transform.NewParentPivotPoint(0, 1, .5))
		world.Groups.Component().Set(settingsEntity, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())

		world.Render.Mesh().Set(settingsEntity, render.NewMesh(world.GameAssets.SquareMesh))
		world.Render.Texture().Set(settingsEntity, render.NewTexture(world.GameAssets.Hud.Settings))

		world.Inputs.LeftClick().Set(settingsEntity, inputs.NewLeftClick(settings.EnterSettingsEvent{}))
		world.Inputs.KeepSelected().Set(settingsEntity, inputs.KeepSelectedComponent{})
		world.Collider.Component().Set(settingsEntity, collider.NewCollider(world.GameAssets.SquareCollider))
	}

	{
		bgCamera := world.NewEntity()
		world.Hierarchy.SetParent(bgCamera, sceneParent)
		world.Camera.Priority().Set(bgCamera, camera.NewPriority(-1))
		world.Camera.Ortho().Set(bgCamera, camera.NewOrtho(-1000, +1000))
		world.Groups.Component().Set(bgCamera, groups.EmptyGroups().Ptr().Enable(BgGroup).Val())

		bg := world.NewEntity()
		world.Hierarchy.SetParent(bg, bgCamera)
		world.Transform.Parent().Set(bg, transform.NewParent(transform.RelativePos|transform.RelativeSizeXY))
		world.Groups.Inherit().Set(bg, groups.InheritGroupsComponent{})
		world.Ui.AnimatedBackground().Set(bg, ui.AnimatedBackgroundComponent{})
	}

	gameCamera := world.NewEntity()
	world.Hierarchy.SetParent(gameCamera, sceneParent)
	world.UUID.Component().Set(gameCamera, uuid.New([16]byte{48}))
	world.Camera.Ortho().Set(gameCamera, camera.NewOrtho(-1000, +1000))
	world.Groups.Component().Set(gameCamera, groups.EmptyGroups().Ptr().Enable(GameGroup).Val())
	world.Camera.Mobile().Set(gameCamera, camera.NewMobileCamera())
	world.Camera.Limits().Set(gameCamera, camera.NewCameraLimits(
		mgl32.Vec3{50 * float32(-cols), 50 * float32(-rows), -1000},
		mgl32.Vec3{50 * float32(cols), 50 * float32(rows), 1000},
	))

	if isServer {
		gridEntity := world.NewEntity()

		world.Hierarchy.SetParent(gridEntity, sceneParent)
		world.Groups.Component().Set(gridEntity, groups.EmptyGroups().Ptr().Enable(GameGroup).Val())

		task := world.Generation.Generate(generation.NewConfig(
			gridEntity,
			seed.New(world.Clock.Now().Unix()),
			// seed.New(21377137),
			grid.NewCoords(cols, rows),
		))
		world.Batcher.Queue(task)

		farm := world.NewEntity()
		world.Hierarchy.SetParent(farm, gridEntity)
		world.Groups.Component().Set(farm, groups.EmptyGroups().Ptr().Enable(GameGroup).Val())

		world.Construct.ID().Set(farm, construct.NewID(registry.ConstructFarm))
		world.Construct.Coords().Set(farm, construct.NewCoords(grid.NewCoords(499, 500)))
	}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.GameBuilder {
		return func(sceneParent ecs.EntityID) {
			world := ioc.GetServices[gamescenes.World](c)
			addScene(
				world,
				sceneParent,
				true, // is server
			)
		}
	})
}
