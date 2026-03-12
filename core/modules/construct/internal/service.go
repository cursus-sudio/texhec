package internal

import (
	"core/modules/construct"
	"core/modules/definitions"
	"core/modules/tile"
	"engine"
	"engine/modules/render"
	"engine/modules/transform"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type service struct {
	engine.World `inject:"1"`
	Tile         tile.Service            `inject:"1"`
	GameAssets   definitions.Definitions `inject:"1"`

	dirtySet ecs.DirtySet

	constructs      ecs.ComponentsArray[construct.ConstructComponent]
	constructCoords ecs.ComponentsArray[construct.CoordsComponent]
}

func NewService(c ioc.Dic) construct.Service {
	s := ioc.GetServices[*service](c)

	s.dirtySet = ecs.NewDirtySet()

	s.constructs = ecs.GetComponentsArray[construct.ConstructComponent](s)
	s.constructCoords = ecs.GetComponentsArray[construct.CoordsComponent](s)

	s.constructCoords.AddDirtySet(s.dirtySet)

	s.Transform.Pos().AddDependency(s.constructCoords)
	s.Transform.Size().AddDependency(s.constructCoords)

	s.Render.Mesh().BeforeGet(s.BeforeGet)
	s.Render.Texture().BeforeGet(s.BeforeGet)
	s.Transform.Pos().BeforeGet(s.BeforeGet)
	s.Transform.Size().BeforeGet(s.BeforeGet)

	return s
}

func (s *service) BeforeGet() {
	for _, entity := range s.dirtySet.Get() {
		construct, ok := s.constructs.Get(entity)
		if !ok {
			continue
		}
		coords, ok := s.constructCoords.Get(entity)
		if !ok {
			continue
		}

		pos := s.Tile.GetPos(coords.Coords)
		pos.Pos[2] += 1
		s.Render.Mesh().Set(entity, render.NewMesh(s.GameAssets.SquareMesh))
		s.Render.Texture().Set(entity, render.NewTexture(construct.Construct))

		s.Transform.ParentPivotPoint().Set(entity, transform.NewParentPivotPoint(0, 0, .5))
		s.Transform.Pos().Set(entity, pos)
		s.Transform.Size().Set(entity, s.Tile.GetTileSize())

	}
}

func (s *service) Construct() ecs.ComponentsArray[construct.ConstructComponent] {
	return s.constructs
}
func (s *service) Coords() ecs.ComponentsArray[construct.CoordsComponent] {
	return s.constructCoords
}
