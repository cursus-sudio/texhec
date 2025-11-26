package uimodule

import (
	gameassets "core/assets"
	"core/modules/ui"
	"engine/modules/animation"
	"engine/modules/collider"
	"engine/modules/genericrenderer"
	"engine/modules/groups"
	"engine/modules/hierarchy"
	"engine/modules/inputs"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
	"errors"
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

func (s *uiSys) EnsureInit() {
	if s.menu != nil {
		return
	}
	// start transactions
	s.transformTransaction = s.transformTool.Transaction()
	s.renderTransaction = s.renderTool.Transaction()
	s.textTransaction = s.textTool.Transaction()
	s.hierarchyTransaction = s.hierarchyTool.Transaction()

	s.animationTransaction = s.animationArray.Transaction()
	s.groupInheritTransaction = s.groupInheritArray.Transaction()
	s.pipelineTransaction = s.pipelineArray.Transaction()
	s.colliderTransaction = s.colliderArray.Transaction()
	s.leftClickTransaction = s.leftClickArray.Transaction()
	s.keepSelectedTransaction = s.keepSelectedArray.Transaction()

	// initialize main entities
	menu := s.world.NewEntity()
	menuTransform := s.transformTransaction.GetObject(menu)
	menuTransform.ParentPivotPoint().Set(transform.NewParentPivotPoint(1, 1, .5))
	menuTransform.Size().Set(transform.NewSize(.2, 1, 1))
	menuTransform.PivotPoint().Set(transform.NewPivotPoint(1, 1, .5))

	menuRender := s.renderTransaction.GetObject(menu)
	menuRender.Color().Set(render.NewColor(mgl32.Vec4{1, 1, 1, .5}))
	menuRender.Mesh().Set(render.NewMesh(gameassets.SquareMesh))
	menuRender.Texture().Set(render.NewTexture(gameassets.WaterTileTextureID))
	s.pipelineTransaction.SaveComponent(menu, genericrenderer.PipelineComponent{})

	s.groupInheritTransaction.SaveComponent(menu, groups.InheritGroupsComponent{})
	s.colliderTransaction.SaveComponent(menu, collider.NewCollider(gameassets.SquareColliderID))

	menuText := s.textTransaction.GetObject(menu)
	menuText.FontSize().Set(text.FontSizeComponent{FontSize: 32})
	menuText.TextAlign().Set(text.TextAlignComponent{Vertical: .5, Horizontal: .5})

	childrenContainer := s.world.NewEntity()
	s.hierarchyTransaction.GetObject(childrenContainer).Parent().Set(hierarchy.NewParent(menu))
	childrenContainerTransform := s.transformTransaction.GetObject(childrenContainer)
	childrenContainerTransform.Parent().Set(transform.NewParent(transform.RelativePos | transform.RelativeSize))

	quit := s.world.NewEntity()
	s.hierarchyTransaction.GetObject(quit).Parent().Set(hierarchy.NewParent(menu))
	quitTransform := s.transformTransaction.GetObject(quit)
	quitTransform.Parent().Set(transform.NewParent(transform.RelativePos))
	quitTransform.ParentPivotPoint().Set(transform.NewParentPivotPoint(1, 1, .5))
	quitTransform.Size().Set(transform.NewSize(25, 25, 2))
	quitTransform.PivotPoint().Set(transform.NewPivotPoint(1, 1, .5))

	quitText := s.textTransaction.GetObject(quit)
	quitText.Text().Set(text.TextComponent{Text: "X"})
	quitText.FontSize().Set(text.FontSizeComponent{FontSize: 25})
	quitText.TextAlign().Set(text.TextAlignComponent{Vertical: .5, Horizontal: .5})

	quitRender := s.renderTransaction.GetObject(quit)
	quitRender.Color().Set(render.NewColor(mgl32.Vec4{1, 0, 0, 1}))
	quitRender.Mesh().Set(render.NewMesh(gameassets.SquareMesh))
	quitRender.Texture().Set(render.NewTexture(gameassets.WaterTileTextureID))
	s.pipelineTransaction.SaveComponent(quit, genericrenderer.PipelineComponent{})
	s.groupInheritTransaction.SaveComponent(quit, groups.InheritGroupsComponent{})

	s.leftClickTransaction.SaveComponent(quit, inputs.NewMouseLeftClick(ui.UnselectEvent{}))
	s.keepSelectedTransaction.SaveComponent(quit, inputs.KeepSelectedComponent{})
	s.colliderTransaction.SaveComponent(quit, collider.NewCollider(gameassets.SquareColliderID))

	s.menu = &menuData{
		menu:              menu,
		childrenContainer: childrenContainer,
		visible:           false,
	}
}

func (s *uiSys) Render() error {
	s.EnsureInit()
	for _, entity := range s.hierarchyTransaction.
		GetObject(s.menu.childrenContainer).
		FlatChildren().GetIndices() {
		s.world.RemoveEntity(entity)
	}

	if s.currentState == nil {
		if s.menu.visible {
			s.menu.visible = false
			s.animationTransaction.SaveComponent(s.menu.menu, animation.NewAnimationComponent(gameassets.HideMenuAnimation, s.animationDuration))
		}
		// mark parrent as hidden
		return s.Flush()
	}

	state := *s.currentState

	cameras := s.uiCameraArray.GetEntities()
	if len(cameras) != 1 {
		return errors.New("expected one ui camera")
	}
	camera := cameras[0]

	if !s.menu.visible {
		s.menu.visible = true
		s.animationTransaction.SaveComponent(s.menu.menu, animation.NewAnimationComponent(gameassets.ShowMenuAnimation, s.animationDuration))
	}

	menu := s.menu.menu
	s.hierarchyTransaction.GetObject(menu).Parent().Set(hierarchy.NewParent(camera))
	s.transformTransaction.GetObject(menu).Parent().Set(transform.NewParent(transform.RelativePos | transform.RelativeSize))

	menuText := s.textTransaction.GetObject(menu)
	menuText.Text().Set(text.TextComponent{Text: fmt.Sprintf("pos is %v", state.Tile)})
	menuText.FontSize().Set(text.FontSizeComponent{FontSize: 32})
	menuText.TextAlign().Set(text.TextAlignComponent{Vertical: .5, Horizontal: .5})

	return s.Flush()
}
