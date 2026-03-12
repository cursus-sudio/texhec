package construct

import (
	"engine/modules/grid"
	"engine/services/ecs"
)

//

type ConstructComponent struct {
	Construct ecs.EntityID
}

func NewConstruct(construct ecs.EntityID) ConstructComponent { return ConstructComponent{construct} }

//

type CoordsComponent struct {
	Coords grid.Coords
}

func NewCoords(coords grid.Coords) CoordsComponent { return CoordsComponent{coords} }

//

// type BlueprintComponent struct {
// 	Construct string
// 	// Size int
//
// 	// complexity 1:
// 	// texture
// 	// click event
// 	// size (1x1 or 2x2 for example)
//
// 	// complexity 2:
// 	// healh
// 	// profits
// 	// other features like defense
// }
//
// func NewBlueprint(construct string) BlueprintComponent {
// 	return BlueprintComponent{construct}
// }

//

type Service interface {
	// adds mesh, texture, mouse click event, adds to grid
	Construct() ecs.ComponentsArray[ConstructComponent]
	Coords() ecs.ComponentsArray[CoordsComponent]
}
