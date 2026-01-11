package tilecollider

import (
	"core/modules/tile"
	"engine"
	"engine/modules/collider"
	"engine/modules/groups"
	"engine/modules/inputs"
	"engine/modules/transform"
	"engine/modules/uuid"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type system struct {
	engine.World `inject:"1"`
	Tile         tile.Service `inject:"1"`
}

func TileColliderSystem(c ioc.Dic,
	tileSize int32, // transform
	gridDepth float32,
	tileGroups groups.GroupsComponent, // groups
	colliderComponent collider.Component, // collider
) tile.System {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[system](c)
		tilePosDirtySet := ecs.NewDirtySet()
		s.Tile.Pos().AddDirtySet(tilePosDirtySet)
		s.UUID.Component().BeforeGet(func() {
			ei := tilePosDirtySet.Get()
			if len(ei) == 0 {
				return
			}
			for _, entity := range ei {
				if _, ok := s.UUID.Component().Get(entity); ok {
					continue
				}
				comp := uuid.New(s.UUID.NewUUID())
				s.UUID.Component().Set(entity, comp)
			}
		})

		//

		tileDirtySet := ecs.NewDirtySet()
		s.Tile.Pos().AddDirtySet(tileDirtySet)
		applyTileCollider := func() {
			ei := tileDirtySet.Get()
			if len(ei) == 0 {
				return
			}

			for _, entity := range ei {

				pos, ok := s.Tile.Pos().Get(entity)
				if !ok {
					continue
				}
				transformPos := transform.NewPos(
					float32(tileSize)*float32(pos.X)+float32(tileSize)/2,
					float32(tileSize)*float32(pos.Y)+float32(tileSize)/2,
					gridDepth+float32(pos.Layer),
				)
				s.Transform.Pos().Set(entity, transformPos)
				s.Transform.Size().Set(entity, transform.NewSize(float32(tileSize), float32(tileSize), 1))
				s.Inputs.LeftClick().Set(entity, inputs.NewLeftClick(tile.NewTileClickEvent(pos)))
				s.Inputs.Stack().Set(entity, inputs.StackComponent{})
				s.Collider.Component().Set(entity, colliderComponent)
				s.Groups.Component().Set(entity, tileGroups)
			}
		}

		s.Transform.Size().BeforeGet(applyTileCollider)
		s.Transform.Pos().BeforeGet(applyTileCollider)
		s.Collider.Component().BeforeGet(applyTileCollider)
		s.Inputs.LeftClick().BeforeGet(applyTileCollider)
		s.Inputs.Stack().BeforeGet(applyTileCollider)
		s.Groups.Component().BeforeGet(applyTileCollider)

		return nil
	})
}
