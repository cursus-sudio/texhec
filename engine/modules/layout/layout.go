package layout

import (
	"engine/modules/hierarchy"
	"engine/modules/transform"
	"engine/services/ecs"
)

type ToolFactory ecs.ToolFactory[World, LayoutTool]
type LayoutTool interface {
	// doesn't work for nested layout
	Layout() Interface
}
type World interface {
	ecs.World
	transform.TransformTool
	hierarchy.HierarchyTool
}

type Interface interface {
	Align() ecs.ComponentsArray[AlignComponent]
	Order() ecs.ComponentsArray[OrderComponent]
	// Wrap() ecs.ComponentsArray[WrapComponent]
	Gap() ecs.ComponentsArray[GapComponent]
}

// all components are components for parent:

// centering
// Y axis is reversed for primary axis
type AlignComponent struct {
	// value between 0 and 1 where 0 means aligned to left and 1 aligned to right
	Primary, Secondary float32 // default is 0
}

func NewAlign(primary, secondary float32) AlignComponent {
	return AlignComponent{primary, secondary}
}

// order
type Order uint8

const (
	OrderHorizontal Order = iota
	OrderVectical
)

type OrderComponent struct {
	Order Order // default horizontal
}

func (order *OrderComponent) Primary() Order   { return order.Order }
func (order *OrderComponent) Secondary() Order { return 1 - order.Order }

func NewOrder(order Order) OrderComponent {
	return OrderComponent{order}
}

// wrapping
// type WrapComponent struct {
// 	Wrap bool
// }
//
// func NewWrap(wrap bool) WrapComponent {
// 	return WrapComponent{wrap}
// }

// gaps
type GapComponent struct {
	Gap float32
}

func NewGap(x float32) GapComponent {
	return GapComponent{x}
}
