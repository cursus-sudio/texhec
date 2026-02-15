package internal

import (
	"core/modules/construct"
	"core/modules/registry"
	"core/modules/tile"
	"engine"
	"engine/modules/render"
	"engine/modules/transform"
	"engine/services/datastructures"
	"engine/services/ecs"
	"fmt"

	"github.com/ogiusek/ioc/v2"
)

type service struct {
	engine.World `inject:"1"`
	Tile         tile.Service    `inject:"1"`
	GameAssets   registry.Assets `inject:"1"`

	dirtySet ecs.DirtySet

	constructID     ecs.ComponentsArray[construct.IDComponent]
	constructCoords ecs.ComponentsArray[construct.CoordsComponent]

	blueprints datastructures.SparseArray[construct.ID, construct.Blueprint]
}

func NewService(c ioc.Dic) construct.Service {
	s := ioc.GetServices[*service](c)

	s.dirtySet = ecs.NewDirtySet()

	s.constructID = ecs.GetComponentsArray[construct.IDComponent](s)
	s.constructCoords = ecs.GetComponentsArray[construct.CoordsComponent](s)

	s.blueprints = datastructures.NewSparseArray[construct.ID, construct.Blueprint]()

	s.constructID.AddDirtySet(s.dirtySet)
	s.constructCoords.AddDirtySet(s.dirtySet)

	s.Render.Mesh().AddDependency(s.constructID)
	s.Render.Texture().AddDependency(s.constructID)
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
		id, ok := s.constructID.Get(entity)
		if !ok {
			continue
		}
		blueprint, ok := s.blueprints.Get(id.ID)
		if !ok {
			s.Logger.Warn(fmt.Errorf("construct doesn't have registered blueprint with id %v", id))
			continue
		}
		coords, ok := s.constructCoords.Get(entity)
		if !ok {
			continue
		}

		pos := s.Tile.GetPos(coords.Coords)
		pos.Pos[2] += 1
		s.Render.Mesh().Set(entity, render.NewMesh(s.GameAssets.SquareMesh))
		s.Render.Texture().Set(entity, render.NewTexture(blueprint.Texture))

		s.Transform.ParentPivotPoint().Set(entity, transform.NewParentPivotPoint(0, 0, .5))
		s.Transform.Pos().Set(entity, pos)
		s.Transform.Size().Set(entity, s.Tile.GetTileSize())

	}
}

func (s *service) RegisterConstruct(
	id construct.ID,
	blueprint construct.Blueprint,
) {
	s.blueprints.Set(id, blueprint)
}

func (s *service) ID() ecs.ComponentsArray[construct.IDComponent] {
	return s.constructID
}

func (s *service) Coords() ecs.ComponentsArray[construct.CoordsComponent] {
	return s.constructCoords
}
