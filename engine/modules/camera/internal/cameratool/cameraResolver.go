package cameratool

import (
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/media/window"

	"github.com/go-gl/mathgl/mgl32"
)

// type cameraDataID
type projectionID uint8

type projectionComponent struct {
	projectionID
}

type tool struct {
	camera.World

	cameraArray      ecs.ComponentsArray[camera.Component]
	projectionsArray ecs.ComponentsArray[projectionComponent]

	*toolFactory
	projections datastructures.SparseArray[projectionID, ProjectionData]

	dirtySet ecs.DirtySet

	mobileCamera       ecs.ComponentsArray[camera.MobileCameraComponent]
	cameraLimits       ecs.ComponentsArray[camera.CameraLimitsComponent]
	viewport           ecs.ComponentsArray[camera.ViewportComponent]
	normalizedViewport ecs.ComponentsArray[camera.NormalizedViewportComponent]

	ortho              ecs.ComponentsArray[camera.OrthoComponent]
	orthoResolution    ecs.ComponentsArray[camera.OrthoResolutionComponent]
	perspective        ecs.ComponentsArray[camera.PerspectiveComponent]
	dynamicPerspective ecs.ComponentsArray[camera.DynamicPerspectiveComponent]
}

func (t tool) Camera() camera.Interface { return t }

func (t tool) Component() ecs.ComponentsArray[camera.Component] {
	return t.cameraArray
}

func (t tool) Mobile() ecs.ComponentsArray[camera.MobileCameraComponent] {
	return t.mobileCamera
}
func (t tool) Limits() ecs.ComponentsArray[camera.CameraLimitsComponent] {
	return t.cameraLimits
}
func (t tool) Viewport() ecs.ComponentsArray[camera.ViewportComponent] {
	return t.viewport
}
func (t tool) NormalizedViewport() ecs.ComponentsArray[camera.NormalizedViewportComponent] {
	return t.normalizedViewport
}

func (t tool) Ortho() ecs.ComponentsArray[camera.OrthoComponent] {
	return t.ortho
}
func (t tool) OrthoResolution() ecs.ComponentsArray[camera.OrthoResolutionComponent] {
	return t.orthoResolution
}
func (t tool) Perspective() ecs.ComponentsArray[camera.PerspectiveComponent] {
	return t.perspective
}
func (t tool) DynamicPerspective() ecs.ComponentsArray[camera.DynamicPerspectiveComponent] {
	return t.dynamicPerspective
}

//

func (t tool) GetViewport(entity ecs.EntityID) (x, y, w, h int32) {
	viewportComponent, ok := t.viewport.Get(entity)
	if ok {
		return viewportComponent.Viewport()
	}
	normalizedViewportComponent, ok := t.normalizedViewport.Get(entity)
	if ok {
		return normalizedViewportComponent.Viewport(t.window.Window().GetSize())
	}

	w, h = t.window.Window().GetSize()
	return 0, 0, w, h
}
func (t tool) Mat4(entity ecs.EntityID) mgl32.Mat4 {
	comp, ok := t.projectionsArray.Get(entity)
	if !ok {
		return mgl32.Mat4{}
	}
	data, ok := t.projections.Get(comp.projectionID)
	if !ok {
		return mgl32.Mat4{}
	}
	return data.Mat4(entity)
}
func (t tool) ShootRay(camera ecs.EntityID, mousePos window.MousePos) collider.Ray {
	comp, ok := t.projectionsArray.Get(camera)
	if !ok {
		return collider.Ray{}
	}
	data, ok := t.projections.Get(comp.projectionID)
	if !ok {
		return collider.Ray{}
	}

	ray := data.ShootRay(camera, mousePos)
	groups, _ := t.Groups().Component().Get(camera)
	ray.Groups = groups
	return ray
}

//

func (t tool) BeforeGet() {
	dirtyEntities := t.dirtySet.Get()
	if len(dirtyEntities) == 0 {
		return
	}

	for _, entity := range dirtyEntities {
		cam, ok := t.cameraArray.Get(entity)
		if !ok {
			t.projectionsArray.Remove(entity)
			continue
		}
		projID, ok := t.projectionIDs[cam.Projection]
		if !ok {
			t.projectionsArray.Remove(entity)
			continue
		}
		projComp := projectionComponent{projID}
		t.projectionsArray.Set(entity, projComp)
	}
}
