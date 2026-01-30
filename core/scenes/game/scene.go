package gamescene

import (
	gameassets "core/assets"
	"core/modules/generation"
	"core/modules/settings"
	"core/modules/ui"
	gamescenes "core/scenes"
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/connection"
	"engine/modules/grid"
	"engine/modules/groups"
	"engine/modules/inputs"
	"engine/modules/netsync"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/modules/uuid"
	"engine/services/ecs"
	"errors"

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
)

func addScene(
	world gamescenes.World,
	sceneParent ecs.EntityID,
	gameAssets gameassets.GameAssets,
	isServer bool,
) {
	// biggest maps on mods in rusted warfare 2560x1440
	// - all tiles are rendered at once
	// - strategic map is used at some point
	// biggest zoom out in factorio is 448x256 (in 4k)
	rows := 1000
	cols := 1000

	uiCamera := world.NewEntity()
	world.Hierarchy.SetParent(uiCamera, sceneParent)
	world.Camera.Ortho().Set(uiCamera, camera.NewOrtho(-1000, +1000))
	world.Groups.Component().Set(uiCamera, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())
	world.Ui.UiCamera().Set(uiCamera, ui.UiCameraComponent{})

	gameCamera := world.NewEntity()
	world.Hierarchy.SetParent(gameCamera, sceneParent)
	world.UUID.Component().Set(gameCamera, uuid.New([16]byte{48}))
	world.Camera.Ortho().Set(gameCamera, camera.NewOrtho(-1000, +1000))
	world.Groups.Component().Set(gameCamera, groups.EmptyGroups().Ptr().Enable(GameGroup).Val())
	world.Camera.Mobile().Set(gameCamera, camera.NewMobileCamera())
	world.Camera.Limits().Set(gameCamera, camera.NewCameraLimits(
		mgl32.Vec3{50 * float32(-rows), 50 * float32(-cols), -1000},
		mgl32.Vec3{50 * float32(rows), 50 * float32(cols), 1000},
	))

	signature := world.NewEntity()
	world.Hierarchy.SetParent(signature, uiCamera)
	world.Transform.Pos().Set(signature, transform.NewPos(5, 5, 1))
	world.Transform.Size().Set(signature, transform.NewSize(100, 50, 1))
	world.Transform.PivotPoint().Set(signature, transform.NewPivotPoint(0, .5, 0))
	world.Transform.Parent().Set(signature, transform.NewParent(transform.RelativePos))
	world.Transform.ParentPivotPoint().Set(signature, transform.NewParentPivotPoint(0, 0, .5))
	world.Groups.Component().Set(signature, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())

	world.Text.FontSize().Set(signature, text.FontSizeComponent{FontSize: 32})
	world.Text.Break().Set(signature, text.BreakComponent{Break: text.BreakNone})

	settingsEntity := world.NewEntity()
	world.Hierarchy.SetParent(settingsEntity, uiCamera)
	world.Transform.Pos().Set(settingsEntity, transform.NewPos(10, -10, 5))
	world.Transform.Size().Set(settingsEntity, transform.NewSize(50, 50, 1))
	world.Transform.PivotPoint().Set(settingsEntity, transform.NewPivotPoint(0, 1, .5))
	world.Transform.Parent().Set(settingsEntity, transform.NewParent(transform.RelativePos))
	world.Transform.ParentPivotPoint().Set(settingsEntity, transform.NewParentPivotPoint(0, 1, .5))
	world.Groups.Component().Set(settingsEntity, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())

	world.Render.Mesh().Set(settingsEntity, render.NewMesh(gameAssets.SquareMesh))
	world.Render.Texture().Set(settingsEntity, render.NewTexture(gameAssets.Hud.Settings))

	world.Inputs.LeftClick().Set(settingsEntity, inputs.NewLeftClick(settings.EnterSettingsEvent{}))
	world.Inputs.KeepSelected().Set(settingsEntity, inputs.KeepSelectedComponent{})
	world.Collider.Component().Set(settingsEntity, collider.NewCollider(gameAssets.SquareCollider))

	if isServer {
		gridEntity := world.NewEntity()

		world.Hierarchy.SetParent(gridEntity, sceneParent)
		world.Transform.Size().Set(gridEntity, transform.NewSize(float32(cols)*100, float32(rows)*100, 1))
		world.Groups.Component().Set(gridEntity, groups.EmptyGroups().Ptr().Enable(GameGroup).Val())

		world.Collider.Component().Set(gridEntity, collider.NewCollider(gameAssets.SquareCollider))
		world.Inputs.Stack().Set(gridEntity, inputs.StackComponent{})

		task := world.Generation.Generate(generation.NewConfiguration(
			gridEntity, 21377137,
			grid.NewCoords(cols, rows),
		))
		world.Batcher.Queue(task)
	}

	if isServer {
		listenerEntity := world.NewEntity()
		world.Hierarchy.SetParent(listenerEntity, sceneParent)
		listener, err := world.Connection.Host(":8000", func(cc connection.ConnectionComponent) {
			entity := world.NewEntity()
			world.NetSync.Client().Set(entity, netsync.ClientComponent{})
			world.Connection.Component().Set(entity, cc)
		})
		if err != nil {
			world.Logger.Warn(err)
			return
		}
		world.Connection.Listener().Set(listenerEntity, listener)
	} else {
		comp, err := world.Connection.Connect(":8000")
		if err != nil {
			world.Logger.Warn(errors.New("there is no server"))
		}
		entity := world.NewEntity()
		world.Hierarchy.SetParent(entity, sceneParent)
		world.NetSync.Server().Set(entity, netsync.ServerComponent{})
		world.Connection.Component().Set(entity, comp)
	}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.GameBuilder {
		return func(sceneParent ecs.EntityID) {
			world := ioc.GetServices[gamescenes.World](c)
			addScene(
				world,
				sceneParent,
				ioc.Get[gameassets.GameAssets](c),
				true, // is server
			)
		}
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.GameClientBuilder {
		return func(sceneParent ecs.EntityID) {
			world := ioc.GetServices[gamescenes.World](c)
			addScene(
				world,
				sceneParent,
				ioc.Get[gameassets.GameAssets](c),
				false, // is server
			)
		}
	})
}
