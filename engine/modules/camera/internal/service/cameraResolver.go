package service

import (
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/groups"
	"engine/modules/transform"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/media/window"
	"reflect"

	"github.com/go-gl/mathgl/mgl32"
)

// extra data

type Service interface {
	Register(
		reflect.Type,
		ProjectionData,
	)
	camera.Service
}

type ProjectionData struct {
	Mat4     func(entity ecs.EntityID) mgl32.Mat4
	ShootRay func(entity ecs.EntityID, mousePos window.MousePos) collider.Ray
}

//

// type cameraDataID
type projectionID uint8

type projectionComponent struct {
	projectionID
}

type service struct {
	world     ecs.World
	transform transform.Service
	groups    groups.Service

	cameraArray      ecs.ComponentsArray[camera.Component]
	projectionsArray ecs.ComponentsArray[projectionComponent]

	window        window.Api
	projectionIDs map[reflect.Type]projectionID
	projections   datastructures.SparseArray[projectionID, ProjectionData]

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

func NewSerivce(
	world ecs.World,
	transform transform.Service,
	groups groups.Service,
	window window.Api,
) Service {
	s := &service{
		world,
		transform,
		groups,

		ecs.GetComponentsArray[camera.Component](world),
		ecs.GetComponentsArray[projectionComponent](world),

		window,
		make(map[reflect.Type]projectionID),
		datastructures.NewSparseArray[projectionID, ProjectionData](),
		ecs.NewDirtySet(),

		ecs.GetComponentsArray[camera.MobileCameraComponent](world),
		ecs.GetComponentsArray[camera.CameraLimitsComponent](world),
		ecs.GetComponentsArray[camera.ViewportComponent](world),
		ecs.GetComponentsArray[camera.NormalizedViewportComponent](world),

		ecs.GetComponentsArray[camera.OrthoComponent](world),
		ecs.GetComponentsArray[camera.OrthoResolutionComponent](world),
		ecs.GetComponentsArray[camera.PerspectiveComponent](world),
		ecs.GetComponentsArray[camera.DynamicPerspectiveComponent](world),
	}

	s.projectionsArray.BeforeGet(s.BeforeGet)
	s.cameraArray.AddDirtySet(s.dirtySet)
	return s
}

func (t *service) Component() ecs.ComponentsArray[camera.Component] {
	return t.cameraArray
}

func (t *service) Mobile() ecs.ComponentsArray[camera.MobileCameraComponent] {
	return t.mobileCamera
}
func (t *service) Limits() ecs.ComponentsArray[camera.CameraLimitsComponent] {
	return t.cameraLimits
}
func (t *service) Viewport() ecs.ComponentsArray[camera.ViewportComponent] {
	return t.viewport
}
func (t *service) NormalizedViewport() ecs.ComponentsArray[camera.NormalizedViewportComponent] {
	return t.normalizedViewport
}

func (t *service) Ortho() ecs.ComponentsArray[camera.OrthoComponent] {
	return t.ortho
}
func (t *service) OrthoResolution() ecs.ComponentsArray[camera.OrthoResolutionComponent] {
	return t.orthoResolution
}
func (t *service) Perspective() ecs.ComponentsArray[camera.PerspectiveComponent] {
	return t.perspective
}
func (t *service) DynamicPerspective() ecs.ComponentsArray[camera.DynamicPerspectiveComponent] {
	return t.dynamicPerspective
}

//

func (t *service) GetViewport(entity ecs.EntityID) (x, y, w, h int32) {
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
func (t *service) Mat4(entity ecs.EntityID) mgl32.Mat4 {
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
func (t *service) ShootRay(camera ecs.EntityID, mousePos window.MousePos) collider.Ray {
	comp, ok := t.projectionsArray.Get(camera)
	if !ok {
		return collider.Ray{}
	}
	data, ok := t.projections.Get(comp.projectionID)
	if !ok {
		return collider.Ray{}
	}

	ray := data.ShootRay(camera, mousePos)
	groups, _ := t.groups.Component().Get(camera)
	ray.Groups = groups
	return ray
}

//

func (t *service) Register(
	componentType reflect.Type,
	data ProjectionData,
) {
	if _, ok := t.projectionIDs[componentType]; ok {
		return
	}
	i := projectionID(len(t.projections.GetIndices()))
	t.projectionIDs[componentType] = i
	t.projections.Set(i, data)
}

//

func (t *service) BeforeGet() {
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
