package construct

import (
	"engine/modules/grid"
	"engine/services/assets"
	"engine/services/ecs"
)

type ID uint32

//

type IDComponent struct {
	ID ID
}

type CoordsComponent struct {
	Coords grid.Coords
}

func NewID(id ID) IDComponent                      { return IDComponent{id} }
func NewCoords(coords grid.Coords) CoordsComponent { return CoordsComponent{coords} }

//

type Blueprint struct {
	Texture assets.AssetID
	// Size int

	// complexity 1:
	// texture
	// click event
	// size (1x1 or 2x2 for example)

	// complexity 2:
	// healh
	// profits
	// other features like defense
}

func NewBlueprint(texture assets.AssetID) Blueprint {
	return Blueprint{Texture: texture}
}

//

type Service interface {
	RegisterConstruct(ID, Blueprint)

	// adds mesh, texture, mouse click event, adds to grid
	ID() ecs.ComponentsArray[IDComponent]
	Coords() ecs.ComponentsArray[CoordsComponent]
}
