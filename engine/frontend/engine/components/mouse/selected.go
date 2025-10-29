package mouse

type Hovered struct{}

func NewHovered() Hovered { return Hovered{} }

//

type Dragged struct{}

func NewDragged() Dragged { return Dragged{} }
