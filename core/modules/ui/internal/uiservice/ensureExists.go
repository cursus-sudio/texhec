package uiservice

import (
	"core/modules/ui"
	"engine/modules/collider"
	"engine/modules/groups"
	"engine/modules/inputs"
	"engine/modules/layout"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"

	"github.com/go-gl/mathgl/mgl32"
)

func (t *service) EnsureExists() {
mainLoop:
	for _, camera := range t.uiCameraArray.GetEntities() {
		// objects
		// menu
		for _, child := range t.Hierarchy.Children(camera).GetIndices() {
			if _, ok := t.menuArray.Get(child); ok {
				continue mainLoop
			}
		}
		menu := t.NewEntity()
		t.Hierarchy.SetParent(menu, camera)
		t.Transform.ParentPivotPoint().Set(menu, transform.NewParentPivotPoint(1, 1, .5))
		t.Transform.Pos().Set(menu, transform.NewPos(0, 0, 1))
		t.Transform.Size().Set(menu, transform.NewSize(.2, 1, 1))
		t.Transform.PivotPoint().Set(menu, transform.NewPivotPoint(0, 1, .5))

		t.Render.Color().Set(menu, render.NewColor(mgl32.Vec4{1, 1, 1, .5}))
		t.AnimatedBackground().Set(menu, ui.AnimatedBackgroundComponent{})

		t.Groups.Inherit().Set(menu, groups.InheritGroupsComponent{})
		t.Collider.Component().Set(menu, collider.NewCollider(t.GameAssets.SquareCollider))
		t.Inputs.KeepSelected().Set(menu, inputs.KeepSelectedComponent{})
		t.menuArray.Set(menu, menuComponent{})

		// quit btn
		quit := t.NewEntity()

		t.Hierarchy.SetParent(quit, menu)
		t.Groups.Inherit().Set(quit, groups.InheritGroupsComponent{})

		t.Transform.Parent().Set(quit, transform.NewParent(transform.RelativePos))
		t.Transform.ParentPivotPoint().Set(quit, transform.NewParentPivotPoint(1, 1, 1))
		t.Transform.Size().Set(quit, transform.NewSize(25, 25, 1))
		t.Transform.PivotPoint().Set(quit, transform.NewPivotPoint(1, 1, 0))

		t.Text.Content().Set(quit, text.TextComponent{Text: "X"})
		t.Text.FontSize().Set(quit, text.FontSizeComponent{FontSize: 25})
		t.Text.Align().Set(quit, text.TextAlignComponent{Vertical: .5, Horizontal: .5})

		t.Render.Color().Set(quit, render.NewColor(mgl32.Vec4{1, 0, 0, 1}))
		t.Render.Mesh().Set(quit, render.NewMesh(t.GameAssets.SquareMesh))
		t.Render.Texture().Set(quit, render.NewTexture(t.GameAssets.Blank))

		t.Inputs.LeftClick().Set(quit, inputs.NewLeftClick(ui.HideUiEvent{}))
		t.Inputs.KeepSelected().Set(quit, inputs.KeepSelectedComponent{})
		t.Collider.Component().Set(quit, collider.NewCollider(t.GameAssets.SquareCollider))

		// child wrapper
		childWrapper := t.NewEntity()
		t.Hierarchy.SetParent(childWrapper, menu)
		t.Groups.Inherit().Set(childWrapper, groups.InheritGroupsComponent{})
		t.Transform.Parent().Set(childWrapper, transform.NewParent(transform.RelativePos|transform.RelativeSizeXY))

		t.Layout.Order().Set(childWrapper, layout.NewOrder(layout.OrderVectical))
		t.Layout.Align().Set(childWrapper, layout.NewAlign(0, .5))
		t.Layout.Gap().Set(childWrapper, layout.NewGap(10))
		t.childrenWrapperArray.Set(childWrapper, childrenComponent{})
	}
}
